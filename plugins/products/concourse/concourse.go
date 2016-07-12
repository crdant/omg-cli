package concourse

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/plugins/products/concourse/enaml-gen/atc"
	"github.com/enaml-ops/omg-cli/plugins/products/concourse/enaml-gen/baggageclaim"
	"github.com/enaml-ops/omg-cli/plugins/products/concourse/enaml-gen/garden"
	"github.com/enaml-ops/omg-cli/plugins/products/concourse/enaml-gen/groundcrew"
	"github.com/enaml-ops/omg-cli/plugins/products/concourse/enaml-gen/postgresql"
	"github.com/enaml-ops/omg-cli/plugins/products/concourse/enaml-gen/tsa"
	"github.com/xchapter7x/lo"
)

const (
	concourseReleaseName string = "concourse"
	concourseReleaseVer  string = "1.0.1"
	concourseReleaseURL  string = "https://bosh.io/d/github.com/concourse/concourse"
	concourseReleaseSHA  string = "ef60fe5182a7b09df324d167f36d26481e7f5e01"
	gardenReleaseName    string = "garden-linux"
	gardenReleaseVer     string = "0.337.0"
	gardenReleaseURL     string = "https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release"
	gardenReleaseSHA     string = "d1d81d56c3c07f6f9f04ebddc68e51b8a3cf541d"
	stemcellName         string = "ubuntu-trusty"
	stemcellVer          string = "3232.4"
	stemcellURL          string = "https://bosh.io/d/stemcells/bosh-aws-xen-hvm-ubuntu-trusty-go_agent"
	stemcellSHA          string = "ac920cae17c7159dee3bf1ebac727ce2d01564e9"
)

//Deployment -
type Deployment struct {
	enaml.Deployment
	manifest            *enaml.DeploymentManifest
	ConcourseURL        string
	ConcourseUserName   string
	ConcoursePassword   string
	DirectorUUID        string
	NetworkName         string
	WebIPs              []string
	WebInstances        int
	WebAZs              []string
	DatabaseAZs         []string
	WorkerAZs           []string
	DeploymentName      string
	NetworkRange        string
	NetworkGateway      string
	StemcellAlias       string
	PostgresPassword    string
	ResourcePoolName    string
	WebVMType           string
	WorkerVMType        string
	DatabaseVMType      string
	DatabaseStorageType string
	CloudConfigYml      string
	StemcellURL         string
	StemcellSHA         string
}

//NewDeployment -
func NewDeployment() (d Deployment) {
	d = Deployment{}
	d.manifest = new(enaml.DeploymentManifest)
	return
}

func (d *Deployment) doCloudConfigValidation(data []byte) (err error) {
	lo.G.Debug("Cloud Config:", string(data))
	c := &enaml.CloudConfigManifest{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}

	for _, azName := range d.WebAZs {
		if !c.ContainsAZName(azName) {
			err = fmt.Errorf("WebAZ [%s] is not defined as a AZ in cloud config", azName)
			return
		}
	}
	for _, azName := range d.WorkerAZs {
		if !c.ContainsAZName(azName) {
			err = fmt.Errorf("WorkerAZ[%s] is not defined as a AZ in cloud config", azName)
			return
		}
	}
	for _, azName := range d.DatabaseAZs {
		if !c.ContainsAZName(azName) {
			err = fmt.Errorf("DatabaseAZ[%s] is not defined as a AZ in cloud config", azName)
			return
		}
	}

	if !c.ContainsVMType(d.WebVMType) {
		err = fmt.Errorf("WebVMType[%s] is not defined as a VMType in cloud config", d.WebVMType)
		return
	}
	if !c.ContainsVMType(d.WorkerVMType) {
		err = fmt.Errorf("WorkerVMType[%s] is not defined as a VMType in cloud config", d.WorkerVMType)
		return
	}
	if !c.ContainsVMType(d.DatabaseVMType) {
		err = fmt.Errorf("DatabaseVMType[%s] is not defined as a VMType in cloud config", d.DatabaseVMType)
		return
	}
	if !c.ContainsDiskType(d.DatabaseStorageType) {
		err = fmt.Errorf("DatabaseStorageType[%s] is not defined as a DiskType in cloud config", d.DatabaseStorageType)
		return
	}
	return
}

