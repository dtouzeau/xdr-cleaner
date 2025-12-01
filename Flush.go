package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type FlushManager struct {
	mu            sync.Mutex
	outfile       string
	flushEvery    int
	currentAlerts []Alert
	totalFlushed  int
	debug         bool
	firstWrite    bool
}

func NewFlushManager(outfile string, flushEvery int, debug bool) *FlushManager {
	return &FlushManager{
		outfile:       outfile,
		flushEvery:    flushEvery,
		currentAlerts: make([]Alert, 0, flushEvery),
		debug:         debug,
		firstWrite:    true,
	}
}

func (fm *FlushManager) AddAlerts(alerts []Alert) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.currentAlerts = append(fm.currentAlerts, alerts...)

	if len(fm.currentAlerts) >= fm.flushEvery {
		return fm.flush()
	}

	return nil
}

func (fm *FlushManager) flush() error {
	if len(fm.currentAlerts) == 0 {
		return nil
	}

	if fm.debug {
		fmt.Printf("Flushing %d alerts to disk (total flushed: %d)...\n", len(fm.currentAlerts), fm.totalFlushed+len(fm.currentAlerts))
	}

	var file *os.File
	var err error

	if fm.firstWrite {
		file, err = os.Create(fm.outfile)
		if err != nil {
			return fmt.Errorf("create file error: %w", err)
		}
		_, _ = file.WriteString("{\n  \"Alerts\": [\n")
		fm.firstWrite = false
	} else {
		file, err = os.OpenFile(fm.outfile, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("open file error: %w", err)
		}
	}
	defer file.Close()

	for i, alert := range fm.currentAlerts {
		alertJSON, err := json.MarshalIndent(alert, "    ", "  ")
		if err != nil {
			return fmt.Errorf("marshal error: %w", err)
		}

		if fm.totalFlushed > 0 || i > 0 {
			_, _ = file.WriteString(",\n")
		}

		_, _ = file.WriteString("    ")
		_, _ = file.Write(alertJSON)
	}

	fm.totalFlushed += len(fm.currentAlerts)
	fm.currentAlerts = fm.currentAlerts[:0]

	return nil
}

func (fm *FlushManager) Finalize() error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if err := fm.flush(); err != nil {
		return err
	}

	file, err := os.OpenFile(fm.outfile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("finalize open error: %w", err)
	}
	defer file.Close()

	_, _ = file.WriteString("\n  ]\n}\n")

	if fm.debug {
		fmt.Printf("Finalized: total %d alerts written to %s\n", fm.totalFlushed, fm.outfile)
	}

	return nil
}

func (fm *FlushManager) GetTotalFlushed() int {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	return fm.totalFlushed + len(fm.currentAlerts)
}
