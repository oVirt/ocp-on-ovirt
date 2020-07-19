# ocp-on-rhv CI deployment
- configure the required params
  ```bash
    export OVIRT_ENGINE_USERNAME=admin@internal
    export OVIRT_ENGINE_PASSWORD=123456
    export OVIRT_ENGINE_URL=https://engine2.es.localvms.com/ovirt-engine/api
  ```

- running full  deployment
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -e@ci-bm-params.yml
  ```

- running proxy-vm provision only
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -t proxy-vm -e@ci-bm-params.yml
  ```

- running environment sanity tests throught the proxy-vm
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -t ovirt-ocp-tests -e@ci-bm-params.yml
  ```
