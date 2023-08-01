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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MetricSpec defines the desired state of Metric
type MetricSetSpec struct {

	// The name of the metric (that will be associated with a flavor like storage)
	Metrics []Metric `json:"metrics"`

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

	// For metrics that require an application, we need a container and name (for now)
	// +optional
	Application Application `json:"application"`

	// Number of job completions (e.g., pods)
	// +kubebuilder:default=1
	// +default=1
	// +optional
	Completions int32 `json:"completions"`
}

// Storage that will be monitored
type Storage struct {

	// Volume type to test
	Volume Volume `json:"volume"`
}

// Application that will be monitored
type Application struct {
	Image string `json:"image"`

	// command to execute and monitor
	Command string `json:"command"`

	// Entrypoint of container, if different from command
	//+optional
	Entrypoint string `json:"entrypoint"`

	// A pull secret for the application container
	//+optional
	PullSecret string `json:"pullSecret"`

	// Existing Volumes for the application
	// +optional
	Volumes map[string]Volume `json:"volumes"`
}

// A Volume should correspond with an existing volume, either:
// config map, secret, or claim name. This will be added soon.
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

	// Global attributes shared by all metrics
	// Sampling rate in seconds. Defaults to every 10 seconds
	// +kubebuilder:default=10
	// +default=10
	// +optional
	Rate int32 `json:"rate"`

	// Completions
	// Number of completions to do, more relevant for service type applications
	// that run forever, or a storage metric. If not set (0) then don't set a limit
	// +optional
	Completions int32 `json:"completions"`

	// Custom attributes specific to metrics
	// +optional
	Attributes map[string]string `json:"attributes"`
}

// Get pod labels for a metric set
func (m *MetricSet) GetPodLabels() map[string]string {
	podLabels := map[string]string{}
	podLabels["cluster-name"] = m.Name
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

// Determine if an application is present
func (m *MetricSet) HasApplication() bool {
	return !reflect.DeepEqual(m.Spec.Application, Application{})
}

func (m *MetricSet) HasStorage() bool {
	return !reflect.DeepEqual(m.Spec.Storage, Storage{})
}

// Validate a requested metricset
func (m *MetricSet) Validate() bool {

	// An application or storage setup is required
	if !m.HasApplication() && !m.HasStorage() {
		fmt.Printf("😥️ An application OR storage entry is required.\n")
		return false
	}

	// We don't currently support running both at once
	if m.HasApplication() && m.HasStorage() {
		fmt.Printf("😥️ An application OR storage entry is required, not both.\n")
		return false
	}

	if len(m.Spec.Metrics) == 0 {
		fmt.Printf("😥️ One or more metrics are required.\n")
		return false
	}

	// Storage or an application can have completions (replicas)
	if m.Spec.Completions < 1 {
		fmt.Printf("😥️ Completions must be >= 1.")
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

	// Validation for each metric
	for _, metric := range m.Spec.Metrics {
		if metric.Rate <= 0 {
			metric.Rate = 10
		}
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
