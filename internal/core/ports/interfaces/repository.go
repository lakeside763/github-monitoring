package interfaces

import "github.com/lakeside763/github-repo/internal/core/models"

type Repometadata interface {
	Create(repo *models.Repository) (*models.Repository, error)
	GetByFullName(repoFullName string) (*models.Repository, error)
}