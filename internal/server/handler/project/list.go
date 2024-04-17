package projecthandler

import (
	"net/http"
	"nishojib/gotrello/internal/data/repository"
	"nishojib/gotrello/internal/server/handler"
	projectui "nishojib/gotrello/ui/html/project"

	"github.com/uptrace/bun"
)

func List(db *bun.DB) http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		projectList, err := repository.NewProjectRepository(db).List()
		if err != nil {
			return err
		}

		return handler.Render(w, r, projectui.List(projectList))
	})
}
