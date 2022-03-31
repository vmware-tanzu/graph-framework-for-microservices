ARG BUILDER_TAG

FROM nexus/compiler-builder:$BUILDER_TAG
ADD compiler.tar /go/src/gitlab.eng.vmware.com/nexus/compiler

WORKDIR /go/src/gitlab.eng.vmware.com/nexus/compiler
RUN make init_submodules
RUN go mod download

WORKDIR /go/src/gitlab.eng.vmware.com/nexus/compiler/_deps/github.com/kubernetes/code-generator
RUN go mod download

WORKDIR /go/src/gitlab.eng.vmware.com/nexus/compiler

CMD make generate_code
