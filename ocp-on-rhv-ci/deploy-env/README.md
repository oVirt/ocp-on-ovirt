# ocp-on-rhv CI deployment

the deployment Consists of 4 roles:
- creating OVN networks for internal OCP communication
- spawning proxy VM for outside communication
- uploading RCHOS image template (not must)
- sanity tests
  

## how to use it?

- apply the env for the oVirt engine communication:
  ```bash
    export OVIRT_ENGINE_USERNAME=admin@internal
    export OVIRT_ENGINE_PASSWORD=123456
    export OVIRT_ENGINE_URL=https://engine2.es.localvms.com/ovirt-engine/api
  ```

- prepare your env configuration, examples can be found [here](https://github.com/oVirt/ocp-on-ovirt/tree/master/ocp-on-rhv-ci/deploy-env/config)
  
- running full  deployment
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -e@config/ci-bm-params.yml
  ```

- running proxy-vm provision only
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -t proxy-vm -e@config/ci-bm-params.yml
  ```

- upload latest rchos template
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -t rhcos-template -e@config/ci-bm-params.yml

  ```

- running environment sanity tests throught the proxy-vm
  ```bash
  ansible-playbook ocp_on_rhv-deploy.yml -t ovirt-ocp-tests -e@config/ci-bm-params.yml
  ```
