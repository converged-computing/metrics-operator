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
}

// A benchmark is details for one metric
// The difference between benchmark and metric is subtle.
// A metric is more a measurment, and the benchmark is the comparison value.
// I don't have strong opinions but I like exposing "metric/metrics" to the user
type Metric struct {
	Name string `json:"name"`
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

// Validate a requested metric
func (m *MetricSet) Validate() bool {

	if len(m.Spec.Metrics) == 0 {
		fmt.Printf("😥️ One or more metrics are required.\n")
		return false
	}

	// Validation for each metric
	for i, metric := range m.Spec.Metrics {
		if metric.Name == "" {
			fmt.Printf("😥️ Metric in index %d is missing a name.\n", i)
			return false
		}
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
