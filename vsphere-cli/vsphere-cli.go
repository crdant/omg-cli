package vspherecli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-cli/deployments/bosh-init"
	"github.com/enaml-ops/omg-cli/utils"
)

func deployYaml(myYaml string, boshInitDeploy func(string)) {
	fmt.Println("deploying your bosh")
	content := []byte(myYaml)
	tmpfile, err := ioutil.TempFile("", "bosh-init-deployment")
	defer os.Remove(tmpfile.Name())

	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	boshInitDeploy(tmpfile.Name())
}

func checkRequired(name string, c *cli.Context) {
	if c.String(name) == "" {
		fmt.Println("Sorry you need to provide " + name)
		os.Exit(1)
	}
}

// GetFlags returns the available CLI flags
func GetFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "bosh-release-ver", Value: "256.2", Usage: "the version of the bosh release you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-cpi-release-ver", Value: "52", Usage: "the bosh cpi version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-ver", Value: "3012", Usage: "the go agent version you wish to use (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-release-sha", Value: "ff2f4e16e02f66b31c595196052a809100cfd5a8", Usage: "sha1 of the bosh release being used (found on bosh.io)"},
		cli.StringFlag{Name: "bosh-cpi-release-sha", Value: "dc4a0cca3b33dce291e4fbeb9e9948b6a7be3324", Usage: "sha1 of the cpi release being used (found on bosh.io)"},
		cli.StringFlag{Name: "go-agent-sha", Value: "3380b55948abe4c437dee97f67d2d8df4eec3fc1", Usage: "sha1 of the go agent being use (found on bosh.io)"},
		cli.StringFlag{Name: "director-name", Value: "my-bosh", Usage: "the name of your director"},
		cli.BoolFlag{Name: "print-manifest", Usage: "if you would simply like to output a manifest the set this flag as true."},
		// vsphere specific flags
		cli.StringFlag{Name: "vsphere-address", Value: "", Usage: "IP of the vCenter"},
		cli.StringFlag{Name: "vsphere-user", Value: "", Usage: "vSphere user"},
		cli.StringFlag{Name: "vsphere-password", Value: "", Usage: "vSphere user's password"},
		cli.StringFlag{Name: "vsphere-datacenter-name", Value: "", Usage: "name of the datacenter the Director will use for VM creation"},
		cli.StringFlag{Name: "vsphere-vm-folder", Value: "", Usage: "name of the folder created to hold VMs"},
		cli.StringFlag{Name: "vsphere-template-folder", Value: "", Usage: "the name of the folder created to hold stemcells"},
		cli.StringFlag{Name: "vsphere-datastore-pattern", Value: "", Usage: "name of the datastore the Director will use for storing VMs"},
		cli.StringFlag{Name: "vsphere-persistent-datastore-pattern", Value: "", Usage: "name of the datastore the Director will use for storing persistent disks. Defaults to vsphere-datastore-pattern"},
		cli.StringFlag{Name: "vsphere-disk-path", Value: "", Usage: "name of the VMs folder, disk folder will be automatically created in the chosen datastore."},
		cli.StringSliceFlag{Name: "vsphere-clusters", Value: &cli.StringSlice{""}, Usage: "one or more vSphere datacenter cluster names"},
		cli.StringFlag{Name: "vsphere-network-name", Value: "", Usage: "name of the vSphere network"},
	}
}

// GetAction returns a function action that can be registered with the CLI
func GetAction(boshInitDeploy func(string)) func(c *cli.Context) error {
	return func(c *cli.Context) (e error) {
		checkRequired("vsphere-address", c)
		checkRequired("vsphere-user", c)
		checkRequired("vsphere-password", c)
		checkRequired("vsphere-datacenter-name", c)
		checkRequired("vsphere-vm-folder", c)
		checkRequired("vsphere-template-folder", c)
		checkRequired("vsphere-datastore-pattern", c)
		checkRequired("vsphere-disk-path", c)
		checkRequired("vsphere-clusters", c)
		checkRequired("vsphere-network-name", c)

		manifest := boshinit.NewVSphereBosh(boshinit.BoshInitConfig{
			Name:                  c.String("name"),
			BoshReleaseVersion:    c.String("bosh-release-ver"),
			BoshPrivateIP:         c.String("bosh-private-ip"),
			BoshCPIReleaseVersion: c.String("bosh-cpi-release-ver"),
			GoAgentVersion:        c.String("go-agent-ver"),
			BoshReleaseSHA:        c.String("bosh-release-sha"),
			BoshCPIReleaseSHA:     c.String("bosh-cpi-release-sha"),
			GoAgentSHA:            c.String("go-agent-sha"),
			// vsphere specific
			VSphereAddress:                    c.String("vsphere-address"),
			VSphereUser:                       c.String("vsphere-user"),
			VSpherePassword:                   c.String("vsphere-password"),
			VSphereDatacenterName:             c.String("vsphere-datacenter-name"),
			VSphereVMFolder:                   c.String("vsphere-vm-folder"),
			VSphereTemplateFolder:             c.String("vsphere-template-folder"),
			VSphereDatastorePattern:           c.String("vsphere-datastore-pattern"),
			VSpherePersistentDatastorePattern: c.String("vsphere-persistent-datastore-pattern"),
			VSphereDiskPath:                   c.String("vsphere-disk-path"),
			VSphereClusters:                   utils.ClearDefaultStringSliceValue(c.StringSlice("vsphere-clusters")...),
			VSphereNetworkName:                c.String("vsphere-network-name"),
		})

		if yamlString, err := enaml.Paint(manifest); err == nil {

			if c.Bool("print-manifest") {
				fmt.Println(yamlString)

			} else {
				deployYaml(yamlString, boshInitDeploy)
			}
		} else {
			e = err
		}
		return
	}
}
