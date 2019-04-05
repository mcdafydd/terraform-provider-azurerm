package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var hdInsightStormClusterHeadNodeDefinition = azure.HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         4,
	MaxInstanceCount:         4,
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(2)),
	ValidVmSizes: []string{
		// TODO: are these the same for the other types?
		// only certain sizes are valid for certain machine types for certain kinds
		// otherwise you get an unhelpful error. this is an attempt to make a better UX
		"Standard_A3",
		"Standard_A4",
		"Standard_A4_v2",
		"Standard_A4m_v2",
		"Standard_A6",
		"Standard_A7",
		"Standard_A8_v2",
		"Standard_A8m_v2",
		"Standard_D12_v2",
		"Standard_D13_v2",
		"Standard_D14_v2",
		"Standard_D3_v2",
		"Standard_D4_v2",
		"Standard_D5_v2",
		"Standard_E16_v3",
		"Standard_E20_v3",
		"Standard_E2_v3",
		"Standard_E32_v3",
		"Standard_E4_v3",
		"Standard_E64_v3",
		"Standard_E64i_v3",
		"Standard_E8_v3",
		"Standard_G2",
		"Standard_G3",
		"Standard_G4",
		"Standard_G5",
	},
}

var hdInsightStormClusterWorkerNodeDefinition = azure.HDInsightNodeDefinition{
	CanSpecifyInstanceCount: true,
	MinInstanceCount:        1,
	// can't find a hard limit - appears to be limited by the subscription; setting something sensible for now
	MaxInstanceCount: 9999,
	CanSpecifyDisks:  false,
	ValidVmSizes: []string{
		// TODO: update this
		// TODO: are these the same for the other types?
		// only certain sizes are valid for certain machine types for certain kinds
		// otherwise you get an unhelpful error. this is an attempt to make a better UX
		"Standard_A3",
		"Standard_A4",
		"Standard_A4_v2",
		"Standard_A4m_v2",
		"A6",
		"A7",
		"Standard_A8_v2",
		"Standard_A8m_v2",
		"Standard_D3_v2",
		"Standard_D4_v2",
		"Standard_D5_v2",
		"Standard_D12_v2",
		"Standard_D13_v2",
		"Standard_D14_v2",
		"Standard_E2_v3",
		"Standard_E4_v3",
		"Standard_E8_v3",
		"Standard_E16_v3",
		"Standard_E20_v3",
		"Standard_E32_v3",
		"Standard_E64_v3",
		"Standard_E64i_v3",
		"Standard_G2",
		"Standard_G3",
		"Standard_G4",
		"Standard_G5",
	},
}

var hdInsightStormClusterZookeeperNodeDefinition = azure.HDInsightNodeDefinition{
	CanSpecifyInstanceCount:  false,
	MinInstanceCount:         3,
	MaxInstanceCount:         3,
	CanSpecifyDisks:          false,
	FixedTargetInstanceCount: utils.Int32(int32(3)),
	ValidVmSizes: []string{
		//// this is hard-coded at the API level
		//"Medium",
		"Standard_A4_V2",
	},
}

func resourceArmHDInsightStormCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHDInsightStormClusterCreate,
		Read:   resourceArmHDInsightStormClusterRead,
		Update: hdinsightClusterUpdate("Storm", resourceArmHDInsightStormClusterRead),
		Delete: hdinsightClusterDelete("Storm"),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": azure.SchemaHDInsightName(),

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"cluster_version": azure.SchemaHDInsightClusterVersion(),

			"tier": azure.SchemaHDInsightTier(),

			"component_version": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storm": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"gateway": azure.SchemaHDInsightsGateway(),

			"storage_account": azure.SchemaHDInsightsStorageAccounts(),

			"roles": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"head_node": azure.SchemaHDInsightNodeDefinition("roles.0.head_node", hdInsightStormClusterHeadNodeDefinition),

						"worker_node": azure.SchemaHDInsightNodeDefinition("roles.0.worker_node", hdInsightStormClusterWorkerNodeDefinition),

						"zookeeper_node": azure.SchemaHDInsightNodeDefinition("roles.0.zookeeper_node", hdInsightStormClusterZookeeperNodeDefinition),
					},
				},
			},

			"tags": tagsSchema(),

			"https_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ssh_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmHDInsightStormClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdinsightClustersClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	clusterVersion := d.Get("cluster_version").(string)
	tags := d.Get("tags").(map[string]interface{})
	tier := hdinsight.Tier(d.Get("tier").(string))

	componentVersionsRaw := d.Get("component_version").([]interface{})
	componentVersions := expandHDInsightStormComponentVersion(componentVersionsRaw)

	gatewayRaw := d.Get("gateway").([]interface{})
	gateway := azure.ExpandHDInsightsConfigurations(gatewayRaw)

	storageAccountsRaw := d.Get("storage_account").([]interface{})
	storageAccounts, err := azure.ExpandHDInsightsStorageAccounts(storageAccountsRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `storage_account`: %s", err)
	}

	stormRoles := hdInsightRoleDefinition{
		HeadNodeDef:      hdInsightStormClusterHeadNodeDefinition,
		WorkerNodeDef:    hdInsightStormClusterWorkerNodeDefinition,
		ZookeeperNodeDef: hdInsightStormClusterZookeeperNodeDefinition,
	}
	rolesRaw := d.Get("roles").([]interface{})
	roles, err := expandHDInsightRoles(rolesRaw, stormRoles)
	if err != nil {
		return fmt.Errorf("Error expanding `roles`: %+v", err)
	}

	if requireResourcesToBeImported {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_hdinsight_storm_cluster", *existing.ID)
		}
	}

	params := hdinsight.ClusterCreateParametersExtended{
		Location: utils.String(location),
		Properties: &hdinsight.ClusterCreateProperties{
			Tier:           tier,
			OsType:         hdinsight.Linux,
			ClusterVersion: utils.String(clusterVersion),
			ClusterDefinition: &hdinsight.ClusterDefinition{
				Kind:             utils.String("Storm"),
				ComponentVersion: componentVersions,
				Configurations:   gateway,
			},
			StorageProfile: &hdinsight.StorageProfile{
				Storageaccounts: storageAccounts,
			},
			ComputeProfile: &hdinsight.ComputeProfile{
				Roles: roles,
			},
		},
		Tags: expandTags(tags),
	}
	future, err := client.Create(ctx, resourceGroup, name, params)
	if err != nil {
		return fmt.Errorf("Error creating HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Error reading ID for HDInsight Storm Cluster %q (Resource Group %q)", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHDInsightStormClusterRead(d, meta)
}

func resourceArmHDInsightStormClusterRead(d *schema.ResourceData, meta interface{}) error {
	clustersClient := meta.(*ArmClient).hdinsightClustersClient
	configurationsClient := meta.(*ArmClient).hdinsightConfigurationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["clusters"]

	resp, err := clustersClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] HDInsight Storm Cluster %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	configuration, err := configurationsClient.Get(ctx, resourceGroup, name, "gateway")
	if err != nil {
		return fmt.Errorf("Error retrieving Configuration for HDInsight Storm Cluster %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	// storage_account isn't returned so I guess we just leave it ¯\_(ツ)_/¯
	if props := resp.Properties; props != nil {
		d.Set("cluster_version", props.ClusterVersion)
		d.Set("tier", string(props.Tier))

		if def := props.ClusterDefinition; def != nil {
			if err := d.Set("component_version", flattenHDInsightStormComponentVersion(def.ComponentVersion)); err != nil {
				return fmt.Errorf("Error flattening `component_version`: %+v", err)
			}

			if err := d.Set("gateway", azure.FlattenHDInsightsConfigurations(configuration.Value)); err != nil {
				return fmt.Errorf("Error flattening `gateway`: %+v", err)
			}
		}

		stormRoles := hdInsightRoleDefinition{
			HeadNodeDef:      hdInsightStormClusterHeadNodeDefinition,
			WorkerNodeDef:    hdInsightStormClusterWorkerNodeDefinition,
			ZookeeperNodeDef: hdInsightStormClusterZookeeperNodeDefinition,
		}
		flattenedRoles := flattenHDInsightRoles(d, props.ComputeProfile, stormRoles)
		if err := d.Set("roles", flattenedRoles); err != nil {
			return fmt.Errorf("Error flattening `roles`: %+v", err)
		}

		httpEndpoint := azure.FindHDInsightConnectivityEndpoint("HTTPS", props.ConnectivityEndpoints)
		d.Set("https_endpoint", httpEndpoint)
		sshEndpoint := azure.FindHDInsightConnectivityEndpoint("SSH", props.ConnectivityEndpoints)
		d.Set("ssh_endpoint", sshEndpoint)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func expandHDInsightStormComponentVersion(input []interface{}) map[string]*string {
	vs := input[0].(map[string]interface{})
	return map[string]*string{
		"Storm": utils.String(vs["storm"].(string)),
	}
}

func flattenHDInsightStormComponentVersion(input map[string]*string) []interface{} {
	stormVersion := ""
	if v, ok := input["Storm"]; ok {
		if v != nil {
			stormVersion = *v
		}
	}
	return []interface{}{
		map[string]interface{}{
			"storm": stormVersion,
		},
	}
}