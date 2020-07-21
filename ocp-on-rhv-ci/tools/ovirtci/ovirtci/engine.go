package ovirtci

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	ovirtsdk4 "github.com/oVirt/go-ovirt"
	log "github.com/sirupsen/logrus"
)

//Engine - engine SDK
type Engine struct {
	conn     *ovirtsdk4.Connection
	version  string
	connURL  string
	username string
	password string
}

func (e *Engine) setFromEnvVars() {
	if username, exists := os.LookupEnv("OVIRT_ENGINE_USERNAME"); exists {
		log.Debugf("env exists %s - %s", "OVIRT_ENGINE_USERNAME", username)
		e.username = username
	}

	if password, exists := os.LookupEnv("OVIRT_ENGINE_PASSWORD"); exists {
		log.Debugf("env exists %s - %s", "OVIRT_ENGINE_PASSWORD", password)
		e.password = password
	}

	if connURL, exists := os.LookupEnv("OVIRT_ENGINE_URL"); exists {
		log.Debugf("env exists %s - %s", "OVIRT_ENGINE_URL", connURL)
		e.connURL = connURL
	}
}
func (e *Engine) ConnectToOvirt() {
	e.setFromEnvVars()
	inputRawURL := e.connURL
	log.Printf("Connecting to ovirt on %s", inputRawURL)
	conn, err := ovirtsdk4.NewConnectionBuilder().
		URL(inputRawURL).
		Username(e.username).
		Password(e.password).
		Insecure(true).
		Compress(true).
		Timeout(time.Second * 10).
		Build()
	if err != nil {
		log.Fatalf("Make connection failed, reason: %s", err.Error())
	}

	e.conn = conn

	// To use `Must` methods, you should recover it if panics
}
func (e *Engine) SetOvirtInfo() {
	// Get API information from the root service:
	api := e.conn.SystemService().Get().MustSend().MustApi()
	log.Infof("oVift Cluster Information:")
	log.Infof("%10s %v", "Version:", api.MustProductInfo().MustVersion().MustFullVersion())
	log.Infof("%10s %v", "Hosts:", api.MustSummary().MustHosts().MustTotal())
	log.Infof("%10s %v", "SDs:", api.MustSummary().MustStorageDomains().MustTotal())
	log.Infof("%10s %v", "Users:", api.MustSummary().MustUsers().MustTotal())
	log.Infof("%10s %v", "vms:", api.MustSummary().MustVms().MustTotal())
	e.version = api.MustProductInfo().MustVersion().MustFullVersion()
}

//GetLastPRForCluster - returns the last PR URL for the given cluster
func (e *Engine) GetLastPRForCluster(cluster string) string {
	eventsService := e.conn.SystemService().EventsService()
	var data map[string]interface{}

	events := map[int64]string{}
	keys := []int64{}

	allevents := eventsService.List().Search("origin=openshift-ci").MustSend()
	for _, event := range allevents.MustEvents().Slice() {
		fdata := strings.Split(event.MustDescription(), ";")
		if len(fdata) > 2 {
			json.Unmarshal([]byte(fdata[3]), &data)
			//clusterStatus := strings.TrimSpace(fdata[2])
			clusterid := strings.TrimSpace(fdata[1])
			jobtype := data["type"]

			//log.Printf("AAA %s", data["refs"])
			//if data["refs"] == nil {
			//	continue
			//}

			var prLink interface{}

			if jobtype != "presubmit" {
				prLink = "Release"
			} else {
				prLink = data["refs"].(map[string]interface{})["pulls"].([]interface{})[0].(map[string]interface{})["link"]
			}

			buildId, _ := strconv.ParseInt(data["buildid"].(string), 10, 64)
			//data["buildid"].(int64)

			//log.Printf("AAA %s", buildId)

			if clusterid == cluster {
				events[buildId] = fmt.Sprintf("%v", prLink)
				//log.Printf("AAA %s - %s", clusterid, prLink)
				keys = append(keys, buildId)
				//return fmt.Sprintf("%v", prLink)
			}

		}
	}
	//sort the keys in descend order (we need the recent build)
	sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })
	return fmt.Sprintf("%v", events[keys[0]])

}

func (e *Engine) ListEvents() {
	eventsService := e.conn.SystemService().EventsService()

	var data map[string]interface{}

	events := eventsService.List().Search("origin=openshift-ci").MustSend()
	for _, event := range events.MustEvents().Slice() {
		fdata := strings.Split(event.MustDescription(), ";")
		if len(fdata) > 2 {
			json.Unmarshal([]byte(fdata[3]), &data)
			fmt.Printf("%s - cluster:%s \n", data["refs"].(map[string]interface{})["pulls"].([]interface{})[0].(map[string]interface{})["link"], fdata[1])
		}
	}
}

func (e *Engine) AddComment(vmname string, comment string) {

	// Get the reference to the "vms" service:
	vmsService := e.conn.SystemService().VmsService()
	ovirtsdk4.NewVmBuilder().Comment("test").MustBuild()
	// Retrieve the description of the virtual machine:
	vmsResp, err := vmsService.List().Search(fmt.Sprintf("name=%s", vmname)).Send()
	if err != nil {
		fmt.Printf("Failed to get vm list, reason: %v\n", err)
		return
	}
	vm := vmsResp.MustVms().Slice()[0]

	//In order to update the virtual machine we need a reference to the service
	// the manages it:
	vmService := vmsService.VmService(vm.MustId())
	vmService.Update().
		Vm(
			ovirtsdk4.NewVmBuilder().
				Comment(comment).
				MustBuild()).
		Send()

}
func (e *Engine) ListVMs() map[string]string {
	// Get the reference to the "vms" service:
	vmsService := e.conn.SystemService().VmsService()

	// Use the "list" method of the "vms" service to list all the virtual machines
	vmsResponse, err := vmsService.List().Send()

	if err != nil {
		fmt.Printf("Failed to get vm list, reason: %v\n", err)
		return nil
	}
	vmSlice := map[string]string{}

	if vms, ok := vmsResponse.Vms(); ok {
		// Print the virtual machine names and identifiers:
		for _, vm := range vms.Slice() {
			fmt.Print("VM: (")
			if vmName, ok := vm.Name(); ok {
				fmt.Printf(" name: %v", vmName)
				vmSlice[vmName], _ = vm.Comment()
			}
			if vmID, ok := vm.Id(); ok {
				fmt.Printf(" id: %v", vmID)
			}
			fmt.Println(")")
		}
	}
	return vmSlice
}
