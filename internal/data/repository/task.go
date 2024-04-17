package repository

import (
	"context"
	"database/sql"
	"errors"
	"nishojib/gotrello/internal/data/models"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TaskRepository struct {
	db *bun.DB
}

func NewTaskRepository(db *bun.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (taskRepo *TaskRepository) Insert(task *models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return taskRepo.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var tasks []models.Task
		err := tx.NewSelect().
			Model(&tasks).
			Where("status_id = ?", task.StatusID).
			Order("sort_order ASC").
			Scan(ctx)
		if err != nil {
			return err
		}

		sortOrder := 1.0
		if len(tasks) > 0 {
			sortOrder = tasks[len(tasks)-1].SortOrder + 1
		}

		task.SortOrder = sortOrder
		_, err = tx.NewInsert().Model(task).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func (taskRepo *TaskRepository) Get(id uuid.UUID) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var task models.Task
	err := taskRepo.db.NewSelect().Model(&task).Where("id = ?", id).Scan(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Task{}, ErrRecordNotFound
		default:
			return models.Task{}, err
		}
	}

	return task, nil
}

func (taskRepo *TaskRepository) List(statusID uuid.UUID) ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var tasks []models.Task
	err := taskRepo.db.NewSelect().
		Model(&tasks).
		Where("status_id = ?", statusID).
		Order("sort_order ASC").
		Scan(ctx)
	if err != nil {
		return []models.Task{}, err
	}

	return tasks, nil
}

func (taskRepo *TaskRepository) Update(task *models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := taskRepo.db.NewUpdate().
		Model(task).
		Set("name = ?", task.Name).
		Set("status_id = ?", task.StatusID).
		Set("sort_order = ?", task.SortOrder).
		Set("version = ?", task.Version+1).
		Where("id = ?", task.ID).
		Where("version = ?", task.Version).
		Exec(ctx)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrEditConflict
	}

	return nil
}

func (taskRepo *TaskRepository) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := taskRepo.db.NewDelete().
		Model(&models.Task{}).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
