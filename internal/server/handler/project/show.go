package projecthandler

import (
	"net/http"
	"nishojib/gotrello/internal/data/repository"
	"nishojib/gotrello/internal/server/handler"
	projectui "nishojib/gotrello/ui/html/project"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func Show(db *bun.DB) http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		projectIDStr := chi.URLParam(r, "projectID")

		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			return err
		}

		project, err := repository.NewProjectRepository(db).Get(projectID)
		if err != nil {
			return err
		}

		statuses, err := repository.NewStatusRepository(db).List(project.ID)
		if err != nil {
			return err
		}

		return handler.Render(w, r, projectui.Show(project.Name, statuses))
	})
}
