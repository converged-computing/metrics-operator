/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package metrics

import (
	"fmt"
	"log"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

var (
	Registry = make(map[string]Metric)
)

// A StorageMetric is intended to measure a storage interface
type StorageMetric interface {
	Description() string
	Name() string
	SetOptions(*api.Metric)
	Validate(*api.MetricSet) bool
}

// A Metric defines a generic interface for the operator to interact with
// The functionality of different metric types might vary based on the type
// All metrics return a JobSet of some type (and potentially a replicated job)
type Metric interface {

	// Indicates that the metric requires an application to measure
	RequiresApplication() bool
	RequiresStorage() bool
	Standalone() bool
	Description() string
	Name() string
	SetOptions(*api.Metric)
	Validate(*api.MetricSet) bool

	// Functions for standalone metrics
	ReplicatedJobs(*api.MetricSet, *[]Metric) ([]jobset.ReplicatedJob, error)
	SuccessJobs() []string

	// Container specific attributes!
	// Entrypoint scripts, with
	EntrypointScripts(*api.MetricSet) []EntrypointScript
	WorkingDir() string
	Image() string
}

// GetMetric returns the Component specified by name from `Registry`.
func GetMetric(metric *api.Metric, set *api.MetricSet) (Metric, error) {
	if _, ok := Registry[metric.Name]; ok {
		m := Registry[metric.Name]

		// Validate it's for storage OR application
		if m.RequiresApplication() && m.RequiresStorage() {
			return nil, fmt.Errorf("%s cannot be for an application and storage", metric.Name)
		}
		if !m.Validate(set) {
			return nil, fmt.Errorf("%s is not valid", metric.Name)
		}

		// Set global and custom options on the registry metric from the CRD
		m.SetOptions(metric)
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
