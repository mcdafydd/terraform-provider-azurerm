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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMMonitorScheduledQueryRules_basic(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "0"),
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

func TestAccAzureRMMonitorScheduledQueryRules_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

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
				Config:      testAccAzureRMMonitorScheduledQueryRules_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_monitor_activity_log_alert"),
			},
		},
	})
	return
}

func TestAccAzureRMMonitorScheduledQueryRules_AlertingAction(t *testing.T) {
	resourceName := "azurerm_monitor_activity_log_alert.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMMonitorScheduledQueryRules_singleResource(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operation_name", "Microsoft.Storage/storageAccounts/write"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttrSet(resourceName, "criteria.0.resource_id"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "1"),
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
	config := testAccAzureRMMonitorScheduledQueryRules_complete(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is just a test resource."),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operation_name", "Microsoft.Storage/storageAccounts/write"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_provider", "Microsoft.Storage"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_type", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttrSet(resourceName, "criteria.0.resource_group"),
					resource.TestCheckResourceAttrSet(resourceName, "criteria.0.resource_id"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.caller", "user@example.com"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.level", "Error"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.status", "Failed"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
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
	config := testAccAzureRMMonitorScheduledQueryRules_complete(ri, rs, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is just a test resource."),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operation_name", "Microsoft.Storage/storageAccounts/write"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_provider", "Microsoft.Storage"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_type", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttrSet(resourceName, "criteria.0.resource_group"),
					resource.TestCheckResourceAttrSet(resourceName, "criteria.0.resource_id"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.caller", "user@example.com"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.level", "Error"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.status", "Failed"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
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

func TestAccAzureRMMonitorScheduledQueryRules_basicAndCompleteUpdate(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_id", ""),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.caller", ""),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.level", ""),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.status", ""),
					resource.TestCheckResourceAttr(resourceName, "action.#", "0"),
				),
			},
			{
				Config: completeConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is just a test resource."),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.operation_name", "Microsoft.Storage/storageAccounts/write"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_provider", "Microsoft.Storage"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_type", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttrSet(resourceName, "criteria.0.resource_group"),
					resource.TestCheckResourceAttrSet(resourceName, "criteria.0.resource_id"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.caller", "user@example.com"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.level", "Error"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.status", "Failed"),
					resource.TestCheckResourceAttr(resourceName, "action.#", "2"),
				),
			},
			{
				Config: basicConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRulesExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.category", "Recommendation"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_id", ""),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.caller", ""),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.level", ""),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.status", ""),
					resource.TestCheckResourceAttr(resourceName, "action.#", "0"),
				),
			},
		},
	})
	return
}

func testAccAzureRMMonitorScheduledQueryRules_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestScheduledQueryRules-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  scopes              = ["${azurerm_resource_group.test.id}"]

  criteria {
    category = "Recommendation"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMMonitorScheduledQueryRules_requiresImport(rInt int, location string) string {
	template := testAccAzureRMMonitorScheduledQueryRules_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_activity_log_alert" "import" {
  name                = "${azurerm_monitor_activity_log_alert.test.name}"
  resource_group_name = "${azurerm_monitor_activity_log_alert.test.resource_group_name}"
  scopes              = ["${azurerm_resource_group.test.id}"]

  criteria {
    category = "Recommendation"
  }
}
`, template)
}

func testAccAzureRMMonitorScheduledQueryRules_singleResource(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestScheduledQueryRules-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  scopes              = ["${azurerm_resource_group.test.id}"]

  criteria {
    operation_name = "Microsoft.Storage/storageAccounts/write"
    category       = "Recommendation"
    resource_id    = "${azurerm_storage_account.test.id}"
  }

  action {
    action_group_id = "${azurerm_monitor_action_group.test.id}"
  }
}
`, rInt, location, rInt, rString, rInt)
}

func testAccAzureRMMonitorScheduledQueryRules_complete(rInt int, rString, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_monitor_action_group" "test1" {
  name                = "acctestActionGroup1-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag1"
}

resource "azurerm_monitor_action_group" "test2" {
  name                = "acctestActionGroup2-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag2"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_activity_log_alert" "test" {
  name                = "acctestScheduledQueryRules-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  enabled             = true
  description         = "This is just a test resource."

  scopes = [
    "${azurerm_resource_group.test.id}",
    "${azurerm_storage_account.test.id}",
  ]

  criteria {
    operation_name    = "Microsoft.Storage/storageAccounts/write"
    category          = "Recommendation"
    resource_provider = "Microsoft.Storage"
    resource_type     = "Microsoft.Storage/storageAccounts"
    resource_group    = "${azurerm_resource_group.test.name}"
    resource_id       = "${azurerm_storage_account.test.id}"
    caller            = "user@example.com"
    level             = "Error"
    status            = "Failed"
  }

  action {
    action_group_id = "${azurerm_monitor_action_group.test1.id}"
  }

  action {
    action_group_id = "${azurerm_monitor_action_group.test2.id}"

    webhook_properties = {
      from = "terraform test"
      to   = "microsoft azure"
    }
  }
}
`, rInt, location, rInt, rInt, rString, rInt)
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
