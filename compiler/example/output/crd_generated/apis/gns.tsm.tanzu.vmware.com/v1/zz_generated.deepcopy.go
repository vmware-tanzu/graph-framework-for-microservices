//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

import (
	v1alpha1 "github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdditionalDescription) DeepCopyInto(out *AdditionalDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdditionalDescription.
func (in *AdditionalDescription) DeepCopy() *AdditionalDescription {
	if in == nil {
		return nil
	}
	out := new(AdditionalDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdditionalGnsData) DeepCopyInto(out *AdditionalGnsData) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdditionalGnsData.
func (in *AdditionalGnsData) DeepCopy() *AdditionalGnsData {
	if in == nil {
		return nil
	}
	out := new(AdditionalGnsData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AdditionalGnsData) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdditionalGnsDataList) DeepCopyInto(out *AdditionalGnsDataList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AdditionalGnsData, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdditionalGnsDataList.
func (in *AdditionalGnsDataList) DeepCopy() *AdditionalGnsDataList {
	if in == nil {
		return nil
	}
	out := new(AdditionalGnsDataList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AdditionalGnsDataList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdditionalGnsDataNexusStatus) DeepCopyInto(out *AdditionalGnsDataNexusStatus) {
	*out = *in
	out.Status = in.Status
	out.Nexus = in.Nexus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdditionalGnsDataNexusStatus.
func (in *AdditionalGnsDataNexusStatus) DeepCopy() *AdditionalGnsDataNexusStatus {
	if in == nil {
		return nil
	}
	out := new(AdditionalGnsDataNexusStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdditionalGnsDataSpec) DeepCopyInto(out *AdditionalGnsDataSpec) {
	*out = *in
	out.Description = in.Description
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdditionalGnsDataSpec.
func (in *AdditionalGnsDataSpec) DeepCopy() *AdditionalGnsDataSpec {
	if in == nil {
		return nil
	}
	out := new(AdditionalGnsDataSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdditionalStatus) DeepCopyInto(out *AdditionalStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdditionalStatus.
func (in *AdditionalStatus) DeepCopy() *AdditionalStatus {
	if in == nil {
		return nil
	}
	out := new(AdditionalStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in AliasArr) DeepCopyInto(out *AliasArr) {
	{
		in := &in
		*out = make(AliasArr, len(*in))
		copy(*out, *in)
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AliasArr.
func (in AliasArr) DeepCopy() AliasArr {
	if in == nil {
		return nil
	}
	out := new(AliasArr)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Answer) DeepCopyInto(out *Answer) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Answer.
func (in *Answer) DeepCopy() *Answer {
	if in == nil {
		return nil
	}
	out := new(Answer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BarChild) DeepCopyInto(out *BarChild) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BarChild.
func (in *BarChild) DeepCopy() *BarChild {
	if in == nil {
		return nil
	}
	out := new(BarChild)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BarChild) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BarChildList) DeepCopyInto(out *BarChildList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BarChild, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BarChildList.
func (in *BarChildList) DeepCopy() *BarChildList {
	if in == nil {
		return nil
	}
	out := new(BarChildList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BarChildList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BarChildNexusStatus) DeepCopyInto(out *BarChildNexusStatus) {
	*out = *in
	out.Nexus = in.Nexus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BarChildNexusStatus.
func (in *BarChildNexusStatus) DeepCopy() *BarChildNexusStatus {
	if in == nil {
		return nil
	}
	out := new(BarChildNexusStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BarChildSpec) DeepCopyInto(out *BarChildSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BarChildSpec.
func (in *BarChildSpec) DeepCopy() *BarChildSpec {
	if in == nil {
		return nil
	}
	out := new(BarChildSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Child) DeepCopyInto(out *Child) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Child.
func (in *Child) DeepCopy() *Child {
	if in == nil {
		return nil
	}
	out := new(Child)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Description) DeepCopyInto(out *Description) {
	*out = *in
	if in.TestAns != nil {
		in, out := &in.TestAns, &out.TestAns
		*out = make([]Answer, len(*in))
		copy(*out, *in)
	}
	out.HostPort = in.HostPort
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Description.
func (in *Description) DeepCopy() *Description {
	if in == nil {
		return nil
	}
	out := new(Description)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Dns) DeepCopyInto(out *Dns) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Dns.
func (in *Dns) DeepCopy() *Dns {
	if in == nil {
		return nil
	}
	out := new(Dns)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Dns) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DnsList) DeepCopyInto(out *DnsList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Dns, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DnsList.
func (in *DnsList) DeepCopy() *DnsList {
	if in == nil {
		return nil
	}
	out := new(DnsList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DnsList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DnsNexusStatus) DeepCopyInto(out *DnsNexusStatus) {
	*out = *in
	out.Nexus = in.Nexus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DnsNexusStatus.
func (in *DnsNexusStatus) DeepCopy() *DnsNexusStatus {
	if in == nil {
		return nil
	}
	out := new(DnsNexusStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Gns) DeepCopyInto(out *Gns) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Gns.
func (in *Gns) DeepCopy() *Gns {
	if in == nil {
		return nil
	}
	out := new(Gns)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Gns) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GnsList) DeepCopyInto(out *GnsList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Gns, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GnsList.
func (in *GnsList) DeepCopy() *GnsList {
	if in == nil {
		return nil
	}
	out := new(GnsList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GnsList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GnsNexusStatus) DeepCopyInto(out *GnsNexusStatus) {
	*out = *in
	out.State = in.State
	out.Nexus = in.Nexus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GnsNexusStatus.
func (in *GnsNexusStatus) DeepCopy() *GnsNexusStatus {
	if in == nil {
		return nil
	}
	out := new(GnsNexusStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GnsSpec) DeepCopyInto(out *GnsSpec) {
	*out = *in
	in.Description.DeepCopyInto(&out.Description)
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int)
		**out = **in
	}
	if in.OtherDescription != nil {
		in, out := &in.OtherDescription, &out.OtherDescription
		*out = new(Description)
		(*in).DeepCopyInto(*out)
	}
	if in.MapPointer != nil {
		in, out := &in.MapPointer, &out.MapPointer
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
	if in.SlicePointer != nil {
		in, out := &in.SlicePointer, &out.SlicePointer
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	in.WorkloadSpec.DeepCopyInto(&out.WorkloadSpec)
	if in.DifferentSpec != nil {
		in, out := &in.DifferentSpec, &out.DifferentSpec
		*out = new(v1alpha1.WorkloadSpec)
		(*in).DeepCopyInto(*out)
	}
	out.ServiceSegmentRef = in.ServiceSegmentRef
	if in.ServiceSegmentRefPointer != nil {
		in, out := &in.ServiceSegmentRefPointer, &out.ServiceSegmentRefPointer
		*out = new(ServiceSegmentRef)
		**out = **in
	}
	if in.ServiceSegmentRefs != nil {
		in, out := &in.ServiceSegmentRefs, &out.ServiceSegmentRefs
		*out = make([]ServiceSegmentRef, len(*in))
		copy(*out, *in)
	}
	if in.ServiceSegmentRefMap != nil {
		in, out := &in.ServiceSegmentRefMap, &out.ServiceSegmentRefMap
		*out = make(map[string]ServiceSegmentRef, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.GnsServiceGroupsGvk != nil {
		in, out := &in.GnsServiceGroupsGvk, &out.GnsServiceGroupsGvk
		*out = make(map[string]Child, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.GnsAccessControlPolicyGvk != nil {
		in, out := &in.GnsAccessControlPolicyGvk, &out.GnsAccessControlPolicyGvk
		*out = new(Child)
		**out = **in
	}
	if in.FooChildGvk != nil {
		in, out := &in.FooChildGvk, &out.FooChildGvk
		*out = new(Child)
		**out = **in
	}
	if in.IgnoreChildGvk != nil {
		in, out := &in.IgnoreChildGvk, &out.IgnoreChildGvk
		*out = new(Child)
		**out = **in
	}
	if in.DnsGvk != nil {
		in, out := &in.DnsGvk, &out.DnsGvk
		*out = new(Link)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GnsSpec.
func (in *GnsSpec) DeepCopy() *GnsSpec {
	if in == nil {
		return nil
	}
	out := new(GnsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GnsState) DeepCopyInto(out *GnsState) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GnsState.
func (in *GnsState) DeepCopy() *GnsState {
	if in == nil {
		return nil
	}
	out := new(GnsState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostPort) DeepCopyInto(out *HostPort) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostPort.
func (in *HostPort) DeepCopy() *HostPort {
	if in == nil {
		return nil
	}
	out := new(HostPort)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IgnoreChild) DeepCopyInto(out *IgnoreChild) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IgnoreChild.
func (in *IgnoreChild) DeepCopy() *IgnoreChild {
	if in == nil {
		return nil
	}
	out := new(IgnoreChild)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IgnoreChild) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IgnoreChildList) DeepCopyInto(out *IgnoreChildList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IgnoreChild, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IgnoreChildList.
func (in *IgnoreChildList) DeepCopy() *IgnoreChildList {
	if in == nil {
		return nil
	}
	out := new(IgnoreChildList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IgnoreChildList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IgnoreChildNexusStatus) DeepCopyInto(out *IgnoreChildNexusStatus) {
	*out = *in
	out.Nexus = in.Nexus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IgnoreChildNexusStatus.
func (in *IgnoreChildNexusStatus) DeepCopy() *IgnoreChildNexusStatus {
	if in == nil {
		return nil
	}
	out := new(IgnoreChildNexusStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IgnoreChildSpec) DeepCopyInto(out *IgnoreChildSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IgnoreChildSpec.
func (in *IgnoreChildSpec) DeepCopy() *IgnoreChildSpec {
	if in == nil {
		return nil
	}
	out := new(IgnoreChildSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Link) DeepCopyInto(out *Link) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Link.
func (in *Link) DeepCopy() *Link {
	if in == nil {
		return nil
	}
	out := new(Link)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NexusStatus) DeepCopyInto(out *NexusStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NexusStatus.
func (in *NexusStatus) DeepCopy() *NexusStatus {
	if in == nil {
		return nil
	}
	out := new(NexusStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RandomDescription) DeepCopyInto(out *RandomDescription) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RandomDescription.
func (in *RandomDescription) DeepCopy() *RandomDescription {
	if in == nil {
		return nil
	}
	out := new(RandomDescription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RandomGnsData) DeepCopyInto(out *RandomGnsData) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RandomGnsData.
func (in *RandomGnsData) DeepCopy() *RandomGnsData {
	if in == nil {
		return nil
	}
	out := new(RandomGnsData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RandomGnsData) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RandomGnsDataList) DeepCopyInto(out *RandomGnsDataList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RandomGnsData, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RandomGnsDataList.
func (in *RandomGnsDataList) DeepCopy() *RandomGnsDataList {
	if in == nil {
		return nil
	}
	out := new(RandomGnsDataList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RandomGnsDataList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RandomGnsDataNexusStatus) DeepCopyInto(out *RandomGnsDataNexusStatus) {
	*out = *in
	out.Status = in.Status
	out.Nexus = in.Nexus
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RandomGnsDataNexusStatus.
func (in *RandomGnsDataNexusStatus) DeepCopy() *RandomGnsDataNexusStatus {
	if in == nil {
		return nil
	}
	out := new(RandomGnsDataNexusStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RandomGnsDataSpec) DeepCopyInto(out *RandomGnsDataSpec) {
	*out = *in
	out.Description = in.Description
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RandomGnsDataSpec.
func (in *RandomGnsDataSpec) DeepCopy() *RandomGnsDataSpec {
	if in == nil {
		return nil
	}
	out := new(RandomGnsDataSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RandomStatus) DeepCopyInto(out *RandomStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RandomStatus.
func (in *RandomStatus) DeepCopy() *RandomStatus {
	if in == nil {
		return nil
	}
	out := new(RandomStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicationSource) DeepCopyInto(out *ReplicationSource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicationSource.
func (in *ReplicationSource) DeepCopy() *ReplicationSource {
	if in == nil {
		return nil
	}
	out := new(ReplicationSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceSegmentRef) DeepCopyInto(out *ServiceSegmentRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceSegmentRef.
func (in *ServiceSegmentRef) DeepCopy() *ServiceSegmentRef {
	if in == nil {
		return nil
	}
	out := new(ServiceSegmentRef)
	in.DeepCopyInto(out)
	return out
}
