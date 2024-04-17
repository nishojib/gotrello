package models

import (
	"nishojib/gotrello/internal/validator"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name      string
	CreatedAt time.Time `bun:"default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`
	Version   int
}

func (project *Project) Validate(validator *validator.Validator) {
	validator.Check(project.Name != "", "title", "must be provided")
	validator.Check(len(project.Name) <= 100, "title", "must not be more than 100 bytes long")
}
