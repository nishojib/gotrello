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

type ProjectRepository struct {
	db *bun.DB
}

func NewProjectRepository(db *bun.DB) *ProjectRepository {
	return &ProjectRepository{db}
}

func (projectRepo *ProjectRepository) Insert(project *models.Project) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := projectRepo.db.NewInsert().Model(project).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (projectRepo *ProjectRepository) Get(id uuid.UUID) (models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var project models.Project
	err := projectRepo.db.NewSelect().Model(&project).Where("id = ?", id).Scan(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Project{}, ErrRecordNotFound
		default:
			return models.Project{}, err
		}
	}

	return project, nil
}

func (projectRepo *ProjectRepository) List() ([]models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var projects []models.Project
	err := projectRepo.db.NewSelect().Model(&projects).Scan(ctx)
	if err != nil {
		return []models.Project{}, err
	}

	return projects, nil
}

func (projectRepo *ProjectRepository) Update(project *models.Project) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := projectRepo.db.NewUpdate().
		Model(project).
		Set("name = ?", project.Name).
		Set("version = ?", project.Version+1).
		Where("id = ?", project.ID).
		Where("version = ?", project.Version).
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

func (projectRepo *ProjectRepository) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := projectRepo.db.NewDelete().
		Model(&models.Project{}).
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
