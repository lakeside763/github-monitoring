package interfaces

import (
	"github.com/lakeside763/github-repo/internal/core/models"
	"gorm.io/gorm"
)

type Commits interface {
	Create(commit *models.Commit) (*models.Commit, error)
	BeginTransaction() (*gorm.DB, error)
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error
}