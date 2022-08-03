#! /bin/bash

set -e

go mod edit -require github.com/elliotchance/orderedmap@v1.4.0
go mod edit -require k8s.io/apimachinery@v0.22.2
go mod edit -require k8s.io/client-go@v0.22.2

go mod edit -replace k8s.io/api=k8s.io/api@v0.22.2
go mod edit -replace k8s.io/apiextensions-apiserver=k8s.io/apiextensions-apiserver@v0.22.2
go mod edit -replace k8s.io/apimachinery=k8s.io/apimachinery@v0.22.2
go mod edit -replace k8s.io/client-go=k8s.io/client-go@v0.22.2
