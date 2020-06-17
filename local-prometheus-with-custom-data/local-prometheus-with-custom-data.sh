#!/bin/bash
set -euo pipefail

prometheus_tar_url=$1
port="${2:-9090}"

echo "Using port ${port}"

dir=$(mktemp -d -p /tmp/)
pushd $dir
wget $prometheus_tar_url
tar xvf prometheus.tar

sudo chown -R 65534:65534 ./*
sudo chmod 777 $dir
sudo chmod -R 777 ./*
rm -rf prometheus.tar
docker run -p ${port}:9090 -v `pwd`:/prometheus prom/prometheus
popd
