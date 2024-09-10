package services

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/lakeside763/github-repo/config"
	"github.com/lakeside763/github-repo/internal/adapters/clients"
	"github.com/lakeside763/github-repo/internal/core/models"
	"github.com/lakeside763/github-repo/internal/core/ports/interfaces"
	log "github.com/sirupsen/logrus"
)

type MonitoringService struct {
	GithubAPIClient clients.GithubAPIClient
	RedisClient     interfaces.Redis
	Commit          interfaces.Commits
}

func NewMonitoringService(
	githubAPIClient clients.GithubAPIClient, redisClient interfaces.Redis, commitRepo interfaces.Commits,
) *MonitoringService {
	return &MonitoringService{
		GithubAPIClient: githubAPIClient,
		RedisClient:     redisClient,
		Commit:          commitRepo,
	}
}

func (c *MonitoringService) StartMonitoring(interval time.Duration, repo string, limit int) {
	var mu sync.Mutex
	s := gocron.NewScheduler(time.UTC)
	config := config.LoadConfig()
	db := config.DatabaseName

	s.Every(interval).Do(func() {
		// Lock the critical section
		mu.Lock()
		defer mu.Unlock()

		// Start a transaction
		tx, err := c.Commit.BeginTransaction()
		if err != nil {
			log.Errorf("Failed to start transaction: %v", tx.Error)
			return
		}

		pagination, err := c.RedisClient.GetPagination(repo)
		if err != nil {
			log.Errorf("Error retrieving pagination: %v", err)
			tx.Rollback()
			return
		}

		params := models.FetchCommitInput{
			Repo:            repo,
			Since:           pagination.Since,
			Until:           pagination.Until,
			Page:            pagination.Page,
			Limit:           limit,
			NextWindowSince: pagination.NextWindowSince,
		}

		commits, newPagination, err := c.GithubAPIClient.FetchCommitWithRetry(params)
		if err != nil {
			log.Errorf("Error fetching commits: %v", err)
			tx.Rollback()
			return
		}

		// Retrieve the repository ID from Redis
		repoIdKey := fmt.Sprintf("%v-%v-id", repo, db)
		repoIdStr, err := c.RedisClient.Get(repoIdKey)
		if err != nil || repoIdStr == "" {
			log.Errorf("Repo ID not found in Redis: %v", err)
			tx.Rollback()
			return
		}

		repoId, _ := strconv.Atoi(repoIdStr)
		processedSHAsKey := fmt.Sprintf("%v-%v-shas", repo, db)

		// Save commits to the database
		for _, data := range commits {
			commit := models.Commit{
				SHA:     data.SHA,
				Message: data.Commit.Message,
				Author:  data.Commit.Author.Name,
				Date:    data.Commit.Author.Date,
				URL:     data.HTMLURL,
				RepoId:  repoId,
			}

			// Check if commit already processed
			processed, err := c.RedisClient.SIsMember(processedSHAsKey, commit.SHA)
			if err != nil {
				log.Errorf("Error checking SHA in Redis: %v", err)
				tx.Rollback()
				return
			}

			if processed {
				log.Infof("SHA %v already processed, skipping...", commit.SHA)
				continue
			}

			if err := tx.Create(&commit).Error; err != nil {
				log.Errorf("Failed to save commit to database: %v", err)
				tx.Rollback()
				return
			}

			// Add commit SHA to Redis (processedSHAsKey)
			if err := c.RedisClient.SAdd(processedSHAsKey, commit.SHA); err != nil {
				log.Errorf("Error adding SHA to Redis: %v", err)
				tx.Rollback()
				return
			}
		}

		// Save the new pagination data to Redis
		if err := c.RedisClient.SavePagination(newPagination, repo); err != nil {
			log.Errorf("Error saving pagination to Redis: %v", err)
			tx.Rollback()
			return
		}

		// Commit the transaction if all operations succeed
		if err := tx.Commit().Error; err != nil {
			log.Errorf("Failed to commit transaction: %v", err)
			tx.Rollback()
			return
		}

		log.Infof("Fetch commit task completed successfully at %v", time.Now())
	})

	s.StartAsync()
}
