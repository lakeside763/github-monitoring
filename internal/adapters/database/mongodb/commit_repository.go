package mongodb

import (
	"github.com/lakeside763/github-repo/internal/core/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type MongoCommitRepository struct {
	collection *mongo.Collection
}

func NewMongoCommitRepository(client *mongo.Client, dbName string) *MongoCommitRepository {
	return &MongoCommitRepository{
		collection: client.Database(dbName).Collection("commits"),
	}
}

func (r *MongoCommitRepository) Create(commit *models.Commit) (*models.Commit, error) {
	return commit, nil
}

// BeginTransaction implements interfaces.Commits.
func (r *MongoCommitRepository) BeginTransaction() (*gorm.DB, error) {
	panic("unimplemented")
}

// CommitTransaction implements interfaces.Commits.
func (r *MongoCommitRepository) CommitTransaction(tx *gorm.DB) error {
	panic("unimplemented")
}

// RollbackTransaction implements interfaces.Commits.
func (r *MongoCommitRepository) RollbackTransaction(tx *gorm.DB) error {
	panic("unimplemented")
}

