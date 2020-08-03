#!/bin/bash

help="
Tool for generting ovirt e2e job template from the command line.


Usage:
e2e-runner OPTIONS

Flags:
-m --minimal            Create a minimal cluster and run a minimal test suite, this is used to run a regular PR job
-c --csi                Create a minimal cluster and run the csi test suite, this is used to test the csi driver
--no-lease              Don't acquire a boskos lease, this is very danguros since it can lead to overloading CI resources, use this flag only when testing new leases configuration.
--no-template           CI job will create a tempalte, currently we are not providing all the configuration needed on install config and don't clean old template VM.
--pause-after-setup     Pause for 3h after setup is successful, after you are done and want to continue log into the setup container and kill the sleep process
--pause-before-setup    Pause for 3h before test is successful, after you are done and want to continue log into the test container and kill the sleep process
--pause-after-test      Pause for 3h after test is successful, after you are done and want to continue log into the test container and kill the sleep process
--tag
--release
--job-name
"

#4.4 || 4.5 || 4.6 || latest
RELEASE_IMAGE_TAG="4.6"
RELEASE_IMAGE_REPO="registry.svc.ci.openshift.org/ocp/release"
OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE="${RELEASE_IMAGE_REPO}:${RELEASE_IMAGE_TAG}"
#conformance || minimal
LEASE_TYPE="minimal"
#conformance || minimal || csi
TEST="minimal"
JOB_NAME="e2e-job-${RELEASE_IMAGE_TAG}-${LEASE_TYPE}"
#true || false
USE_LEASE="true"
USE_OVIRT_TEMPLATE="false"
STOP_AFTER_SETUP="false"
STOP_AFTER_TEST="false"
STOP_BEFORE_TEST="true"

# while [[ "$1" =~ ^- && ! "$1" == "--" ]]; do case $1 in
#   -m | --minimal )
#     LEASE_TYPE="minimal"
#     TEST="minimal"
#     ;;
#   -c | --csi )
#     LEASE_TYPE="minimal"
#     TEST="csi"
#     ;;
#   --no-lease )
#     USE_LEASE="false"
#     ;;
#   --no-template )
#     USE_OVIRT_TEMPLATE="false"
#     ;;
#   --pause-after-setup )
#     STOP_AFTER_SETUP="true"
#     ;;
#   --pause-before-test )
#     STOP_BEFORE_TEST="true"
#     ;;
#   --pause-after-test )
#     STOP_AFTER_TEST="true"
#     ;;
#   --tag )
#     shift
#     RELEASE_IMAGE_TAG="${1}"
#   --release )
#     shift
#     RELEASE_IMAGE_REPO="${1}"
#     ;;
#   --job-name )
#     shift
#     JOB_NAME="${1}"
#     ;;
# esac; shift; done

sed -e "s|PLACEHOLDER_LEASE_TYPE|\"${LEASE_TYPE}\"|" \
    -e "s|PLACEHOLDER_OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE|\"${OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE}\"|" \
    -e "s|PLACEHOLDER_USE_LEASE|\"${USE_LEASE}\"|" \
    -e "s|PLACEHOLDER_USE_OVIRT_TEMPLATE|\"${USE_OVIRT_TEMPLATE}\"|" \
    -e "s|PLACEHOLDER_STOP_AFTER_SETUP|\"${STOP_AFTER_SETUP}\"|" \
    -e "s|PLACEHOLDER_TEST|\"${TEST}\"|" \
    -e "s|PLACEHOLDER_STOP_AFTER_TEST|\"${STOP_AFTER_TEST}\"|" \
    -e "s|PLACEHOLDER_JOB_NAME|\"${JOB_NAME}\"|" \
    -e "s|PLACEHOLDER_STOP_BEFORE_TEST|\"${STOP_BEFORE_TEST}\"|" /home/gzaidman/workspace/upstream/ocp/ocp-ci-tools/ci-jobs/e2e-runner/e2e-suite.in.yaml > ${JOB_NAME}.yaml
