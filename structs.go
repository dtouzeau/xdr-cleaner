package main

type AlertsFile struct {
	Alerts []Alert `json:"Alerts"`
}

type Alert struct {
	Assets                interface{}     `json:"Assets"`
	Assignee              Assignee        `json:"Assignee"`
	CreatedAt             string          `json:"CreatedAt"`
	DetectionTechnologies interface{}     `json:"DetectionTechnologies"`
	ExternalRef           string          `json:"ExternalRef"`
	Extra                 interface{}     `json:"Extra"`
	FirstEventTime        string          `json:"FirstEventTime"`
	HistoryRecords        interface{}     `json:"HistoryRecords"`
	ID                    int             `json:"ID"`
	IncidentID            string          `json:"IncidentID"`
	IncidentLinkType      string          `json:"IncidentLinkType"`
	InternalID            string          `json:"InternalID"`
	IsCII                 bool            `json:"IsCII"`
	LastEventTime         string          `json:"LastEventTime"`
	MITRETactics          interface{}     `json:"MITRETactics"`
	MITRETechniques       interface{}     `json:"MITRETechniques"`
	Name                  string          `json:"Name"`
	Observables           []Observable    `json:"Observables"`
	OriginalEvents        []OriginalEvent `json:"OriginalEvents"`
	Rules                 []Rule          `json:"Rules"`
	Severity              string          `json:"Severity"`
	SourceCreatedAt       string          `json:"SourceCreatedAt"`
	SourceID              string          `json:"SourceID"`
	Status                string          `json:"Status"`
	StatusChangedAt       string          `json:"StatusChangedAt"`
	StatusResolution      string          `json:"StatusResolution"`
	TenantID              string          `json:"TenantID"`
	UpdatedAt             string          `json:"UpdatedAt"`
}

type Assignee struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	Type string `json:"Type"`
}

type Observable struct {
	Details string `json:"Details"`
	Type    string `json:"Type"`
	Value   string `json:"Value"`
}

type OriginalEvent struct {
	N          map[string]interface{} `json:"N"`
	ID         string                 `json:"ID"`
	Name       string                 `json:"Name"`
	Type       int                    `json:"Type"`
	EndTime    int64                  `json:"EndTime"`
	Message    string                 `json:"Message"`
	Priority   int                    `json:"Priority"`
	Severity   string                 `json:"Severity"`
	TenantID   string                 `json:"TenantID"`
	GroupedBy  []string               `json:"GroupedBy"`
	ServiceID  string                 `json:"ServiceID"`
	StartTime  int64                  `json:"StartTime"`
	Timestamp  int64                  `json:"Timestamp"`
	BaseEvents []BaseEvent            `json:"BaseEvents"`
}

type BaseEvent struct {
	ID                  string `json:"ID"`
	Name                string `json:"Name"`
	Type                int    `json:"Type"`
	EndTime             int64  `json:"EndTime"`
	Message             string `json:"Message"`
	Priority            int    `json:"Priority"`
	Severity            string `json:"Severity"`
	TenantID            string `json:"TenantID"`
	ServiceID           string `json:"ServiceID"`
	Timestamp           int64  `json:"Timestamp"`
	ExternalID          string `json:"ExternalID"`
	SourcePort          int    `json:"SourcePort"`
	ServiceName         string `json:"ServiceName"`
	DeviceAction        string `json:"DeviceAction"`
	DeviceVendor        string `json:"DeviceVendor"`
	DeviceAddress       string `json:"DeviceAddress"`
	DeviceProduct       string `json:"DeviceProduct"`
	SourceAddress       string `json:"SourceAddress"`
	DeviceHostName      string `json:"DeviceHostName"`
	DeviceTimeZone      string `json:"DeviceTimeZone"`
	DestinationPort     int    `json:"DestinationPort"`
	DeviceExternalID    string `json:"DeviceExternalID"`
	DeviceReceiptTime   int64  `json:"DeviceReceiptTime"`
	TransportProtocol   string `json:"TransportProtocol"`
	DestinationAddress  string `json:"DestinationAddress"`
	DeviceEventClassID  string `json:"DeviceEventClassID"`
	ApplicationProtocol string `json:"ApplicationProtocol"`
	DeviceEventCategory string `json:"DeviceEventCategory"`
}

type Rule struct {
	Confidence string `json:"Confidence"`
	Custom     bool   `json:"Custom"`
	ID         string `json:"ID"`
	InternalID string `json:"InternalID"`
	Name       string `json:"Name"`
	Severity   string `json:"Severity"`
	Type       string `json:"Type"`
}
