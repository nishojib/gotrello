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

type StatusRepository struct {
	db *bun.DB
}

func NewStatusRepository(db *bun.DB) *StatusRepository {
	return &StatusRepository{db}
}

func (statusRepo *StatusRepository) Insert(status *models.Status) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := statusRepo.db.NewInsert().Model(status).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (statusRepo *StatusRepository) Get(id uuid.UUID) (models.Status, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var status models.Status
	err := statusRepo.db.NewSelect().
		Model(&status).
		Relation("Tasks", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("sort_order ASC")
		}).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Status{}, ErrRecordNotFound
		default:
			return models.Status{}, err
		}
	}

	return status, nil
}

func (statusRepo *StatusRepository) List(projectID uuid.UUID) ([]models.Status, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var statuses []models.Status
	err := statusRepo.db.NewSelect().
		Model(&statuses).
		Relation("Tasks", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Order("sort_order ASC")
		}).
		Where("project_id = ?", projectID).
		Scan(ctx)
	if err != nil {
		return []models.Status{}, err
	}

	return statuses, nil
}

func (statusRepo *StatusRepository) Update(status *models.Status) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := statusRepo.db.NewUpdate().
		Model(status).
		Set("name = ?", status.Name).
		Set("version = ?", status.Version+1).
		Where("id = ?", status.ID).
		Where("version = ?", status.Version).
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

func (statusRepo *StatusRepository) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := statusRepo.db.NewDelete().
		Model(&models.Status{}).
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
