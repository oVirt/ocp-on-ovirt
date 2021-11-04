

# set vault address & token from https://vault.ci.openshift.org/
```bash
export VAULT_ADDR=https://vault.ci.openshift.org
export VAULT_TOKEN=change_me
export VAULT_SECRETS_PATH="kv/selfservice/ovirt/cluster-secrets-ovirt"
```

# backup CI secrets from vault to a local folder
`./backup_ci_secrets.sh`

# restore CI secrects from local folder

`./restore_ci_secrets.sh`