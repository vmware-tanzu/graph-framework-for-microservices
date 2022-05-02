package nexus_client

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"

	baseClientset "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/client/clientset/versioned"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/helper"

	baseconfigtsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/apis/config.tsm.tanzu.vmware.com/v1"
	basegnstsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/apis/gns.tsm.tanzu.vmware.com/v1"
	basepolicytsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/apis/policy.tsm.tanzu.vmware.com/v1"
	baseroottsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/apis/root.tsm.tanzu.vmware.com/v1"
	baseservicegrouptsmtanzuvmwarecomv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/_generated/apis/servicegroup.tsm.tanzu.vmware.com/v1"
)

type Clientset struct {
	baseClient        *baseClientset.Clientset
	rootTsmV1         *RootTsmV1
	configTsmV1       *ConfigTsmV1
	gnsTsmV1          *GnsTsmV1
	servicegroupTsmV1 *ServicegroupTsmV1
	policyTsmV1       *PolicyTsmV1
}

func NewForConfig(config *rest.Config) (*Clientset, error) {
	baseClient, err := baseClientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client := &Clientset{}
	client.baseClient = baseClient
	client.rootTsmV1 = newRootTsmV1(client)
	client.configTsmV1 = newConfigTsmV1(client)
	client.gnsTsmV1 = newGnsTsmV1(client)
	client.servicegroupTsmV1 = newServicegroupTsmV1(client)
	client.policyTsmV1 = newPolicyTsmV1(client)

	return client, nil
}

type PatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type Patch []PatchOp

func (p Patch) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (c *Clientset) Root() *RootTsmV1 {
	return c.rootTsmV1
}
func (c *Clientset) Config() *ConfigTsmV1 {
	return c.configTsmV1
}
func (c *Clientset) Gns() *GnsTsmV1 {
	return c.gnsTsmV1
}
func (c *Clientset) Servicegroup() *ServicegroupTsmV1 {
	return c.servicegroupTsmV1
}
func (c *Clientset) Policy() *PolicyTsmV1 {
	return c.policyTsmV1
}

type RootTsmV1 struct {
	roots *rootRootTsmV1
}

func newRootTsmV1(client *Clientset) *RootTsmV1 {
	return &RootTsmV1{
		roots: &rootRootTsmV1{
			client: client,
		},
	}
}

type rootRootTsmV1 struct {
	client *Clientset
}

func (obj *RootTsmV1) Roots() *rootRootTsmV1 {
	return obj.roots
}

type ConfigTsmV1 struct {
	configs *configConfigTsmV1
}

func newConfigTsmV1(client *Clientset) *ConfigTsmV1 {
	return &ConfigTsmV1{
		configs: &configConfigTsmV1{
			client: client,
		},
	}
}

type configConfigTsmV1 struct {
	client *Clientset
}

func (obj *ConfigTsmV1) Configs() *configConfigTsmV1 {
	return obj.configs
}

type GnsTsmV1 struct {
	gnses *gnsGnsTsmV1
	dnses *dnsGnsTsmV1
}

func newGnsTsmV1(client *Clientset) *GnsTsmV1 {
	return &GnsTsmV1{
		gnses: &gnsGnsTsmV1{
			client: client,
		},
		dnses: &dnsGnsTsmV1{
			client: client,
		},
	}
}

type gnsGnsTsmV1 struct {
	client *Clientset
}

func (obj *GnsTsmV1) Gnses() *gnsGnsTsmV1 {
	return obj.gnses
}

type dnsGnsTsmV1 struct {
	client *Clientset
}

func (obj *GnsTsmV1) Dnses() *dnsGnsTsmV1 {
	return obj.dnses
}

type ServicegroupTsmV1 struct {
	svcgroups *svcgroupServicegroupTsmV1
}

func newServicegroupTsmV1(client *Clientset) *ServicegroupTsmV1 {
	return &ServicegroupTsmV1{
		svcgroups: &svcgroupServicegroupTsmV1{
			client: client,
		},
	}
}

type svcgroupServicegroupTsmV1 struct {
	client *Clientset
}

func (obj *ServicegroupTsmV1) SvcGroups() *svcgroupServicegroupTsmV1 {
	return obj.svcgroups
}

type PolicyTsmV1 struct {
	accesscontrolpolicies *accesscontrolpolicyPolicyTsmV1
	acpconfigs            *acpconfigPolicyTsmV1
}

func newPolicyTsmV1(client *Clientset) *PolicyTsmV1 {
	return &PolicyTsmV1{
		accesscontrolpolicies: &accesscontrolpolicyPolicyTsmV1{
			client: client,
		},
		acpconfigs: &acpconfigPolicyTsmV1{
			client: client,
		},
	}
}

type accesscontrolpolicyPolicyTsmV1 struct {
	client *Clientset
}

func (obj *PolicyTsmV1) AccessControlPolicies() *accesscontrolpolicyPolicyTsmV1 {
	return obj.accesscontrolpolicies
}

type acpconfigPolicyTsmV1 struct {
	client *Clientset
}

func (obj *PolicyTsmV1) ACPConfigs() *acpconfigPolicyTsmV1 {
	return obj.acpconfigs
}

// Get hashes object's name and returns stored kubernetes object with all children and softlinks.
// To resolve a hashed name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *rootRootTsmV1) Get(ctx context.Context, name string, parents map[string]string) (result *baseroottsmtanzuvmwarecomv1.Root, err error) {
	hashedName := helper.GetHashedName("roots.root.tsm.tanzu.vmware.com", parents, name)
	return obj.GetByName(ctx, hashedName)
}

// GetByName works as Get but without hashing a name
func (obj *rootRootTsmV1) GetByName(ctx context.Context, name string) (result *baseroottsmtanzuvmwarecomv1.Root, err error) {
	result, err = obj.client.baseClient.RootTsmV1().Roots().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return
}

// Delete hashes object's name and deletes the object and all it's children
// To resolve a hash names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *rootRootTsmV1) Delete(ctx context.Context, name string, parents map[string]string) (err error) {
	if parents == nil {
		parents = map[string]string{}
	}
	parents[helper.IS_NAME_HASHED_LABEL] = "true"
	hashedName := helper.GetHashedName("roots.root.tsm.tanzu.vmware.com", parents, name)
	return obj.DeleteByName(ctx, hashedName, parents)
}

// DeleteByName works as Delete but without hashing a name
func (obj *rootRootTsmV1) DeleteByName(ctx context.Context, name string, parents map[string]string) (err error) {

	result, err := obj.client.baseClient.RootTsmV1().Roots().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if parents == nil {
		parents = make(map[string]string, 1)
	}

	if _, ok := result.GetLabels()[helper.DISPLAY_NAME_LABEL]; ok {
		parents["roots.root.tsm.tanzu.vmware.com"] = result.GetLabels()[helper.DISPLAY_NAME_LABEL]
	} else {
		parents["roots.root.tsm.tanzu.vmware.com"] = name
	}

	if result.Spec.ConfigGvk != nil {
		err := obj.client.Config().Configs().DeleteByName(ctx, result.Spec.ConfigGvk.Name, parents)
		if err != nil {
			return err
		}
	}

	err = obj.client.baseClient.RootTsmV1().Roots().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return
}

// Create hashes object's name and creates an object in the apiserver. Only spec fields can be provided, links and
// children can't be added using this function.
// To hash object's name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *rootRootTsmV1) Create(ctx context.Context, objToCreate *baseroottsmtanzuvmwarecomv1.Root, parents map[string]string) (result *baseroottsmtanzuvmwarecomv1.Root, err error) {
	if objToCreate.Labels == nil {
		objToCreate.Labels = map[string]string{}
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
		objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("roots.root.tsm.tanzu.vmware.com", parents, objToCreate.GetName())
		objToCreate.Name = hashedName
	}
	return obj.CreateByName(ctx, objToCreate, parents)
}

