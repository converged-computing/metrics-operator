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

package v1alpha2

import (
	"fmt"

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

	// Pod spec for the application, standalone, or storage metrics
	//+optional
	Pod Pod `json:"pod"`

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

	//+optional
	Privileged bool `json:"privileged"`

	//+optional
	AllowPtrace bool `json:"allowPtrace"`

	//+optional
	AllowAdmin bool `json:"allowAdmin"`
}

// A Metric addon is an interface that exposes extra volumes for a metric. Examples include:
// A storage volume to be mounted on one or more of the replicated jobs
// A single application container.
type MetricAddon struct {
	Name string `json:"name"`

	// Metric Addon Options
	// +optional
	Options map[string]intstr.IntOrString `json:"options"`

	// Addon List Options
	// +optional
	ListOptions map[string][]intstr.IntOrString `json:"listOptions"`

	// Addon Map Options
	// +optional
	MapOptions map[string]map[string]intstr.IntOrString `json:"mapOptions"`
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

// The difference between benchmark and metric is subtle.
// A metric is more a measurment, and the benchmark is the comparison value.
// I don't have strong opinions but I think we are doing more measurment
// not necessarily with benchmarks

// A metric is basically a container. It minimally provides:
// In the simplest case, a sidecar container (e.g., service or similar)
// Optionally: Possibly additional volumes that can be mounted in
// With a shared process namespace, ability to monitor

type Metric struct {
	Name string `json:"name"`

	// Metric Options
	// Metric specific options
	// +optional
	Options map[string]intstr.IntOrString `json:"options"`

	// Use a custom container image (advanced users only)
	// +optional
	Image string `json:"image,omitempty"`

	// A Metric addon can be storage (volume) or an application,
	// It's an additional entity that can customize a replicated job,
	// either adding assets / features or entire containers to the pod
	//+optional
	Addons []MetricAddon `json:"addons"`

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

// Validate a requested metricset
func (m *MetricSet) Validate() bool {

	if len(m.Spec.Metrics) == 0 {
		fmt.Printf("😥️ One or more metrics are required.\n")
		return false
	}
	if m.Spec.Pods < 1 {
		fmt.Printf("😥️ Pods must be >= 1.")
		return false
	}
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
