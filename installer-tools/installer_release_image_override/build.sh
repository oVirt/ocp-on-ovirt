# OPENSHIFT_RELEASE_VERSION can be 4.4, 4.5, 4.6 etc
export OPENSHIFT_RELEASE_VERSION=4.6
export MIRROR="mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp-dev-preview/latest-${OPENSHIFT_RELEASE_VERSION}"
export OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE="$(curl -s -l "${MIRROR}/release.txt" | sed -n 's/^Pull From: //p')"
./bin/openshift-install create cluster --dir=ocp --log-level=debug
