resource "azurerm_resource_group" "example" {
  name     = format("%s-resources", var.prefix)
  location = var.location
}

resource "azurerm_resource_group" "example" {
  name     = "monitoring-resources"
  location = "West US"
}
resource "azurerm_application_insights" "example" {
  name                = format("%s-insights", var.prefix)
  location            = var.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}
resource "azurerm_application_insights" "example2" {
  name                = format("%s-insights2", var.prefix)
  location            = var.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

# Example: AlertingAction
resource "azurerm_scheduled_query_rule" "example" {
  name                   = format("%s-queryrule", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name
  "azns_action": {
    "action_group": [],
    "email_subject": "Email Header",
    "custom_webhook_payload": "{}"
  },
  "severity": "1",
  "trigger": {
    "threshold_operator": "GreaterThan",
    "threshold": 3
  }
  action_type            = "AlertingAction"
  data_source_id         = azurerm_application_insights.example.id
  description            = "Scheduled query rule AlertingAction example"
  enabled                = true
  frequency              = 5
  time_window            = 30
}

# Example: AlertingAction Cross-Resource
resource "azurerm_scheduled_query_rule" "example" {
  name                   = format("%s-queryrule", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name
  "azns_action": {
    "action_group": [],
    "email_subject": "Email Header",
    "custom_webhook_payload": "{}"
  },
  "severity": "1",
  "trigger": {
    "threshold_operator": "GreaterThan",
    "threshold": 3
  }
  action_type            = "AlertingAction"
  authorized_resources   = [azurerm_application_insights.example.id]
  data_source_id         = azurerm_application_insights.example.id
  description            = "Scheduled query rule AlertingAction cross-resource example"
  enabled                = true
  frequency              = 5
  query                  = "requests | where status_code >= 500 | summarize AggregatedValue = count() by bin(TimeGenerated, 5m)"
  query_type             = "ResultCount"
  time_window            = 30
}

# Example: LogToMetricAction
resource "azurerm_scheduled_query_rule" "example" {
  name                   = format("%s-queryrule", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name
  "criteria": [
    {
      "metric_name": "Average_% Idle Time",
      "dimensions": [{
        name             = "dimension
        operator         = "GreaterThan"
        values           = ["latency"]
      }]
    }
  ]
  action_type            = "LogToMetricAction"
  data_source_id         = azurerm_application_insights.example.id
  description            = "Scheduled query rule LogToMetric example"
  enabled                = true
}
