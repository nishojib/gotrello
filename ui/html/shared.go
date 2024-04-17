package ui

import (
	"context"
	"nishojib/gotrello/internal/data/models"
)

func GetAuthenticatedUser(ctx context.Context) models.AuthenticatedUser {
	user, ok := ctx.Value(models.UserContextKey).(models.AuthenticatedUser)
	if !ok {
		return models.AuthenticatedUser{}
	}

	return user
}
