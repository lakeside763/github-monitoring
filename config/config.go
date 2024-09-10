package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL        string
	RedisURL           string
	DatabaseType       string
	DatabaseName       string
	GithubAPIBaseURL   string
	GithubToken        string
	MonitoringInterval time.Duration
	GithubAccountOwner string
	RepoFullName       string
	RepoName           string
	FetchLimit         int
}

func LoadConfig() *Config {
	monitoringIntervalStr := os.Getenv("MONITORING_INTERVAL")
	monitoringInterval, _ := strconv.Atoi(monitoringIntervalStr)

	fetchLimitStr := os.Getenv("FETCH_LIMIT")
	fetchLimit, _ := strconv.Atoi(fetchLimitStr)

	databaseType := os.Getenv("DATABASE_TYPE")
	var databaseURL string
	if databaseType == "mongodb" {
		databaseURL = os.Getenv("DATABASE_URL")
	} else if databaseType == "postgresql" {
		databaseURL = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	}

	return &Config{
		DatabaseURL:        databaseURL,
		RedisURL:           os.Getenv("REDIS_URL"),
		DatabaseType:       databaseType,
		DatabaseName:       os.Getenv("DB_NAME"),
		GithubAPIBaseURL:   os.Getenv("GITHUB_API_BASE_URL"),
		GithubToken:        os.Getenv("GITHUB_TOKEN"),
		GithubAccountOwner: os.Getenv("GITHUB_ACCOUNT_OWNER"),
		RepoName:           os.Getenv("GITHUB_REPO_NAME"),
		RepoFullName:       os.Getenv("GITHUB_REPO_FULL_NAME"),
		MonitoringInterval: time.Duration(monitoringInterval) * time.Second,
		FetchLimit:         fetchLimit,
	}
}
