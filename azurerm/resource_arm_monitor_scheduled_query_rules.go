package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorScheduledQueryRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorScheduledQueryRulesCreateUpdate,
		Read:   resourceArmMonitorScheduledQueryRulesRead,
		Update: resourceArmMonitorScheduledQueryRulesCreateUpdate,
		Delete: resourceArmMonitorScheduledQueryRulesDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"action_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AlertingAction",
					"LogToMetricAction",
				}, false),
			},
			"azns_action": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_group": {
							Type:         schema.TypeSet,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"custom_webhook_payload": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.URLIsHTTPOrHTTPS,
						},
						"email_subject": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"authorized_resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
			"data_source_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"frequency": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"lastUpdatedTime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioningState": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ResultCount",
				ValidateFunc: validation.StringInSlice([]string{
					"ResultCount",
				}, false),
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"0",
					"1",
					"2",
					"3",
					"4",
				}, false),
			},
			"throttling": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"time_window": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"trigger": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operator": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GreaterThan",
								"LessThan",
								"Equal",
							}, false),
						},
						"threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMonitorScheduledQueryRulesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Monitor.ScheduledQueryRulesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Monitor Scheduled Query Rules %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_scheduled_query_rules", *existing.ID)
		}
	}

	actionType := d.Get("action_type").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(insights.Enabled)

	location := azure.NormalizeLocation(d.Get("location").(string))

	action := insights.BasicAction{}
	switch actionType {
	case "AlertingAction":
		action = expandMonitorScheduledQueryRulesAlertingAction(d)
	case "LogToMetricAction":
		action = expandMonitorScheduledQueryRulesLogToMetricAction(d)
	default:
		return fmt.Errorf("Invalid action_type %q. Value must be either 'AlertingAction' or 'LogToMetricAction'", action_type)
	}

	source, _ := expandMonitorScheduledQueryRulesSource(d)
	schedule, _ := expandMonitorScheduledQueryRulesSchedule(d)

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := insights.LogSearchRuleResource{
		Location: utils.String(location),
		LogSearchRule: &insights.LogSearchRule{
			Description: utils.String(description),
			Enabled:     enabled,
			Source:      &source,
			Schedule:    &schedule,
			Action:      action,
		},
		Tags: expandedTags,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating or updating scheduled query rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Scheduled query rule %q (resource group %q) ID is empty", name, resourceGroup)
	}
	d.SetId(*read.ID)

	return resourceArmMonitorScheduledQueryRulesRead(d, meta)
}

func resourceArmMonitorScheduledQueryRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Monitor.ScheduledQueryRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ScheduledQueryRules"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Scheduled Query Rule %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting scheduled query rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if rule := resp.LogSearchRule; rule != nil {
		d.Set("enabled", rule.Enabled)
		d.Set("description", rule.Description)
		if err := d.Set("source", flattenAzureRmScheduledQueryRulesSource(rule.Source)); err != nil {
			return fmt.Errorf("Error setting `source`: %+v", err)
		}
		if err := d.Set("schedule", flattenAzureRmScheduledQueryRulesSchedule(rule.Schedule)); err != nil {
			return fmt.Errorf("Error setting `schedule`: %+v", err)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMonitorScheduledQueryRulesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Monitor.ScheduledQueryRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ScheduledQueryRules"]

	if resp, err := client.Delete(ctx, resourceGroup, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting scheduled query rule %q (resource group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandMonitorScheduledQueryRulesAlertingAction(d *schema.ResourceData) *insights.AlertingAction {
	aznsAction := d.Get("azns_action").(string).List()
	severity := d.Get("severity").(string)
	throttling := d.Get("throttling").(int)
	trigger := d.Get("trigger").(*schema.Set).List()

	action := insights.AlertingAction{
		AznsAction: &aznsAction,
		Severity: severity,
		ThrottlingInMin: utils.Int(throttling),
		Trigger: &trigger,
		OdataType: insights.OdataTypeMicrosoftWindowsAzureManagementMonitoringAlertsModelsMicrosoftAppInsightsNexusDataContractsResourcesScheduledQueryRulesAlertingAction,
	}

	return &action
}

func expandMonitorScheduledQueryRulesCriteria(input []interface{}) (*[]insights.Criteria, error) {
	actionType := d.Get("action_type").(string)

	if actionType != "LogToMetricAction" {
		return nil, fmt.Errorf("Criteria only supported if action_type is 'LogToMetricAction'")
	}

	criteria := make([]insights.Criteria, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		dimensions := make([]insights.Dimension, 0)
		for _, dimension := range v["dimension"].([]interface{}) {
			dVal := dimension.(map[string]interface{})
			dimensions = append(dimensions, insights.Dimension{
				Name:     utils.String(dVal["name"].(string)),
				Operator: utils.String(dVal["operator"].(string)),
				Values:   utils.ExpandStringSlice(dVal["values"].([]interface{})),
			})
		}

		criteria = append(criteria, insights.Criteria{
			MetricName: utils.String(v["metric_name"].(string)),
			Dimensions: &dimensions,
		})
	}
	return &criteria
}

func expandMonitorScheduledQueryRulesLogToMetricAction(d *schema.ResourceData) insights.LogToMetricAction {
	criteria := d.Get("criteria").(string).List()

	action := insights.LogToMetricAction{
		Criteria: &criteria,
		OdataType: insights.OdataTypeMicrosoftWindowsAzureManagementMonitoringAlertsModelsMicrosoftAppInsightsNexusDataContractsResourcesScheduledQueryRulesLogToMetricAction,
	}

	return &action
}

func expandMonitorScheduledQueryRulesSchedule(d *schema.ResourceData) *insights.Schedule {
	actionType := d.Get("action_type").(string)

	if actionType != "AlertingAction" {
		return nil, fmt.Errorf("'frequency' and 'time_window' only supported if action_type is 'AlertingAction'")
	}

	frequency := d.Get("frequency").(*schema.Set).List()
	timeWindow := d.Get("time_window").(int)

	schedule := insights.Schedule{
		FrequencyInMin: utils.Int(frequency),
		TimeWindowInMin: utils.Int(timeWindow),
	}

	return &schedule
}

func expandMonitorScheduledQueryRulesSource(d *schema.ResourceData) *insights.Source {
	authorizedResources := d.Get("authorized_resources").(*schema.Set).List()
	dataSourceID := d.Get("data_source_id").(string)
	query := d.Get("query").(string)
	queryType := d.Get("query_type").(string)

	source := insights.Source{
		AuthorizedResources: utils.ExpandStringSlice(authorizedResources.([]interface{})),
		DataSourceID: utils.String(dataSourceID),
		Query: utils.String(query),
		QueryType: queryType,
	}

	return &source
}

func flattenAzureRmScheduledQueryRulesAznsAction(input *insights.AzNsActionGroup) []interface{} {
	result := make([]interface{}, 0)
	v := make(map[string]interface{})

	if input != nil {
		if input.ActionGroup != nil {
			v["action_group"] = *input.ActionGroup
		}
		v["custom_webhook_payload"] = *input.CustomWebhookPayload
		v["email_subject"] = *input.EmailSubject
	}
	result = append(result, v)

	return result
}

func flattenAzureRmScheduledQueryRulesCriteria(input *[]insights.Criteria) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, criteria := range *input {
			v := make(map[string]interface{})
			dimension := make(map[string]interface{})

			v["dimension"] = dimension
			v["metric_name"] = *criteria.MetricName

			result = append(result, v)
		}
	}

	return result
}

func flattenAzureRmScheduledQueryRulesSchedule(input *insights.Schedule) []interface{} {
	result := make(map[string]interface{})

	if input == nil {
		return []interface{}{}
	}

	if input.FrequencyInMinutes != nil {
		result["frequency_in_minutes"] = *input.FrequencyInMinutes
	}

	if input.TimeWindowInMinutes != nil {
		result["time_window_in_minutes"] = *input.TimeWindowInMinutes
	}

	return []interface{}{result}
}

func flattenAzureRmScheduledQueryRulesSource(input *insights.Source) []interface{} {
	result := make(map[string]interface{})

	if input.AuthorizedResources != nil {
		result["authorized_resources"] = *input.AuthorizedResources
	}
	if input.DataSourceID != nil {
		result["data_source_id"] = *input.DataSourceID
	}
	if input.Query != nil {
		result["query"] = *input.Query
	}
	result["query_type"] = input.QueryType

	return []interface{}{result}
}

func flattenAzureRmScheduledQueryRulesTrigger(input *insights.TriggerCondition) []interface{} {
	result := make(map[string]interface{})

	if input.MetricTrigger != nil {
		result["metric_trigger"] = *input.MetricTrigger
	}
	if input.ThresholdOperator != "" {
		result["operator"] = input.ThresholdOperator
	}
	if input.Threshold != nil {
		result["threshold"] = *input.Threshold
	}

	return []interface{}{result}
}
