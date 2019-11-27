---
subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules"
sidebar_current: "docs-azurerm-datasource-monitor-scheduled-query-rules"
description: |-
  Get information about the specified Scheduled Query Rule.
---

# Data Source: azurerm_monitor_scheduled_query_rules

Use this data source to access the properties of a Scheduled Query Rule.

## Example Usage

```hcl
data "azurerm_monitor_scheduled_query_rules" "example" {
  resource_group_name = "terraform-example-rg"
  name                = "tfex-queryrule"
}

output "query_rule_id" {
  value = "${data.azurerm_monitor_scheduled_query_rules.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Scheduled Query Rule.
* `resource_group_name` - (Required) Specifies the name of the resource group the Scheduled Query Rule is located in.

## Attributes Reference

* `id` - The ID of the Scheduled Query Rule.
* `action` - An `action` block as defined below. Defines the action to be taken when the rule is run.
* `action_type` - Must equal ether `AlertingAction` or `LogToMetricAction`.
* `authorized_resources` - List of Resource IDs referred into query.
* `custom_webhook_payload` - Custom payload to be sent for all webhook URI in Azure action group.
* `data_source_id` - The resource uri over which log search query is to be run.
* `description` - The description of the Scheduled Query Rule.
* `email_subject` - Custom subject override for all email ids in Azure action group.
* `enabled` - Whether this scheduled query rule is enabled.
* `frequency` - Frequency (in minutes) at which rule condition should be evaluated.
* `query` - Log search query. Required for action type - `alerting_action`.
* `query_type` - Must equal "ResultCount".
* `time_window` - Time window for which data needs to be fetched for query (should be greater than or equal to frequency_in_minutes).

---

`action` supports the following if `action_type` is `AlertingAction`:

* `azns_action` - An `azns_action` block as defined below.
* `severity` - Severity of the alert. Possible values include: 'Zero', 'One', 'Two', 'Three', 'Four'.
* `throttling` - Time (in minutes) for which Alerts should be throttled or suppressed.
* `trigger` - A `trigger` block as defined below. The condition that results in the alert rule being run.

---

`action` supports the following if `action_type` is `LogToMetricAction`:

* `criteria` - A `criteria` block as defined below.
* `metric_name` - Name of the metric.

---

`azns_action` supports the following:

* `action_group` - List of action group reference resource IDs.
* `custom_webhook_payload` - (Optional) Custom payload to be sent for all webhook payloads in alerting action.
* `email_subject` - Email subject line.

---

`criteria` supports the following:

* `dimension` - A `dimension` block as defined below.
* `metric_name` - Name of the metric

---

`dimension` supports the following:

* `name` - Name of the dimension.
* `operator` - Operator for dimension values, - 'Exclude' or 'Include'.
* `values` - List of dimension values.

---

`metricTrigger` supports the following:

* `metricColumn` - Evaluation of metric on a particular column.
* `metricTriggerType` - Metric Trigger Type - 'Consecutive' or 'Total'.
* `operator` - Evaluation operation for rule - 'Equal', 'GreaterThan' or 'LessThan'.
* `threshold` - The threshold of the metric trigger.

---

`trigger` supports the following:

* `metricTrigger` - A `metricTrigger` block as defined above. Trigger condition for metric query rule.
* `operator` - Evaluation operation for rule - 'Equal', 'GreaterThan' or 'LessThan'.
* `threshold` - Result or count threshold based on which rule should be triggered.
