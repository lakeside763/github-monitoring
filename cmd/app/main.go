package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .evn file: %v", err)
	} else {
		log.Info("Loading configuration")
	}
}

func main() {
	// Initialize dotenv
	config := config.LoadConfig()

	// Initialize repositories
	dataRepo := repository.NewDataRepository()

	redisClient := redis.NewRedisCache(config.RedisURL)

	// Initialize Github API Client
	githubAPIClient := clients.NewGithubAPIClient(config.GithubAPIBaseURL, config.GithubToken)


	// Create a context for cancellation and graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup

	// Initialize Repometadata service
	repoMetadataService := services.NewRepoMetadataService(githubAPIClient, redisClient, dataRepo.Repometadata)
	_, err := repoMetadataService.FetchRepoMetadata(config.RepoFullName)
	if err != nil {
		log.Fatalf("Failed to fetch repo metadata: %v", err)
	}

	// Initialize Commit Monitoring Service
	monitoringService := services.NewMonitoringService(*githubAPIClient, redisClient, dataRepo.Commit)
	wg.Add(1)
	go func() {
		defer wg.Done()
		monitoringService.StartMonitoring(config.MonitoringInterval, config.RepoFullName, config.FetchLimit)
	}()

	router := httprouter.New()
	routes.SetupRoutes(router, dataRepo)

	server := &http.Server{
		Addr: ":5200",
		Handler: router,
	}

	// Run the HTTP server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Starting server on :5200")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to serve: %v", err)
		}
	}()
	
	// log.Info("Starting server on :5200")
	// log.Fatal(http.ListenAndServe(":5200", router))
	
	// Handle graceful shutdown
	waitForShutdown(ctx, cancel, server, &wg)
}

func waitForShutdown(ctx context.Context, cancel context.CancelFunc, server *http.Server, wg *sync.WaitGroup) {
	// Capture termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// wait for a termination signal
	<-sigChan
	log.Info("Shutting down gracefully")

	// Cancel the context to signal goroutines to stop
	cancel()

	// Gracefully shut down the HTTP server with a timeout
	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	} else {
		log.Info("HTTP server stopped gracefully")
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Info("All services stopped")
}