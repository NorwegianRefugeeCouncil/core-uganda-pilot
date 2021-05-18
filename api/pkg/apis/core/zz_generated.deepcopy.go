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

package core

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormDefinition) DeepCopyInto(out *FormDefinition) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormDefinition.
func (in *FormDefinition) DeepCopy() *FormDefinition {
	if in == nil {
		return nil
	}
	out := new(FormDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FormDefinition) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormDefinitionList) DeepCopyInto(out *FormDefinitionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FormDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormDefinitionList.
func (in *FormDefinitionList) DeepCopy() *FormDefinitionList {
	if in == nil {
		return nil
	}
	out := new(FormDefinitionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FormDefinitionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormDefinitionNames) DeepCopyInto(out *FormDefinitionNames) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormDefinitionNames.
func (in *FormDefinitionNames) DeepCopy() *FormDefinitionNames {
	if in == nil {
		return nil
	}
	out := new(FormDefinitionNames)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormDefinitionSchema) DeepCopyInto(out *FormDefinitionSchema) {
	*out = *in
	in.Root.DeepCopyInto(&out.Root)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormDefinitionSchema.
func (in *FormDefinitionSchema) DeepCopy() *FormDefinitionSchema {
	if in == nil {
		return nil
	}
	out := new(FormDefinitionSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormDefinitionSpec) DeepCopyInto(out *FormDefinitionSpec) {
	*out = *in
	out.Names = in.Names
	if in.Versions != nil {
		in, out := &in.Versions, &out.Versions
		*out = make([]FormDefinitionVersion, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormDefinitionSpec.
func (in *FormDefinitionSpec) DeepCopy() *FormDefinitionSpec {
	if in == nil {
		return nil
	}
	out := new(FormDefinitionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormDefinitionValidation) DeepCopyInto(out *FormDefinitionValidation) {
	*out = *in
	in.FormSchema.DeepCopyInto(&out.FormSchema)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormDefinitionValidation.
func (in *FormDefinitionValidation) DeepCopy() *FormDefinitionValidation {
	if in == nil {
		return nil
	}
	out := new(FormDefinitionValidation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormDefinitionVersion) DeepCopyInto(out *FormDefinitionVersion) {
	*out = *in
	in.Schema.DeepCopyInto(&out.Schema)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormDefinitionVersion.
func (in *FormDefinitionVersion) DeepCopy() *FormDefinitionVersion {
	if in == nil {
		return nil
	}
	out := new(FormDefinitionVersion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FormElementDefinition) DeepCopyInto(out *FormElementDefinition) {
	*out = *in
	if in.Label != nil {
		in, out := &in.Label, &out.Label
		*out = make(TranslatedStrings, len(*in))
		copy(*out, *in)
	}
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = make(TranslatedStrings, len(*in))
		copy(*out, *in)
	}
	if in.Help != nil {
		in, out := &in.Help, &out.Help
		*out = make(TranslatedStrings, len(*in))
		copy(*out, *in)
	}
	if in.Children != nil {
		in, out := &in.Children, &out.Children
		*out = make([]FormElementDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MaxLength != nil {
		in, out := &in.MaxLength, &out.MaxLength
		*out = new(int64)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FormElementDefinition.
func (in *FormElementDefinition) DeepCopy() *FormElementDefinition {
	if in == nil {
		return nil
	}
	out := new(FormElementDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TranslatedString) DeepCopyInto(out *TranslatedString) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TranslatedString.
func (in *TranslatedString) DeepCopy() *TranslatedString {
	if in == nil {
		return nil
	}
	out := new(TranslatedString)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in TranslatedStrings) DeepCopyInto(out *TranslatedStrings) {
	{
		in := &in
		*out = make(TranslatedStrings, len(*in))
		copy(*out, *in)
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TranslatedStrings.
func (in TranslatedStrings) DeepCopy() TranslatedStrings {
	if in == nil {
		return nil
	}
	out := new(TranslatedStrings)
	in.DeepCopyInto(out)
	return *out
}
