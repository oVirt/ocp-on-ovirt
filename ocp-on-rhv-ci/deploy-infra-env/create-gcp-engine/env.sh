#!/bin/bash


  function terraform() {
      [[ -z "$rhv_host_count" ]] && rhv_host_count=1

      docker run --rm -it --env TF_VAR_rhv_host_count=$rhv_host_count --env GOOGLE_CLOUD_KEYFILE_JSON=mygcloud/ocp-on-rhv-service.json -v $(pwd):/opt/app -v ~/.ssh:/home/terraform/.ssh contino/terraform "$@";
    }


    function ansible-playbook() {
    docker run --rm \
    -e USER=ansible \
    -e MY_UID=$(id -u) \
    -e MY_GID=$(id -g) \
    -v ${HOME}/.ssh/:/home/ansible/.ssh/:ro  -v $(pwd):/data quay.io/eslutsky/ansible:latest-tools ansible-playbook "$@";
    }

    function ansible() {
    docker run --rm \
    -e USER=ansible \
    -e MY_UID=$(id -u) \
    -e MY_GID=$(id -g) \
    -v ${HOME}/.ssh/:/home/ansible/.ssh/:ro  -v $(pwd):/data quay.io/eslutsky/ansible:latest-tools ansible "$@";
    }


    function ansible-inventory() {
    docker run --rm \
    -e USER=ansible \
    -e MY_UID=$(id -u) \
    -e MY_GID=$(id -g) \
    -v ${HOME}/.ssh/:/home/ansible/.ssh/:ro  -v $(pwd):/data quay.io/eslutsky/ansible:latest-tools ansible-inventory "$@";
    }


    function ansible-doc() {
    docker run --rm \
    -e USER=ansible \
    -e MY_UID=$(id -u) \
    -e MY_GID=$(id -g) \
    -v ${HOME}/.ssh/:/home/ansible/.ssh/:ro  -v $(pwd):/data quay.io/eslutsky/ansible:latest-tools ansible-doc "$@";
    }

    function gcloud(){
      docker run --rm \
      -e CLOUDSDK_CORE_PROJECT=openshift-gce-devel \
      -e CLOUDSDK_CONFIG=/config/mygcloud \
      -v $(pwd)/mygcloud:/config/mygcloud \
      -v $(pwd):/certs \
      gcr.io/google.com/cloudsdktool/cloud-sdk:alpine gcloud "$@";

    }

    function get_free_public_ip()
    {
      local regex_filter=$1
      free_ip=""
      free_ip=$(gcloud compute addresses list --filter="name~'${regex_filter}' \
      AND status:RESERVED" \
      --format='value(ADDRESS)' | head -1) >/dev/null
      [[ -z "$free_ip" ]] && return 1
      return 0
    }

    function get_vms_without_public_ips()
    {
      #[ "$DEBUG" == 'true' ] && set -x
      local regex_filter=$1
      vms=()
      vms=(`gcloud compute instances list --filter="name~'${regex_filter}' \
      AND -EXTERNAL_IP:*" --format='value(NAME)'`)
      [[ -z "$vms" ]] && return 1

      return 0

    }
