package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) error {
	apiURL := os.Getenv("API_URL")
	workerSecret := os.Getenv("WORKER_SECRET")

	if apiURL == "" || workerSecret == "" {
		return fmt.Errorf("API_URL or WORKER_SECRET not set")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+workerSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
