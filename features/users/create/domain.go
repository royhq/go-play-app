package create

import (
	"time"

	"github.com/royhq/go-play-app/shared/domain"
)

type User struct {
	ID      domain.UserID
	Name    string
	Age     int
	Created time.Time
}

type CreatedUserEvent struct {
	Date   time.Time `json:"date"`
	UserID string    `json:"user_id"`
}
