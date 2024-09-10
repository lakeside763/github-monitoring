package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lakeside763/github-repo/internal/core/models"
	log "github.com/sirupsen/logrus"
)

type GithubAPIClient struct {
	BaseURL string
	Client  *http.Client
	Token   string
}

func NewGithubAPIClient(baseURL string, token string) *GithubAPIClient {
	return &GithubAPIClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
		Token:   token,
	}
}

func (c *GithubAPIClient) FetchRepoMetadata(repoName string) (*models.Repository, error) {
	url := fmt.Sprintf("%s/repos/%s", c.BaseURL, repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository metadata: %v", err)
	}
	defer resp.Body.Close()

	var repository models.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, err
	}

	return &repository, nil
}

func (c *GithubAPIClient) FetchCommitWithRetry(f models.FetchCommitInput) ([]models.ResponseCommit, models.Pagination, error) {
	httpClient := &http.Client{
		Timeout: time.Minute * 5,
	}

	const maxRetries = 3
	var err error

	var commits []models.ResponseCommit
	for retries := 0; retries < maxRetries; retries++ {
		url := fmt.Sprintf("%s/repos/%s/commits?since=%s&until=%s&per_page=%d&page=%d",
			c.BaseURL, f.Repo, f.Since.Format(time.RFC3339), f.Until.Format(time.RFC3339), f.Limit, f.Page)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, models.Pagination{}, err
		}

		if c.Token != "" {
			req.Header.Set("Authorization", "token "+c.Token)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Errorf("Error during HTTP request: %v. Retry %d/%d", err, retries+1, maxRetries)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("failed to fetch commits, status code: %d", resp.StatusCode)
			log.Errorf("%v. Retry %d/%d", err, retries+1, maxRetries)
			continue
		}

		err = json.NewDecoder(resp.Body).Decode(&commits)
		if err != nil {
			log.Errorf("Error decoding response: %v. Retry %d/%d", err, retries+1, maxRetries)
			continue
		}

		// Initialize nextSince and nextUntil with current values
		nextSince := f.Since
		nextUntil := f.Until
		nextWindowSince := f.NextWindowSince // Keep track of Since for next window

		// If the number of commits is equal to the limit, prepare for the next page
		if len(commits) >= f.Limit {
			if f.Page == 1 {
				lastCommit := commits[0]
				nextWindowSince = lastCommit.Commit.Author.Date.Add(1 * time.Second) // Add one second to prevent duplicate fetching
			}

			nextPage := f.Page + 1
			pagination := models.Pagination{
				Since:           nextSince,
				Until:           nextUntil,
				Page:            nextPage,
				NextWindowSince: nextWindowSince,
			}

			return commits, pagination, nil
		}

		// Initialize pagination for the next window
		timeNow, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
		pagination := models.Pagination{
			Since:           nextWindowSince,
			Until:           timeNow,
			Page:            1,
			NextWindowSince: nextWindowSince,
		}

		return commits, pagination, nil
	}
	return nil, models.Pagination{}, fmt.Errorf("failed to fetch commits after %d attempts: %v", maxRetries, err)
}
