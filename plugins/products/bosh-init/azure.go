package boshinit

import (
	"fmt"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/cloudproperties/azure"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/aws_cpi"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/cpi"
)

func NewAzureBosh(cfg BoshInitConfig) *enaml.DeploymentManifest {
	var ntpProperty = NewNTP("0.pool.ntp.org", "1.pool.ntp.org")
	var cpiTemplate = enaml.Template{Name: "cpi", Release: "bosh-azure-cpi"}
	var manifest = NewBoshDeploymentBase(cfg, "cpi", ntpProperty)

	manifest.AddRelease(enaml.Release{
		Name: "bosh-azure-cpi",
		URL:  "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-azure-cpi-release?v=" + cfg.BoshCPIReleaseVersion,
		SHA1: cfg.BoshCPIReleaseSHA,
	})

	resourcePool := enaml.ResourcePool{
		Name:    "vms",
		Network: "private",
	}
	resourcePool.Stemcell = enaml.Stemcell{
		URL:  "https://bosh.io/d/stemcells/bosh-azure-hyperv-ubuntu-trusty-go_agent?v=" + cfg.GoAgentVersion,
		SHA1: cfg.GoAgentSHA,
	}
	resourcePool.CloudProperties = azurecloudproperties.ResourcePool{
		InstanceType: cfg.BoshInstanceSize,
	}
	manifest.AddResourcePool(resourcePool)
	manifest.AddDiskPool(enaml.DiskPool{
		Name:     "disks",
		DiskSize: 20000,
	})
	net := enaml.NewManualNetwork("private")
	net.AddSubnet(enaml.Subnet{
		Range:   cfg.BoshCIDR,
		Gateway: cfg.BoshGateway,
		DNS:     cfg.BoshDNS,
		CloudProperties: azurecloudproperties.Network{
			VnetName:   cfg.AzureVnet,
			SubnetName: cfg.AzureSubnet,
		},
	})
	manifest.AddNetwork(net)
	manifest.AddNetwork(enaml.NewVIPNetwork("public"))
	boshJob := manifest.Jobs[0]
	boshJob.AddTemplate(cpiTemplate)
	boshJob.AddNetwork(enaml.Network{
		Name:      "public",
		StaticIPs: []string{cfg.AzurePublicIP},
	})
	var agentProperty = aws_cpi.Agent{
		Mbus: "nats://nats:nats-password@" + cfg.BoshPrivateIP + ":4222",
	}
	boshJob.AddProperty(agentProperty)
	azureProperty := NewAzureProperty(
		cfg.AzureEnvironment,
		cfg.AzureSubscriptionID,
		cfg.AzureTenantID,
		cfg.AzureClientID,
		cfg.AzureClientSecret,
		cfg.AzureResourceGroup,
		cfg.AzureStorageAccount,
		cfg.AzureDefaultSecurityGroup,
		cfg.AzureSSHUser,
		cfg.AzureSSHPubKey,
	)
	boshJob.AddProperty(azureProperty)
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(NewAzureCloudProvider(azureProperty, cpiTemplate, cfg.AzurePublicIP, cfg.AzurePrivateKeyPath, ntpProperty))
	return manifest
}

func NewAzureCloudProvider(myazure cpi.Azure, cpiTemplate enaml.Template, pubip, keypath string, ntpProperty []string) enaml.CloudProvider {
	return enaml.CloudProvider{
		Template: cpiTemplate,
		MBus:     fmt.Sprintf("https://mbus:mbus-password@%s:6868", pubip),
		SSHTunnel: enaml.SSHTunnel{
			Host:           pubip,
			Port:           22,
			User:           "vcap",
			PrivateKeyPath: keypath,
		},
		Properties: map[string]interface{}{
			"azure": myazure,
			"agent": map[string]string{
				"mbus": "https://mbus:mbus-password@0.0.0.0:6868",
			},
			"blobstore": map[string]string{
				"provider": "local",
				"path":     "/var/vcap/micro_bosh/data/cache",
			},
			"ntp": ntpProperty,
		},
	}
}

func NewAzureProperty(azureenv, subid, tenantid, clientid, clientsecret, resourcegroup, storageaccount, securitygroup, sshuser, sshkey string) cpi.Azure {
	return cpi.Azure{
		Environment:          azureenv,
		SubscriptionId:       subid,
		TenantId:             tenantid,
		ClientId:             clientid,
		ClientSecret:         clientsecret,
		ResourceGroupName:    resourcegroup,
		StorageAccountName:   storageaccount,
		DefaultSecurityGroup: securitygroup,
		SshUser:              sshuser,
		SshPublicKey:         sshkey,
	}
}
