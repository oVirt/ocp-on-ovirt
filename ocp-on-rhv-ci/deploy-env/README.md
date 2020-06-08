# ocp-on-rhv CI deployment
- configure the required params
  ```bash
    export OVIRT_ENGINE_USERNAME=admin@internal
    export OVIRT_ENGINE_PASSWORD=xxxxyyyzz
    export OVIRT_ENGINE_URL=https://engine-fqdn/ovirt-engine/api
  ```

- running  deployment
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml
  ```