# OCP on RHV infra deployment

## pre-req

- RHV-4.4 installation
- Cluster with at least one host with these specs:
  - minimum of 46Gb RAM:
    - 10Gb for each master x 3
    - 8Gb for each worker x 2
    - 20VCPUs
- DC with Attached Storage at least 100GB (not production)
- create ssh public key under user ~/.ssh/id_rsa.pub

## deployment summary

the deployment Consists of 4 stages:

- creating OVN networks for between VM communication
- creating proxy VM:
  - it will serve as Router,DHCP,NAT,DNS (dnsmasq) and HAProxy.
  - it will have one external nic (ovirtmgmt) and additional nic for each ovn network.
- sanity tests
  - spawn test VM that will try to communicate through the proxy-vm
  
Each stage can be run separately using ansible tag or everything can run in one flow.

## deployment steps

- apply the env for the oVirt engine communication:
  ```bash
    export OVIRT_ENGINE_USERNAME=admin@internal
    export OVIRT_ENGINE_PASSWORD=123456
    export OVIRT_ENGINE_URL=https://engine2.es.localvms.com/ovirt-engine/api
  ```

- prepare your env configuration, examples can be found [here](https://github.com/oVirt/ocp-on-ovirt/tree/master/ocp-on-rhv-ci/deploy-env/config)
  
- running full deployment
  
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -e@config/ci-bm-params.yml
  ```

- provision OVN networks only
  
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -t networks -e@config/ci-bm-params.yml
  ```

- running proxy-vm provision only
  
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -t proxy-vm -e@config/ci-bm-params.yml
  ```

- running environment sanity tests throught the proxy-vm
  
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -t ovirt-ocp-tests -e@config/ci-bm-params.yml
  ```

## post installation optional steps

- reapply OVN configuration for the rhv cluster (works on rhv 4.4)
  
  on the engine run the following:

  ```shell
  cd /usr/share/ovirt-engine/ansible-runner-service-project/project
  cp /usr/share/ovirt-engine/playbooks/ovirt-provider-ovn-driver.yml .
  ansible-playbook --key-file /etc/pki/ovirt-engine/keys/engine_id_rsa -i /usr/share/ovirt-engine-metrics/bin/ovirt-engine-hosts-ansible-inventory --extra-vars " cluster_name=<CLUSTER_NAME> ovn_central=<ENGINE-IP> ovn_tunneling_interface=ovirtmgmt" ovirt-provider-ovn-driver.yml
  ```
  note: this playbook needs to  run every time new hosts is added to the cluster.
  
  to verify that the rhv host is succefully connected to the centralized OVN switch examine the log: `/var/log/openvswitch/ovn-controller.log` on the rhv Host.
  ```
  2021-05-19T13:56:39.844Z|00007|reconnect|INFO|ssl:169.63.244.91:6642: connected
  2021-05-19T13:56:39.872Z|00008|ofctrl|INFO|unix:/var/run/openvswitch/br-int.mgmt: connecting to switch
  2021-05-19T13:56:39.872Z|00009|rconn|INFO|unix:/var/run/openvswitch/br-int.mgmt: connecting...
  2021-05-19T13:56:39.872Z|00010|rconn|WARN|unix:/var/run/openvswitch/br-int.mgmt: connection failed (No such file or directory)
  2021-05-19T13:56:39.872Z|00011|rconn|INFO|unix:/var/run/openvswitch/br-int.mgmt: waiting 1 seconds before reconnect
  2021-05-19T13:56:40.873Z|00012|rconn|INFO|unix:/var/run/openvswitch/br-int.mgmt: connecting...
  ```
