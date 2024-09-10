package mongodb

import (
	"context"
	"time"

	"github.com/lakeside763/github-repo/internal/core/models"
	"github.com/lakeside763/github-repo/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepometadataRepository struct {
	Repository *mongo.Collection
}

func NewMongoRepometadataRepository(client *mongo.Client, dbName string) *MongoRepometadataRepository {
	return &MongoRepometadataRepository{
		Repository: client.Database(dbName).Collection("repositories"),
	}
}

func (r *MongoRepometadataRepository) Create(repo *models.Repository) (*models.Repository, error) {
	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()

	result, err := r.Repository.InsertOne(context.TODO(), repo)
	if err != nil {
		return nil, err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		repo.ID = oid.Hex()
	} else {
		return nil, mongo.ErrInvalidIndexValue
	}

	return repo, nil
}

// GetByFullName implements interfaces.Repometadata.
func (r *MongoRepometadataRepository) GetByFullName(repoFullName string) (*models.Repository, error) {
	var repo models.Repository
	filter := bson.M{"full_name": repoFullName}

	err := r.Repository.FindOne(context.TODO(), filter).Decode(&repo)
	if err != nil {
		return nil, utils.HandleMongoError(err)
	}
	return &repo, nil
}