//Initialize -
func (d *Deployment) Initialize(cloudConfig []byte) (err error) {

	if err = d.doCloudConfigValidation(cloudConfig); err != nil {
		return
	}
	var web *enaml.InstanceGroup
	var db *enaml.InstanceGroup
	var worker *enaml.InstanceGroup
	d.manifest.SetName(d.DeploymentName)
	d.manifest.SetDirectorUUID(d.DirectorUUID)
	d.manifest.AddRemoteRelease(concourseReleaseName, concourseReleaseVer, concourseReleaseURL, concourseReleaseSHA)
	d.manifest.AddRemoteRelease(gardenReleaseName, gardenReleaseVer, gardenReleaseURL, gardenReleaseSHA)
	d.manifest.AddRemoteStemcell(stemcellName, d.StemcellAlias, stemcellVer, d.StemcellURL, d.StemcellSHA)

	update := d.CreateUpdate()
	d.manifest.SetUpdate(update)

	if web, err = d.CreateWebInstanceGroup(); err != nil {
		return
	}
	d.manifest.AddInstanceGroup(web)

	if db, err = d.CreateDatabaseInstanceGroup(); err != nil {
		return
	}
	d.manifest.AddInstanceGroup(db)

	if worker, err = d.CreateWorkerInstanceGroup(); err != nil {
		return
	}
	d.manifest.AddInstanceGroup(worker)

	return
}

//CreateWebInstanceGroup -
func (d *Deployment) CreateWebInstanceGroup() (web *enaml.InstanceGroup, err error) {
	if err = validateInstanceGroup(d.ResourcePoolName, d.StemcellAlias, "WebAZs", d.WebAZs); err == nil {
		web = &enaml.InstanceGroup{
			Name:         "web",
			Instances:    d.WebInstances,
			ResourcePool: d.ResourcePoolName,
			VMType:       d.WebVMType,
			AZs:          d.WebAZs,
			Stemcell:     d.StemcellAlias,
		}
		web.AddNetwork(enaml.Network{
			Name:      d.NetworkName,
			StaticIPs: d.WebIPs,
		})
		web.AddJob(d.CreateAtcJob())
		web.AddJob(d.CreateTsaJob())
	}
	return
}

//CreateAtcJob -
func (d *Deployment) CreateAtcJob() (job *enaml.InstanceJob) {
	job = enaml.NewInstanceJob("atc", concourseReleaseName, atc.Atc{
		ExternalUrl:        d.ConcourseURL,
		BasicAuthUsername:  d.ConcourseUserName,
		BasicAuthPassword:  d.ConcoursePassword,
		PostgresqlDatabase: "atc",
	})
	return
}

//CreateTsaJob -
func (d *Deployment) CreateTsaJob() (job *enaml.InstanceJob) {
	job = enaml.NewInstanceJob("tsa", concourseReleaseName, tsa.Tsa{})
	return
}

//CreateDatabaseInstanceGroup -
func (d *Deployment) CreateDatabaseInstanceGroup() (db *enaml.InstanceGroup, err error) {
	persistenceDisk := 10240
	if d.DatabaseStorageType != "" {
		persistenceDisk = 0
	}
	if err = validateInstanceGroup(d.ResourcePoolName, d.StemcellAlias, "DatabaseAzs", d.DatabaseAZs); err == nil {
		db = &enaml.InstanceGroup{
			Name:               "db",
			Instances:          1,
			ResourcePool:       d.ResourcePoolName,
			PersistentDisk:     persistenceDisk,
			PersistentDiskType: d.DatabaseStorageType,
			VMType:             d.DatabaseVMType,
			AZs:                d.DatabaseAZs,
			Stemcell:           d.StemcellAlias,
		}
		db.AddNetwork(d.CreateNetwork())
		db.AddJob(d.CreatePostgresqlJob())
	}

	return
}

func validateInstanceGroup(resourcePoolName, stemcellAlias, propertyName string, azs []string) (err error) {
	if resourcePoolName == "" {
		if (len(azs) == 0) || (stemcellAlias == "") {
			err = fmt.Errorf("No resource pool name so must provide %s and StemcellAlias property", propertyName)
		}
	} else if (len(azs) > 0) || (stemcellAlias != "") {
		err = fmt.Errorf("ResourcePoolName defined so cannot also define %s (%s) and StemcellAlias (%s) properties", propertyName, azs, stemcellAlias)
	}
	return
}

