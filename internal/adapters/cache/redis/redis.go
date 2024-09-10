package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lakeside763/github-repo/config"
	"github.com/lakeside763/github-repo/internal/core/models"
	log "github.com/sirupsen/logrus"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	config := config.LoadConfig()
	client := redis.NewClient(&redis.Options{
		Addr: config.RedisURL,
	})
	return &RedisCache{Client: client}
}

func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(context.TODO(), key, value, expiration).Err()
}

func (c *RedisCache) Get(key string) (string, error) {
	return c.Client.Get(context.TODO(), key).Result()
}

func (c *RedisCache) SIsMember(key string, member interface{}) (bool, error) {
	return c.Client.SIsMember(context.TODO(), key, member).Result()
}

func (c *RedisCache) SAdd(key string, member interface{}) error {
	return c.Client.SAdd(context.TODO(), key, member).Err()
}

func (c *RedisCache) GetPagination(repo string) (*models.Pagination, error) {
	ctx := context.Background()
	config := config.LoadConfig()
	db := config.DatabaseName
	key := fmt.Sprintf("%v-%v-repo", repo, db)

	exists, err := c.Client.Exists(ctx, key).Result()
	if err != nil {
		log.Errorf("Error checking existence of pagination in Redis: %v", err)
		return nil, err
	}

	if exists == 0 {
		// Return default values if the pagination key doesn't exist
		sinceTime, _ := time.Parse(time.RFC3339, "2010-01-01T00:00:00Z")
		untilTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		return &models.Pagination{
			Since:           sinceTime,
			Until:           untilTime,
			Page:            1,
			NextWindowSince: time.Time{},
		}, nil
	}

	// Get the JSON string from Redis
	paginationData, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// Deserialize the JSON string back to the Pagination struct
	var pagination models.Pagination
	err = json.Unmarshal([]byte(paginationData), &pagination)
	if err != nil {
		log.Errorf("Failed to unmarshall pagination data: %v", err)
		return nil, err
	}

	return &pagination, nil
}

func (c *RedisCache) SavePagination(pagination models.Pagination, repo string) (error) {
	ctx := context.Background()
	config := config.LoadConfig()
	db := config.DatabaseName

	// Serialize the pagination struct to a JSON string
	paginationData, err := json.Marshal(pagination)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%v-%v-repo", repo, db)
	err = c.Client.Set(ctx, key, paginationData, 0).Err()
	if err != nil {
		log.Errorf("Failed to save pagination to Redis: %v", err)
		return err
	}
	return nil
}
