#!/bin/bash

#vault configuration
export localPath=ci-secrets
mkdir -p $localPath

getVaultField(){
local path=$1
local field=$2
echo "saving secret ${localPath}/${field}"
vault kv get -field=$field $path >${localPath}/${field}
}

export IFS=''
files=$(vault kv get -format=json  kv/selfservice/ovirt/cluster-secrets-ovirt  | jq '.data.data | keys[]' |  tr -d '"')
echo $files | while read -r line ; do echo "aa $line" ; getVaultField "$VAULT_SECRETS_PATH" "$line" ;done 


