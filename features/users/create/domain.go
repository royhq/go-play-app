package create

import "time"

type UserID string

type User struct {
	ID      UserID
	Name    string
	Age     int
	Created time.Time
}
