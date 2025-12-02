# Query Filters Documentation

## Overview
The `queryFilters` field in the configuration allows you to add custom query parameters to the API URL. This supports various filter operators including `_contains`, `_ne` (not equal), and other API-specific operators.

## Configuration Example

To use query filters, add the `queryFilters` object to your `config.json`:

```json
{
  "pageNumber": 1,
  "tenantID": "3d6c4203-7328-4fb2-98fa-f37d385ffbde",
  "token": "your-api-token-here",
  "baseURL": "https://api.louni.priv/xdr/api/v1/alerts",
  "debug": true,
  "maxConcurrentPages": 50,
  "queryFilters": {
    "name_contains": "R219_01",
    "observable_value_contains": "62.201.149",
    "status_ne": "Closed"
  }
}
```
## Generated URL

With the above configuration, the application will generate URLs like:
```
https://your.xdr.server/xdr/api/v1/alerts?tenantID=3d6c4203-7328-4fb2-98fa-f37d385ffbde&status_ne=Closed&name_contains=R219_01&observable_value_contains=62.201.149&page=1
```

## Supported Filter Operators

The `queryFilters` field supports any filter operator that your API accepts. Common operators include:

- **`_contains`** - Field contains the specified value (partial match)
  - Example: `"name_contains": "R219_01"`

- **`_ne`** - Field not equal to the specified value
  - Example: `"status_ne": "Closed"`

- **`_gt`** - Field greater than the specified value
  - Example: `"severity_gt": "3"`

- **`_lt`** - Field less than the specified value
  - Example: `"severity_lt": "8"`

- **`_gte`** - Field greater than or equal to
  - Example: `"created_gte": "2025-01-01"`

- **`_lte`** - Field less than or equal to
  - Example: `"created_lte": "2025-12-31"`

## Multiple Filters Example

You can combine multiple filters to create complex queries:

```json
{
  "queryFilters": {
    "name_contains": "suspicious",
    "observable_value_contains": "192.168",
    "severity_gt": "5",
    "status_ne": "Closed",
    "source_contains": "firewall"
  }
}
```

## Notes

1. All filter keys and values are URL-encoded automatically by the application
2. Empty string values are ignored (not added to the URL)
3. The `queryFilters` field is optional - if not present, the application works as before
4. Query filters are appended to the standard parameters (tenantID, page, etc.)
5. The default `status_ne=Closed` is still added automatically unless overridden

## Complete Configuration Example

```json
{
  "pageNumber": 1,
  "ids": "",
  "tenantID": "3d6c4203-7328-4fb2-98fa-f37d385ffbde",
  "token": "your-bearer-token",
  "fromDate": "2025-01-01T00:00:00Z",
  "toDate": "2025-12-31T23:59:59Z",
  "status": "",
  "withEvents": "true",
  "withAffected": "true",
  "withHistory": "false",
  "outfile": "/var/log/xdr/out.log",
  "baseURL": "https://api.louni.priv/xdr/api/v1/alerts",
  "debug": true,
  "maxConcurrentPages": 10,
  "filterMode": true,
  "filteredOutfile": "/var/log/xdr/filtered.json",
  "filters": [
    {
      "field": "name",
      "value": "R219_01"
    }
  ],
  "closeAlerts": false,
  "closeReason": "falsePositive",
  "flushEvery": 1000,
  "queryFilters": {
    "name_contains": "R219_01",
    "observable_value_contains": "62.201.149"
  }
}
```

## Debugging

Enable `"debug": true` in your configuration to see the generated URLs in the output:
```
Fetching page 1: https://api.louni.priv/xdr/api/v1/alerts?name_contains=R219_01&observable_value_contains=62.201.149&page=1&status_ne=Closed&tenantID=3d6c4203-7328-4fb2-98fa-f37d385ffbde
```
