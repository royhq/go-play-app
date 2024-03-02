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
