#!/bin/bash


#vault configuration
export localPath=ci-secrets
mkdir -p $localPath

putVaultField(){
local vaultpath=$1
local localfile=$2

field=$(basename $localfile)
echo $localfile
eval "cat $localfile | vault kv patch $vaultpath $field=-"
}

for FILE in $localPath/*; do 

putVaultField "$VAULT_SECRETS_PATH" $FILE

done



