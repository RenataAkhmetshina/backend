package client

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type FlashcardClient struct {
	client *resty.Client
}

func NewFlashcardClient() *FlashcardClient {
	r := resty.New()

	r.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		fmt.Printf("[Internal Call] Outgoing to: %s %s\n", req.Method, req.URL)
		return nil
	})

	r.SetHeader("X-Internal-Secret", os.Getenv("INTERNAL_SERVICE_KEY"))

	return &FlashcardClient{client: r}
}

func (fc *FlashcardClient) WipeUserData(userID string) error {
	resp, err := fc.client.R().
		Delete(fmt.Sprintf("http://localhost:8082/internal/flashcards/user/%s", userID))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("flashcard service returned error: %d", resp.StatusCode())
	}

	return nil
}
