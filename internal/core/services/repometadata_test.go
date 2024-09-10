package services

import (
	"testing"

	"github.com/lakeside763/github-repo/internal/core/models"
	"github.com/lakeside763/github-repo/internal/core/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFetchRepoMetadata_RepoExists(t *testing.T) {
    mockGithubAPIClient := new(mocks.MockGithubAPIClient)
    mockRedisClient := new(mocks.MockRedisClient)
    mockRepoMetadata := new(mocks.MockRepoMetadata)

    service := NewRepoMetadataService(mockGithubAPIClient, mockRedisClient, mockRepoMetadata)

    repoFullName := "test-repo"
    existingRepo := &models.Repository{FullName: repoFullName, RepoID: 1}

    mockRepoMetadata.On("GetByFullName", repoFullName).Return(existingRepo, nil)

    repo, err := service.FetchRepoMetadata(repoFullName)

    assert.NoError(t, err)
    assert.Equal(t, existingRepo, repo)
    mockRepoMetadata.AssertExpectations(t)
}