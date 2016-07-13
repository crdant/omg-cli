package boshinit

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/aws_cpi"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/director"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/postgres"
	"github.com/enaml-ops/omg-cli/plugins/products/bosh-init/enaml-gen/registry"
)

type BoshInitConfig struct {
	BoshAvailabilityZone              string
	BoshInstanceSize                  string
	AWSSubnet                         string
	AWSPEMFilePath                    string
	AWSAccessKeyID                    string
	AWSSecretKey                      string
	AWSRegion                         string
	AWSSecurityGroups                 []string
	AWSKeyName                        string
	AzureVnet                         string
	AzureSubnet                       string
	AzureSubscriptionID               string
	AzureTenantID                     string
	AzureClientID                     string
	AzureClientSecret                 string
	AzureResourceGroup                string
	AzureStorageAccount               string
	AzureDefaultSecurityGroup         string
	AzureSSHPubKey                    string
	AzureSSHUser                      string
	AzureEnvironment                  string
	AzurePrivateKeyPath               string
	VSphereAddress                    string
	VSphereUser                       string
	VSpherePassword                   string
	VSphereDatacenterName             string
	VSphereVMFolder                   string
	VSphereTemplateFolder             string
	VSphereDatastorePattern           string
	VSpherePersistentDatastorePattern string
	VSphereDiskPath                   string
	VSphereClusters                   []string
	VSphereNetworks                   []Network
}

type Network struct {
	Name    string
	Range   string
	Gateway string
	DNS     []string
}

type Rr registry.Registry
type Ar aws_cpi.Registry

type RegistryProperty struct {
	Rr      `yaml:",inline"`
	Ar      `yaml:",inline"`
	Address string `yaml:"address"`
}
type user struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type DirectorProperty struct {
	director.Director `yaml:",inline"`
	Address           string
}

type PgSql struct {
	User     string
	Host     string
	Password string
	Database string
	Adapter  string
}

type IAASManifestProvider interface {
	CreateCPIRelease() enaml.Release
	CreateCPITemplate() enaml.Template
	CreateDiskPool() enaml.DiskPool
	CreateResourcePool() enaml.ResourcePool
	CreateManualNetwork() enaml.ManualNetwork
	CreateVIPNetwork() enaml.VIPNetwork
	CreateJobNetwork() enaml.Network
	CreateCloudProvider() enaml.CloudProvider
	CreateCPIJobProperty() interface{}
	CreateDeploymentManifest() *enaml.DeploymentManifest
}

type Postgres interface {
	GetDirectorDB() *director.DirectorDb
	GetRegistryDB() *registry.Db
	GetPostgresDB() postgres.Postgres
}

type BoshBase struct {
	NetworkCIDR         string
	NetworkGateway      string
	NetworkDNS          []string
	DirectorName        string
	DirectorPassword    string
	AgentPassword       string
	DBPassword          string
	CPIName             string
	NtpServers          []string
	NatsPassword        string
	MBusPassword        string
	UAAPublicKey        string
	PrivateIP           string
	PublicIP            string
	SSLCert             string
	SSLKey              string
	SigningKey          string
	VerificationKey     string
	HealthMonitorSecret string
	LoginSecret         string
	RegistryPassword    string
	CACert              string
	BoshReleaseSHA      string
	BoshReleaseVersion  string
	CPIReleaseSHA       string
	CPIReleaseVersion   string
	GOAgentVersion      string
	GOAgentSHA          string
	UAAReleaseSHA       string
	UAAReleaseVersion   string
}
type BoshDefaults struct {
	CIDR               string
	Gateway            string
	DNS                *cli.StringSlice
	BoshReleaseVersion string
	BoshReleaseSHA     string
	PrivateIP          string
	CPIReleaseVersion  string
	CPIReleaseSHA      string
	CPIName            string
	GOAgentVersion     string
	GOAgentSHA         string
	NtpServers         *cli.StringSlice
}

//UAAClient - Structure to represent map of client priviledges
type UAAClient struct {
	ID                   string      `yaml:"id,omitempty"`
	Secret               string      `yaml:"secret,omitempty"`
	Scope                string      `yaml:"scope,omitempty"`
	AuthorizedGrantTypes string      `yaml:"authorized-grant-types,omitempty"`
	Authorities          string      `yaml:"authorities,omitempty"`
	AutoApprove          interface{} `yaml:"autoapprove,omitempty"`
	Override             bool        `yaml:"override,omitempty"`
	RedirectURI          string      `yaml:"redirect-uri,omitempty"`
	AccessTokenValidity  int         `yaml:"access-token-validity,omitempty"`
	RefreshTokenValidity int         `yaml:"refresh-token-validity,omitempty"`
	ResourceIDs          string      `yaml:"resource_ids,omitempty"`
	Name                 string      `yaml:"name,omitempty"`
	AppLaunchURL         string      `yaml:"app-launch-url,omitempty"`
	ShowOnHomepage       bool        `yaml:"show-on-homepage,omitempty"`
	AppIcon              string      `yaml:"app-icon,omitempty"`
}
