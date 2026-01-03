package business

import "time"

type Business struct {
	ID        string
	Name      string
	Token     string
	CreatedAt time.Time
}
