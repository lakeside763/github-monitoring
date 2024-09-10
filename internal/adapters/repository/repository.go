package repository

import (
	"context"
	"log"

	"github.com/lakeside763/github-repo/config"
	"github.com/lakeside763/github-repo/internal/adapters/database/mongodb"
	"github.com/lakeside763/github-repo/internal/adapters/database/postgresql"
	"github.com/lakeside763/github-repo/internal/core/ports/interfaces"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DataRepository struct {
	Commit	interfaces.Commits
	Repometadata interfaces.Repometadata
}

func NewDataRepository() *DataRepository {
	config := config.LoadConfig()

	var commitRepo interfaces.Commits
	var repometadataRepo interfaces.Repometadata

	// Initialize database and manage repository
	switch config.DatabaseType {
	case "mongodb":
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.DatabaseURL))
		if err != nil {
			log.Fatal(err)
		}
		commitRepo = mongodb.NewMongoCommitRepository(client, config.DatabaseName)
		repometadataRepo = mongodb.NewMongoRepometadataRepository(client, config.DatabaseName)
	case "postgresql":
		db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
		commitRepo = postgresql.NewPostgressCommitRepository(db)
		repometadataRepo = postgresql.NewPostgressRepometadataRepository(db)
	default:
		log.Fatal("Unsupported database type")
	}

	return &DataRepository{
		Commit: commitRepo,
		Repometadata: repometadataRepo,
	}
}