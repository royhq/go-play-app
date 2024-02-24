package uuid

import "github.com/google/uuid"

func New() string {
	return uuid.NewString()
}

func IsUUID(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}
