# ocp-on-rhv CI  env bringup

---

## bring up the Engine on GCP using terraform

- obtain service key from GCP place it in mygcloud/ocp-on-rhv-service.json.
  https://github.com/openshift/shared-secrets/blob/master/gce/aos-serviceaccount.json

- generate ssh key pair to access the engine
- apply terraform script

  this script will:
  - create network/subnetwork for the engine instance.
  - create the engine instance.
  - create the GCP FW rules to allow internal/outside communication.

  ```shell
    source scripts/env.sh
    terraform init
    terraform apply
  ```

## configure  external access to oVirt GCP engine.

- activate gcloud from service account
  ```shell
  gcloud auth activate-service-account --key-file=/config/mygcloud/ocp-on-rhv-service.json
  ```


- check if the engine public address is available and in status `RESERVED`:

  ```shell
  gcloud compute addresses list --filter="name~'ocp-on-rhv-engine'"
  NAME                  ADDRESS/RANGE  TYPE      PURPOSE  NETWORK  REGION       SUBNET  STATUS
  ocp-on-rhv-engine-vm  35.226.85.87   EXTERNAL                    us-central1          RESERVED
  ```

- if the address not available , new address can be created using command:

  ```
  gcloud compute addresses create ocp-rhv-static-ip-engine --region us-central1
  ```

- assign public IP address to the engine instance:

  ```bash
   gcloud compute instances add-access-config ocp-rhv44-vm-engine --address  35.226.85.87 --zone us-central1-c
  ```

- verify DNS record exist and point to the engine address:

  ```shell
  gcloud dns record-sets list --zone=devcluster | grep -i engine

  # engine.rhv.gcp.devcluster.openshift.com.    A     60     35.226.85.87
  ```

- update /etc/hosts with the engine fqdn ( point to the internal IP )

  ```shell
  [root@ocp-rhv-vm-engine ~]# cat /etc/hosts  | grep engine.rhv.gcp.devcluster.openshift.com
  10.0.0.10 ocp-rhv-vm-engine.c.openshift-gce-devel.internal ocp-rhv-vm-engine engine.rhv.gcp.devcluster.openshift.com # Added by Google
  ```

## oVirt installation

- ssh to the engine with user centos.

- installing ovirt-engine 4.3

  ```
    sudo su
    yum install http://resources.ovirt.org/pub/yum-repo/ovirt-release44.rpm -y
    yum update -y
    yum install ovirt-engine -y
  ```

- run `engine-setup` , accept all the defaults :
    engine-fqdn: engine.rhv44.gcp.devcluster.openshift.com

