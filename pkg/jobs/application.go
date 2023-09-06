/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package jobs

import (
	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// These are common templates for application metrics

// SingleApplication is a Metric base for a simple application metric
// be accessible by other packages (and not conflict with function names)
type SingleApplication struct {
	Identifier    string
	Summary       string
	Container     string
	Workdir       string
	ResourceSpec  *api.ContainerResources
	AttributeSpec *api.ContainerSpec
}

// Name returns the metric name
func (m SingleApplication) Name() string {
	return m.Identifier
}

func (m SingleApplication) GetVolumes() map[string]api.Volume {
	return map[string]api.Volume{}
}

func (m SingleApplication) HasSoleTenancy() bool {
	return false
}

// Description returns the metric description
func (m SingleApplication) Description() string {
	return m.Summary
}

// Default SingleApplication is generic performance family
func (m SingleApplication) Family() string {
	return metrics.PerformanceFamily
}

// Return container resources for the metric container
func (m SingleApplication) Resources() *api.ContainerResources {
	return m.ResourceSpec
}
func (m SingleApplication) Attributes() *api.ContainerSpec {
	return m.AttributeSpec
}

// Validation
func (m SingleApplication) Validate(spec *api.MetricSet) bool {
	return true
}

// If we have an application container, return that plus custom logic
// custom: is any custom code (environment, waiting, etc.)
// prefix: is a wrapper to the actual entrypoint command
func (m SingleApplication) ApplicationEntrypoint(
	spec *api.MetricSet,
	custom string,
	prefix string,
	suffix string,
) metrics.EntrypointScript {

	template := `#!/bin/bash`

	// If we have custom logic (environment, sleep, etc) add it here
	if custom != "" {
		template = template + "\n" + custom
	}
	// Add the actual entrypoint
	template = template + "\n" + prefix + " " + spec.Spec.Application.Entrypoint + "\n" + suffix

	// If we do, add the custom logic first
	return metrics.EntrypointScript{
		Script: template,
		Path:   metrics.DefaultApplicationEntrypoint,
		Name:   metrics.DefaultApplicationName,
	}
}

// Container variables
func (m SingleApplication) Image() string {
	return m.Container
}
func (m SingleApplication) WorkingDir() string {
	return m.Workdir
}

func (m SingleApplication) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {
	return []jobset.ReplicatedJob{}, nil
}

func (m SingleApplication) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

func (m SingleApplication) SuccessJobs() []string {
	return []string{}
}

func (m SingleApplication) Type() string {
	return metrics.ApplicationMetric
}
