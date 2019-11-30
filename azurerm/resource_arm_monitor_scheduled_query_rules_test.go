package azurerm

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

/*func TestAccAzureRMMonitorScheduledQueryRules_basic(t *testing.T) {
	resourceName := "azurerm_monitor_activity_log_alert.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorScheduledQueryRules_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
	return
}*/

func TestAccAzureRMMonitorScheduledQueryRules_AlertingAction(t *testing.T) {
	resourceName := "azurerm_monitor_activity_log_alert.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	config := testAccAzureRMMonitorScheduledQueryRules_alertingAction(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
	return
}

func TestAccAzureRMMonitorScheduledQueryRules_AlertingActionCrossResource(t *testing.T) {
	resourceName := "azurerm_monitor_activity_log_alert.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	config := testAccAzureRMMonitorScheduledQueryRules_alertingActionCrossResource(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
	return
}

func TestAccAzureRMMonitorScheduledQueryRules_LogToMetricAction(t *testing.T) {
	resourceName := "azurerm_monitor_activity_log_alert.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	config := testAccAzureRMMonitorScheduledQueryRules_logToMetricAction(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
				),
				PreventDiskCleanup: true,
			},
			{
				PreventDiskCleanup: true,
				ResourceName:       resourceName,
				ImportState:        true,
				ImportStateVerify:  true,
			},
		},
	})
	return
}

/*func TestAccAzureRMMonitorScheduledQueryRules_basicAndCompleteUpdate(t *testing.T) {
	resourceName := "azurerm_monitor_activity_log_alert.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	basicConfig := testAccAzureRMMonitorScheduledQueryRules_basic(ri, location)
	completeConfig := testAccAzureRMMonitorScheduledQueryRules_complete(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorActionGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
				),
			},
			{
				Config: completeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
				),
			},
			{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
				),
			},
		},
	})
	return
}*/
/*
func testAccAzureRMMonitorScheduledQueryRules_basic(rInt int, location string) string {
	return fmt.Sprintf(`
	resource "azurerm_monitor_scheduled_query_rules" "import" {
		name                = "acctestSqr-%d"
		description         = "test alerting action"
		enabled             = true
		action_type         = "Alerting"

		query          = "let data=datatable(id:int, value:string) [1, 'test1', 2, 'testtwo']; data | extend strlen = strlen(value)"
		data_source_id = "${azurerm_log_analytics_workspace.test.id}"
		query_type     = "ResultCount"

		frequency   = 60
		time_window = 60


		severity    = 3
		azns_action {
			action_group = ["${azurerm_monitor_action_group.test.id}"]
			email_subject = "Custom alert email subject"
		}

		trigger {
			operator  = "GreaterThan"
			threshold = 5000
		}
	}

`, rInt, location, rInt)
}*/
func testAccAzureRMMonitorScheduledQueryRules_alertingAction(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestAppInsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "test" {
	name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}

resource "azurerm_monitor_scheduled_query_rules" "test" {
  name                = "acctestsqr-%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
	description         = "test log to metric action"
	enabled             = true
	action_type         = "LogToMetric"

	data_source_id = "${azurerm_application_insights.test.id}"
  query          = "let data=datatable(id:int, value:string) [1, 'test1', 2, 'testtwo']; data | extend strlen = strlen(value)"
	query_type     = "ResultCount"

	frequency   = 60
  time_window = 60

	severity     = 3
	azns_action {
		action_group = ["${azurerm_monitor_action_group.test.id}"]
		email_subject = "Custom alert email subject"
	}

	trigger {
		operator = "GreaterThan"
		threshold         = 5000
		metric_trigger {
			operator            = "GreaterThan"
			threshold           = 5
			metric_trigger_type = "Consecutive"
			metric_column       = "Computer"
		}
	}
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMMonitorScheduledQueryRules_alertingActionCrossResource(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestAppInsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestWorkspace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}

resource "azurerm_monitor_scheduled_query_rules" "test" {
  name                = "acctestsqr-%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
	description         = "test log to metric action"
	enabled             = true
	action_type         = "LogToMetric"

	authorized_resources = ["${azurerm_application_insights.test.id}", "${azurerm_log_analytics_workspace.test.id}"]
	data_source_id       = "${azurerm_application_insights.test.id}"
  query                = "union requests, workspace(${azurerm_log_analytics_workspace.test.name}).Heartbeat"
	query_type           = "ResultCount"

	frequency   = 60
  time_window = 60

	severity     = 3
	azns_action {
		action_group = ["${azurerm_monitor_action_group.test.id}"]
		email_subject = "Custom alert email subject"
	}

	trigger {
		operator = "GreaterThan"
		threshold         = 5000
		metric_trigger {
			operator            = "GreaterThan"
			threshold           = 5
			metric_trigger_type = "Consecutive"
			metric_column       = "Computer"
		}
	}
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMMonitorScheduledQueryRules_logToMetricAction(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestAppInsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "test" {
	name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}

resource "azurerm_monitor_scheduled_query_rules" "test" {
  name                = "acctestsqr-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
	description         = "test log to metric action"
	enabled             = true
	action_type         = "LogToMetric"

	data_source_id = "${azurerm_application_insights.test.id}"

	criteria {
		metric_name        = "Average_percent Idle Time"
		dimension {
			name             = "dimension"
			operator         = "GreaterThan"
			values           = ["latency"]
		}
	}
}
`, rInt, location, rInt, rInt, rInt)
}

func testCheckAzureRMMonitorScheduledQueryRulesDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).Monitor.ScheduledQueryRulesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_activity_log_alert" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Activity log alert still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Activity Log Alert Instance: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).Monitor.ScheduledQueryRulesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on monitorScheduledQueryRulesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Activity Log Alert Instance %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}
