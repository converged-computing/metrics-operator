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

	// For metrics that require an application, we need a container and name (for now)
	// +optional
	Application Application `json:"application"`
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

	// Do we need to run more than one completion (pod)?
	//+optional
	Completions int32 `json:"completions"`
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
	return m.Spec.Application.Image != ""
}

// Validate a requested metricset
func (m *MetricSet) Validate() bool {

	if len(m.Spec.Metrics) == 0 {
		fmt.Printf("😥️ One or more metrics are required.\n")
		return false
	}

	// Validation for application
	if m.Spec.Application.Command == "" {
		fmt.Printf("😥️ Application is missing a command.")
		return false
	}
	if m.Spec.Application.Entrypoint == "" {
		m.Spec.Application.Entrypoint = m.Spec.Application.Command
	}

	// Validation for each metric
	//for i, metric := range m.Spec.Metrics {}
	return true
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
