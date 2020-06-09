# ocp-on-rhv CI deployment
- configure the required params
  ```bash
    export OVIRT_ENGINE_USERNAME=admin@internal
    export OVIRT_ENGINE_PASSWORD=xxxxyyyzz
    export OVIRT_ENGINE_URL=https://engine-fqdn/ovirt-engine/api
  ```

- running full  deployment
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml
  ```

- running proxy-vm provision only
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -t proxy-vm
  ```

- running environment sanity tests throught the proxy-vm
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -t ovirt-ocp-tests
  ```