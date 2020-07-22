# This Dockerfile is used by CI to publish the installer image.
# It builds an image containing only the openshift-install.

FROM registry.svc.ci.openshift.org/openshift/release:golang-1.14 AS builder
#ADD . /go/src/github.com/ovirt/ocp-on-ovirt/ocp-on-rhv-ci/tools/ovirtci
WORKDIR /go/src/github.com/oVirt/ocp-on-ovirt/ocp-on-rhv-ci/tools/ovirtci
COPY . .
#RUN go build -mod vendor cli/prFinder.go
RUN hack/build.sh


FROM registry.svc.ci.openshift.org/origin/4.1:base
COPY --from=builder /go/src/github.com/oVirt/ocp-on-ovirt/ocp-on-rhv-ci/tools/ovirtci/prFinder /bin/prFinder
RUN mkdir /output && chown 1000:1000 /output
USER 1000:1000
ENV PATH /bin
ENV HOME /output
WORKDIR /output
ENTRYPOINT ["/bin/prFinder"]
