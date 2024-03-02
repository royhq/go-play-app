package get

import "github.com/royhq/go-play-app/shared/domain"

type Query struct {
	UserID domain.UserID
}

type QueryOutput struct {
	User User
}