// CreateByName works as Create but without hashing the name
func (obj *rootRootTsmV1) CreateByName(ctx context.Context, objToCreate *baseroottsmtanzuvmwarecomv1.Root, parents map[string]string) (result *baseroottsmtanzuvmwarecomv1.Root, err error) {
	for k, v := range parents {
		objToCreate.Labels[k] = v
	}
	if _, ok := objToCreate.Labels[helper.DISPLAY_NAME_LABEL]; !ok {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
	}

	objToCreate.Spec.ConfigGvk = nil

	result, err = obj.client.baseClient.RootTsmV1().Roots().Create(ctx, objToCreate, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// Update hashes object's name and updates an object in the apiserver. Only spec fields and metadata can be updated,
// links and children can't be added or updated using this function.
// To hash the name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *rootRootTsmV1) Update(ctx context.Context, objToUpdate *baseroottsmtanzuvmwarecomv1.Root, parents map[string]string) (result *baseroottsmtanzuvmwarecomv1.Root, err error) {
	if objToUpdate.Labels == nil {
		objToUpdate.Labels = map[string]string{}
	}
	if objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToUpdate.Labels[helper.DISPLAY_NAME_LABEL] = objToUpdate.GetName()
		objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("roots.root.tsm.tanzu.vmware.com", parents, objToUpdate.GetName())
		objToUpdate.Name = hashedName
	}
	return obj.UpdateByName(ctx, objToUpdate)
}

// UpdateByName works as Update but without hashing the name
func (obj *rootRootTsmV1) UpdateByName(ctx context.Context, objToUpdate *baseroottsmtanzuvmwarecomv1.Root) (result *baseroottsmtanzuvmwarecomv1.Root, err error) {
	// ResourceVersion must be set for update
	if objToUpdate.ResourceVersion == "" {
		current, err := obj.client.baseClient.RootTsmV1().Roots().Get(ctx, objToUpdate.GetName(), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		objToUpdate.ResourceVersion = current.ResourceVersion
	}

	var patch Patch
	patchOpMeta := PatchOp{
		Op:    "replace",
		Path:  "/metadata",
		Value: objToUpdate.ObjectMeta,
	}
	patch = append(patch, patchOpMeta)

	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.RootTsmV1().Roots().Patch(ctx, objToUpdate.GetName(), types.JSONPatchType, marshaled, metav1.PatchOptions{}, "")
	if err != nil {
		return nil, err
	}

	return
}

func (obj *rootRootTsmV1) List(ctx context.Context, opts metav1.ListOptions) (result *baseroottsmtanzuvmwarecomv1.RootList, err error) {
	return obj.client.baseClient.RootTsmV1().Roots().List(ctx, opts)
}

func (obj *rootRootTsmV1) GetConfig(ctx context.Context, srcObj *baseroottsmtanzuvmwarecomv1.Root) (result *baseconfigtsmtanzuvmwarecomv1.Config, err error) {
	return obj.GetConfigByName(ctx, srcObj)
}

func (obj *rootRootTsmV1) GetConfigByName(ctx context.Context, srcObj *baseroottsmtanzuvmwarecomv1.Root) (result *baseconfigtsmtanzuvmwarecomv1.Config, err error) {
	if srcObj.Spec.ConfigGvk != nil {
		return obj.client.baseClient.ConfigTsmV1().Configs().Get(ctx, srcObj.Spec.ConfigGvk.Name, metav1.GetOptions{})
	} else {
		var parentName string
		if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
			parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
		} else {
			parentName = srcObj.GetName()
		}
		return nil, fmt.Errorf("There's no child Config for parent %s", parentName)
	}

}

// Get hashes object's name and returns stored kubernetes object with all children and softlinks.
// To resolve a hashed name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *configConfigTsmV1) Get(ctx context.Context, name string, parents map[string]string) (result *baseconfigtsmtanzuvmwarecomv1.Config, err error) {
	hashedName := helper.GetHashedName("configs.config.tsm.tanzu.vmware.com", parents, name)
	return obj.GetByName(ctx, hashedName)
}

// GetByName works as Get but without hashing a name
func (obj *configConfigTsmV1) GetByName(ctx context.Context, name string) (result *baseconfigtsmtanzuvmwarecomv1.Config, err error) {
	result, err = obj.client.baseClient.ConfigTsmV1().Configs().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return
}

// Delete hashes object's name and deletes the object and all it's children
// To resolve a hash names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *configConfigTsmV1) Delete(ctx context.Context, name string, parents map[string]string) (err error) {
	if parents == nil {
		parents = map[string]string{}
	}
	parents[helper.IS_NAME_HASHED_LABEL] = "true"
	hashedName := helper.GetHashedName("configs.config.tsm.tanzu.vmware.com", parents, name)
	return obj.DeleteByName(ctx, hashedName, parents)
}

// DeleteByName works as Delete but without hashing a name
func (obj *configConfigTsmV1) DeleteByName(ctx context.Context, name string, parents map[string]string) (err error) {

	result, err := obj.client.baseClient.ConfigTsmV1().Configs().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if parents == nil {
		parents = make(map[string]string, 1)
	}

	if _, ok := result.GetLabels()[helper.DISPLAY_NAME_LABEL]; ok {
		parents["configs.config.tsm.tanzu.vmware.com"] = result.GetLabels()[helper.DISPLAY_NAME_LABEL]
	} else {
		parents["configs.config.tsm.tanzu.vmware.com"] = name
	}

	if result.Spec.GNSGvk != nil {
		err := obj.client.Gns().Gnses().DeleteByName(ctx, result.Spec.GNSGvk.Name, parents)
		if err != nil {
			return err
		}
	}

	err = obj.client.baseClient.ConfigTsmV1().Configs().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	var patch Patch

	patchOp := PatchOp{
		Op:   "remove",
		Path: "/spec/configGvk",
	}

	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return err
	}
	parentName, ok := parents["roots.root.tsm.tanzu.vmware.com"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if parents[helper.IS_NAME_HASHED_LABEL] == "true" {
		parentName = helper.GetHashedName("roots.root.tsm.tanzu.vmware.com", parents, parentName)
	}
	_, err = obj.client.baseClient.RootTsmV1().Roots().Patch(ctx, parentName, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return
}

// Create hashes object's name and creates an object in the apiserver. Only spec fields can be provided, links and
// children can't be added using this function.
// To hash object's name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *configConfigTsmV1) Create(ctx context.Context, objToCreate *baseconfigtsmtanzuvmwarecomv1.Config, parents map[string]string) (result *baseconfigtsmtanzuvmwarecomv1.Config, err error) {
	if objToCreate.Labels == nil {
		objToCreate.Labels = map[string]string{}
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
		objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("configs.config.tsm.tanzu.vmware.com", parents, objToCreate.GetName())
		objToCreate.Name = hashedName
	}
	return obj.CreateByName(ctx, objToCreate, parents)
}

