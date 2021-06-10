# IPI tools

## Setup user and group on oVirt

The ansible playbook [ovirt_ocpadmin_user-setup.yml](ovirt_ocpadmin_user-setup.yml) helps to setup
a user and group with the minimum privileges to run the *openshift-install* on ovirt.

Before running the playbook the following variables should be exported

```bash
export OVIRT_ENGINE_FQDN=engine.example.com
export OVIRT_ENGINE_USERNAME=johndoe
export OVIRT_ENGINE_PASSWORD=xxxyyyzzz
export OCP_ADMIN_PASSWORD="rrrsssttt"
```

In the playbook you can customize the ocp username, groupname and the datacenter name

```YAML
    ocp_admin_username: ocpadmin
    ocp_admin_groupname: ocp-administrator

    data_center_name: Default
```


Run the playbook from the oVirt engine node or from a machine with the ovirt.infra roles
installed

```bash
ansible-playbook ovirt_ocpadmin_user-setup.yml
```

