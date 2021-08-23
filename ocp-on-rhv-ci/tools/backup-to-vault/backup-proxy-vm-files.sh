#!/bin/bash


filename=proxyvm-backup-$(date +%s).tar.gz
remote_path=/root/backups
remote_ssh_host=root@ovirt-proxy-vm-2.rhv44.gcp.devcluster.openshift.com 

#vault configuration
export VAULT_ADDR=https://vault.ci.openshift.org
export VAULT_TOKEN=changeme


ssh $remote_ssh_host "tar cvfz $remote_path/$filename \
/etc/sysconfig/network-scripts/ifcfg-net* \
/etc/haproxy/haproxy.cfg \
/var/lib/dnsmasq/*.conf \
/usr/lib/udev/rules.d/60-net.rules \
/etc/systemd/system/multi-user.target.wants/dnsmasq@net-*.service \
/usr/lib/systemd/system/dnsmasq@.service \
/etc/firewalld/zones/*"

scp $remote_ssh_host:$remote_path/$filename .
cat $filename | base64 >$filename.b64

cat $filename.b64 | \
vault kv put  kv/selfservice/ovirt/proxy-vm-backup proxy-vm-config=-

rm -rf $filename
