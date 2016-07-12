package boshinit

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/cloudproperties/aws"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/aws_cpi"
)

func NewAWSBosh(cfg BoshInitConfig) *enaml.DeploymentManifest {
	var ntpProperty = NewNTP("0.pool.ntp.org", "1.pool.ntp.org")
	var manifest = NewBoshDeploymentBase(cfg, "aws_cpi", ntpProperty)
	var awsProperty = aws_cpi.Aws{
		AccessKeyId:           cfg.AWSAccessKeyID,
		SecretAccessKey:       cfg.AWSSecretKey,
		DefaultKeyName:        cfg.AWSKeyName,
		DefaultSecurityGroups: cfg.AWSSecurityGroups,
		Region:                cfg.AWSRegion,
	}

	var agentProperty = aws_cpi.Agent{
		Mbus: "nats://nats:nats-password@" + cfg.BoshPrivateIP + ":4222",
	}

	manifest.AddRelease(enaml.Release{
		Name: "bosh-aws-cpi",
		URL:  "https://bosh.io/d/github.com/cloudfoundry-incubator/bosh-aws-cpi-release?v=" + cfg.BoshCPIReleaseVersion,
		SHA1: cfg.BoshCPIReleaseSHA,
	})

	resourcePool := enaml.ResourcePool{
		Name:    "vms",
		Network: "private",
	}
	resourcePool.Stemcell = enaml.Stemcell{
		URL:  "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent?v=" + cfg.GoAgentVersion,
		SHA1: cfg.GoAgentSHA,
	}
	resourcePool.CloudProperties = awscloudproperties.ResourcePool{
		InstanceType: cfg.BoshInstanceSize,
		EphemeralDisk: awscloudproperties.EphemeralDisk{
			Size:     25000,
			DiskType: "gp2",
		},
		AvailabilityZone: cfg.BoshAvailabilityZone,
	}
	manifest.AddResourcePool(resourcePool)
	manifest.AddDiskPool(enaml.DiskPool{
		Name:     "disks",
		DiskSize: 20000,
		CloudProperties: awscloudproperties.RootDisk{
			DiskType: "gp2",
		},
	})
	net := enaml.NewManualNetwork("private")
	net.AddSubnet(enaml.Subnet{
		Range:   cfg.BoshCIDR,
		Gateway: cfg.BoshGateway,
		DNS:     cfg.BoshDNS,
		CloudProperties: awscloudproperties.Network{
			Subnet: cfg.AWSSubnet,
		},
	})
	manifest.AddNetwork(net)
	manifest.AddNetwork(enaml.NewVIPNetwork("public"))
	boshJob := manifest.Jobs[0]
	boshJob.AddTemplate(enaml.Template{Name: "aws_cpi", Release: "bosh-aws-cpi"})
	boshJob.AddNetwork(enaml.Network{
		Name:      "public",
		StaticIPs: []string{cfg.AWSElasticIP},
	})
	boshJob.AddProperty(agentProperty)
	boshJob.AddProperty(awsProperty)
	manifest.Jobs[0] = boshJob
	manifest.SetCloudProvider(NewAWSCloudProvider(cfg.AWSElasticIP, cfg.AWSPEMFilePath, awsProperty, ntpProperty))
	return manifest
}