// CreateByName works as Create but without hashing the name
func (obj *configConfigTsmV1) CreateByName(ctx context.Context, objToCreate *baseconfigtsmtanzuvmwarecomv1.Config, parents map[string]string) (result *baseconfigtsmtanzuvmwarecomv1.Config, err error) {
	for k, v := range parents {
		objToCreate.Labels[k] = v
	}
	if _, ok := objToCreate.Labels[helper.DISPLAY_NAME_LABEL]; !ok {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
	}

	objToCreate.Spec.GNSGvk = nil

	result, err = obj.client.baseClient.ConfigTsmV1().Configs().Create(ctx, objToCreate, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	parentName, ok := parents["roots.root.tsm.tanzu.vmware.com"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] == "true" {
		parentName = helper.GetHashedName("roots.root.tsm.tanzu.vmware.com", parents, parentName)
	}

	var patch Patch
	patchOp := PatchOp{
		Op:   "replace",
		Path: "/spec/configGvk",
		Value: baseconfigtsmtanzuvmwarecomv1.Child{
			Group: "config.tsm.tanzu.vmware.com",
			Kind:  "Config",
			Name:  objToCreate.Name,
		},
	}
	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	_, err = obj.client.baseClient.RootTsmV1().Roots().Patch(ctx, parentName, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// Update hashes object's name and updates an object in the apiserver. Only spec fields and metadata can be updated,
// links and children can't be added or updated using this function.
// To hash the name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *configConfigTsmV1) Update(ctx context.Context, objToUpdate *baseconfigtsmtanzuvmwarecomv1.Config, parents map[string]string) (result *baseconfigtsmtanzuvmwarecomv1.Config, err error) {
	if objToUpdate.Labels == nil {
		objToUpdate.Labels = map[string]string{}
	}
	if objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToUpdate.Labels[helper.DISPLAY_NAME_LABEL] = objToUpdate.GetName()
		objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("configs.config.tsm.tanzu.vmware.com", parents, objToUpdate.GetName())
		objToUpdate.Name = hashedName
	}
	return obj.UpdateByName(ctx, objToUpdate)
}

// UpdateByName works as Update but without hashing the name
func (obj *configConfigTsmV1) UpdateByName(ctx context.Context, objToUpdate *baseconfigtsmtanzuvmwarecomv1.Config) (result *baseconfigtsmtanzuvmwarecomv1.Config, err error) {
	// ResourceVersion must be set for update
	if objToUpdate.ResourceVersion == "" {
		current, err := obj.client.baseClient.ConfigTsmV1().Configs().Get(ctx, objToUpdate.GetName(), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		objToUpdate.ResourceVersion = current.ResourceVersion
	}

	var patch Patch
	patchOpMeta := PatchOp{
		Op:    "replace",
		Path:  "/metadata",
		Value: objToUpdate.ObjectMeta,
	}
	patch = append(patch, patchOpMeta)

	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.ConfigTsmV1().Configs().Patch(ctx, objToUpdate.GetName(), types.JSONPatchType, marshaled, metav1.PatchOptions{}, "")
	if err != nil {
		return nil, err
	}

	return
}

func (obj *configConfigTsmV1) List(ctx context.Context, opts metav1.ListOptions) (result *baseconfigtsmtanzuvmwarecomv1.ConfigList, err error) {
	return obj.client.baseClient.ConfigTsmV1().Configs().List(ctx, opts)
}

func (obj *configConfigTsmV1) GetGNS(ctx context.Context, srcObj *baseconfigtsmtanzuvmwarecomv1.Config) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {
	return obj.GetGNSByName(ctx, srcObj)
}

func (obj *configConfigTsmV1) GetGNSByName(ctx context.Context, srcObj *baseconfigtsmtanzuvmwarecomv1.Config) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {
	if srcObj.Spec.GNSGvk != nil {
		return obj.client.baseClient.GnsTsmV1().Gnses().Get(ctx, srcObj.Spec.GNSGvk.Name, metav1.GetOptions{})
	} else {
		var parentName string
		if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
			parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
		} else {
			parentName = srcObj.GetName()
		}
		return nil, fmt.Errorf("There's no child GNS for parent %s", parentName)
	}

}

// Get hashes object's name and returns stored kubernetes object with all children and softlinks.
// To resolve a hashed name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *gnsGnsTsmV1) Get(ctx context.Context, name string, parents map[string]string) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {
	hashedName := helper.GetHashedName("gnses.gns.tsm.tanzu.vmware.com", parents, name)
	return obj.GetByName(ctx, hashedName)
}

// GetByName works as Get but without hashing a name
func (obj *gnsGnsTsmV1) GetByName(ctx context.Context, name string) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {
	result, err = obj.client.baseClient.GnsTsmV1().Gnses().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return
}

// Delete hashes object's name and deletes the object and all it's children
// To resolve a hash names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *gnsGnsTsmV1) Delete(ctx context.Context, name string, parents map[string]string) (err error) {
	if parents == nil {
		parents = map[string]string{}
	}
	parents[helper.IS_NAME_HASHED_LABEL] = "true"
	hashedName := helper.GetHashedName("gnses.gns.tsm.tanzu.vmware.com", parents, name)
	return obj.DeleteByName(ctx, hashedName, parents)
}

// DeleteByName works as Delete but without hashing a name
func (obj *gnsGnsTsmV1) DeleteByName(ctx context.Context, name string, parents map[string]string) (err error) {

	result, err := obj.client.baseClient.GnsTsmV1().Gnses().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if parents == nil {
		parents = make(map[string]string, 1)
	}

	if _, ok := result.GetLabels()[helper.DISPLAY_NAME_LABEL]; ok {
		parents["gnses.gns.tsm.tanzu.vmware.com"] = result.GetLabels()[helper.DISPLAY_NAME_LABEL]
	} else {
		parents["gnses.gns.tsm.tanzu.vmware.com"] = name
	}

	for _, v := range result.Spec.GnsServiceGroupsGvk {
		err := obj.client.Servicegroup().SvcGroups().DeleteByName(ctx, v.Name, parents)
		if err != nil {
			return err
		}
	}

	if result.Spec.GnsAccessControlPolicyGvk != nil {
		err := obj.client.Policy().AccessControlPolicies().DeleteByName(ctx, result.Spec.GnsAccessControlPolicyGvk.Name, parents)
		if err != nil {
			return err
		}
	}

	err = obj.client.baseClient.GnsTsmV1().Gnses().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	var patch Patch

	patchOp := PatchOp{
		Op:   "remove",
		Path: "/spec/gNSGvk",
	}

	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return err
	}
	parentName, ok := parents["configs.config.tsm.tanzu.vmware.com"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if parents[helper.IS_NAME_HASHED_LABEL] == "true" {
		parentName = helper.GetHashedName("configs.config.tsm.tanzu.vmware.com", parents, parentName)
	}
	_, err = obj.client.baseClient.ConfigTsmV1().Configs().Patch(ctx, parentName, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return
}

// Create hashes object's name and creates an object in the apiserver. Only spec fields can be provided, links and
// children can't be added using this function.
// To hash object's name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *gnsGnsTsmV1) Create(ctx context.Context, objToCreate *basegnstsmtanzuvmwarecomv1.Gns, parents map[string]string) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {
	if objToCreate.Labels == nil {
		objToCreate.Labels = map[string]string{}
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
		objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("gnses.gns.tsm.tanzu.vmware.com", parents, objToCreate.GetName())
		objToCreate.Name = hashedName
	}
	return obj.CreateByName(ctx, objToCreate, parents)
}

