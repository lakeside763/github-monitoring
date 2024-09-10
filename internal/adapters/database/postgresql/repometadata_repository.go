package postgresql

import (
	"github.com/lakeside763/github-repo/internal/core/models"
	"github.com/lakeside763/github-repo/pkg/utils"
	"gorm.io/gorm"
)

type PostgresRepometadataRepository struct {
	db *gorm.DB
}

func NewPostgressRepometadataRepository(db *gorm.DB) *PostgresRepometadataRepository {
	return &PostgresRepometadataRepository{db: db}
}

func (r PostgresRepometadataRepository) Create(repo *models.Repository) (*models.Repository, error) {
	if err := r.db.Create(repo).Error; err != nil {
		return nil, err
	}
	return repo, nil
}

// GetByFullName implements interfaces.Repometadata.
func (r *PostgresRepometadataRepository) GetByFullName(repoFullName string) (*models.Repository, error) {
	var repo models.Repository
	err := r.db.Where("full_name = ?", repoFullName).First(&repo).Error
	if err != nil {
		return nil, utils.HandleGormError(err)
	}
	return &repo, nil
}
