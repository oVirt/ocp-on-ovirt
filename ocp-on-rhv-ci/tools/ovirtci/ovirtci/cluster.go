package ovirtci

import (
	"fmt"
	"strings"
)

type Vm struct {
	macAddress string
	IpAddress  string
	Name       string
	Alive      bool
}

type Cluster struct {
	VmList    []Vm
	Name      string
	Raw       string
	LeaseFile string
}

func (c *Cluster) SetRaw(s string) {
	c.Raw = s
}

func (c *Cluster) parseVmList() {
	Vms := strings.Split(c.Raw, "\n")
	for _, val := range Vms {
		fields := strings.Split(val, " ")
		if len(fields) > 3 {
			vm := Vm{
				macAddress: fields[1],
				IpAddress:  fields[2],
				Name:       fields[3],
			}
			//fmt.Print(vm.name)
			c.VmList = append(c.VmList, vm)
		}
		//c.vmList = append(c.vmList, vm)
	}
}

type Clusters map[string]Cluster
type LeaseFiles []string

//ParseLeaseFiles - parse file list and populate the slice
func (c *LeaseFiles) ParseLeaseFiles(fileLists ChannelLease) {

	str := strings.Split(fileLists.Content, "\n")
	for _, val := range str[:len(str)-1] {
		*c = append(*c, val)
	}
}

//SetCluster - add new cluster to the map
func (c *Clusters) SetCluster(key string, value string) {
	cluster := Cluster{Raw: value}
	cluster.parseVmList()
	(*c)[key] = cluster
}

func (c *Clusters) PrintClusters() {
	for key, val := range *c {
		fmt.Println(key, "=", val.VmList[0].IpAddress)
	}

}
