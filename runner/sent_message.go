package runner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/curefatih/message-sender/model"
	"github.com/spf13/viper"
)

func runMessageTask(ctx context.Context, cfg *viper.Viper, messageTask model.MessageTask) error {
	// Prepare the request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		cfg.GetString("process.task.message.url"),
		bytes.NewBuffer([]byte(messageTask.MessageContent)),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the header
	req.Header.Set("x-ins-auth-key", cfg.GetString("process.task.message.auth.key"))
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check for a non-200 status code
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, body)
	}

	fmt.Printf("response: OK for task %s\n", messageTask.ID)
	return nil
}
