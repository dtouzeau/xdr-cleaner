package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func BuildURL(TheConf JsonConfig, currentPage int) string {
	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", currentPage))
	//params.Set("timestampField", *timestampField)

	if TheConf.Ids != "" {
		for _, v := range strings.Split(TheConf.Ids, ",") {
			params.Add("id", strings.TrimSpace(v))
		}
	}

	for _, v := range strings.Split(TheConf.TenantID, ",") {
		params.Add("tenantID", strings.TrimSpace(v))
	}

	if TheConf.FromDate != "" {
		if _, err := time.Parse(time.RFC3339, TheConf.FromDate); err == nil {
			params.Set("from", TheConf.FromDate)
		}
	}

	if TheConf.ToDate != "" {
		if _, err := time.Parse(time.RFC3339, TheConf.ToDate); err == nil {
			params.Set("to", TheConf.ToDate)
		}
	}

	if TheConf.Status != "" {
		for _, v := range strings.Split(TheConf.Status, ",") {
			params.Add("status", strings.TrimSpace(v))
		}
	}

	if TheConf.WithEvents != "" {
		params.Set("withEvents", TheConf.WithEvents)
	}

	if TheConf.WithAffected != "" {
		params.Set("withAffected", TheConf.WithAffected)
	}

	if TheConf.WithHistory != "" {
		params.Set("withHistory", TheConf.WithHistory)
	}

	return TheConf.BaseURL + "?" + params.Encode()
}
func BuilClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}