// CreateByName works as Create but without hashing the name
func (obj *gnsGnsTsmV1) CreateByName(ctx context.Context, objToCreate *basegnstsmtanzuvmwarecomv1.Gns, parents map[string]string) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {
	for k, v := range parents {
		objToCreate.Labels[k] = v
	}
	if _, ok := objToCreate.Labels[helper.DISPLAY_NAME_LABEL]; !ok {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
	}

	objToCreate.Spec.GnsServiceGroupsGvk = nil
	objToCreate.Spec.GnsAccessControlPolicyGvk = nil
	objToCreate.Spec.DnsGvk = nil

	result, err = obj.client.baseClient.GnsTsmV1().Gnses().Create(ctx, objToCreate, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	parentName, ok := parents["configs.config.tsm.tanzu.vmware.com"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] == "true" {
		parentName = helper.GetHashedName("configs.config.tsm.tanzu.vmware.com", parents, parentName)
	}

	var patch Patch
	patchOp := PatchOp{
		Op:   "replace",
		Path: "/spec/gNSGvk",
		Value: basegnstsmtanzuvmwarecomv1.Child{
			Group: "gns.tsm.tanzu.vmware.com",
			Kind:  "Gns",
			Name:  objToCreate.Name,
		},
	}
	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	_, err = obj.client.baseClient.ConfigTsmV1().Configs().Patch(ctx, parentName, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// Update hashes object's name and updates an object in the apiserver. Only spec fields and metadata can be updated,
// links and children can't be added or updated using this function.
// To hash the name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *gnsGnsTsmV1) Update(ctx context.Context, objToUpdate *basegnstsmtanzuvmwarecomv1.Gns, parents map[string]string) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {
	if objToUpdate.Labels == nil {
		objToUpdate.Labels = map[string]string{}
	}
	if objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToUpdate.Labels[helper.DISPLAY_NAME_LABEL] = objToUpdate.GetName()
		objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("gnses.gns.tsm.tanzu.vmware.com", parents, objToUpdate.GetName())
		objToUpdate.Name = hashedName
	}
	return obj.UpdateByName(ctx, objToUpdate)
}

// UpdateByName works as Update but without hashing the name
func (obj *gnsGnsTsmV1) UpdateByName(ctx context.Context, objToUpdate *basegnstsmtanzuvmwarecomv1.Gns) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {
	// ResourceVersion must be set for update
	if objToUpdate.ResourceVersion == "" {
		current, err := obj.client.baseClient.GnsTsmV1().Gnses().Get(ctx, objToUpdate.GetName(), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		objToUpdate.ResourceVersion = current.ResourceVersion
	}

	var patch Patch
	patchOpMeta := PatchOp{
		Op:    "replace",
		Path:  "/metadata",
		Value: objToUpdate.ObjectMeta,
	}
	patch = append(patch, patchOpMeta)

	patchValueDomain := objToUpdate.Spec.Domain
	patchOpDomain := PatchOp{
		Op:    "replace",
		Path:  "/spec/domain",
		Value: patchValueDomain,
	}
	patch = append(patch, patchOpDomain)

	patchValueUseSharedGateway := objToUpdate.Spec.UseSharedGateway
	patchOpUseSharedGateway := PatchOp{
		Op:    "replace",
		Path:  "/spec/useSharedGateway",
		Value: patchValueUseSharedGateway,
	}
	patch = append(patch, patchOpUseSharedGateway)

	patchValueDescription := objToUpdate.Spec.Description
	patchOpDescription := PatchOp{
		Op:    "replace",
		Path:  "/spec/description",
		Value: patchValueDescription,
	}
	patch = append(patch, patchOpDescription)

	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.GnsTsmV1().Gnses().Patch(ctx, objToUpdate.GetName(), types.JSONPatchType, marshaled, metav1.PatchOptions{}, "")
	if err != nil {
		return nil, err
	}

	return
}

func (obj *gnsGnsTsmV1) List(ctx context.Context, opts metav1.ListOptions) (result *basegnstsmtanzuvmwarecomv1.GnsList, err error) {
	return obj.client.baseClient.GnsTsmV1().Gnses().List(ctx, opts)
}

// AddDns updates srcObj with linkToAdd object
func (obj *gnsGnsTsmV1) AddDns(ctx context.Context, srcObj *basegnstsmtanzuvmwarecomv1.Gns, linkToAdd *basegnstsmtanzuvmwarecomv1.Dns) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {

	var patch Patch
	patchOp := PatchOp{
		Op:   "replace",
		Path: "/spec/dnsGvk",
		Value: basegnstsmtanzuvmwarecomv1.Child{
			Group: "gns.tsm.tanzu.vmware.com",
			Kind:  "Dns",
			Name:  linkToAdd.Name,
		},
	}
	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.GnsTsmV1().Gnses().Patch(ctx, srcObj.Name, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// RemoveDns removes linkToRemove object from srcObj
func (obj *gnsGnsTsmV1) RemoveDns(ctx context.Context, srcObj *basegnstsmtanzuvmwarecomv1.Gns, linkToRemove *basegnstsmtanzuvmwarecomv1.Dns) (result *basegnstsmtanzuvmwarecomv1.Gns, err error) {
	var patch Patch

	patchOp := PatchOp{
		Op:   "remove",
		Path: "/spec/dnsGvk",
	}

	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.GnsTsmV1().Gnses().Patch(ctx, srcObj.Name, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

func (obj *gnsGnsTsmV1) GetGnsServiceGroups(ctx context.Context, srcObj *basegnstsmtanzuvmwarecomv1.Gns, GnsServiceGroupsObjectName string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	parents := srcObj.GetLabels()
	if parents == nil {
		parents = make(map[string]string, 1)
	}
	if _, ok := srcObj.GetLabels()[helper.DISPLAY_NAME_LABEL]; ok {
		parents["gnses.gns.tsm.tanzu.vmware.com"] = srcObj.GetLabels()[helper.DISPLAY_NAME_LABEL]
	} else {
		parents["gnses.gns.tsm.tanzu.vmware.com"] = srcObj.GetName()
	}
	hashedName := helper.GetHashedName(result.CRDName(), parents, GnsServiceGroupsObjectName)
	return obj.GetGnsServiceGroupsByName(ctx, srcObj, hashedName)

}

func (obj *gnsGnsTsmV1) GetGnsServiceGroupsByName(ctx context.Context, srcObj *basegnstsmtanzuvmwarecomv1.Gns, GnsServiceGroupsObjectName string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	if srcObj.Spec.GnsServiceGroupsGvk != nil {
		l, ok := srcObj.Spec.GnsServiceGroupsGvk[GnsServiceGroupsObjectName]
		if ok {
			return obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Get(ctx, l.Name, metav1.GetOptions{})
		}
	}
	var parentName string
	if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
		parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
	} else {
		parentName = srcObj.GetName()
	}
	return nil, fmt.Errorf("there's no child %s for parent %s", GnsServiceGroupsObjectName, parentName)

}

func (obj *gnsGnsTsmV1) GetAllGnsServiceGroups(ctx context.Context, srcObj *basegnstsmtanzuvmwarecomv1.Gns) (result []*baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	result = make([]*baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, 0, len(srcObj.Spec.GnsServiceGroupsGvk))
	for _, v := range srcObj.Spec.GnsServiceGroupsGvk {
		l, err := obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Get(ctx, v.Name, metav1.GetOptions{})
		if err != nil {
			var parentName string
			if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
				parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
			} else {
				parentName = srcObj.GetName()
			}
			return nil, fmt.Errorf("couldn't get child %s for parent %s: %v", v.Name, parentName, err)
		}
		result = append(result, l)
	}
	return
}

func (obj *gnsGnsTsmV1) GetGnsAccessControlPolicy(ctx context.Context, srcObj *basegnstsmtanzuvmwarecomv1.Gns) (result *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	return obj.GetGnsAccessControlPolicyByName(ctx, srcObj)
}

func (obj *gnsGnsTsmV1) GetGnsAccessControlPolicyByName(ctx context.Context, srcObj *basegnstsmtanzuvmwarecomv1.Gns) (result *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	if srcObj.Spec.GnsAccessControlPolicyGvk != nil {
		return obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().Get(ctx, srcObj.Spec.GnsAccessControlPolicyGvk.Name, metav1.GetOptions{})
	} else {
		var parentName string
		if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
			parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
		} else {
			parentName = srcObj.GetName()
		}
		return nil, fmt.Errorf("There's no child GnsAccessControlPolicy for parent %s", parentName)
	}

}

func (obj *gnsGnsTsmV1) GetDns(ctx context.Context, srcObj *basegnstsmtanzuvmwarecomv1.Gns) (result *basegnstsmtanzuvmwarecomv1.Dns, err error) {
	return obj.GetDnsByName(ctx, srcObj)
}

func (obj *gnsGnsTsmV1) GetDnsByName(ctx context.Context, srcObj *basegnstsmtanzuvmwarecomv1.Gns) (result *basegnstsmtanzuvmwarecomv1.Dns, err error) {
	if srcObj.Spec.DnsGvk != nil {
		return obj.client.baseClient.GnsTsmV1().Dnses().Get(ctx, srcObj.Spec.DnsGvk.Name, metav1.GetOptions{})
	} else {
		var parentName string
		if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
			parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
		} else {
			parentName = srcObj.GetName()
		}
		return nil, fmt.Errorf("There's no child Dns for parent %s", parentName)
	}

}

// Get hashes object's name and returns stored kubernetes object with all children and softlinks.
// To resolve a hashed name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *dnsGnsTsmV1) Get(ctx context.Context, name string, parents map[string]string) (result *basegnstsmtanzuvmwarecomv1.Dns, err error) {
	hashedName := helper.GetHashedName("dnses.gns.tsm.tanzu.vmware.com", parents, name)
	return obj.GetByName(ctx, hashedName)
}

// GetByName works as Get but without hashing a name
func (obj *dnsGnsTsmV1) GetByName(ctx context.Context, name string) (result *basegnstsmtanzuvmwarecomv1.Dns, err error) {
	result, err = obj.client.baseClient.GnsTsmV1().Dnses().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return
}

// Delete hashes object's name and deletes the object and all it's children
// To resolve a hash names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *dnsGnsTsmV1) Delete(ctx context.Context, name string, parents map[string]string) (err error) {
	if parents == nil {
		parents = map[string]string{}
	}
	parents[helper.IS_NAME_HASHED_LABEL] = "true"
	hashedName := helper.GetHashedName("dnses.gns.tsm.tanzu.vmware.com", parents, name)
	return obj.DeleteByName(ctx, hashedName, parents)
}

