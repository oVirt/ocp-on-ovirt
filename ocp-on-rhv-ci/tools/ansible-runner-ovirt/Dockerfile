FROM quay.io/centos/centos:stream8
LABEL maintainer="Evgeny Slutsky"

ENV container docker

# Update packages to the latest version
RUN dnf -y update 

# Configure systemd.
# See https://hub.docker.com/_/centos/ for details.
RUN (cd /lib/systemd/system/sysinit.target.wants/; for i in *; do [ $i == systemd-tmpfiles-setup.service ] || rm -f $i; done) ;\
rm -f /lib/systemd/system/multi-user.target.wants/* ;\
rm -f /etc/systemd/system/*.wants/* ;\
rm -f /lib/systemd/system/local-fs.target.wants/* ;\
rm -f /lib/systemd/system/sockets.target.wants/*udev* ;\
rm -f /lib/systemd/system/sockets.target.wants/*initctl* ;\
rm -f /lib/systemd/system/basic.target.wants/* ;\
rm -f /lib/systemd/system/anaconda.target.wants/*

# Install required packages.
# Remove packages that are nolonger requried.
# Clean the dnf cache.
RUN dnf -y install \
epel-release \
initscripts \
&& dnf -y update \
&& dnf module -y install python39 \
&& dnf -y install sudo 

# install ovirt stuff
RUN dnf -y install centos-release-ovirt45
RUN dnf -y install python3-ovirt-engine-sdk4
RUN dnf -y install ansible-core ansible-collection-ansible-netcommon ovirt-ansible-collection ansible-collection-ansible-utils

RUN dnf -y autoremove \
&& dnf clean all \
&& rm -rf /var/cache/dnf/*

VOLUME [ "/sys/fs/cgroup" ]
CMD ["/usr/lib/systemd/systemd"]