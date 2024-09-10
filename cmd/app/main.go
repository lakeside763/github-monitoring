package main

import (
	"net/http"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/lakeside763/github-repo/config"
	"github.com/lakeside763/github-repo/internal/adapters/cache/redis"
	"github.com/lakeside763/github-repo/internal/adapters/clients"
	"github.com/lakeside763/github-repo/internal/adapters/http/routes"
	"github.com/lakeside763/github-repo/internal/adapters/repository"
	"github.com/lakeside763/github-repo/internal/core/services"
	log "github.com/sirupsen/logrus"
)

func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .evn file: %v", err)
	}
}

func main() {
	// Initialize dotenv
	initEnv()
	config := config.LoadConfig()

	// Initialize repositories
	dataRepo := repository.NewDataRepository()

	redisClient := redis.NewRedisCache(config.RedisURL)

	// Initialize Github API Client
	githubAPIClient := clients.NewGithubAPIClient(config.GithubAPIBaseURL, config.GithubToken)

	// Initialize Repometadata service
	repoMetadataService := services.NewRepoMetadataService(githubAPIClient, redisClient, dataRepo.Repometadata,)
	_, err := repoMetadataService.FetchRepoMetadata(config.RepoFullName)
	if err != nil {
		log.Fatalf("Failed to fetch repo metadata: %v", err)
	}

	// Initialize Commit Monitoring Service
	monitoringService := services.NewMonitoringService(*githubAPIClient, redisClient, dataRepo.Commit)
	go func() {
		monitoringService.StartMonitoring(config.MonitoringInterval, config.RepoFullName, config.FetchLimit)
	}()

	router := httprouter.New()
	routes.SetupRoutes(router, dataRepo)
	
	log.Info("Starting server on :5200")
	log.Fatal(http.ListenAndServe(":5200", router))
}