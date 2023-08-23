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
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

var (
	RegistrySet = make(map[string]MetricSet)
)

const (
	// Metric Design Types
	ApplicationMetric = "application"
	StorageMetric     = "storage"
	StandaloneMetric  = "standalone"

	// Metric Family Types (these likely can be changed)
	StorageFamily         = "storage"
	MachineLearningFamily = "machine-learning"
	NetworkFamily         = "network"
	SimulationFamily      = "simulation"
	SolverFamily          = "solver"

	// Generic (more than one type, CPU/io, etc)
	ProxyAppFamily    = "proxyapp"
	PerformanceFamily = "performance"
)

// A MetricSet interface holds one or more Metrics
// and exposes the JobSet
type MetricSet interface {

	// Metric Set Type (string)
	Type() string
	Add(m *Metric)
	Exists(m *Metric) bool
	Metrics() []*Metric
	EntrypointScripts(*api.MetricSet) []EntrypointScript
	ReplicatedJobs(*api.MetricSet) ([]jobset.ReplicatedJob, error)
	HasSoleTenancy() bool
}

// ConsolidateEntrypointScripts from a metric set into one list
func consolidateEntrypointScripts(metrics []*Metric, set *api.MetricSet) []EntrypointScript {
	scripts := []EntrypointScript{}
	for _, metric := range metrics {
		for _, script := range (*metric).EntrypointScripts(set, metric) {
			scripts = append(scripts, script)
		}
	}
	return scripts
}

// BaseMetricSet
type BaseMetricSet struct {
	name        string
	metrics     []*Metric
	metricNames map[string]bool
}

func (m BaseMetricSet) Metrics() []*Metric {
	return m.metrics
}
func (m BaseMetricSet) Type() string {
	return m.name
}
func (m BaseMetricSet) Exists(metric *Metric) bool {
	_, ok := m.metricNames[(*metric).Name()]
	return ok
}

// Determine if any metrics in the set need sole tenancy
// This is defined on the level of the jobset for now
func (m BaseMetricSet) HasSoleTenancy() bool {
	for _, m := range m.metrics {
		if (*m).HasSoleTenancy() {
			return true
		}
	}
	return false
}

func (m *BaseMetricSet) Add(metric *Metric) {
	if !m.Exists(metric) {
		m.metrics = append(m.metrics, metric)
		m.metricNames[(*metric).Name()] = true
	}
}
func (m *BaseMetricSet) EntrypointScripts(set *api.MetricSet) []EntrypointScript {
	return consolidateEntrypointScripts(m.metrics, set)
}

// Types of Metrics: Storage, Application, and Standalone

// StorageMetricSet defines a MetricSet to measure storage interfaces
type StorageMetricSet struct {
	BaseMetricSet
}

// ApplicationMetricSet defines a MetricSet to measure application performance
type ApplicationMetricSet struct {
	BaseMetricSet
}
type StandaloneMetricSet struct {
	BaseMetricSet
}

// Register a new Metric type, adding it to the Registry
func RegisterSet(m MetricSet) {
	name := m.Type()
	if _, ok := RegistrySet[name]; ok {
		log.Fatalf("MetricSet: %s has already been added to the registry", name)
	}
	RegistrySet[name] = m
}

// GetMetric returns the Component specified by name from `Registry`.
func GetMetricSet(name string) (MetricSet, error) {
	if _, ok := RegistrySet[name]; ok {
		m := RegistrySet[name]
		return m, nil
	}
	return nil, fmt.Errorf("%s is not a registered MetricSet type", name)
}

func init() {
	RegisterSet(&StorageMetricSet{BaseMetricSet{name: StorageMetric, metricNames: map[string]bool{}}})
	RegisterSet(&ApplicationMetricSet{BaseMetricSet{name: ApplicationMetric, metricNames: map[string]bool{}}})
	RegisterSet(&StandaloneMetricSet{BaseMetricSet{name: StandaloneMetric, metricNames: map[string]bool{}}})
}
