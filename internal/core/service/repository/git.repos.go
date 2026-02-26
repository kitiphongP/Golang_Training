package repository

import (
	"context"
	"golang/internal/core/models"
)

type SkillRepository interface {
	Save(ctx context.Context, skill *models.UserSkills) error
	FindByUsername(ctx context.Context, username string) (*models.UserSkills, error)
}