package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/curefatih/message-sender/model"
	"github.com/curefatih/message-sender/model/dto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func runMessageTask(ctx context.Context, cfg *viper.Viper, messageTask model.MessageTask) (*dto.MessageTaskSendResponse, error) {
	jsonPayload, err := json.Marshal(dto.MessageTaskSendPayload{
		MessageContent: messageTask.MessageContent,
		To:             messageTask.To,
	})
	if err != nil {
		log.Info().Msgf("Error marshaling struct to JSON: %v\n", err)
		return nil, err
	}

	// Prepare the request
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		cfg.GetString("process.task.message.url"),
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		log.Error().Err(err)
		return nil, fmt.Errorf("failed to create request: %w", err)
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
		log.Error().Err(err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for a non-200 status code
	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, body)
	}

	// Unmarshal the response into the struct
	var response dto.MessageTaskSendResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error().Err(err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	log.Info().Msgf("response: OK for task %d", messageTask.ID)

	return &response, nil
}
