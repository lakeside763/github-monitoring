package interfaces

import "github.com/lakeside763/github-repo/internal/core/models"

type Github interface {
	FetchRepoMetadata(repoName string) (*models.Repository, error)
}