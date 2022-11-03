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
	HELLOWORLD_URL         = "https://storage.googleapis.com/nexus-template-downloads/%s/helloworld-example.tar"
	DATAMODEL_TEMPLATE_URL = "https://storage.googleapis.com/nexus-template-downloads/%s/datamodel-templatedir.tar"
	NEXUS_TEMPLATE_URL     = "https://storage.googleapis.com/nexus-template-downloads/%s/nexus-template.tar"
	RUNTIME_MANIFESTS_URL  = "https://storage.googleapis.com/nexus-template-downloads/%s/runtime-manifests.tar"

	//For helm only
	HELM_REPO_URL = "https://harbor-repo.vmware.com/chartrepo/nexus"
	HELM_REPO     = "harbor-vmware"
)

var HELM_CHART_NAME = fmt.Sprintf("%s/nexus-runtime", HELM_REPO)

const TEMPLATE_URL = "https://storage.googleapis.com/nexus-template-downloads/%s/app-template.tar"
const Filename = "app-template.tar"

type NexusConfig struct {
	Name string
}

var NexusDMPropertiesFile string = "nexus.yaml"
var NexusConfFile = "NEXUSDATAMODEL"

//go:embed values.yaml
var TemplateFs embed.FS

//go:embed runtime_installer.yaml
var RuntimeTemplate embed.FS

var WaitTimeout = 2 * time.Minute

var NexusGroupSuffix string = ".tsm.vmware.com"

func GetEnvList() []string {
	return []string{
		fmt.Sprintf("GOPRIVATE=%s", getGoPrivate()),
	}
}

func getGoPrivate() string {
	return "*.eng.vmware.com," + os.Getenv("GOPRIVATE")
}

type Tags struct {
	FieldName  string
	VersionEnv string
	ImageName  string
}

var TagsList = map[string]Tags{
	"controller": {
		VersionEnv: "NEXUS_CONTROLLER_TEMPLATE_VERSION",
		FieldName:  "controller",
		ImageName:  "controller",
	},
	"validation": {
		VersionEnv: "NEXUS_VALIDATION_TEMPLATE_VERSION",
		FieldName:  "validation",
		ImageName:  "nexus-validation",
	},
	"api-gateway": {
		VersionEnv: "NEXUS_API_GATEWAY_TEMPLATE_VERSION",
		FieldName:  "api_gateway",
		ImageName:  "api-gateway",
	},
	"connector": {
		VersionEnv: "NEXUS_CONNECTOR_TEMPLATE_VERSION",
		FieldName:  "connector",
		ImageName:  "connector",
	},
	"api": {
		VersionEnv: "NEXUS_API_DATAMODEL_CRD_VERSION",
		FieldName:  "api",
		ImageName:  "api",
	},
}

type Datamodel struct {
	DatamodelInstaller  DatamodelInstaller
	IsImagePullSecret   bool
	ImagePullSecret     string
	SkipCRDInstallation string
	DatamodelTitle      string
	GraphqlPath         string
}
type DatamodelInstaller struct {
	Image string
	Name  string
}

type RuntimeInstaller struct {
	Image   string
	Name    string
	Args    []string
	Command []string
}
