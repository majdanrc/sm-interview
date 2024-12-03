package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qwark97/interview/fetcher/model"
	"io"
	"net/http"
	"time"
)

type Fetcher struct {
	baseUrl string
	client  *http.Client
	retries int
}

func New(url string, client *http.Client) Fetcher {
	return Fetcher{
		baseUrl: url,
		client:  client,
		retries: 3,
	}
}

func (f *Fetcher) Users(ctx context.Context, processor func(user model.User) error) error {
	url := fmt.Sprintf("%s/users", f.baseUrl)

	for url != "" {
		resp, err := f.requestUsers(ctx, url)
		if err != nil {
			return err
		}

		for _, user := range resp.Users {
			if err := processor(user); err != nil {
				return err
			}
		}

		url = resp.NextLink
	}

	return nil
}

func (f *Fetcher) requestUsers(ctx context.Context, url string) (*model.Response, error) {
	var response *model.Response
	var err error

	for i := 0; i < f.retries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := f.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		return response, nil
	}

	return nil, fmt.Errorf("failed after %d retries: %w", f.retries, err)
}
