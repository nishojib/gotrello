package models

import (
	"nishojib/gotrello/internal/validator"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Status struct {
	bun.BaseModel `bun:"table:statuses"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name      string
	CreatedAt time.Time `bun:"default:current_timestamp"`
	DeletedAt time.Time `bun:",soft_delete"`
	Version   int
	ProjectID uuid.UUID `bun:"type:uuid"`
	Project   Project   `bun:"rel:belongs-to,join:project_id=id"`
	Tasks     []Task    `bun:"rel:has-many,join:id=status_id"`
}

func (status *Status) Validate(validator *validator.Validator) {
	validator.Check(status.Name != "", "title", "must be provided")
	validator.Check(len(status.Name) <= 100, "title", "must not be more than 100 bytes long")
}
