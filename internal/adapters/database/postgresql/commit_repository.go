package postgresql

import (
	"github.com/lakeside763/github-repo/internal/core/models"
	"gorm.io/gorm"
)

type PostgresCommitRepository struct {
	db *gorm.DB
}

func NewPostgressCommitRepository(db *gorm.DB) *PostgresCommitRepository {
	return &PostgresCommitRepository{db: db}
}

func (r *PostgresCommitRepository) Create(commit *models.Commit) (*models.Commit, error) {
	if err := r.db.Create(commit).Error; err != nil {
		return nil, err
	}
	return commit, nil
}

func (r *PostgresCommitRepository) BeginTransaction() (*gorm.DB, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
			return nil, tx.Error
	}
	return tx, nil
}

func (r *PostgresCommitRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *PostgresCommitRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}