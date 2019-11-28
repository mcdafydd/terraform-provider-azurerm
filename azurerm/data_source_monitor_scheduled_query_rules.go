package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMonitorScheduledQueryRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMonitorScheduledQueryRulesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"action_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azns_action": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_group": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     schema.TypeString,
						},
						"custom_webhook_payload": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_subject": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"criteria": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dimension": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     schema.TypeString,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"data_source_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"frequency": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"throttling": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"time_window": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"trigger": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_trigger": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_column": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     schema.TypeString,
									},
									"metric_trigger_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmMonitorScheduledQueryRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Monitor.ScheduledQueryRulesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Scheduled Query Rule %q was not found", name)
		}
		return fmt.Errorf("Error reading Scheduled Query Rule: %+v", err)
	}

	d.SetId(*resp.ID)
	d.Set("description", resp.Description)
	d.Set("enabled", resp.Enabled)

	if action, ok := resp.Action.(*insights.AlertingAction); ok {
		d.Set("action_type", "AlertingAction")
		d.Set("severity", string(action.Severity))
		d.Set("throttling", *action.ThrottlingInMin)

		flattenAzureRmScheduledQueryRulesAlertingAction(action)
		if err := d.Set("trigger", flattenAzureRmScheduledQueryRulesAlertingAction(action.Trigger)); err != nil {
			return fmt.Errorf("Error setting `trigger`: %+v", err)
		}
	}

	if action, ok := resp.Action.(*insights.LogToMetricAction); ok {
		d.Set("action_type", "LogToMetricAction")

		flattenAzureRmScheduledQueryRulesAlertingAction(action)

		if err := d.Set("criteria", flattenAzureRmScheduledQueryRulesCriteria(action.Criteria)); err != nil {
			return fmt.Errorf("Error setting `criteria`: %+v", err)
		}
	}

	if schedule := resp.Schedule; schedule != nil {
		if schedule.FrequencyInMinutes != nil {
			d.Set("frequency", *schedule.FrequencyInMinutes)
		}
		if schedule.TimeWindowInMinutes != nil {
			d.Set("time_window", *schedule.TimeWindowInMinutes)
		}
	}

	if source := resp.Source; source != nil {
		if err := d.Set("source", flattenAzureRmScheduledQueryRulesSource(resp.Source)); err != nil {
			return fmt.Errorf("Error setting `source`: %+v", err)
		}
	}

	// read-only props
	d.Set("last_updated_time", resp.LastUpdatedTime)
	d.Set("provisioning_state", resp.ProvisioningState)

	return nil
}
