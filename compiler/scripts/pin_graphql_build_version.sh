set -e

### Pinning this dependency version for graphql and compiler libraries using similar version to build plugin
COMPILER_SRC_DIRECTORY=$1
go mod edit -require sigs.k8s.io/yaml@v1.3.0
go mod edit -replace github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen=${COMPILER_SRC_DIRECTORY}/../gqlgen
go mod edit -replace github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi=${COMPILER_SRC_DIRECTORY}/../kube-openapi
go mod edit -replace github.com/vmware-tanzu/graph-framework-for-microservices/nexus=${COMPILER_SRC_DIRECTORY}/../nexus
go mod edit -require github.com/cespare/xxhash/v2@v2.1.2
go mod edit -require github.com/google/gofuzz@v1.1.0
go mod edit -require github.com/imdario/mergo@v0.3.12
go mod edit -require k8s.io/apimachinery@v0.26.0
go mod edit -require k8s.io/client-go@v0.26.0
go mod edit -require golang.org/x/sys@v0.2.0
go mod edit -require golang.org/x/time@v0.3.0
go mod edit -require golang.org/x/term@v0.2.0
go mod edit -require golang.org/x/text@v0.4.0
go mod edit -require golang.org/x/oauth2@v0.0.0-20220411215720-9780585627b5
go mod edit -require k8s.io/klog/v2@v2.70.1
go mod edit -require golang.org/x/net@v0.2.0
go mod edit -require google.golang.org/grpc@v1.51.0
go mod edit -require k8s.io/utils@v0.0.0-20221128185143-99ec85e7a448
go mod edit -require sigs.k8s.io/controller-runtime@v0.14.1
go mod edit -require k8s.io/api@v0.26.0
