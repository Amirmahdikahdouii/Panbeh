package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/panbeh/otp-backend/internal/config"
	"github.com/panbeh/otp-backend/internal/domain/business"
	"github.com/panbeh/otp-backend/internal/domain/otp"
	businessRepo "github.com/panbeh/otp-backend/internal/repository/businessRepo"
	"github.com/panbeh/otp-backend/internal/repository/databases"
	oTPRepo "github.com/panbeh/otp-backend/internal/repository/otpRepo"
	"github.com/panbeh/otp-backend/internal/service"
	transport "github.com/panbeh/otp-backend/internal/transport/http"
	loggerPkg "github.com/panbeh/otp-backend/pkg/logger"
)

func main() {
	ctx := context.Background()

	err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configs: %v", err)
	}

	logger := loggerPkg.New(config.GetLogLevel())
	logger.Info("starting application")

	postgresDB := databases.ConnectPostgres(config.GetPostgres())
	redisDB := databases.ConnectRedisStandAlone(config.GetRedis().RedisAddr)

	businessRepo := businessRepo.NewBusinessRepository(postgresDB)
	otpRepo := oTPRepo.NewOTPRepository(redisDB)

	businessDomainSvc := business.NewService(business.ServiceConfig{})
	// TODO: Should we have a unique type for OTP TTL? if yes, store it in config or domain?
	otpTTL, err := otp.NewCodeTTL(config.GetOTPTTL())
	if err != nil {
		panic("OTP TTL is not valid")
	}

	otpDomainSvc := otp.NewService(otp.ServiceConfig{
		TTL: otpTTL,
	})

	businessAppSvc := service.NewBusinessAppService(businessRepo, businessDomainSvc)
	otpSender := service.NewLogOTPSender(logger)
	otpAppSvc := service.NewOTPAppService(otpRepo, otpDomainSvc, otpSender)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(loggerPkg.EchoMiddleware(logger))

	router := transport.NewRouter(
		logger,
		transport.RouterDeps{
			BusinessService: businessAppSvc,
			OTPService:      otpAppSvc,
			AuthResolver:    businessAppSvc,
		},
	)
	router.Register(e)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.GetRestServerAddr()),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Note: Echo's StartServer blocks, so we run it in a goroutine and handle shutdown via signals.
	go func() {
		if err := e.StartServer(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", slog.Any("err", err))
	}
}
