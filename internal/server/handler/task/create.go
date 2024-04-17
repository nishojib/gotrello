package taskhandler

import (
	"fmt"
	"log/slog"
	"net/http"
	"nishojib/gotrello/internal/data/models"
	"nishojib/gotrello/internal/data/repository"
	"nishojib/gotrello/internal/server/handler"
	"nishojib/gotrello/internal/validator"
	projectui "nishojib/gotrello/ui/html/project"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func Create(db *bun.DB) http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		statusIDStr := r.FormValue("statusID")

		slog.Info(fmt.Sprintf("statusIDStr: %s", statusIDStr))

		statusID, err := uuid.Parse(statusIDStr)
		if err != nil {
			return err
		}

		task := models.Task{Name: r.FormValue("name"), StatusID: statusID}

		v := validator.New()
		if task.Validate(v); !v.Valid() {
			return handler.Render(
				w,
				r,
				projectui.TaskForm(projectui.TaskParams{Task: task}, projectui.TaskErrors{
					Name: v.Errors["name"],
				}),
			)
		}

		if err = repository.NewTaskRepository(db).Insert(&task); err != nil {
			return err
		}

		return handler.Render(
			w,
			r,
			projectui.TaskForm(projectui.TaskParams{Task: task}, projectui.TaskErrors{}),
		)
	})
}
