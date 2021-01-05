package main

import (
	"fmt"

	ovirtci "github.com/ovirt/ocp-on-ovirt/ocp-on-rhv-ci/tools/ovirtci/ovirtci"

	log "github.com/sirupsen/logrus"
)

func main() {

	oengine := ovirtci.Engine{}

	subnetCIDR := "192.168.0.0/16"
	oengine.ConnectToOvirt()

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panics occurs %v, try the non-Must methods to find the reason", err)
		}
	}()

	engineVms := oengine.ListVMs()
	oengine.SetOvirtInfo()

	if len(engineVms)>0{
		log.Printf("#{engineVms}")
	}

	for vmName, _ := range engineVms {
		vmip, err  := oengine.GetVmIp(vmName,subnetCIDR)
		if err == nil {
			fmt.Printf("vm: %s , got address: %s \n", vmName, vmip)
		}
	}


}
