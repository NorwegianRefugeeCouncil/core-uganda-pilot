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

package discovery

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIGroup) DeepCopyInto(out *APIGroup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Versions != nil {
		in, out := &in.Versions, &out.Versions
		*out = make([]GroupVersionForDiscovery, len(*in))
		copy(*out, *in)
	}
	out.PreferredVersion = in.PreferredVersion
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIGroup.
func (in *APIGroup) DeepCopy() *APIGroup {
	if in == nil {
		return nil
	}
	out := new(APIGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *APIGroup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIGroupList) DeepCopyInto(out *APIGroupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Groups != nil {
		in, out := &in.Groups, &out.Groups
		*out = make([]APIGroup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIGroupList.
func (in *APIGroupList) DeepCopy() *APIGroupList {
	if in == nil {
		return nil
	}
	out := new(APIGroupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *APIGroupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIResource) DeepCopyInto(out *APIResource) {
	*out = *in
	if in.Verbs != nil {
		in, out := &in.Verbs, &out.Verbs
		*out = make(Verbs, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIResource.
func (in *APIResource) DeepCopy() *APIResource {
	if in == nil {
		return nil
	}
	out := new(APIResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIResourceList) DeepCopyInto(out *APIResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.APIResources != nil {
		in, out := &in.APIResources, &out.APIResources
		*out = make([]APIResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIResourceList.
func (in *APIResourceList) DeepCopy() *APIResourceList {
	if in == nil {
		return nil
	}
	out := new(APIResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *APIResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIService) DeepCopyInto(out *APIService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIService.
func (in *APIService) DeepCopy() *APIService {
	if in == nil {
		return nil
	}
	out := new(APIService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *APIService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIServiceList) DeepCopyInto(out *APIServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]APIService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIServiceList.
func (in *APIServiceList) DeepCopy() *APIServiceList {
	if in == nil {
		return nil
	}
	out := new(APIServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *APIServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIServiceSpec) DeepCopyInto(out *APIServiceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIServiceSpec.
func (in *APIServiceSpec) DeepCopy() *APIServiceSpec {
	if in == nil {
		return nil
	}
	out := new(APIServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GroupVersionForDiscovery) DeepCopyInto(out *GroupVersionForDiscovery) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GroupVersionForDiscovery.
func (in *GroupVersionForDiscovery) DeepCopy() *GroupVersionForDiscovery {
	if in == nil {
		return nil
	}
	out := new(GroupVersionForDiscovery)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in Verbs) DeepCopyInto(out *Verbs) {
	{
		in := &in
		*out = make(Verbs, len(*in))
		copy(*out, *in)
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Verbs.
func (in Verbs) DeepCopy() Verbs {
	if in == nil {
		return nil
	}
	out := new(Verbs)
	in.DeepCopyInto(out)
	return *out
}