// DeleteByName works as Delete but without hashing a name
func (obj *dnsGnsTsmV1) DeleteByName(ctx context.Context, name string, parents map[string]string) (err error) {

	err = obj.client.baseClient.GnsTsmV1().Dnses().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return
}

// Create hashes object's name and creates an object in the apiserver. Only spec fields can be provided, links and
// children can't be added using this function.
// To hash object's name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *dnsGnsTsmV1) Create(ctx context.Context, objToCreate *basegnstsmtanzuvmwarecomv1.Dns, parents map[string]string) (result *basegnstsmtanzuvmwarecomv1.Dns, err error) {
	if objToCreate.Labels == nil {
		objToCreate.Labels = map[string]string{}
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
		objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("dnses.gns.tsm.tanzu.vmware.com", parents, objToCreate.GetName())
		objToCreate.Name = hashedName
	}
	return obj.CreateByName(ctx, objToCreate, parents)
}

// CreateByName works as Create but without hashing the name
func (obj *dnsGnsTsmV1) CreateByName(ctx context.Context, objToCreate *basegnstsmtanzuvmwarecomv1.Dns, parents map[string]string) (result *basegnstsmtanzuvmwarecomv1.Dns, err error) {
	for k, v := range parents {
		objToCreate.Labels[k] = v
	}
	if _, ok := objToCreate.Labels[helper.DISPLAY_NAME_LABEL]; !ok {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
	}

	result, err = obj.client.baseClient.GnsTsmV1().Dnses().Create(ctx, objToCreate, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// Update hashes object's name and updates an object in the apiserver. Only spec fields and metadata can be updated,
// links and children can't be added or updated using this function.
// To hash the name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *dnsGnsTsmV1) Update(ctx context.Context, objToUpdate *basegnstsmtanzuvmwarecomv1.Dns, parents map[string]string) (result *basegnstsmtanzuvmwarecomv1.Dns, err error) {
	if objToUpdate.Labels == nil {
		objToUpdate.Labels = map[string]string{}
	}
	if objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToUpdate.Labels[helper.DISPLAY_NAME_LABEL] = objToUpdate.GetName()
		objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("dnses.gns.tsm.tanzu.vmware.com", parents, objToUpdate.GetName())
		objToUpdate.Name = hashedName
	}
	return obj.UpdateByName(ctx, objToUpdate)
}

// UpdateByName works as Update but without hashing the name
func (obj *dnsGnsTsmV1) UpdateByName(ctx context.Context, objToUpdate *basegnstsmtanzuvmwarecomv1.Dns) (result *basegnstsmtanzuvmwarecomv1.Dns, err error) {
	// ResourceVersion must be set for update
	if objToUpdate.ResourceVersion == "" {
		current, err := obj.client.baseClient.GnsTsmV1().Dnses().Get(ctx, objToUpdate.GetName(), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		objToUpdate.ResourceVersion = current.ResourceVersion
	}

	var patch Patch
	patchOpMeta := PatchOp{
		Op:    "replace",
		Path:  "/metadata",
		Value: objToUpdate.ObjectMeta,
	}
	patch = append(patch, patchOpMeta)

	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.GnsTsmV1().Dnses().Patch(ctx, objToUpdate.GetName(), types.JSONPatchType, marshaled, metav1.PatchOptions{}, "")
	if err != nil {
		return nil, err
	}

	return
}

func (obj *dnsGnsTsmV1) List(ctx context.Context, opts metav1.ListOptions) (result *basegnstsmtanzuvmwarecomv1.DnsList, err error) {
	return obj.client.baseClient.GnsTsmV1().Dnses().List(ctx, opts)
}

// Get hashes object's name and returns stored kubernetes object with all children and softlinks.
// To resolve a hashed name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *svcgroupServicegroupTsmV1) Get(ctx context.Context, name string, parents map[string]string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	hashedName := helper.GetHashedName("svcgroups.servicegroup.tsm.tanzu.vmware.com", parents, name)
	return obj.GetByName(ctx, hashedName)
}

// GetByName works as Get but without hashing a name
func (obj *svcgroupServicegroupTsmV1) GetByName(ctx context.Context, name string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	result, err = obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return
}

// Delete hashes object's name and deletes the object and all it's children
// To resolve a hash names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *svcgroupServicegroupTsmV1) Delete(ctx context.Context, name string, parents map[string]string) (err error) {
	if parents == nil {
		parents = map[string]string{}
	}
	parents[helper.IS_NAME_HASHED_LABEL] = "true"
	hashedName := helper.GetHashedName("svcgroups.servicegroup.tsm.tanzu.vmware.com", parents, name)
	return obj.DeleteByName(ctx, hashedName, parents)
}

// DeleteByName works as Delete but without hashing a name
func (obj *svcgroupServicegroupTsmV1) DeleteByName(ctx context.Context, name string, parents map[string]string) (err error) {

	err = obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	var patch Patch

	patchOp := PatchOp{
		Op:   "remove",
		Path: "/spec/gnsServiceGroupsGvk/" + name,
	}

	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return err
	}
	parentName, ok := parents["gnses.gns.tsm.tanzu.vmware.com"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if parents[helper.IS_NAME_HASHED_LABEL] == "true" {
		parentName = helper.GetHashedName("gnses.gns.tsm.tanzu.vmware.com", parents, parentName)
	}
	_, err = obj.client.baseClient.GnsTsmV1().Gnses().Patch(ctx, parentName, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return
}

// Create hashes object's name and creates an object in the apiserver. Only spec fields can be provided, links and
// children can't be added using this function.
// To hash object's name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *svcgroupServicegroupTsmV1) Create(ctx context.Context, objToCreate *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, parents map[string]string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	if objToCreate.Labels == nil {
		objToCreate.Labels = map[string]string{}
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
		objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("svcgroups.servicegroup.tsm.tanzu.vmware.com", parents, objToCreate.GetName())
		objToCreate.Name = hashedName
	}
	return obj.CreateByName(ctx, objToCreate, parents)
}

// CreateByName works as Create but without hashing the name
func (obj *svcgroupServicegroupTsmV1) CreateByName(ctx context.Context, objToCreate *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, parents map[string]string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	for k, v := range parents {
		objToCreate.Labels[k] = v
	}
	if _, ok := objToCreate.Labels[helper.DISPLAY_NAME_LABEL]; !ok {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
	}

	result, err = obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Create(ctx, objToCreate, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	parentName, ok := parents["gnses.gns.tsm.tanzu.vmware.com"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] == "true" {
		parentName = helper.GetHashedName("gnses.gns.tsm.tanzu.vmware.com", parents, parentName)
	}

	payload := "{\"spec\": {\"gnsServiceGroupsGvk\": {\"" + objToCreate.Name + "\": {\"name\": \"" + objToCreate.Name + "\",\"kind\": \"SvcGroup\", \"group\": \"servicegroup.tsm.tanzu.vmware.com\"}}}}"
	_, err = obj.client.baseClient.GnsTsmV1().Gnses().Patch(ctx, parentName, types.MergePatchType, []byte(payload), metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// Update hashes object's name and updates an object in the apiserver. Only spec fields and metadata can be updated,
// links and children can't be added or updated using this function.
// To hash the name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *svcgroupServicegroupTsmV1) Update(ctx context.Context, objToUpdate *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, parents map[string]string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	if objToUpdate.Labels == nil {
		objToUpdate.Labels = map[string]string{}
	}
	if objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToUpdate.Labels[helper.DISPLAY_NAME_LABEL] = objToUpdate.GetName()
		objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("svcgroups.servicegroup.tsm.tanzu.vmware.com", parents, objToUpdate.GetName())
		objToUpdate.Name = hashedName
	}
	return obj.UpdateByName(ctx, objToUpdate)
}

