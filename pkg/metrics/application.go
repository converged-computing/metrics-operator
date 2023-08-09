/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// Get ReplicatedJobs intended to run a performance metric for an application
// For this setup, we expect to create a container for each metric
func (m *ApplicationMetricSet) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {
	rjs := []jobset.ReplicatedJob{}
	for _, metric := range m.Metrics() {
		jobs, err := GetApplicationReplicatedJobs(spec, metric, spec.Spec.Application.Volumes, true)
		if err != nil {
			return rjs, err
		}
		rjs = append(rjs, jobs...)
	}
	return rjs, nil
}

// Create a standalone JobSet, one without volumes or application
// This will be definition be a JobSet for only one metric
func GetApplicationReplicatedJobs(
	spec *api.MetricSet,
	metric *Metric,
	volumes map[string]api.Volume,
	shareProcessNamespace bool,
) ([]jobset.ReplicatedJob, error) {

	// Prepare a replicated job
	rjs := []jobset.ReplicatedJob{}

	// We currently don't expose applications to allow custom replicated jobs
	// If we return no replicated jobs, fall back to default
	m := (*metric)

	// This defaults to one replicated job, named "m"
	job, err := GetReplicatedJob(spec, shareProcessNamespace, spec.Spec.Pods, spec.Spec.Completions, "")
	if err != nil {
		return rjs, err
	}

	// Add volumes expecting an application. GetVolumes creates metric entrypoint volumes
	// and adds existing volumes (application) to our set of mounts. We need both
	// for the jobset.
	runnerScripts := GetMetricsKeyToPath([]*Metric{metric})
	job.Template.Spec.Template.Spec.Volumes = GetVolumes(spec, runnerScripts, volumes)

	// Derive the containers for the metric
	// TODO fnish this
	containerSpec := ContainerSpec{
		Image:      m.Image(),
		Command:    []string{"/bin/bash", "/metrics_operator/entrypoint-0.sh"},
		WorkingDir: m.WorkingDir(),
		Name:       m.Name(),
	}

	// This is for the metric and application containers
	// Metric containers have metric entrypoint volumes
	// Application containers have existing volumes
	containers, err := GetContainers(
		spec,
		[]ContainerSpec{containerSpec},
		volumes,

		// Allow ptrace
		true,
	)

	if err != nil {
		logger.Errorf("There was an error getting containers for %s: %s\n", m.Name(), err)
		return rjs, err
	}
	job.Template.Spec.Template.Spec.Containers = containers
	rjs = append(rjs, *job)
	return rjs, nil
}
