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

resource "azurerm_log_analytics_workspace" "example" {
  name                = format("%s-logs", var.prefix)
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

# Example: AlertingAction
resource "azurerm_scheduled_query_rule" "example" {
  name                   = format("%s-queryrule", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name

  action_type              = "AlertingAction"
  azns_action              = {
    action_group           = []
    email_subject          = "Email Header"
    custom_webhook_payload = {}
  }
  data_source_id           = azurerm_application_insights.example.id
  description              = "Scheduled query rule AlertingAction example"
  enabled                  = true
  frequency                = 5
  query                    = "requests | where status_code >= 500 | summarize AggregatedValue = count() by bin(TimeGenerated, 5m)"
  query_type               = "ResultCount"
  severity                 = "1"
  time_window              = 30
  trigger                  = {
    threshold_operator     = "GreaterThan"
    threshold              = 3
  }
}

# Example: AlertingAction Cross-Resource
resource "azurerm_scheduled_query_rule" "example2" {
  name                   = format("%s-queryrule2", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name

  action_type              = "AlertingAction"
  authorized_resources     = [azurerm_application_insights.example.id,
                              azurerm_log_analytics_workspace.example.id]
  azns_action              = {
    action_group           = []
    email_subject          = "Email Header"
    custom_webhook_payload = {}
  }
  data_source_id           = azurerm_application_insights.example.id
  description              = "Scheduled query rule AlertingAction cross-resource example"
  enabled                  = true
  frequency                = 5
  query                    = "union requests, workspace(${azurerm_log_analytics_workspace.example.name}).Heartbeat"
  query_type               = "ResultCount"
  severity                 = "1"
  time_window              = 30
  trigger                  = {
    threshold_operator     = "GreaterThan"
    threshold              = 3
  }
}

# Example: LogToMetricAction
resource "azurerm_scheduled_query_rule" "example3" {
  name                   = format("%s-queryrule3", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name

  action_type            = "LogToMetricAction"
  criteria               = [{
      metric_name        = "Average_% Idle Time"
      dimensions         = [{
        name             = "dimension"
        operator         = "GreaterThan"
        values           = ["latency"]
      }]
  }]
  data_source_id         = azurerm_application_insights.example.id
  description            = "Scheduled query rule LogToMetric example"
  enabled                = true
}
