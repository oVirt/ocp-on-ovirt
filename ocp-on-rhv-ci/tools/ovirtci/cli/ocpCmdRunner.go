package main

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/oVirt/ocp-on-ovirt/ocp-on-rhv-ci/tools/ovirtci/ovirtci"

	log "github.com/sirupsen/logrus"
)

const proxyvm string = "ovirt-proxy-vm.rhv44.gcp.devcluster.openshift.com:22"

func main() {

	oengine := ovirtci.Engine{}

	//get current logged user
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	sshKeyFile := fmt.Sprintf("/home/%s/.ssh/id_rsa", user.Username)
	sshVMKeyFile := fmt.Sprintf("/home/%s/id_rsa", user.Username)

	fmt.Println(fmt.Sprintf("welcome %s , ocp-on-rhv CI infra cluster debugging tool , you key %s", user.Username, sshKeyFile))

	proxyVm := ovirtci.Proxyvm{
		Address:  proxyvm,
		SshUser:  user.Username,
		SshKey:   sshKeyFile,
		VmSshKey: sshVMKeyFile,
	}

	log.SetLevel(log.DebugLevel)

	oengine.ConnectToOvirt()

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panics occurs %v, try the non-Must methods to find the reason", err)
		}
	}()

	engineVms := oengine.ListVMs()
	oengine.SetOvirtInfo()

	res := ovirtci.RunProxyVM("find /var/lib/dnsmasq/net-1*.leases -not -empty", proxyVm)
	log.Debugln("getting VMs info")

	leaseFiles := ovirtci.LeaseFiles{}
	leaseFiles.ParseLeaseFiles(res)
	clusters := ovirtci.Clusters{}

	for _, val := range leaseFiles {
		res = ovirtci.RunProxyVM(fmt.Sprintf("cat %s", val), proxyVm)
		if strings.HasPrefix(res.FileName, "cat ") {
			res.FileName = filepath.Base(strings.Trim(res.FileName, "cat "))
		}
		clusters.SetCluster(res.FileName, res.Content)
	}
	clusters.PrintClusters()

	//connect to each VM and run command
	log.Debugln("connecting to OCP VMs")

	opts := map[string][]string{}

	for name, cluster := range clusters {
		for _, vm := range cluster.VmList {
			vm.Alive = false

			//only choose the VMs that reported by the engine
			if _, ok := engineVms[vm.Name]; ok {
				log.Debugf("name: %s , address: %s", vm.Name, vm.IpAddress)
				opts[name] = append(opts[name], vm.IpAddress)
			}

		}
	}

	for name := range clusters {
		log.Debugln("cluster", opts[name])
		res, err := ovirtci.RunVMMany(opts[name], proxyVm, "sudo uptime")
		if err != nil {
			log.Fatal()
		}

		for addr, resl := range res {
			log.Debugln(addr, resl)

		}

	}

	keys := make([]string, 0, len(opts))
	for k := range opts {
		keys = append(keys, k)
	}

}
