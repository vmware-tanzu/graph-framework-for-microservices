/*
Copyright (c) 2021 VMware, Inc. All rights reserved.

Proprietary and confidential.

Unauthorized copying or use of this file, in any medium or form,
is strictly prohibited.
*/

package common

import (
	"embed"
	"fmt"
	"os"
	"time"
)

// VERSION ...Version set at compile time.
var VERSION string

// OS ...OS set at compile time.
var OS string

// BUILT ...Architecture set at compile time.
var BUILT string

// GIT_BRANCH ...Git branch set at compile time.
var GIT_BRANCH string

// GIT_COMMIT ...Git commit SHA set at compile time.
var GIT_COMMIT string

const NEXUS_DIR = "nexus"

var HarborRepo = "harbor-repo.vmware.com/nexus"

const (
	HELLOWORLD_URL                   = "https://storage.googleapis.com/nexus-template-downloads/%s/helloworld-example.tar"
	DATAMODEL_TEMPLATE_URL           = "https://storage.googleapis.com/nexus-template-downloads/%s/datamodel-templatedir.tar"
	NEXUS_TEMPLATE_URL               = "https://storage.googleapis.com/nexus-template-downloads/%s/nexus-template.tar"
	RUNTIME_MANIFESTS_URL            = "https://storage.googleapis.com/nexus-template-downloads/%s/runtime-manifests.tar"
	VALIDATION_MANIFESTS_URL         = "https://storage.googleapis.com/nexus-template-downloads/%s/validation-manifests.tar"
	API_GATEWAY_MANIFESTS_URL        = "https://storage.googleapis.com/nexus-template-downloads/%s/api-gw-manifests.tar"
	API_DATAMODEL_CRD_URL            = "https://storage.googleapis.com/nexus-template-downloads/%s/api-datamodel-crds.tar"
	CONNECT_CONTROLLER_MANIFESTS_URL = "https://storage.googleapis.com/nexus-template-downloads/%s/controller-manifests.tar"
)

const TEMPLATE_URL = "https://storage.googleapis.com/nexus-template-downloads/%s/app-template.tar"
const Filename = "app-template.tar"

type NexusConfig struct {
	Name string
}

var NexusConfFile = "NEXUSDATAMODEL"

//go:embed values.yaml
var TemplateFs embed.FS

var WaitTimeout = 2 * time.Minute

var RuntimePodLabels = []string{
	"-lapp=nexus-etcd",
	"-lapp=nexus-kube-apiserver",
	"-lname=nexus-kube-controllermanager",
	"-lapp=nexus-validation",
	"-lcontrol-plane=connect-controller",
	"-lapp.kubernetes.io/component=controller,app.kubernetes.io/name=ingress-nginx",
	"-lcontrol-plane=api-gw",
}

func GetEnvList() []string {
	return []string{
		fmt.Sprintf("GOPRIVATE=%s", getGoPrivate()),
	}
}

func getGoPrivate() string {
	return "*.eng.vmware.com," + os.Getenv("GOPRIVATE")
}

type ImageTemplate struct {
	Image                string
	Tag                  string
	IsImagePullSecret    bool
	ImagePullSecret      string
	Namespace            string
	NetworkingAPIVersion string
	GatewayMode          string
}

type Manifest struct {
	URL            string
	Directory      string
	VersionStrName string
	FileName       string
	Templatized    bool
	VersionEnv     string
	ImageName      string
	Image          ImageTemplate
}

var RuntimeManifests = map[string]Manifest{
	"controller": {
		URL:            CONNECT_CONTROLLER_MANIFESTS_URL,
		Directory:      "nexus-manifests-controller",
		VersionEnv:     "NEXUS_CONTROLLER_TEMPLATE_VERSION",
		VersionStrName: "NexusConnectController",
		FileName:       "controller-manifests.tar",
		Templatized:    true,
		ImageName:      "controller",
	},
	"runtime": {
		URL:            RUNTIME_MANIFESTS_URL,
		Directory:      "nexus-manifests-runtime",
		VersionEnv:     "NEXUS_RUNTIME_MANIFESTS_VERSION",
		VersionStrName: "NexusRuntime",
		FileName:       "runtime-manifests.tar",
		Templatized:    true,
		ImageName:      "",
	},
	"validation": {
		URL:            VALIDATION_MANIFESTS_URL,
		Directory:      "nexus-manifests-validation",
		VersionEnv:     "NEXUS_VALIDATION_TEMPLATE_VERSION",
		VersionStrName: "NexusValidationTemplates",
		FileName:       "validation-manifets.tar",
		ImageName:      "nexus-validation",
		Templatized:    true,
	},
	"api-gateway": {
		URL:            API_GATEWAY_MANIFESTS_URL,
		Directory:      "nexus-manifests-api-gateway",
		VersionEnv:     "NEXUS_API_GATEWAY_TEMPLATE_VERSION",
		VersionStrName: "NexusApiGatewayTemplates",
		FileName:       "api-gw-manifests.tar",
		ImageName:      "api-gateway",
		Templatized:    true,
	},
}

var NexusApiDatamodelManifest = Manifest{
	URL:            API_DATAMODEL_CRD_URL,
	Directory:      "nexus-manifests-api-datamodel-crds",
	VersionEnv:     "NEXUS_API_DATAMODEL_CRD",
	VersionStrName: "NexusApiDatamodelCrds",
	FileName:       "api-datamodel-crds.tar",
	ImageName:      "",
	Templatized:    false,
}
