#!/usr/bin/bash

#Link in the form of https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/origin-ci-test/logs/release-openshift-ocp-installer-e2e-ovirt-upgrade-4.5-stable-to-4.6-ci/1306062535597756416/artifacts/e2e-ovirt/
link="$1"


proto="$(echo $link | grep :// | sed -e's,^\(.*://\).*,\1,g')"
# remove the protocol
url="$(echo ${1/$proto/})"

#echo $url
path="$(echo $url | grep / | cut -d/ -f3-)"



job_name="$(echo $path | grep / | cut -d/ -f8)"
build_id="$(echo $path | grep / | cut -d/ -f9)"

echo $path
echo $job_name
echo $build_id

cd ~/dev/tmp/aa
gsutil -m cp -R gs://${path} .