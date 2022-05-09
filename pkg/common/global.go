/*
Copyright (c) 2021 VMware, Inc. All rights reserved.

Proprietary and confidential.

Unauthorized copying or use of this file, in any medium or form,
is strictly prohibited.
*/

package common

import (
	"embed"
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

const (
	HELLOWORLD_URL           = "https://storage.googleapis.com/nexus-template-downloads/%s/helloworld-example.tar"
	DATAMODEL_TEMPLATE_URL   = "https://storage.googleapis.com/nexus-template-downloads/%s/datamodel-templatedir.tar"
	NEXUS_TEMPLATE_URL       = "https://storage.googleapis.com/nexus-template-downloads/%s/nexus-template.tar"
	RUNTIME_MANIFESTS_URL    = "https://storage.googleapis.com/nexus-template-downloads/%s/runtime-manifests.tar"
	VALIDATION_MANIFESTS_URL = "https://storage.googleapis.com/nexus-template-downloads/%s/validation-manifests.tar"
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

var PodLabels [5]string = [5]string{
	"-lapp=nexus-etcd",
	"-lapp=nexus-kube-apiserver",
	"-lname=nexus-kube-controllermanager",
	"-lcontrol-plane=api-gw",
	"-lapp=nexus-validation",
}

var EnvList = []string{
	"GOPRIVATE=gitlab.eng.vmware.com",
}

func GetEnvList() []string {
	envList := []string{}
	envList = append(envList, EnvList...)
	return envList

}
