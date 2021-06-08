#!/bin/sh -x

base_url="https://fqdn/ovirt-engine/api"
user="admin@internal"
password="pass"


listStorageConnectionExtensions() {
    host_id=$1
    curl \
    --insecure \
    --request GET \
    --header "Accept: application/xml" \
    --header "Content-Type: application/xml" \
    --user "${user}:${password}" \
    $base_url/hosts/${host_id}/storageconnectionextensions
}

addStorageConnectionExtensions(){
    hostId=$1
    target=$2
    iscsiUser=$3
    iscsiPass=$4
    

    curl \
    --insecure \
    --request POST \
    --header "Accept: application/xml" \
    --header "Content-Type: application/xml" \
    --user "${user}:${password}" \
    --data "<storage_connection_extension>
  <target>${target}</target>
  <username>${iscsiUser}</username>
  <password>${iscsiPass}</password>
</storage_connection_extension>" \
    $base_url/hosts/${hostId}/storageconnectionextensions
}

updateStorageConnectionExtensions(){
    hostId=$1
    target=$2
    iscsiUser=$3
    iscsiPass=$4
    extId=$5

    curl \
    --insecure \
    --request PUT \
    --header "Accept: application/xml" \
    --header "Content-Type: application/xml" \
    --user "${user}:${password}" \
    --data "<storage_connection_extension>
  <target>${target}</target>
  <username>${iscsiUser}</username>
  <password>${iscsiPass}</password>
</storage_connection_extension>" \
    $base_url/hosts/${hostId}/storageconnectionextensions/${extId}
}

deleteStorageConnectionExtensions(){
    hostId=$1
    storageConnectionExtensionId=$2

    curl \
    --insecure \
    --request DELETE \
    --header "Accept: application/xml" \
    --header "Content-Type: application/xml" \
    --user "${user}:${password}" \
    $base_url/hosts/${hostId}/storageconnectionextensions/${storageConnectionExtensionId}
}

#example
addStorageConnectionExtensions "host-uuid" \
 "iqn.xxxxxxxxxxxxxxx1:test" \
 "iscsiuser" \
 "secretpass"
