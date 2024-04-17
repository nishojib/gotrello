package models

import "github.com/google/uuid"

type contextKey string

const UserContextKey = contextKey("user")

type AuthenticatedUser struct {
	ID          uuid.UUID
	Email       string
	LoggedIn    bool
	AccessToken string
}
