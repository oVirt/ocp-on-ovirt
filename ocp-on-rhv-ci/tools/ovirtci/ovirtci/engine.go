package ovirtci

import (
	"fmt"
	"os"
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
				vmSlice[vmName] = "Found"
			}
			if vmID, ok := vm.Id(); ok {
				fmt.Printf(" id: %v", vmID)
			}
			fmt.Println(")")
		}
	}
	return vmSlice
}
