/*
Copyright 2022.

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

package v1beta1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// CloudChain
// +k8s:openapi-gen=true
type CloudChain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CloudChainSpec   `json:"spec,omitempty"`
	Status CloudChainStatus `json:"status,omitempty"`
}

func (in *CloudChain) DeepCopyObject() runtime.Object {
	//TODO implement me
	panic("implement me")
}

// CloudChainList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CloudChainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CloudChain `json:"items"`
}

// CloudChainSpec defines the desired state of CloudChain
type CloudChainSpec struct {
}

var _ resource.Object = &CloudChain{}
var _ resourcestrategy.Validater = &CloudChain{}

func (in *CloudChain) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *CloudChain) NamespaceScoped() bool {
	return false
}

func (in *CloudChain) New() runtime.Object {
	return &CloudChain{}
}

func (in *CloudChain) NewList() runtime.Object {
	return &CloudChainList{}
}

func (in *CloudChain) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "chain.cloud.io",
		Version:  "v1beta1",
		Resource: "cloudchains",
	}
}

func (in *CloudChain) IsStorageVersion() bool {
	return true
}

func (in *CloudChain) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &CloudChainList{}

func (in *CloudChainList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

// CloudChainStatus defines the observed state of CloudChain
type CloudChainStatus struct {
}

func (in CloudChainStatus) SubResourceName() string {
	return "status"
}

// CloudChain implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &CloudChain{}

func (in *CloudChain) GetStatus() resource.StatusSubResource {
	return in.Status
}

// CloudChainStatus{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &CloudChainStatus{}

func (in CloudChainStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*CloudChain).Status = in
}

var _ resource.ObjectWithArbitrarySubResource = &CloudChain{}

func (in *CloudChain) GetArbitrarySubResources() []resource.ArbitrarySubResource {
	return []resource.ArbitrarySubResource{
		// +kubebuilder:scaffold:subresource
		&CloudChainFabric{},
	}
}
