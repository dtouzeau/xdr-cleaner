package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func main() {

	sPath := ConfigPath()
	TheConf := LoadConfig()

	// Required args
	if TheConf.TenantID == "" {
		fmt.Println("ERROR: is required see:" + sPath)
		os.Exit(1)
	}
	if TheConf.Token == "" {
		fmt.Println("ERROR: token is required see:" + sPath)
		os.Exit(1)
	}
	if len(TheConf.BaseURL) < 3 {
		fmt.Println("ERROR: URL is required see:" + sPath)
		os.Exit(1)
	}

	if TheConf.Debug {
		fmt.Println("Using token:", TheConf.Token)
		fmt.Printf("Max concurrent pages: %d\n", TheConf.MaxConcurrentPages)
		fmt.Printf("Flush every: %d alerts\n", TheConf.FlushEvery)
	}

	client := BuilClient()

	flushMgr := NewFlushManager(TheConf.Outfile, TheConf.FlushEvery, TheConf.Debug)
	allAlerts := fetchAllAlertsParallelWithFlush(TheConf, client, flushMgr)

	err := flushMgr.Finalize()
	if err != nil {
		fmt.Println("FLUSH FINALIZE ERROR:", err)
		os.Exit(1)
	}

	fmt.Printf("Saved %d alerts to %s\n", len(allAlerts), TheConf.Outfile)

	// Apply filters if enabled
	if TheConf.FilterMode {
		fmt.Println("\n=== Filtering Alerts ===")
		filteredAlerts := FilterAlerts(allAlerts, TheConf)

		if len(filteredAlerts) > 0 {
			err := SaveFilteredAlerts(filteredAlerts, TheConf.FilteredOutfile)
			if err != nil {
				fmt.Println("FILTER SAVE ERROR:", err)
				os.Exit(1)
			}

			// Close filtered alerts if enabled
			if TheConf.CloseAlerts {
				fmt.Println("\n=== Closing Filtered Alerts ===")
				CloseAlerts(filteredAlerts, TheConf, client)
			}
		} else {
			fmt.Println("No alerts matched the filters")
		}
	}
}

type PageResult struct {
	PageNum int
	Alerts  []Alert
	Err     error
}

func fetchPage(client *http.Client, config JsonConfig, pageNum int) PageResult {
	result := PageResult{PageNum: pageNum}

	fullURL := BuildURL(config, pageNum)
	if config.Debug {
		fmt.Printf("Fetching page %d: %s\n", pageNum, fullURL)
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		result.Err = fmt.Errorf("request error on page %d: %w", pageNum, err)
		return result
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		result.Err = fmt.Errorf("HTTP error on page %d: %w", pageNum, err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		result.Err = fmt.Errorf("HTTP %d on page %d: %s", resp.StatusCode, pageNum, string(body))
		return result
	}

	var chunk AlertsFile
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&chunk); err != nil {
		result.Err = fmt.Errorf("JSON error on page %d: %w", pageNum, err)
		return result
	}

	result.Alerts = chunk.Alerts
	return result
}

func fetchAllAlertsParallel(config JsonConfig, client *http.Client) []Alert {
	var allAlerts []Alert
	var mu sync.Mutex
	var wg sync.WaitGroup

	currentPage := config.PageNumber
	semaphore := make(chan struct{}, config.MaxConcurrentPages)
	results := make(chan PageResult, config.MaxConcurrentPages)
	done := make(chan bool)

	go func() {
		for result := range results {
			if result.Err != nil {
				fmt.Printf("Error: %v\n", result.Err)
				continue
			}

			mu.Lock()
			allAlerts = append(allAlerts, result.Alerts...)
			mu.Unlock()

			if config.Debug {
				fmt.Printf("Page %d completed: %d alerts\n", result.PageNum, len(result.Alerts))
			}

			if len(result.Alerts) == 100 {
				wg.Add(1)
				go func(page int) {
					semaphore <- struct{}{}
					fetchResult := fetchPage(client, config, page)
					results <- fetchResult
					<-semaphore
					wg.Done()
				}(result.PageNum + config.MaxConcurrentPages)
			}
		}
		done <- true
	}()

	for i := 0; i < config.MaxConcurrentPages; i++ {
		wg.Add(1)
		go func(page int) {
			semaphore <- struct{}{}
			result := fetchPage(client, config, page)
			results <- result
			<-semaphore
			wg.Done()
		}(currentPage + i)
	}

	wg.Wait()
	close(results)
	<-done

	fmt.Printf("Total alerts fetched: %d\n", len(allAlerts))
	return allAlerts
}

func fetchAllAlertsParallelWithFlush(config JsonConfig, client *http.Client, flushMgr *FlushManager) []Alert {
	var allAlerts []Alert
	var mu sync.Mutex
	var wg sync.WaitGroup

	currentPage := config.PageNumber
	semaphore := make(chan struct{}, config.MaxConcurrentPages)
	results := make(chan PageResult, config.MaxConcurrentPages)
	done := make(chan bool)

	go func() {
		for result := range results {
			if result.Err != nil {
				fmt.Printf("Error: %v\n", result.Err)
				continue
			}

			mu.Lock()
			allAlerts = append(allAlerts, result.Alerts...)
			mu.Unlock()

			if err := flushMgr.AddAlerts(result.Alerts); err != nil {
				fmt.Printf("Flush error: %v\n", err)
			}

			if config.Debug {
				fmt.Printf("Page %d completed: %d alerts (total in memory: %d)\n", result.PageNum, len(result.Alerts), len(allAlerts))
			}

			if len(result.Alerts) == 100 {
				wg.Add(1)
				go func(page int) {
					semaphore <- struct{}{}
					fetchResult := fetchPage(client, config, page)
					results <- fetchResult
					<-semaphore
					wg.Done()
				}(result.PageNum + config.MaxConcurrentPages)
			}
		}
		done <- true
	}()

	for i := 0; i < config.MaxConcurrentPages; i++ {
		wg.Add(1)
		go func(page int) {
			semaphore <- struct{}{}
			result := fetchPage(client, config, page)
			results <- result
			<-semaphore
			wg.Done()
		}(currentPage + i)
	}

	wg.Wait()
	close(results)
	<-done

	fmt.Printf("Total alerts fetched: %d\n", len(allAlerts))
	return allAlerts
}
