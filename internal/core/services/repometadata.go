package services

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/lakeside763/github-repo/config"
	"github.com/lakeside763/github-repo/internal/core/models"
	"github.com/lakeside763/github-repo/internal/core/ports/interfaces"
)

type RepoMetadataService struct {
	GithubAPIClient interfaces.Github
	RedisClient	 interfaces.Redis
	Repository interfaces.Repometadata
}

func NewRepoMetadataService(
	githubAPIClient interfaces.Github, redisClient interfaces.Redis, repository interfaces.Repometadata,
) *RepoMetadataService {
	return &RepoMetadataService{
		GithubAPIClient: 	githubAPIClient,
		RedisClient: redisClient,
		Repository: repository,
	}
}

func (r *RepoMetadataService) FetchRepoMetadata(repoFullName string) (*models.Repository, error) {
	config := config.LoadConfig()
	dbName := config.DatabaseName

	// Check if the repo already exists in the database
	existingRepo, err := r.Repository.GetByFullName(repoFullName)
	if err != nil {
		log.Errorf("Error checking existing repo: %v", err)
		return nil, err
	}
	if existingRepo != nil {
		log.Info("Repo already fetched from Github")
		return existingRepo, nil
	}

	// Fetch repo from Github
	repo, err := r.GithubAPIClient.FetchRepoMetadata(repoFullName)
	if err != nil {
		log.Errorf("Error fetching repo metadata from Github: %v", err)
		return nil, err
	}
	
	if _, err := r.Repository.Create(repo); err != nil {
		log.Fatalf("Error to save repo metadata to database %v", err)
		return nil, err
	}

	repoKey := fmt.Sprintf("%s-%s-id", repo.FullName, dbName)
	if err := r.RedisClient.Set(repoKey, repo.RepoID, 0); err != nil {
		log.Errorf("Failed to save repo ID in Redis: %v", err)
		return nil, err
	}

	log.Infof("Fetched and saved repo metadata successfully %v", repoFullName)
	return repo, nil
}