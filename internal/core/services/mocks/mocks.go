package mocks

import (
	"time"

	"github.com/lakeside763/github-repo/internal/core/models"
	"github.com/stretchr/testify/mock"
)

// MockGithubClient is a mock implementation of the Github interface
type MockGithubAPIClient struct {
	mock.Mock
}

func (m *MockGithubAPIClient) FetchRepoMetadata(repoFullName string) (*models.Repository, error) {
	args := m.Called(repoFullName)
	return args.Get(0).(*models.Repository), args.Error(1)
}

// MockRedisClient is a mock implementation of the Redis interface
type MockRedisClient struct {
	mock.Mock
}

// GetPagination implements interfaces.Redis.
func (m *MockRedisClient) GetPagination(repo string) (*models.Pagination, error) {
	panic("unimplemented")
}

// SAdd implements interfaces.Redis.
func (m *MockRedisClient) SAdd(key string, member interface{}) error {
	panic("unimplemented")
}

// SIsMember implements interfaces.Redis.
func (m *MockRedisClient) SIsMember(key string, member interface{}) (bool, error) {
	panic("unimplemented")
}

// SavePagination implements interfaces.Redis.
func (m *MockRedisClient) SavePagination(pagination models.Pagination, repo string) error {
	panic("unimplemented")
}

func (m *MockRedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

// MockRepoMetadata is a mock implementation of the Repometadata interface
type MockRepoMetadata struct {
	mock.Mock
}

func (m *MockRepoMetadata) GetByFullName(fullName string) (*models.Repository, error) {
	args := m.Called(fullName)
	return args.Get(0).(*models.Repository), args.Error(1)
}

func (m *MockRepoMetadata) Create(repo *models.Repository) (*models.Repository, error) {
	args := m.Called(repo)
	return args.Get(0).(*models.Repository), args.Error(1)
}