// UpdateByName works as Update but without hashing the name
func (obj *svcgroupServicegroupTsmV1) UpdateByName(ctx context.Context, objToUpdate *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	// ResourceVersion must be set for update
	if objToUpdate.ResourceVersion == "" {
		current, err := obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Get(ctx, objToUpdate.GetName(), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		objToUpdate.ResourceVersion = current.ResourceVersion
	}

	var patch Patch
	patchOpMeta := PatchOp{
		Op:    "replace",
		Path:  "/metadata",
		Value: objToUpdate.ObjectMeta,
	}
	patch = append(patch, patchOpMeta)

	patchValueDisplayName := objToUpdate.Spec.DisplayName
	patchOpDisplayName := PatchOp{
		Op:    "replace",
		Path:  "/spec/displayName",
		Value: patchValueDisplayName,
	}
	patch = append(patch, patchOpDisplayName)

	patchValueDescription := objToUpdate.Spec.Description
	patchOpDescription := PatchOp{
		Op:    "replace",
		Path:  "/spec/description",
		Value: patchValueDescription,
	}
	patch = append(patch, patchOpDescription)

	patchValueColor := objToUpdate.Spec.Color
	patchOpColor := PatchOp{
		Op:    "replace",
		Path:  "/spec/color",
		Value: patchValueColor,
	}
	patch = append(patch, patchOpColor)

	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Patch(ctx, objToUpdate.GetName(), types.JSONPatchType, marshaled, metav1.PatchOptions{}, "")
	if err != nil {
		return nil, err
	}

	return
}

func (obj *svcgroupServicegroupTsmV1) List(ctx context.Context, opts metav1.ListOptions) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroupList, err error) {
	return obj.client.baseClient.ServicegroupTsmV1().SvcGroups().List(ctx, opts)
}

// Get hashes object's name and returns stored kubernetes object with all children and softlinks.
// To resolve a hashed name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *accesscontrolpolicyPolicyTsmV1) Get(ctx context.Context, name string, parents map[string]string) (result *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	hashedName := helper.GetHashedName("accesscontrolpolicies.policy.tsm.tanzu.vmware.com", parents, name)
	return obj.GetByName(ctx, hashedName)
}

// GetByName works as Get but without hashing a name
func (obj *accesscontrolpolicyPolicyTsmV1) GetByName(ctx context.Context, name string) (result *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	result, err = obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return
}

// Delete hashes object's name and deletes the object and all it's children
// To resolve a hash names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *accesscontrolpolicyPolicyTsmV1) Delete(ctx context.Context, name string, parents map[string]string) (err error) {
	if parents == nil {
		parents = map[string]string{}
	}
	parents[helper.IS_NAME_HASHED_LABEL] = "true"
	hashedName := helper.GetHashedName("accesscontrolpolicies.policy.tsm.tanzu.vmware.com", parents, name)
	return obj.DeleteByName(ctx, hashedName, parents)
}

// DeleteByName works as Delete but without hashing a name
func (obj *accesscontrolpolicyPolicyTsmV1) DeleteByName(ctx context.Context, name string, parents map[string]string) (err error) {

	result, err := obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if parents == nil {
		parents = make(map[string]string, 1)
	}

	if _, ok := result.GetLabels()[helper.DISPLAY_NAME_LABEL]; ok {
		parents["accesscontrolpolicies.policy.tsm.tanzu.vmware.com"] = result.GetLabels()[helper.DISPLAY_NAME_LABEL]
	} else {
		parents["accesscontrolpolicies.policy.tsm.tanzu.vmware.com"] = name
	}

	for _, v := range result.Spec.PolicyConfigsGvk {
		err := obj.client.Policy().ACPConfigs().DeleteByName(ctx, v.Name, parents)
		if err != nil {
			return err
		}
	}

	err = obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	var patch Patch

	patchOp := PatchOp{
		Op:   "remove",
		Path: "/spec/gnsAccessControlPolicyGvk",
	}

	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return err
	}
	parentName, ok := parents["gnses.gns.tsm.tanzu.vmware.com"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if parents[helper.IS_NAME_HASHED_LABEL] == "true" {
		parentName = helper.GetHashedName("gnses.gns.tsm.tanzu.vmware.com", parents, parentName)
	}
	_, err = obj.client.baseClient.GnsTsmV1().Gnses().Patch(ctx, parentName, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return
}

// Create hashes object's name and creates an object in the apiserver. Only spec fields can be provided, links and
// children can't be added using this function.
// To hash object's name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *accesscontrolpolicyPolicyTsmV1) Create(ctx context.Context, objToCreate *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, parents map[string]string) (result *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	if objToCreate.Labels == nil {
		objToCreate.Labels = map[string]string{}
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
		objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("accesscontrolpolicies.policy.tsm.tanzu.vmware.com", parents, objToCreate.GetName())
		objToCreate.Name = hashedName
	}
	return obj.CreateByName(ctx, objToCreate, parents)
}

// CreateByName works as Create but without hashing the name
func (obj *accesscontrolpolicyPolicyTsmV1) CreateByName(ctx context.Context, objToCreate *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, parents map[string]string) (result *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	for k, v := range parents {
		objToCreate.Labels[k] = v
	}
	if _, ok := objToCreate.Labels[helper.DISPLAY_NAME_LABEL]; !ok {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
	}

	objToCreate.Spec.PolicyConfigsGvk = nil

	result, err = obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().Create(ctx, objToCreate, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	parentName, ok := parents["gnses.gns.tsm.tanzu.vmware.com"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] == "true" {
		parentName = helper.GetHashedName("gnses.gns.tsm.tanzu.vmware.com", parents, parentName)
	}

	var patch Patch
	patchOp := PatchOp{
		Op:   "replace",
		Path: "/spec/gnsAccessControlPolicyGvk",
		Value: basepolicytsmtanzuvmwarecomv1.Child{
			Group: "policy.tsm.tanzu.vmware.com",
			Kind:  "AccessControlPolicy",
			Name:  objToCreate.Name,
		},
	}
	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	_, err = obj.client.baseClient.GnsTsmV1().Gnses().Patch(ctx, parentName, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// Update hashes object's name and updates an object in the apiserver. Only spec fields and metadata can be updated,
// links and children can't be added or updated using this function.
// To hash the name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *accesscontrolpolicyPolicyTsmV1) Update(ctx context.Context, objToUpdate *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, parents map[string]string) (result *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	if objToUpdate.Labels == nil {
		objToUpdate.Labels = map[string]string{}
	}
	if objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToUpdate.Labels[helper.DISPLAY_NAME_LABEL] = objToUpdate.GetName()
		objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("accesscontrolpolicies.policy.tsm.tanzu.vmware.com", parents, objToUpdate.GetName())
		objToUpdate.Name = hashedName
	}
	return obj.UpdateByName(ctx, objToUpdate)
}

// UpdateByName works as Update but without hashing the name
func (obj *accesscontrolpolicyPolicyTsmV1) UpdateByName(ctx context.Context, objToUpdate *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy) (result *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, err error) {
	// ResourceVersion must be set for update
	if objToUpdate.ResourceVersion == "" {
		current, err := obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().Get(ctx, objToUpdate.GetName(), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		objToUpdate.ResourceVersion = current.ResourceVersion
	}

	var patch Patch
	patchOpMeta := PatchOp{
		Op:    "replace",
		Path:  "/metadata",
		Value: objToUpdate.ObjectMeta,
	}
	patch = append(patch, patchOpMeta)

	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().Patch(ctx, objToUpdate.GetName(), types.JSONPatchType, marshaled, metav1.PatchOptions{}, "")
	if err != nil {
		return nil, err
	}

	return
}

func (obj *accesscontrolpolicyPolicyTsmV1) List(ctx context.Context, opts metav1.ListOptions) (result *basepolicytsmtanzuvmwarecomv1.AccessControlPolicyList, err error) {
	return obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().List(ctx, opts)
}

func (obj *accesscontrolpolicyPolicyTsmV1) GetPolicyConfigs(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, PolicyConfigsObjectName string) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	parents := srcObj.GetLabels()
	if parents == nil {
		parents = make(map[string]string, 1)
	}
	if _, ok := srcObj.GetLabels()[helper.DISPLAY_NAME_LABEL]; ok {
		parents["accesscontrolpolicies.policy.tsm.tanzu.vmware.com"] = srcObj.GetLabels()[helper.DISPLAY_NAME_LABEL]
	} else {
		parents["accesscontrolpolicies.policy.tsm.tanzu.vmware.com"] = srcObj.GetName()
	}
	hashedName := helper.GetHashedName(result.CRDName(), parents, PolicyConfigsObjectName)
	return obj.GetPolicyConfigsByName(ctx, srcObj, hashedName)

}

