/*
Copyright 2023.

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
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MetricSpec defines the desired state of Metric
type MetricSetSpec struct {

	// The name of the metric (that will be associated with a flavor like storage)
	// +optional
	Metrics []Metric `json:"metrics"`

	// Don't set JobSet FQDN
	// +optional
	DontSetFQDN bool `json:"dontSetFQDN"`

	// Service name for the JobSet (MetricsSet) cluster network
	// +kubebuilder:default="ms"
	// +default="ms"
	// +optional
	ServiceName string `json:"serviceName"`

	// Should the job be limited to a particular number of seconds?
	// Approximately one year. This cannot be zero or job won't start
	// +kubebuilder:default=31500000
	// +default=31500000
	// +optional
	DeadlineSeconds int64 `json:"deadlineSeconds,omitempty"`

	// A storage setup that we want to measure performance for.
	// and binding to storage metrics
	// +optional
	Storage Storage `json:"storage"`

	// Pod spec for the application, standalone, or storage metrics
	//+optional
	Pod Pod `json:"pod"`

	// For metrics that require an application, we need a container and name (for now)
	// +optional
	Application Application `json:"application"`

	// Parallelism (e.g., pods)
	// +kubebuilder:default=1
	// +default=1
	// +optional
	Pods int32 `json:"pods"`

	// Resources include limits and requests for each pod (that include a JobSet)
	// +optional
	Resources ContainerResource `json:"resources"`

	// Logging spec, preparing for other kinds of logging
	// Right now we just include an interactive option
	//+optional
	Logging Logging `json:"logging"`

	// Single pod completion, meaning the jobspec completions is unset
	// and we only require one main completion
	// +optional
	Completions int32 `json:"completions"`
}

type Logging struct {

	// Don't allow the application, metric, or storage test to finish
	// This adds sleep infinity at the end to allow for interactive mode.
	// +optional
	Interactive bool `json:"interactive"`
}

// Pod attributes that can be given to an application or metric
type Pod struct {

	// name of service account to associate with pod
	//+optional
	ServiceAccountName string `json:"serviceAccountName"`

	// NodeSelector labels
	//+optional
	NodeSelector map[string]string `json:"nodeSelector"`
}

// A container spec can belong to a metric or application
type ContainerSpec struct {

	// Security context for the pod
	//+optional
	SecurityContext SecurityContext `json:"securityContext"`
}

type SecurityContext struct {
	Privileged bool `json:"privileged"`
}

// Storage that will be monitored, or storage alongside a standalone metric
type Storage struct {

	// Volume type to test (not all storage interfaces require one explicitly)
	//+optional
	Volume Volume `json:"volume"`

	// Commands to run (pre is supported to make bind)
	// +optional
	Commands Commands `json:"commands"`
}

// Application that will be monitored
type Application struct {
	Image string `json:"image"`

	// command to execute and monitor (if consistent across pods)
	Command string `json:"command"`

	// Working Directory
	//+optional
	WorkingDir string `json:"workingDir"`

	// Entrypoint of container, if different from command
	//+optional
	Entrypoint string `json:"entrypoint"`

	// A pull secret for the application container
	//+optional
	PullSecret string `json:"pullSecret"`

	// Resources include limits and requests for the application
	// +optional
	Resources ContainerResources `json:"resources"`

	// Container Spec has attributes for the container
	//+optional
	Attributes ContainerSpec `json:"attributes"`

	// Existing Volumes for the application
	// +optional
	Volumes map[string]Volume `json:"volumes"`
}

// ContainerResources include limits and requests
type ContainerResources struct {

	// +optional
	Limits ContainerResource `json:"limits"`

	// +optional
	Requests ContainerResource `json:"requests"`
}

type Commands struct {

	// pre command happens at start (before anything else)
	// +optional
	Pre string `json:"pre"`

	// Command prefix to put in front of a metric main command (not applicable for all)
	//+optional
	Prefix string `json:"prefix"`

	// post happens at end (after collection end)
	// +optional
	Post string `json:"post"`
}

type ContainerResource map[string]intstr.IntOrString

// A Volume should correspond with an existing volume, either:
// config map, secret, or claim name.
type Volume struct {

	// Path and claim name are always required if a secret isn't defined
	// +optional
	Path string `json:"path,omitempty"`

	// Hostpath volume on the host to bind to path
	// +optional
	HostPath string `json:"hostPath"`

	// Config map name if the existing volume is a config map
	// You should also define items if you are using this
	// +optional
	ConfigMapName string `json:"configMapName,omitempty"`

	// Items (key and paths) for the config map
	// +optional
	Items map[string]string `json:"items"`

	// Claim name if the existing volume is a PVC
	// +optional
	ClaimName string `json:"claimName,omitempty"`

	// An existing secret
	// +optional
	SecretName string `json:"secretName,omitempty"`

	// EmptyVol if true generates an empty volume at the path
	// +kubebuilder:default=false
	// +default=false
	// +optional
	EmptyVol bool `json:"emptyVol,omitempty"`

	// +kubebuilder:default=false
	// +default=false
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`
}

// The difference between benchmark and metric is subtle.
// A metric is more a measurment, and the benchmark is the comparison value.
// I don't have strong opinions but I think we are doing more measurment
// not necessarily with benchmarks
type Metric struct {
	Name string `json:"name"`

	// Metric Options
	// Metric specific options
	// +optional
	Options map[string]intstr.IntOrString `json:"options"`

	// Metric List Options
	// Metric specific options
	// +optional
	ListOptions map[string][]intstr.IntOrString `json:"listOptions"`

	// Metric Map Options
	// +optional
	MapOptions map[string]map[string]intstr.IntOrString `json:"mapOptions"`

	// Container Spec has attributes for the container
	//+optional
	Attributes ContainerSpec `json:"attributes"`

	// Resources include limits and requests for the metric container
	// +optional
	Resources ContainerResources `json:"resources"`
}

// Get pod labels for a metric set
func (m *MetricSet) GetPodLabels() map[string]string {

	podLabels := map[string]string{}
	// This is for autoscaling, although haven't used yet
	podLabels["cluster-name"] = m.Name
	// This is for the headless service
	podLabels["metricset-name"] = m.Name
	podLabels["namespace"] = m.Namespace
	podLabels["app.kubernetes.io/name"] = m.Name
	return podLabels
}

// MetricStatus defines the observed state of Metric
type MetricSetStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MetricSet is the Schema for the metrics API
type MetricSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetricSetSpec   `json:"spec,omitempty"`
	Status MetricSetStatus `json:"status,omitempty"`
}

// Determine if an application or storage is present, or standalone
func (m *MetricSet) HasApplication() bool {
	return !reflect.DeepEqual(m.Spec.Application, Application{})
}
func (m *MetricSet) HasStorage() bool {
	return !reflect.DeepEqual(m.Spec.Storage, Storage{})
}
func (m *MetricSet) HasStorageVolume() bool {
	return !reflect.DeepEqual(m.Spec.Storage.Volume, Volume{})
}
func (m *MetricSet) IsStandalone() bool {
	return !m.HasStorage() && !m.HasApplication()
}

// Validate a requested metricset
func (m *MetricSet) Validate() bool {

	// An application or storage setup is required
	if !m.HasApplication() && !m.HasStorage() && !m.IsStandalone() {
		fmt.Printf("😥️ An application OR storage OR standalone entry is required.\n")
		return false
	}

	// We don't currently support running both at once
	// (but should be fine to allow extra standalone)
	if m.HasApplication() && m.HasStorage() {
		fmt.Printf("😥️ An application OR storage entry is required, not both.\n")
		return false
	}

	if len(m.Spec.Metrics) == 0 {
		fmt.Printf("😥️ One or more metrics are required.\n")
		return false
	}
	if m.Spec.Pods < 1 {
		fmt.Printf("😥️ Pods must be >= 1.")
		return false
	}

	// Validation for application
	if m.HasApplication() {
		if m.Spec.Application.Command == "" {
			fmt.Printf("😥️ Application is missing a command.")
			return false
		}

		if m.Spec.Application.Entrypoint == "" {
			m.Spec.Application.Entrypoint = m.Spec.Application.Command
		}

		// For existing volumes, if it's a claim, a path is required.
		if !m.validateVolumes(m.Spec.Application.Volumes) {
			fmt.Printf("😥️ Application container volumes are not valid\n")
			return false
		}
	}

	// Validate for storage
	if m.HasStorage() && !m.validateVolumes(map[string]Volume{"storage": m.Spec.Storage.Volume}) {
		fmt.Printf("😥️ Storage volumes are not valid\n")
		return false
	}

	// If completions unset, set to parallelism
	if m.Spec.Completions == 0 {
		m.Spec.Completions = m.Spec.Pods
	}

	// A standalone metric by definition runs alone
	if len(m.Spec.Metrics) > 1 && m.IsStandalone() {
		fmt.Printf("😥️ A standalone metric by definition runs on its own\n")
		return false
	}

	return true
}

// validateExistingVolumes ensures secret names vs. volume paths are valid
func (m *MetricSet) validateVolumes(volumes map[string]Volume) bool {

	valid := true
	for key, volume := range volumes {

		// Case 1: it's a secret and we only need that
		if volume.SecretName != "" {
			continue
		}

		// Case 2: it's a config map (and will have items too, but we don't hard require them)
		if volume.ConfigMapName != "" {
			continue
		}

		// Case 3: Hostpath volume (mostly for testing)
		if volume.HostPath != "" {
			continue
		}

		// Case 4: claim desired without path
		if volume.ClaimName == "" && volume.Path != "" {
			fmt.Printf("😥️ Found existing volume %s with path %s that is missing a claim name\n", key, volume.Path)
			valid = false
		}
		// Case 5: reverse of the above
		if volume.ClaimName != "" && volume.Path == "" {
			fmt.Printf("😥️ Found existing volume %s with claimName %s that is missing a path\n", key, volume.ClaimName)
			valid = false
		}

		// Case 6: empty volume needs path
		if volume.EmptyVol && volume.Path == "" {
			fmt.Printf("😥️ Found empty volume %s that is missing a path\n", key)
			valid = false
		}
	}
	return valid
}

//+kubebuilder:object:root=true

// MetricSetList contains a list of MetricSet
type MetricSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetricSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MetricSet{}, &MetricSetList{})
}
