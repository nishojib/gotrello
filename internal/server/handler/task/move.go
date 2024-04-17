package taskhandler

import (
	"errors"
	"net/http"
	"nishojib/gotrello/internal/data/models"
	"nishojib/gotrello/internal/data/repository"
	"nishojib/gotrello/internal/server/handler"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func Move(db *bun.DB) http.HandlerFunc {
	return handler.Make(func(w http.ResponseWriter, r *http.Request) error {
		newColumnIDStr := strings.Replace(r.FormValue("newColumnID"), "column-", "", 1)
		newColumnID, err := uuid.Parse(newColumnIDStr)
		if err != nil {
			return err
		}

		oldColumnIDStr := strings.Replace(r.FormValue("oldColumnID"), "column-", "", 1)
		oldColumnID, err := uuid.Parse(oldColumnIDStr)
		if err != nil {
			return err
		}

		newColumn, err := repository.NewStatusRepository(db).Get(newColumnID)
		if err != nil {
			return err
		}

		oldColumn, err := repository.NewStatusRepository(db).Get(oldColumnID)
		if err != nil {
			return err
		}

		itemIDStr := strings.Replace(r.FormValue("itemID"), "item-", "", 1)
		itemID, err := uuid.Parse(itemIDStr)
		if err != nil {
			return err
		}

		itemIndex := slices.IndexFunc(oldColumn.Tasks, func(task models.Task) bool {
			return task.ID == itemID
		})

		if itemIndex == -1 {
			return errors.New("item not found")
		}

		newColumnItems := newColumn.Tasks

		prevItemIDStr := strings.Replace(r.FormValue("prevItemID"), "item-", "", 1)
		nextItemIDStr := strings.Replace(r.FormValue("nextItemID"), "item-", "", 1)

		movedTask := oldColumn.Tasks[itemIndex]
		movedTask.StatusID = newColumnID

		if prevItemIDStr == "" && nextItemIDStr == "" {
			if len(newColumnItems) != 0 {
				return errors.New("bad request")
			}

			movedTask.SortOrder = 1
			if err := repository.NewTaskRepository(db).Update(&movedTask); err != nil {
				return err
			}

			return nil
		}

		if prevItemIDStr == "" && nextItemIDStr != "" {
			nextItemID, err := uuid.Parse(nextItemIDStr)
			if err != nil {
				return err
			}

			if newColumnItems[0].ID != nextItemID {
				return errors.New("bad request")
			}

			movedTask.SortOrder = newColumnItems[0].SortOrder / 2
			if err := repository.NewTaskRepository(db).Update(&movedTask); err != nil {
				return err
			}

			return nil
		}

		if prevItemIDStr != "" && nextItemIDStr == "" {
			prevItemID, err := uuid.Parse(prevItemIDStr)
			if err != nil {
				return err
			}

			if newColumnItems[len(newColumnItems)-1].ID != prevItemID {
				return errors.New("bad request")
			}

			movedTask.SortOrder = newColumnItems[len(newColumnItems)-1].SortOrder + 1
			if err := repository.NewTaskRepository(db).Update(&movedTask); err != nil {
				return err
			}

			return nil
		}

		prevItemID, err := uuid.Parse(prevItemIDStr)
		if err != nil {
			return err
		}

		nextItemID, err := uuid.Parse(nextItemIDStr)
		if err != nil {
			return err
		}

		prevItemIndex := slices.IndexFunc(newColumnItems, func(task models.Task) bool {
			return task.ID == prevItemID
		})

		if prevItemIndex == -1 {
			return errors.New("item not found")
		}

		if newColumnItems[prevItemIndex+1].ID != nextItemID {
			return errors.New("item not found")
		}

		movedTask.SortOrder =
			(newColumnItems[prevItemIndex].SortOrder + newColumnItems[prevItemIndex+1].SortOrder) / 2

		if err := repository.NewTaskRepository(db).Update(&movedTask); err != nil {
			return err
		}

		return nil
	})
}
