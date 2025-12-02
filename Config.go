package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type JsonConfig struct {
	PageNumber         int               `json:"pageNumber"`
	Ids                string            `json:"ids"`
	TenantID           string            `json:"tenantID"`
	Token              string            `json:"token"`
	FromDate           string            `json:"fromDate"`
	ToDate             string            `json:"toDate"`
	Status             string            `json:"status"`
	WithEvents         string            `json:"withEvents"`
	WithAffected       string            `json:"withAffected"`
	WithHistory        string            `json:"withHistory"`
	Outfile            string            `json:"outfile"`
	BaseURL            string            `json:"baseURL"`
	Debug              bool              `json:"debug"`
	MaxConcurrentPages int               `json:"maxConcurrentPages"`
	FilterMode         bool              `json:"filterMode"`
	FilteredOutfile    string            `json:"filteredOutfile"`
	Filters            []Filter          `json:"filters"`
	CloseAlerts        bool              `json:"closeAlerts"`
	CloseReason        string            `json:"closeReason"`
	FlushEvery         int               `json:"flushEvery"`
	QueryFilters       map[string]string `json:"queryFilters"`
}

type Filter struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func ConfigPath() string {
	exePath, _ := os.Executable()
	CurDir := DirName(exePath)
	return filepath.Join(CurDir, "config.json")
}

func LoadConfig() JsonConfig {

	ConfPath := ConfigPath()
	ebytes, _ := fileGetContentsBytes(ConfPath)
	var config JsonConfig
	_ = json.Unmarshal(ebytes, &config)
	if config.PageNumber == 0 {
		config.PageNumber = 1
	}
	if len(config.Outfile) == 0 {
		config.Outfile = DirName(ConfPath) + "/out.log"
	}
	if len(config.BaseURL) == 0 {
		config.BaseURL = "https://mydomain.com/xdr/api/v1"
	}
	if config.MaxConcurrentPages == 0 {
		config.MaxConcurrentPages = 50
	}
	if len(config.FilteredOutfile) == 0 {
		config.FilteredOutfile = DirName(ConfPath) + "/filtered.json"
	}
	if len(config.CloseReason) == 0 {
		config.CloseReason = "falsePositive"
	}
	if config.FlushEvery == 0 {
		config.FlushEvery = 1000
	}
	if !FileExists(ConfPath) {
		sbytes, _ := json.MarshalIndent(config, "", "\t")
		_ = FilePutContentsBytes(ConfPath, sbytes)
	}
	return config
}
func DirName(sFilepath string) string {
	return filepath.Dir(sFilepath)
}
func fileGetContentsBytes(filename string) ([]byte, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return content, nil
}
func FileExists(spath string) bool {
	spath = strings.TrimSpace(spath)
	if IsLink(spath) {
		return true
	}

	if _, err := os.Stat(spath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
func IsLink(path string) bool {

	info, err := os.Lstat(path)
	if err != nil {
		return false
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return true
	}
	return false
}
func FilePutContentsBytes(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}
