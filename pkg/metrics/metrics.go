/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"fmt"
	"log"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

var (
	Registry = make(map[string]Metric)
)

// A general metric produces a JobSet with one or more replicated Jobs
type Metric interface {
	Name() string
	Description() string
	Url() string

	SetOptions(*api.Metric)
	Validate(*api.MetricSet) bool

	// Attributes to expose for containers
	WorkingDir() string
	Image() string
	Family() string

	// One or more replicated jobs to populate a JobSet
	ReplicatedJobs(*api.MetricSet) ([]jobset.ReplicatedJob, error)
	SuccessJobs() []string
	Resources() *api.ContainerResources
	Attributes() *api.ContainerSpec

	// Metric type to know how to add to MetricSet
	Type() string

	// Exportable attributes
	Options() map[string]intstr.IntOrString
	ListOptions() map[string][]intstr.IntOrString

	// EntrypointScripts are required to generate ConfigMaps
	EntrypointScripts(*api.MetricSet, *Metric) []EntrypointScript
}

// GetMetric returns the Component specified by name from `Registry`.
func GetMetric(metric *api.Metric, set *api.MetricSet) (Metric, error) {
	if _, ok := Registry[metric.Name]; ok {
		m := Registry[metric.Name]

		// Ensure the type is one acceptable
		if !(m.Type() == ApplicationMetric || m.Type() == StorageMetric || m.Type() == StandaloneMetric) {
			return nil, fmt.Errorf("%s is not a valid type", metric.Name)
		}

		// Set global and custom options on the registry metric from the CRD
		m.SetOptions(metric)

		// After options are set, final validation
		if !m.Validate(set) {
			return nil, fmt.Errorf("%s did not validate", metric.Name)
		}

		return m, nil
	}
	return nil, fmt.Errorf("%s is not a registered Metric type", metric.Name)
}

// Register a new Metric type, adding it to the Registry
func Register(m Metric) {
	name := m.Name()
	if _, ok := Registry[name]; ok {
		log.Fatalf("Metric: %s has already been added to the registry", name)
	}
	Registry[name] = m
}
