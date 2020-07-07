#!/bin/bash

help="Tool for running ovirt e2e jobs from the command line,\n
you have to be connected to ovirt openshift ci namespace before running this tool\n

Usage:
e2e-runner OPTIONS
"

LEASE_TYPE="conformance"
OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE="registry.svc.ci.openshift.org/ocp/release:4.6"
USE_LEASE="true"
USE_OVIRT_TEMPLATE="true"
STOP_AFTER_SETUP="false"
TEST="conformance"
STOP_AFTER_TEST="false"
STOP_BEFORE_TEST="true"

echo $PWD

sed -e "s|PLACEHOLDER_LEASE_TYPE|${LEASE_TYPE}|" \
    -e "s|PLACEHOLDER_OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE|${OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE}|" \
    -e "s|PLACEHOLDER_USE_LEASE|${USE_LEASE}|" \
    -e "s|PLACEHOLDER_USE_OVIRT_TEMPLATE|${USE_OVIRT_TEMPLATE}|" \
    -e "s|PLACEHOLDER_STOP_AFTER_SETUP|${STOP_AFTER_SETUP}|" \
    -e "s|PLACEHOLDER_TEST|${TEST}|" \
    -e "s|PLACEHOLDER_STOP_AFTER_TEST|${STOP_AFTER_TEST}|" \
    -e "s|PLACEHOLDER_STOP_BEFORE_TEST|${STOP_BEFORE_TEST}|" ./e2e-suite.in.yaml > e2e-job.yaml