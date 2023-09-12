/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package jobs

import (
	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// These are common templates for storage apps.
// They define the interface of a Metric.

type StorageGeneric struct {
	Identifier string
	Summary    string
	Container  string
	Workdir    string

	ResourceSpec  *api.ContainerResources
	AttributeSpec *api.ContainerSpec
}

// Name returns the metric name
func (m StorageGeneric) Name() string {
	return m.Identifier
}

// Family returns the storage family
func (m StorageGeneric) Family() string {
	return metrics.StorageFamily
}

func (m StorageGeneric) GetVolumes() map[string]api.Volume {
	return map[string]api.Volume{}
}

// Description returns the metric description
func (m StorageGeneric) Description() string {
	return m.Summary
}

// By default assume storage does not have sole tenancy
func (m StorageGeneric) HasSoleTenancy() bool {
	return false
}

// Container
func (m StorageGeneric) Image() string {
	return m.Container
}

// WorkingDir does not matter
func (m StorageGeneric) WorkingDir() string {
	return m.Workdir
}

// Return container resources for the metric container
func (m StorageGeneric) Resources() *api.ContainerResources {
	return m.ResourceSpec
}
func (m StorageGeneric) Attributes() *api.ContainerSpec {
	return m.AttributeSpec
}

// Validation
func (m StorageGeneric) Validate(set *api.MetricSet) bool {
	return true
}

func (m StorageGeneric) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Jobs required for success condition (n is the netmark run)
func (m StorageGeneric) SuccessJobs() []string {
	return []string{}
}

func (m StorageGeneric) Type() string {
	return metrics.StorageMetric
}

func (m StorageGeneric) ReplicatedJobs(set *api.MetricSet) ([]jobset.ReplicatedJob, error) {
	return []jobset.ReplicatedJob{}, nil
}
