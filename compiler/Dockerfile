ARG BUILDER_TAG

FROM nexus/compiler-builder:$BUILDER_TAG
ADD compiler.tar /go/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler
ADD cmd/nexus-openapi-gen /go/bin/nexus-openapi-gen
ADD cmd/gqlgen /go/bin/gqlgen

WORKDIR /go/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler
RUN mv .git .gitmodules kube-openapi gqlgen nexus ..
RUN git config --global --add safe.directory '*' && git config --global --add url."https://github.com/".insteadOf "git@github.com:"
RUN make init_submodules

WORKDIR /go/src/
RUN mkdir nexustmp && cd nexustmp && go mod init && \
    /go/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler/scripts/pin_deps.sh && \
    /go/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler/scripts/pin_graphql_build_version.sh && \
    go mod download

WORKDIR /go/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler/_deps/github.com/kubernetes/code-generator
RUN go mod download


WORKDIR /go/src/github.com/vmware-tanzu/graph-framework-for-microservices/compiler
RUN go mod download

CMD make docker.gitlab_credentials && make preparser
