package cloudfoundry

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	mysqllib "github.com/enaml-ops/omg-cli/plugins/products/cf-mysql/enaml-gen/mysql"
)

//NewMySQLPartition -
func NewMySQLPartition(c *cli.Context) (igf InstanceGrouper) {
	igf = &MySQL{
		AZs:                    c.StringSlice("az"),
		StemcellName:           c.String("stemcell-name"),
		NetworkIPs:             c.StringSlice("mysql-ip"),
		NetworkName:            c.String("mysql-network"),
		VMTypeName:             c.String("mysql-vm-type"),
		PersistentDiskType:     c.String("mysql-disk-type"),
		AdminPassword:          c.String("mysql-admin-password"),
		BootstrapUsername:      c.String("mysql-bootstrap-username"),
		BootstrapPassword:      c.String("mysql-bootstrap-password"),
		DatabaseStartupTimeout: 1200,
		InnodbBufferPoolSize:   2147483648,
		MaxConnections:         1500,
		SyslogAddress:          c.String("syslog-address"),
		SyslogPort:             c.Int("syslog-port"),
		SyslogTransport:        c.String("syslog-transport"),
		MySQLSeededDatabases:   parseSeededDBs(c),
	}
	return
}

func parseSeededDBs(c *cli.Context) (dbs []MySQLSeededDatabase) {
	//TODO GOT TO BE A BETTER WAY
	var dbName string
	dbMap := make(map[string]MySQLSeededDatabase)
	for _, flag := range c.FlagNames() {
		if strings.HasPrefix(flag, "db-") {
			if c.IsSet(flag) {
				baseName := strings.Replace(flag, "db-", "", 1)
				if strings.HasSuffix(flag, "-password") {
					dbName = strings.Replace(baseName, "-password", "", 1)
					pwd := c.String(flag)
					if seededDatabase, ok := dbMap[dbName]; ok {
						seededDatabase.Password = pwd
						dbMap[dbName] = seededDatabase
					} else {
						seededDatabase = MySQLSeededDatabase{
							Name:     dbName,
							Password: pwd,
						}
						dbMap[dbName] = seededDatabase
					}
				} else if strings.HasSuffix(flag, "-username") {
					dbName = strings.Replace(baseName, "-username", "", 1)
					userName := c.String(flag)
					if seededDatabase, ok := dbMap[dbName]; ok {
						seededDatabase.Username = userName
						dbMap[dbName] = seededDatabase
					} else {
						seededDatabase = MySQLSeededDatabase{
							Name:     dbName,
							Username: userName,
						}
						dbMap[dbName] = seededDatabase
					}
				}
			}
		}
	}
	for _, value := range dbMap {
		dbs = append(dbs, value)
	}
	return
}

//ToInstanceGroup -
func (s *MySQL) ToInstanceGroup() (ig *enaml.InstanceGroup) {
	ig = &enaml.InstanceGroup{
		Name:      "mysql-partition",
		Instances: len(s.NetworkIPs),
		VMType:    s.VMTypeName,
		AZs:       s.AZs,
		Stemcell:  s.StemcellName,
		Jobs: []enaml.InstanceJob{
			s.newMySQLJob(),
		},
		Networks: []enaml.Network{
			enaml.Network{Name: s.NetworkName, StaticIPs: s.NetworkIPs},
		},
		Update: enaml.Update{
			MaxInFlight: 1,
		},
	}
	return
}

func (s *MySQL) newMySQLJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "mysql",
		Release: "cf-mysql",
		Properties: &mysqllib.Mysql{
			AdminPassword:          s.AdminPassword,
			ClusterIps:             s.NetworkIPs,
			DatabaseStartupTimeout: s.DatabaseStartupTimeout,
			InnodbBufferPoolSize:   s.InnodbBufferPoolSize,
			MaxConnections:         s.MaxConnections,
			BootstrapEndpoint: &mysqllib.BootstrapEndpoint{
				Username: s.BootstrapUsername,
				Password: s.BootstrapPassword,
			},
			SeededDatabases: s.MySQLSeededDatabases,
			SyslogAggregator: &mysqllib.SyslogAggregator{
				Address:   s.SyslogAddress,
				Port:      s.SyslogPort,
				Transport: s.SyslogTransport,
			},
		},
	}
}

//HasValidValues -
func (s *MySQL) HasValidValues() bool {
	return (len(s.AZs) > 0 &&
		s.StemcellName != "" &&
		s.VMTypeName != "" &&
		s.NetworkName != "" &&
		len(s.NetworkIPs) > 0 &&
		s.PersistentDiskType != "" &&
		s.AdminPassword != "" &&
		s.BootstrapUsername != "" &&
		s.BootstrapPassword != "")
}