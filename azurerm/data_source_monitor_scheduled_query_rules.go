package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frequency_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"time_window_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity": {
							Type:     schema.TypeInt,
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
									"email_subject": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"trigger": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
	// set required props for creation
	d.Set("description", resp.Description)
	d.Set("enabled", resp.Enabled)

	// read-only props
	d.Set("type", *resp.Type)
	d.Set("last_updated_time", resp.LastUpdatedTime)
	d.Set("provisioning_state", resp.ProvisioningState)

	//optional props
	if err := d.Set("action", flattenAzureRmScheduledQueryRulesAction(resp.Action)); err != nil {
		return fmt.Errorf("Error setting `action`: %+v", err)
	}

	if err := d.Set("schedule", flattenAzureRmScheduledQueryRulesSchedule(resp.Schedule)); err != nil {
		return fmt.Errorf("Error setting `schedule`: %+v", err)
	}

	if err := d.Set("source", flattenAzureRmScheduledQueryRulesSource(resp.Source)); err != nil {
		return fmt.Errorf("Error setting `source`: %+v", err)
	}

	return nil
}