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
	RegistrySet = make(map[string]MetricSet)
)

const (
	ApplicationMetric = "application"
	StorageMetric     = "storage"
	StandaloneMetric  = "standalone"
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
}

// ConsolidateEntrypointScripts from a metric set into one list
func consolidateEntrypointScripts(metrics []*Metric, set *api.MetricSet) []EntrypointScript {
	scripts := []EntrypointScript{}
	for _, metric := range metrics {
		for _, script := range (*metric).EntrypointScripts(set) {
			scripts = append(scripts, script)
		}
	}
	return scripts
}

// Types of Metrics: Storage, Application, and Standalone

// StorageMetricSet defines a MetricSet to measure storage interfaces
type StorageMetricSet struct {
	name        string
	metrics     []*Metric
	metricNames map[string]bool
}

func (m StorageMetricSet) Metrics() []*Metric {
	return m.metrics
}
func (m StorageMetricSet) Type() string {
	return m.name
}
func (m StorageMetricSet) Exists(metric *Metric) bool {
	_, ok := m.metricNames[(*metric).Name()]
	return ok
}
func (m *StorageMetricSet) Add(metric *Metric) {
	if !m.Exists(metric) {
		m.metrics = append(m.metrics, metric)
		m.metricNames[(*metric).Name()] = true
	}
}
func (m *StorageMetricSet) EntrypointScripts(set *api.MetricSet) []EntrypointScript {
	return consolidateEntrypointScripts(m.metrics, set)
}

// ApplicationMetricSet defines a MetricSet to measure application performance
type ApplicationMetricSet struct {
	name        string
	metrics     []*Metric
	metricNames map[string]bool
}

func (m ApplicationMetricSet) Metrics() []*Metric {
	return m.metrics
}
func (m ApplicationMetricSet) Type() string {
	return m.name
}
func (m *ApplicationMetricSet) EntrypointScripts(set *api.MetricSet) []EntrypointScript {
	return consolidateEntrypointScripts(m.metrics, set)
}
func (m ApplicationMetricSet) Exists(metric *Metric) bool {
	_, ok := m.metricNames[(*metric).Name()]
	return ok
}
func (m *ApplicationMetricSet) Add(metric *Metric) {
	if !m.Exists(metric) {
		m.metrics = append(m.metrics, metric)
		m.metricNames[(*metric).Name()] = true
	}
}

type StandaloneMetricSet struct {
	name        string
	metrics     []*Metric
	metricNames map[string]bool
}

func (m StandaloneMetricSet) Metrics() []*Metric {
	return m.metrics
}
func (m StandaloneMetricSet) Type() string {
	return m.name
}
func (m *StandaloneMetricSet) EntrypointScripts(set *api.MetricSet) []EntrypointScript {
	return consolidateEntrypointScripts(m.metrics, set)
}
func (m StandaloneMetricSet) Exists(metric *Metric) bool {
	_, ok := m.metricNames[(*metric).Name()]
	return ok
}
func (m *StandaloneMetricSet) Add(metric *Metric) {
	if !m.Exists(metric) {
		m.metrics = append(m.metrics, metric)
		m.metricNames[(*metric).Name()] = true
	}
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
	RegisterSet(&StorageMetricSet{name: StorageMetric, metricNames: map[string]bool{}})
	RegisterSet(&ApplicationMetricSet{name: ApplicationMetric, metricNames: map[string]bool{}})
	RegisterSet(&StandaloneMetricSet{name: StandaloneMetric, metricNames: map[string]bool{}})
}
