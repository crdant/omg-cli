package director 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Openstack struct {

	/*StateTimeout - Descr: Timeout (in seconds) for OpenStack resources desired state (optional, by default 300) Default: 300
*/
	StateTimeout interface{} `yaml:"state_timeout,omitempty"`

	/*ApiKey - Descr: OpenStack API key Default: <nil>
*/
	ApiKey interface{} `yaml:"api_key,omitempty"`

	/*EndpointType - Descr: OpenStack endpoint type (optional, by default publicURL) Default: publicURL
*/
	EndpointType interface{} `yaml:"endpoint_type,omitempty"`

	/*WaitResourcePollInterval - Descr: Changes the delay (in seconds) between each status check to OpenStack when creating a resource (optional, by default 5) Default: 5
*/
	WaitResourcePollInterval interface{} `yaml:"wait_resource_poll_interval,omitempty"`

	/*Project - Descr: OpenStack project name (required for Keystone API version 3) Default: <nil>
*/
	Project interface{} `yaml:"project,omitempty"`

	/*IgnoreServerAvailabilityZone - Descr: When creating disks do not use the servers AZ, default to openstack default Default: false
*/
	IgnoreServerAvailabilityZone interface{} `yaml:"ignore_server_availability_zone,omitempty"`

	/*Tenant - Descr: OpenStack tenant name (required for Keystone API version 2) Default: <nil>
*/
	Tenant interface{} `yaml:"tenant,omitempty"`

	/*BootVolumeCloudProperties - Descr: Volume type for the boot volume (optional) Default: <nil>
*/
	BootVolumeCloudProperties *BootVolumeCloudProperties `yaml:"boot_volume_cloud_properties,omitempty"`

	/*UseDhcp - Descr: Whether to use DHCP when configuring networking on VM (for both manual and dynamic) Default: true
*/
	UseDhcp interface{} `yaml:"use_dhcp,omitempty"`

	/*DefaultKeyName - Descr: Default OpenStack keypair to use when spinning up new vms Default: <nil>
*/
	DefaultKeyName interface{} `yaml:"default_key_name,omitempty"`

	/*ConfigDrive - Descr: Config drive device (cdrom or disk) to use as metadata service on OpenStack (optional, nil by default) Default: <nil>
*/
	ConfigDrive interface{} `yaml:"config_drive,omitempty"`

	/*AuthUrl - Descr: URL of the OpenStack Identity endpoint to connect to Default: <nil>
*/
	AuthUrl interface{} `yaml:"auth_url,omitempty"`

	/*ConnectionOptions - Descr: Hash containing optional connection parameters to the OpenStack API Default: <nil>
*/
	ConnectionOptions interface{} `yaml:"connection_options,omitempty"`

	/*BootFromVolume - Descr: Boot from volume (optional, false by default) Default: false
*/
	BootFromVolume interface{} `yaml:"boot_from_volume,omitempty"`

	/*Region - Descr: OpenStack region (optional) Default: <nil>
*/
	Region interface{} `yaml:"region,omitempty"`

	/*DefaultSecurityGroups - Descr: Default OpenStack security groups to use when spinning up new vms Default: <nil>
*/
	DefaultSecurityGroups interface{} `yaml:"default_security_groups,omitempty"`

	/*Username - Descr: OpenStack user name Default: <nil>
*/
	Username interface{} `yaml:"username,omitempty"`

	/*Domain - Descr: OpenStack domain (required for Keystone API version 3) Default: <nil>
*/
	Domain interface{} `yaml:"domain,omitempty"`

	/*StemcellPublicVisibility - Descr: Set public visibility for stemcells (optional, false by default) Default: false
*/
	StemcellPublicVisibility interface{} `yaml:"stemcell_public_visibility,omitempty"`

}