---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules"
sidebar_current: "docs-azurerm-resource-monitor-scheduled-query-rules"
description: |-
  Manages a Scheduled Query Rule within Azure Monitor
---

# azurerm_monitor_action_group

Manages a Scheduled Query Rule within Azure Monitor.

## Example Usage

```hcl
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

# AlertingAction
resource "azurerm_scheduled_query_rule" "example" {
  name                   = format("%s-queryrule", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name

  action                 = {
    "azns_action": {
      "action_group": [],
      "email_subject": "Email Header",
      "custom_webhook_payload": "{}"
    },
    "severity": "1",
    "trigger": {
      "threshold_operator": "GreaterThan",
      "threshold": 3
    },
  }
  action_type            = "AlertingAction"
  data_source_id         = azurerm_application_insights.example.id
  description            = "Scheduled query rule AlertingAction example"
  enabled                = true
  frequency              = 5
  time_window            = 30
}

# AlertingAction Cross-Resource
resource "azurerm_scheduled_query_rule" "example" {
  name                   = format("%s-queryrule", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name

  action                 = {
    "azns_action": {
      "action_group": [],
      "email_subject": "Email Header",
      "custom_webhook_payload": "{}"
    },
    "severity": "1",
    "trigger": {
      "threshold_operator": "GreaterThan",
      "threshold": 3
    },
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

# LogToMetricAction
resource "azurerm_scheduled_query_rule" "example" {
  name                   = format("%s-queryrule", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name

  action                 = {
    "criteria": [
      {
        "metric_name": "Average_% Idle Time",
        "dimensions": []
      }
    ],
  }
  action_type            = "AlertingAction"
  data_source_id         = azurerm_application_insights.example.id
  description            = "Scheduled query rule LogToMetric example"
  enabled                = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Action Group. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the Action Group instance.
* `action` - (Required) An `action` block as defined below.
* `action_type` - (Required) Must equal ether `AlertingAction` or `LogToMetricAction`.
* `authorized_resources` - (Optional) List of Resource IDs referred into query.
* `data_source_id` - (Required) The resource uri over which log search query is to be run.
* `description` - (Optional) The description of the Scheduled Query Rule.
* `enabled` - (Optional) Whether this scheduled query rule is enabled.  Default is `true`.
* `frequency` - (Optional) Frequency (in minutes) at which rule condition should be evaluated.
* `query` - (Required) Log search query. Required for action type - `alerting_action`.
* `query_type` - (Required) Must equal "ResultCount" for now.
* `time_window` - (Optional) Time window for which data needs to be fetched for query (should be greater than or equal to frequency_in_minutes).

---

`action` supports the following if `action_type` is `AlertingAction`:

* `azns_action` - (Required) An `azns_action` block as defined below.
* `severity` - (Optional) Severity of the alert. Possible values include: 'Zero', 'One', 'Two', 'Three', 'Four'.
* `throttling` - (Optional) Time (in minutes) for which Alerts should be throttled or suppressed.
* `trigger` - (Required) The condition that results in the alert rule being run.

---

`action` supports the following if `action_type` is `LogToMetricAction`:

* `criteria` - (Required) A `criteria` block as defined below.
* `metric_name` - (Required) Name of the metric.

---

* `azns_action` supports the following:

* `action_group` - (Optional) List of action group reference resource IDs.
* `custom_webhook_payload` - (Optional) Custom payload to be sent for all webhook payloads in alerting action.
* `email_subject` - (Optional) Custom subject override for all email ids in Azure action group.

---

`criteria` supports the following:

* `dimension` - (Required) A `dimension` block as defined below.
* `metric_name` - (Required) Name of the metric

---

`dimension` supports the following:

* `name` - (Required) Name of the dimension.
* `operator` - (Required) Operator for dimension values, - 'Exclude' or 'Include'.
* `values` - (Required) List of dimension values.

---

`metricTrigger` supports the following:

* `metricColumn` - (Required) Evaluation of metric on a particular column.
* `metricTriggerType` - (Required) Metric Trigger Type - 'Consecutive' or 'Total'.
* `operator` - (Required) Evaluation operation for rule - 'Equal', 'GreaterThan' or 'LessThan'.
* `threshold` - (Required) The threshold of the metric trigger.

---

`trigger` supports the following:

* `metricTrigger` - (Optional) A `metricTrigger` block as defined above. Trigger condition for metric query rule.
* `operator` - (Required) Evaluation operation for rule - 'Equal', 'GreaterThan' or 'LessThan'.
* `threshold` - (Required) Result or count threshold based on which rule should be triggered.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Scheduled Query Rule.
* `last_updated_time` - Last time the rule was updated in IS08601 format.
* `provisioning_state` - Provisioning state of the scheduled query rule. Possible values include: 'Succeeded', 'Deploying', 'Canceled', 'Failed'

## Import

Scheduled Query Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_scheduled_query_rules.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/scheduledQueryRules/myrulename
```
