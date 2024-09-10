# Github-Repo-Monitoring-Service

### Author
- By Moses Idowu

### Overview
Github monitoring services: The services fetched all the existing commits, watches and fetch new commits

### Prerequisites
- 1.22 Go version


### Installation
- Clone the repo
```
git clone https://github.com/lakeside763/github-monitor.git
cd github-monitor
```

- Setup env configuration
```
# Database configuration
DB_HOST=localhost
DB_NAME=db
DB_NAME_TEST=test_db
DB_USER=postgres
DB_PASSWORD=
DB_PORT=5432
DATABASE_TYPE=postgresql
DATABASE_URL=mongodb://localhost:27017
REDIS_URL
GITHUB_API_BASE_URL=https://api.github.com
GITHUB_TOKEN=
GITHUB_REPO_FULL_NAME=lakeside763/github-monitor
GITHUB_ACCOUNT_OWNER=lakeside763
GITHUB_REPO_NAME=github-monitor
MONITORING_INTERVAL=3600
FETCH_LIMIT=20
```

- Run with CLI
```
go run cmd/app/main.go
```

- Run with docker
```
docker compose build
docker compose up
```

- Run a test
```
go test -v ./internal/core/services
go test -v ./...
```

- The API documentation in details
- [https://documenter.getpostman.com/view/1194460/2sAXjDfFwK]

### Things Achieved
- Github commits monitoring service
- Github Fetch repo metadata
- Implement interfaces on core services component
- Service Retry mechanism incase of timeout or service failure
- Use of Redis for caching
- Prevent duplicate processing
- Get commit API
- Data persistent

#### Things I would have love to add
- Write mote test cases
- Add more APIs for managing Commits an Repo Metadata from the database


- Folder Structures
```
/cmd
  /app
    - main.go                # Entry point for the application

/internal
  /domain
    /models
      - commit_model.go          # Domain model for Commit
      - repository_model.go      # Domain model for Repository
  
  /adapters
    /cache/redis
        - redis.go            # Connect to redis DB
    /clients
      - github_client.go    # Implementation of GitHub API client
    /database
      /mongodb
        - commit_repository.go
        - repometadata_repository.go
      /postgresql
        - migrations
        - commit_repository.go
        - repometadata_repository.go
      /http
        /handlers
          - commit.go
          - repometadata.go
        /routes
          - commit.go
          - repometadata
          - setup.go
        /repository
          - repository.go

  /core
    /models
      - commit.go
      - repository.go

    /services
      /mocks
        - mocks.go                # Mock for testing repo metadata
      - commits_monitoring.go     # Business logic for repository metadata operations
      - repometadata.go           # Business logic for repository metadata operations
      - repometadata_test.go      # Business logic for repository test metadata operations

  /config
    - config.go                   # application configuration

/pkg
  /utils
    - gorm_utils.go
    - mongo_utils.go
    - utils.go             # Utility functions and helpers

```
 


