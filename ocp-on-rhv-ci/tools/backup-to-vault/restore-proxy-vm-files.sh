#!/bin/bash

filename="proxy-vm-config-files.tar.gz"

#vault configuration
export VAULT_ADDR=https://vault.ci.openshift.org
export VAULT_TOKEN=changeme

#download the recent backup from the vault
#vault kv put  kv/selfservice/ovirt/proxy-vm-backup @$filename.b64
vault kv get -field=proxy-vm-config kv/selfservice/ovirt/proxy-vm-backup >$filename.b64

#decode base64
base64 -d $filename.b64 >$filename

#copy to proxy-VM
scp $filename root@ovirt-proxy-vm-2.rhv44.gcp.devcluster.openshift.com:/tmp/$filename

#restore config files on proxy-vm
ssh root@ovirt-proxy-vm-2.rhv44.gcp.devcluster.openshift.com tar xvfz /tmp/$filename -C /

rm -rf $filename.b64