//CreatePostgresqlJob -
func (d *Deployment) CreatePostgresqlJob() (job *enaml.InstanceJob) {
	dbs := make([]DBName, 1)
	dbs[0] = DBName{
		Name:     "atc",
		Role:     "atc",
		Password: d.PostgresPassword,
	}
	job = enaml.NewInstanceJob("postgresql", concourseReleaseName, postgresql.Postgresql{
		Databases: dbs,
	})
	return
}

//CreateWorkerInstanceGroup -
func (d *Deployment) CreateWorkerInstanceGroup() (worker *enaml.InstanceGroup, err error) {
	if err = validateInstanceGroup(d.ResourcePoolName, d.StemcellAlias, "WorkerAZs", d.WorkerAZs); err == nil {
		worker = &enaml.InstanceGroup{
			Name:         "worker",
			Instances:    1,
			ResourcePool: d.ResourcePoolName,
			VMType:       d.WorkerVMType,
			AZs:          d.WorkerAZs,
			Stemcell:     d.StemcellAlias,
		}

		worker.AddNetwork(d.CreateNetwork())
		worker.AddJob(d.CreateGroundCrewJob())
		worker.AddJob(d.CreateBaggageClaimJob())
		worker.AddJob(d.CreateGardenJob())
	}
	return
}

//CreateGardenJob -
func (d *Deployment) CreateGardenJob() (job *enaml.InstanceJob) {
	job = enaml.NewInstanceJob("garden", gardenReleaseName, Garden{
		garden.Garden{
			ListenAddress:   "0.0.0.0:7777",
			ListenNetwork:   "tcp",
			AllowHostAccess: true,
		},
	})
	return
}

//CreateBaggageClaimJob -
func (d *Deployment) CreateBaggageClaimJob() (job *enaml.InstanceJob) {
	job = enaml.NewInstanceJob("baggageclaim", concourseReleaseName, baggageclaim.Baggageclaim{})
	return
}

//CreateGroundCrewJob -
func (d *Deployment) CreateGroundCrewJob() (job *enaml.InstanceJob) {
	job = enaml.NewInstanceJob("groundcrew", concourseReleaseName, groundcrew.Groundcrew{})
	return
}

//CreateNetwork -
func (d *Deployment) CreateNetwork() (network enaml.Network) {
	network = enaml.Network{
		Name: d.NetworkName,
	}
	return
}

//CreateManualDeploymentNetwork -
func (d *Deployment) CreateManualDeploymentNetwork(networkName, networkRange, networkGateway string, webIPs []string) (network *enaml.ManualNetwork) {
	network = &enaml.ManualNetwork{
		Name: networkName,
		Type: "manual",
	}
	subnets := make([]enaml.Subnet, 1)
	subnet := enaml.Subnet{
		Range:   networkRange,
		Gateway: networkGateway,
		Static:  webIPs,
	}
	subnets[0] = subnet
	network.Subnets = subnets

	return
}

//CreateUpdate -
func (d *Deployment) CreateUpdate() (update enaml.Update) {
	update = enaml.Update{
		Canaries:        1,
		MaxInFlight:     3,
		Serial:          false,
		CanaryWatchTime: "1000-60000",
		UpdateWatchTime: "1000-60000",
	}

	return
}

//CreateResourcePool -
func (d *Deployment) CreateResourcePool(networkName string) (resourcePool enaml.ResourcePool) {
	const resourcePoolName = "concourse"
	resourcePool = enaml.ResourcePool{
		Name:    resourcePoolName,
		Network: networkName,
		Stemcell: enaml.Stemcell{
			Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
			Version: "latest",
		},
	}

	return
}

//CreateCompilation -
func (d *Deployment) CreateCompilation(networkName string) (compilation enaml.Compilation) {
	compilation = enaml.Compilation{
		Network: networkName,
		Workers: 3,
	}

	return
}

func (d Deployment) isStrongPass(pass string) (ok bool) {
	ok = false
	if len(pass) > 8 {
		ok = true
	}
	return
}

func insureHAInstanceCount(instances int) int {
	if instances < 2 {
		instances = 2
	}
	return instances
}

//GetDeployment -
func (d Deployment) GetDeployment() enaml.DeploymentManifest {
	return *d.manifest
}
