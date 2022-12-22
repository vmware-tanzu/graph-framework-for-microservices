module github.com/vmware-tanzu/graph-framework-for-microservices/compiler

go 1.17

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/elliotchance/orderedmap v1.4.0
	github.com/fatih/structtag v1.2.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.2
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.23.0
	github.com/sirupsen/logrus v1.8.1
	github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi v0.0.0-20220603123335-7416bd4754d3
	github.com/vmware-tanzu/graph-framework-for-microservices/nexus v0.0.0-20221114103953-d31f6492d7fc
	golang.org/x/mod v0.6.0
	golang.org/x/text v0.5.0
	golang.org/x/tools v0.2.0
	k8s.io/api v0.26.0
	k8s.io/apiextensions-apiserver v0.26.0
	k8s.io/apimachinery v0.26.0
	k8s.io/gengo v0.0.0-20220902162205-c0856e24416d
	k8s.io/utils v0.0.0-20221107191617-1a15be271d1d
	sigs.k8s.io/gateway-api v0.6.0
)

require (
	github.com/emicklei/go-restful v2.15.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	golang.org/x/net v0.3.1-0.20221206200815-1e63c2f08a10 // indirect
	golang.org/x/sys v0.3.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

replace github.com/vmware-tanzu/graph-framework-for-microservices/kube-openapi => ../kube-openapi

replace github.com/vmware-tanzu/graph-framework-for-microservices/gqlgen => ../gqlgen

replace github.com/vmware-tanzu/graph-framework-for-microservices/nexus => ../nexus
