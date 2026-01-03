package databases

import (
	"github.com/panbeh/otp-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedisStandAlone(redisAddr config.RedisAddr) redis.UniversalClient {
	opt, err := redis.ParseURL(string(redisAddr))
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}
