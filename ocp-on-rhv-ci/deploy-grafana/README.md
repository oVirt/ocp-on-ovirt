## deployment steps

- ### deploy prometheus + grafana on CI ovirt Namespace

  - `oc apply -f ci-prometheus-deployment.yaml`
   
- ### deploy node-exporters on oVirt nodes.
   - install ansible pre-req role:

      `ansible-galaxy role install -r ansible-requirements.yaml`

   - run ansible deployment:
  
        `ansible-playbook deploy-node-exporters.yaml`