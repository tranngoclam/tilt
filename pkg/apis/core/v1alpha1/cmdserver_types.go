/*
Copyright 2020 The Tilt Dev Authors

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

package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/tilt-dev/tilt-apiserver/pkg/server/builder/resource"
	"github.com/tilt-dev/tilt-apiserver/pkg/server/builder/resource/resourcestrategy"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CmdServer
// +k8s:openapi-gen=true
type CmdServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   CmdServerSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status CmdServerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// CmdServerList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CmdServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []CmdServer `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// CmdServerSpec defines the desired state of CmdServer
type CmdServerSpec struct {
	// Command-line arguments. Must have length at least 1.
	Args []string `json:"args,omitempty" protobuf:"bytes,1,rep,name=args"`

	// Process working directory.
	//
	// If the working directory is not specified, the command is run
	// in the default Tilt working directory.
	//	// +optional
	// +tilt:local-path=true
	Dir string `json:"dir" protobuf:"bytes,2,opt,name=dir"`

	// Additional variables process environment.
	//
	// Expressed as a C-style array of strings of the form ["KEY1=VALUE1", "KEY2=VALUE2", ...].
	//
	// Environment variables are layered on top of the environment variables
	// that Tilt runs with.
	//
	// +optional
	Env []string `json:"env" protobuf:"bytes,3,rep,name=env"`

	// Periodic probe of service readiness.
	//
	// +optional
	ReadinessProbe *Probe `json:"readinessProbe" protobuf:"bytes,4,opt,name=readinessProbe"`

	// Kubernetes tends to represent this as a "generation" field
	// to force an update.
	//
	// +optional
	TriggerTime metav1.MicroTime `json:"triggerTime" protobuf:"bytes,5,opt,name=triggerTime"`

	// Specifies how to disable this.
	//
	// +optional
	DisableSource *DisableSource `json:"disableSource,omitempty" protobuf:"bytes,6,opt,name=disableSource"`
}

var _ resource.Object = &CmdServer{}
var _ resourcestrategy.Validater = &CmdServer{}

func (in *CmdServer) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *CmdServer) NamespaceScoped() bool {
	return false
}

func (in *CmdServer) New() runtime.Object {
	return &CmdServer{}
}

func (in *CmdServer) NewList() runtime.Object {
	return &CmdServerList{}
}

func (in *CmdServer) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "tilt.dev",
		Version:  "v1alpha1",
		Resource: "cmdservers",
	}
}

func (in *CmdServer) IsStorageVersion() bool {
	return true
}

func (in *CmdServer) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &CmdServerList{}

func (in *CmdServerList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

// CmdServerStatus defines the observed state of CmdServer
type CmdServerStatus struct {
	// Details about whether/why this is disabled.
	// +optional
	DisableStatus *DisableStatus `json:"disableStatus,omitempty" protobuf:"bytes,1,opt,name=disableStatus"`
}

// CmdServer implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &CmdServer{}

func (in *CmdServer) GetStatus() resource.StatusSubResource {
	return in.Status
}

// CmdServerStatus{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &CmdServerStatus{}

func (in CmdServerStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*CmdServer).Status = in
}
