ARG BUILDER_TAG

FROM nexus/compiler-builder:$BUILDER_TAG
ADD compiler.tar /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git
ADD cmd/nexus-openapi-gen /go/bin/nexus-openapi-gen

WORKDIR /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git
RUN make init_submodules

WORKDIR /go/src/
RUN mkdir nexustmp && cd nexustmp && go mod init && \
    /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/scripts/pin_deps.sh && \
    go mod download

WORKDIR /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_deps/github.com/kubernetes/code-generator
RUN go mod download

WORKDIR /go/src/gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git

CMD make docker.gitlab_credentials && make generate_code
