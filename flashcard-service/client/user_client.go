package client

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type UserClient struct {
	client *resty.Client
}

func NewUserClient() *UserClient {
	r := resty.New()

	baseURL := os.Getenv("USER_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}
	r.SetBaseURL(baseURL)

	r.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		fmt.Printf("[UserClient] → %s %s\n", req.Method, req.URL)
		return nil
	})

	r.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		fmt.Printf("[UserClient] ← %d\n", resp.StatusCode())
		return nil
	})

	return &UserClient{client: r}
}

func (uc *UserClient) CheckUserExists(userID uint, token string) (bool, error) {
	resp, err := uc.client.R().
		SetHeader("Authorization", "Bearer "+token).
		Get(fmt.Sprintf("/api/users/%d", userID))

	if err != nil {
		return false, err
	}
	if resp.StatusCode() == 404 {
		return false, nil
	}
	if resp.IsError() {
		return false, fmt.Errorf("user-service error: %d", resp.StatusCode())
	}
	return true, nil
}