func (obj *accesscontrolpolicyPolicyTsmV1) GetPolicyConfigsByName(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy, PolicyConfigsObjectName string) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	if srcObj.Spec.PolicyConfigsGvk != nil {
		l, ok := srcObj.Spec.PolicyConfigsGvk[PolicyConfigsObjectName]
		if ok {
			return obj.client.baseClient.PolicyTsmV1().ACPConfigs().Get(ctx, l.Name, metav1.GetOptions{})
		}
	}
	var parentName string
	if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
		parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
	} else {
		parentName = srcObj.GetName()
	}
	return nil, fmt.Errorf("there's no child %s for parent %s", PolicyConfigsObjectName, parentName)

}

func (obj *accesscontrolpolicyPolicyTsmV1) GetAllPolicyConfigs(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.AccessControlPolicy) (result []*basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	result = make([]*basepolicytsmtanzuvmwarecomv1.ACPConfig, 0, len(srcObj.Spec.PolicyConfigsGvk))
	for _, v := range srcObj.Spec.PolicyConfigsGvk {
		l, err := obj.client.baseClient.PolicyTsmV1().ACPConfigs().Get(ctx, v.Name, metav1.GetOptions{})
		if err != nil {
			var parentName string
			if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
				parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
			} else {
				parentName = srcObj.GetName()
			}
			return nil, fmt.Errorf("couldn't get child %s for parent %s: %v", v.Name, parentName, err)
		}
		result = append(result, l)
	}
	return
}

// Get hashes object's name and returns stored kubernetes object with all children and softlinks.
// To resolve a hashed name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *acpconfigPolicyTsmV1) Get(ctx context.Context, name string, parents map[string]string) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	hashedName := helper.GetHashedName("acpconfigs.policy.tsm.tanzu.vmware.com", parents, name)
	return obj.GetByName(ctx, hashedName)
}

// GetByName works as Get but without hashing a name
func (obj *acpconfigPolicyTsmV1) GetByName(ctx context.Context, name string) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	result, err = obj.client.baseClient.PolicyTsmV1().ACPConfigs().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return
}

// Delete hashes object's name and deletes the object and all it's children
// To resolve a hash names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *acpconfigPolicyTsmV1) Delete(ctx context.Context, name string, parents map[string]string) (err error) {
	if parents == nil {
		parents = map[string]string{}
	}
	parents[helper.IS_NAME_HASHED_LABEL] = "true"
	hashedName := helper.GetHashedName("acpconfigs.policy.tsm.tanzu.vmware.com", parents, name)
	return obj.DeleteByName(ctx, hashedName, parents)
}

// DeleteByName works as Delete but without hashing a name
func (obj *acpconfigPolicyTsmV1) DeleteByName(ctx context.Context, name string, parents map[string]string) (err error) {

	err = obj.client.baseClient.PolicyTsmV1().ACPConfigs().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	var patch Patch

	patchOp := PatchOp{
		Op:   "remove",
		Path: "/spec/policyConfigsGvk/" + name,
	}

	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return err
	}
	parentName, ok := parents["accesscontrolpolicies.policy.tsm.tanzu.vmware.com"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if parents[helper.IS_NAME_HASHED_LABEL] == "true" {
		parentName = helper.GetHashedName("accesscontrolpolicies.policy.tsm.tanzu.vmware.com", parents, parentName)
	}
	_, err = obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().Patch(ctx, parentName, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	return
}

// Create hashes object's name and creates an object in the apiserver. Only spec fields can be provided, links and
// children can't be added using this function.
// To hash object's name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *acpconfigPolicyTsmV1) Create(ctx context.Context, objToCreate *basepolicytsmtanzuvmwarecomv1.ACPConfig, parents map[string]string) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	if objToCreate.Labels == nil {
		objToCreate.Labels = map[string]string{}
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
		objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("acpconfigs.policy.tsm.tanzu.vmware.com", parents, objToCreate.GetName())
		objToCreate.Name = hashedName
	}
	return obj.CreateByName(ctx, objToCreate, parents)
}

// CreateByName works as Create but without hashing the name
func (obj *acpconfigPolicyTsmV1) CreateByName(ctx context.Context, objToCreate *basepolicytsmtanzuvmwarecomv1.ACPConfig, parents map[string]string) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	for k, v := range parents {
		objToCreate.Labels[k] = v
	}
	if _, ok := objToCreate.Labels[helper.DISPLAY_NAME_LABEL]; !ok {
		objToCreate.Labels[helper.DISPLAY_NAME_LABEL] = objToCreate.GetName()
	}

	objToCreate.Spec.DestSvcGroupsGvk = nil
	objToCreate.Spec.SourceSvcGroupsGvk = nil

	result, err = obj.client.baseClient.PolicyTsmV1().ACPConfigs().Create(ctx, objToCreate, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	parentName, ok := parents["accesscontrolpolicies.policy.tsm.tanzu.vmware.com"]
	if !ok {
		parentName = helper.DEFAULT_KEY
	}
	if objToCreate.Labels[helper.IS_NAME_HASHED_LABEL] == "true" {
		parentName = helper.GetHashedName("accesscontrolpolicies.policy.tsm.tanzu.vmware.com", parents, parentName)
	}

	payload := "{\"spec\": {\"policyConfigsGvk\": {\"" + objToCreate.Name + "\": {\"name\": \"" + objToCreate.Name + "\",\"kind\": \"ACPConfig\", \"group\": \"policy.tsm.tanzu.vmware.com\"}}}}"
	_, err = obj.client.baseClient.PolicyTsmV1().AccessControlPolicies().Patch(ctx, parentName, types.MergePatchType, []byte(payload), metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// Update hashes object's name and updates an object in the apiserver. Only spec fields and metadata can be updated,
// links and children can't be added or updated using this function.
// To hash the name names of all consecutive parents must be provided in parents param in form of:
// {'object_crd_definition_name': 'object_name'}
func (obj *acpconfigPolicyTsmV1) Update(ctx context.Context, objToUpdate *basepolicytsmtanzuvmwarecomv1.ACPConfig, parents map[string]string) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	if objToUpdate.Labels == nil {
		objToUpdate.Labels = map[string]string{}
	}
	if objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] != "true" {
		objToUpdate.Labels[helper.DISPLAY_NAME_LABEL] = objToUpdate.GetName()
		objToUpdate.Labels[helper.IS_NAME_HASHED_LABEL] = "true"
		hashedName := helper.GetHashedName("acpconfigs.policy.tsm.tanzu.vmware.com", parents, objToUpdate.GetName())
		objToUpdate.Name = hashedName
	}
	return obj.UpdateByName(ctx, objToUpdate)
}

// UpdateByName works as Update but without hashing the name
func (obj *acpconfigPolicyTsmV1) UpdateByName(ctx context.Context, objToUpdate *basepolicytsmtanzuvmwarecomv1.ACPConfig) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	// ResourceVersion must be set for update
	if objToUpdate.ResourceVersion == "" {
		current, err := obj.client.baseClient.PolicyTsmV1().ACPConfigs().Get(ctx, objToUpdate.GetName(), metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		objToUpdate.ResourceVersion = current.ResourceVersion
	}

	var patch Patch
	patchOpMeta := PatchOp{
		Op:    "replace",
		Path:  "/metadata",
		Value: objToUpdate.ObjectMeta,
	}
	patch = append(patch, patchOpMeta)

	patchValueDisplayName := objToUpdate.Spec.DisplayName
	patchOpDisplayName := PatchOp{
		Op:    "replace",
		Path:  "/spec/displayName",
		Value: patchValueDisplayName,
	}
	patch = append(patch, patchOpDisplayName)

	patchValueGns := objToUpdate.Spec.Gns
	patchOpGns := PatchOp{
		Op:    "replace",
		Path:  "/spec/gns",
		Value: patchValueGns,
	}
	patch = append(patch, patchOpGns)

	patchValueDescription := objToUpdate.Spec.Description
	patchOpDescription := PatchOp{
		Op:    "replace",
		Path:  "/spec/description",
		Value: patchValueDescription,
	}
	patch = append(patch, patchOpDescription)

	patchValueTags := objToUpdate.Spec.Tags
	patchOpTags := PatchOp{
		Op:    "replace",
		Path:  "/spec/tags",
		Value: patchValueTags,
	}
	patch = append(patch, patchOpTags)

	patchValueProjectId := objToUpdate.Spec.ProjectId
	patchOpProjectId := PatchOp{
		Op:    "replace",
		Path:  "/spec/projectId",
		Value: patchValueProjectId,
	}
	patch = append(patch, patchOpProjectId)

	patchValueConditions := objToUpdate.Spec.Conditions
	patchOpConditions := PatchOp{
		Op:    "replace",
		Path:  "/spec/conditions",
		Value: patchValueConditions,
	}
	patch = append(patch, patchOpConditions)

	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.PolicyTsmV1().ACPConfigs().Patch(ctx, objToUpdate.GetName(), types.JSONPatchType, marshaled, metav1.PatchOptions{}, "")
	if err != nil {
		return nil, err
	}

	return
}

