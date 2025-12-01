package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func FilterAlerts(allAlerts []Alert, config JsonConfig) []Alert {
	if !config.FilterMode || len(config.Filters) == 0 {
		return allAlerts
	}

	var filtered []Alert
	for _, alert := range allAlerts {
		if matchesFilters(alert, config.Filters, config.Debug) {
			filtered = append(filtered, alert)
		}
	}

	if config.Debug {
		fmt.Printf("Filtered %d alerts from %d total\n", len(filtered), len(allAlerts))
	}

	return filtered
}

func matchesFilters(alert Alert, filters []Filter, debug bool) bool {
	for _, filter := range filters {
		if !matchFilter(alert, filter, debug) {
			return false
		}
	}
	return true
}

func matchFilter(alert Alert, filter Filter, debug bool) bool {
	parts := strings.Split(filter.Field, "|")
	if len(parts) < 2 {
		if debug {
			fmt.Printf("Invalid filter field format: %s (use 'Section|Field')\n", filter.Field)
		}
		return false
	}

	section := strings.TrimSpace(parts[0])
	field := strings.TrimSpace(parts[1])

	switch section {
	case "Observable":
		return matchObservable(alert.Observables, field, filter.Value, debug)
	case "Rule":
		return matchRule(alert.Rules, field, filter.Value, debug)
	case "BaseEvent":
		return matchBaseEvent(alert.OriginalEvents, field, filter.Value, debug)
	case "Alert":
		return matchAlertField(alert, field, filter.Value, debug)
	default:
		if debug {
			fmt.Printf("Unknown section: %s\n", section)
		}
		return false
	}
}

func matchObservable(observables []Observable, field, value string, debug bool) bool {
	for _, obs := range observables {
		var fieldValue string
		switch strings.ToLower(field) {
		case "value":
			fieldValue = obs.Value
		case "type":
			fieldValue = obs.Type
		case "details":
			fieldValue = obs.Details
		default:
			continue
		}

		if strings.Contains(strings.ToLower(fieldValue), strings.ToLower(value)) {
			if debug {
				fmt.Printf("Matched Observable.%s: %s contains %s\n", field, fieldValue, value)
			}
			return true
		}
	}
	return false
}

func matchRule(rules []Rule, field, value string, debug bool) bool {
	for _, rule := range rules {
		var fieldValue string
		switch strings.ToLower(field) {
		case "name":
			fieldValue = rule.Name
		case "id":
			fieldValue = rule.ID
		case "type":
			fieldValue = rule.Type
		case "severity":
			fieldValue = rule.Severity
		case "confidence":
			fieldValue = rule.Confidence
		default:
			continue
		}

		if strings.Contains(strings.ToLower(fieldValue), strings.ToLower(value)) {
			if debug {
				fmt.Printf("Matched Rule.%s: %s contains %s\n", field, fieldValue, value)
			}
			return true
		}
	}
	return false
}

func matchBaseEvent(events []OriginalEvent, field, value string, debug bool) bool {
	for _, event := range events {
		for _, baseEvent := range event.BaseEvents {
			var fieldValue string
			switch strings.ToLower(field) {
			case "destinationaddress":
				fieldValue = baseEvent.DestinationAddress
			case "sourceaddress":
				fieldValue = baseEvent.SourceAddress
			case "deviceaddress":
				fieldValue = baseEvent.DeviceAddress
			case "devicehostname":
				fieldValue = baseEvent.DeviceHostName
			case "deviceaction":
				fieldValue = baseEvent.DeviceAction
			case "devicevendor":
				fieldValue = baseEvent.DeviceVendor
			case "deviceproduct":
				fieldValue = baseEvent.DeviceProduct
			case "transportprotocol":
				fieldValue = baseEvent.TransportProtocol
			case "applicationprotocol":
				fieldValue = baseEvent.ApplicationProtocol
			case "message":
				fieldValue = baseEvent.Message
			case "name":
				fieldValue = baseEvent.Name
			case "destinationport":
				fieldValue = fmt.Sprintf("%d", baseEvent.DestinationPort)
			case "sourceport":
				fieldValue = fmt.Sprintf("%d", baseEvent.SourcePort)
			default:
				continue
			}

			if strings.Contains(strings.ToLower(fieldValue), strings.ToLower(value)) {
				if debug {
					fmt.Printf("Matched BaseEvent.%s: %s contains %s\n", field, fieldValue, value)
				}
				return true
			}
		}
	}
	return false
}

func matchAlertField(alert Alert, field, value string, debug bool) bool {
	var fieldValue string
	switch strings.ToLower(field) {
	case "name":
		fieldValue = alert.Name
	case "severity":
		fieldValue = alert.Severity
	case "status":
		fieldValue = alert.Status
	case "internalid":
		fieldValue = alert.InternalID
	case "incidentid":
		fieldValue = alert.IncidentID
	case "externalref":
		fieldValue = alert.ExternalRef
	default:
		return false
	}

	if strings.Contains(strings.ToLower(fieldValue), strings.ToLower(value)) {
		if debug {
			fmt.Printf("Matched Alert.%s: %s contains %s\n", field, fieldValue, value)
		}
		return true
	}
	return false
}

func SaveFilteredAlerts(alerts []Alert, filename string) error {
	out := AlertsFile{Alerts: alerts}
	fileData, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON marshal error: %w", err)
	}

	err = FilePutContentsBytes(filename, fileData)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	fmt.Printf("Saved %d filtered alerts to %s\n", len(alerts), filename)
	return nil
}
