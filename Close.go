package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type CloseRequest struct {
	ID       string `json:"ID"`
	TenantID string `json:"TenantID"`
	Reason   string `json:"Reason"`
}

type CloseResult struct {
	AlertID   string
	AlertName string
	Success   bool
	Error     error
}

func CloseAlerts(alerts []Alert, config JsonConfig, client *http.Client) {
	if !config.CloseAlerts {
		fmt.Println("CloseAlerts is disabled in config")
		return
	}

	if len(alerts) == 0 {
		fmt.Println("No alerts to close")
		return
	}

	fmt.Printf("Starting to close %d alerts...\n", len(alerts))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)
	results := make(chan CloseResult, len(alerts))

	for _, alert := range alerts {
		wg.Add(1)
		go func(a Alert) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := closeAlert(a, config, client)
			results <- result
		}(alert)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	successCount := 0
	failCount := 0

	for result := range results {
		if result.Success {
			successCount++
			if config.Debug {
				fmt.Printf("✓ Closed: %s (%s)\n", result.AlertID, result.AlertName)
			}
		} else {
			failCount++
			fmt.Printf("✗ Failed to close %s (%s): %v\n", result.AlertID, result.AlertName, result.Error)
		}
	}

	fmt.Printf("\nClose Summary:\n")
	fmt.Printf("  Success: %d\n", successCount)
	fmt.Printf("  Failed:  %d\n", failCount)
	fmt.Printf("  Total:   %d\n", len(alerts))
}

func closeAlert(alert Alert, config JsonConfig, client *http.Client) CloseResult {
	result := CloseResult{
		AlertID:   alert.InternalID,
		AlertName: alert.Name,
	}

	closeReq := CloseRequest{
		ID:       alert.InternalID,
		TenantID: alert.TenantID,
		Reason:   config.CloseReason,
	}

	jsonData, err := json.Marshal(closeReq)
	if err != nil {
		result.Error = fmt.Errorf("JSON marshal error: %w", err)
		return result
	}

	url := fmt.Sprintf("%s/alerts/close?tenantID=%s", config.BaseURL, alert.TenantID)

	if config.Debug {
		fmt.Printf("Closing alert: POST %s\n", url)
		fmt.Printf("Body: %s\n", string(jsonData))
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		result.Error = fmt.Errorf("request creation error: %w", err)
		return result
	}

	req.Header.Set("Authorization", "Bearer "+config.Token)
	req.Header.Set("Content-Type", "application/json")

	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := client.Do(req)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(time.Second * time.Duration(attempt))
				continue
			}
			result.Error = fmt.Errorf("HTTP error after %d attempts: %w", maxRetries, err)
			return result
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			result.Success = true
			return result
		}

		if attempt < maxRetries && resp.StatusCode >= 500 {
			time.Sleep(time.Second * time.Duration(attempt))
			continue
		}

		result.Error = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		return result
	}

	return result
}