func (obj *acpconfigPolicyTsmV1) List(ctx context.Context, opts metav1.ListOptions) (result *basepolicytsmtanzuvmwarecomv1.ACPConfigList, err error) {
	return obj.client.baseClient.PolicyTsmV1().ACPConfigs().List(ctx, opts)
}

// AddDestSvcGroups updates srcObj with linkToAdd object
func (obj *acpconfigPolicyTsmV1) AddDestSvcGroups(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.ACPConfig, linkToAdd *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {

	payload := "{\"spec\": {\"destSvcGroupsGvk\": {\"" + linkToAdd.Name + "\": {\"name\": \"" + linkToAdd.Name + "\",\"kind\": \"SvcGroup\", \"group\": \"servicegroup.tsm.tanzu.vmware.com\"}}}}"
	result, err = obj.client.baseClient.PolicyTsmV1().ACPConfigs().Patch(ctx, srcObj.Name, types.MergePatchType, []byte(payload), metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// RemoveDestSvcGroups removes linkToRemove object from srcObj
func (obj *acpconfigPolicyTsmV1) RemoveDestSvcGroups(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.ACPConfig, linkToRemove *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	var patch Patch

	patchOp := PatchOp{
		Op:   "remove",
		Path: "/spec/destSvcGroupsGvk/" + linkToRemove.Name,
	}

	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.PolicyTsmV1().ACPConfigs().Patch(ctx, srcObj.Name, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// AddSourceSvcGroups updates srcObj with linkToAdd object
func (obj *acpconfigPolicyTsmV1) AddSourceSvcGroups(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.ACPConfig, linkToAdd *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {

	payload := "{\"spec\": {\"sourceSvcGroupsGvk\": {\"" + linkToAdd.Name + "\": {\"name\": \"" + linkToAdd.Name + "\",\"kind\": \"SvcGroup\", \"group\": \"servicegroup.tsm.tanzu.vmware.com\"}}}}"
	result, err = obj.client.baseClient.PolicyTsmV1().ACPConfigs().Patch(ctx, srcObj.Name, types.MergePatchType, []byte(payload), metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

// RemoveSourceSvcGroups removes linkToRemove object from srcObj
func (obj *acpconfigPolicyTsmV1) RemoveSourceSvcGroups(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.ACPConfig, linkToRemove *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup) (result *basepolicytsmtanzuvmwarecomv1.ACPConfig, err error) {
	var patch Patch

	patchOp := PatchOp{
		Op:   "remove",
		Path: "/spec/sourceSvcGroupsGvk/" + linkToRemove.Name,
	}

	patch = append(patch, patchOp)
	marshaled, err := patch.Marshal()
	if err != nil {
		return nil, err
	}
	result, err = obj.client.baseClient.PolicyTsmV1().ACPConfigs().Patch(ctx, srcObj.Name, types.JSONPatchType, marshaled, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	return
}

func (obj *acpconfigPolicyTsmV1) GetDestSvcGroups(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.ACPConfig, DestSvcGroupsObjectName string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	parents := srcObj.GetLabels()
	if parents == nil {
		parents = make(map[string]string, 1)
	}
	if _, ok := srcObj.GetLabels()[helper.DISPLAY_NAME_LABEL]; ok {
		parents["acpconfigs.policy.tsm.tanzu.vmware.com"] = srcObj.GetLabels()[helper.DISPLAY_NAME_LABEL]
	} else {
		parents["acpconfigs.policy.tsm.tanzu.vmware.com"] = srcObj.GetName()
	}
	hashedName := helper.GetHashedName(result.CRDName(), parents, DestSvcGroupsObjectName)
	return obj.GetDestSvcGroupsByName(ctx, srcObj, hashedName)

}

func (obj *acpconfigPolicyTsmV1) GetDestSvcGroupsByName(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.ACPConfig, DestSvcGroupsObjectName string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	if srcObj.Spec.DestSvcGroupsGvk != nil {
		l, ok := srcObj.Spec.DestSvcGroupsGvk[DestSvcGroupsObjectName]
		if ok {
			return obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Get(ctx, l.Name, metav1.GetOptions{})
		}
	}
	var parentName string
	if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
		parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
	} else {
		parentName = srcObj.GetName()
	}
	return nil, fmt.Errorf("there's no child %s for parent %s", DestSvcGroupsObjectName, parentName)

}

func (obj *acpconfigPolicyTsmV1) GetAllDestSvcGroups(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.ACPConfig) (result []*baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	result = make([]*baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, 0, len(srcObj.Spec.DestSvcGroupsGvk))
	for _, v := range srcObj.Spec.DestSvcGroupsGvk {
		l, err := obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Get(ctx, v.Name, metav1.GetOptions{})
		if err != nil {
			var parentName string
			if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
				parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
			} else {
				parentName = srcObj.GetName()
			}
			return nil, fmt.Errorf("couldn't get child %s for parent %s: %v", v.Name, parentName, err)
		}
		result = append(result, l)
	}
	return
}

func (obj *acpconfigPolicyTsmV1) GetSourceSvcGroups(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.ACPConfig, SourceSvcGroupsObjectName string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	parents := srcObj.GetLabels()
	if parents == nil {
		parents = make(map[string]string, 1)
	}
	if _, ok := srcObj.GetLabels()[helper.DISPLAY_NAME_LABEL]; ok {
		parents["acpconfigs.policy.tsm.tanzu.vmware.com"] = srcObj.GetLabels()[helper.DISPLAY_NAME_LABEL]
	} else {
		parents["acpconfigs.policy.tsm.tanzu.vmware.com"] = srcObj.GetName()
	}
	hashedName := helper.GetHashedName(result.CRDName(), parents, SourceSvcGroupsObjectName)
	return obj.GetSourceSvcGroupsByName(ctx, srcObj, hashedName)

}

func (obj *acpconfigPolicyTsmV1) GetSourceSvcGroupsByName(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.ACPConfig, SourceSvcGroupsObjectName string) (result *baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	if srcObj.Spec.SourceSvcGroupsGvk != nil {
		l, ok := srcObj.Spec.SourceSvcGroupsGvk[SourceSvcGroupsObjectName]
		if ok {
			return obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Get(ctx, l.Name, metav1.GetOptions{})
		}
	}
	var parentName string
	if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
		parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
	} else {
		parentName = srcObj.GetName()
	}
	return nil, fmt.Errorf("there's no child %s for parent %s", SourceSvcGroupsObjectName, parentName)

}

func (obj *acpconfigPolicyTsmV1) GetAllSourceSvcGroups(ctx context.Context, srcObj *basepolicytsmtanzuvmwarecomv1.ACPConfig) (result []*baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, err error) {
	result = make([]*baseservicegrouptsmtanzuvmwarecomv1.SvcGroup, 0, len(srcObj.Spec.SourceSvcGroupsGvk))
	for _, v := range srcObj.Spec.SourceSvcGroupsGvk {
		l, err := obj.client.baseClient.ServicegroupTsmV1().SvcGroups().Get(ctx, v.Name, metav1.GetOptions{})
		if err != nil {
			var parentName string
			if srcObj.Labels != nil && srcObj.Labels[helper.DISPLAY_NAME_LABEL] != "" {
				parentName = srcObj.Labels[helper.DISPLAY_NAME_LABEL]
			} else {
				parentName = srcObj.GetName()
			}
			return nil, fmt.Errorf("couldn't get child %s for parent %s: %v", v.Name, parentName, err)
		}
		result = append(result, l)
	}
	return
}
