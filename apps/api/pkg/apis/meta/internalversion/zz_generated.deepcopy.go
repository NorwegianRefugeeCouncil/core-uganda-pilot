// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package internalversion

import (
	runtime "github.com/nrc-no/core/apps/api/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ListOptions) DeepCopyInto(out *ListOptions) {
	*out = *in
	out.TypeMeta = in.TypeMeta
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListOptions.
func (in *ListOptions) DeepCopy() *ListOptions {
	if in == nil {
		return nil
	}
	out := new(ListOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ListOptions) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
