/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/pkg/cli"
	cfgv2 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/pkg/config/v2"
	cfgv3 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/pkg/config/v3"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/pkg/plugin"
	kustomizecommonv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/pkg/plugins/common/kustomize/v1"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/pkg/plugins/golang"
	declarativev1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/pkg/plugins/golang/declarative/v1"
	golangv2 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/pkg/plugins/golang/v2"
	golangv3 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/pkg/plugins/golang/v3"
)

func main() {

	// Bundle plugin which built the golang projects scaffold by Kubebuilder go/v3
	gov3Bundle, _ := plugin.NewBundle(golang.DefaultNameQualifier, plugin.Version{Number: 3},
		kustomizecommonv1.Plugin{},
		golangv3.Plugin{},
	)

	c, err := cli.New(
		cli.WithCommandName("kubebuilder"),
		cli.WithVersion(versionString()),
		cli.WithPlugins(
			golangv2.Plugin{},
			golangv3.Plugin{},
			gov3Bundle,
			&kustomizecommonv1.Plugin{},
			&declarativev1.Plugin{},
		),
		cli.WithDefaultPlugins(cfgv2.Version, golangv2.Plugin{}),
		cli.WithDefaultPlugins(cfgv3.Version, gov3Bundle),
		cli.WithDefaultProjectVersion(cfgv3.Version),
		cli.WithCompletion(),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
