#  Tool for extracting OCP must-gather log information behind proxy installation

##  example configuration  files
- `conf/env.sh`
```bash
# engine user name - keycloak 
OVIRT_ENGINE_USERNAME=admin@ovirt@internalsso

# engine password
OVIRT_ENGINE_PASSWORD=<secretpass>

# engine USR
OVIRT_ENGINE_URL=https://ovirt-engine-fqdn/ovirt-engine/api

# Proxy VM address
BASTION_PROXY_ADDRESS=Proxy_Address

# path to the boostrap terraform file.
BOOTSTRAP_TFVARS=/home/runner/bootstrap.tfvars.json

# path with the must-gather bundle file
MUST_GATHER_PATH=/home/runner/output/log-bundle-bootstrap.tar.gz
```

- `conf/bootstrap.tfvars.json` - terraform var file copied from the OCP installation directory
```
{"bootstrap_vm_id":"105950dd-ccfd-4e6c-b414-3315d212a66a"}
```

- `conf/ssh-privatekey` -  this private key which should be able to access bastion VM and the internal ocp VMs.

## running playbook from priviliged container
```bash
podman run -ti --privileged  \
--env-file conf/env.sh \
-u root:root \
-v $(pwd)/src/:/home/runner/must-gather \
-v $(pwd)/conf/ssh-privatekey:/home/runner/.ssh/id_rsa \
-v $(pwd)/conf/bootstrap.tfvars.json:/home/runner/bootstrap.tfvars.json \
-v $(pwd)/output/:/home/runner/output/ \
-w /home/runner/must-gather/ \
quay.io/ovirt/ansible-runner:ovirt-45  ansible-playbook must-gather.yaml
```
