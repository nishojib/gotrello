package models

import (
	"nishojib/gotrello/internal/validator"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name        string
	Description string
	CreatedAt   time.Time `bun:"default:current_timestamp"`
	DeletedAt   time.Time `bun:",soft_delete"`
	SortOrder   float64
	Version     int
	StatusID    uuid.UUID `bun:"type:uuid"`
}

func (task *Task) Validate(validator *validator.Validator) {
	validator.Check(task.Name != "", "name", "Name must be provided")
	validator.Check(len(task.Name) <= 500, "name", "Name must not be more than 500 bytes long")
}
