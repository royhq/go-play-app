package create

import "time"

type UserID string

type NewUser struct {
	Name string
	Age  int
}

type User struct {
	ID      UserID
	Name    string
	Age     int
	Created time.Time
}
