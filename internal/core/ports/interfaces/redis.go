package interfaces

import (
	"time"

	"github.com/lakeside763/github-repo/internal/core/models"
)


type Redis interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	SIsMember(key string, member interface{}) (bool, error)
	SAdd(key string, member interface{}) error
	GetPagination(repo string) (*models.Pagination, error)
	SavePagination(pagination models.Pagination, repo string) error
}