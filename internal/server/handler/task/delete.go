package taskhandler

import (
	"net/http"
	"nishojib/gotrello/internal/data/repository"
	"nishojib/gotrello/internal/server/handler"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func Delete(db *bun.DB) http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		taskIDStr := chi.URLParam(r, "taskID")

		taskID, err := uuid.Parse(taskIDStr)
		if err != nil {
			return err
		}

		if err = repository.NewTaskRepository(db).Delete(taskID); err != nil {
			return err
		}

		return nil
	})
}
