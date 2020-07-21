package main

import (
	"strings"

	"github.com/oVirt/ocp-on-ovirt/ocp-on-rhv-ci/tools/ovirtci/ovirtci"
	"github.com/prometheus/common/log"
)

func main() {

	oengine := ovirtci.Engine{}
	oengine.ConnectToOvirt()

	//oengine.ListEvents()
	engineVms := oengine.ListVMs()

	// iterate over VM list
	for vmName, comment := range engineVms {
		fields := strings.Split(vmName, "-")
		log.Infof("Found vm %s with comment %s", vmName, comment)

		//if a VM already have a comment then skip it
		if comment != "" {
			continue
		}

		//skip VMs that dont comply with our naming convention
		if len(fields) > 2 {
			cluster_id := fields[0]

			//find the PR Link from the last received event for that cluster
			prLink := oengine.GetLastPRForCluster(cluster_id)

			if prLink != "" {
				log.Infof("Updating vm %s - cluster - %s with PR %s", vmName, cluster_id, prLink)
				oengine.AddComment(vmName, prLink)
			} else {
				log.Infof("Couldnt find PR for cluster %s", cluster_id)
			}

		}

	}

}
