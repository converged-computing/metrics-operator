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

// Get ReplicatedJobs intended to run storage
// For this setup, we expect to create a container for each storage metric
// And then add the volume bind to it
func (m *StorageMetricSet) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {

	// Prepare replicated jobs list to return
	rjs := []jobset.ReplicatedJob{}

	// Storage metrics do not need to share the process namespace
	// The jobname empty string will use the default
	job, err := GetReplicatedJob(spec, false, spec.Spec.Pods, spec.Spec.Completions, "")
	if err != nil {
		return rjs, err
	}

	// Only add storage volume if we have it! Not all storage interfaces require
	// A Kubernetes abstraction, some are created via a command.
	volumes := map[string]api.Volume{}
	if spec.HasStorageVolume() {
		// Add volumes expecting an application.
		// A storage app is required to have a volume
		volumes = map[string]api.Volume{"storage": spec.Spec.Storage.Volume}
	}

	// Derive running scripts from the metric
	runnerScripts := GetMetricsKeyToPath(m.Metrics())
	job.Template.Spec.Template.Spec.Volumes = GetVolumes(spec, runnerScripts, volumes)

	// Derive the containers, one per metric
	// This will also include mounts for volumes
	containers, err := getContainers(spec, m.Metrics(), volumes)
	if err != nil {
		return rjs, err
	}
	job.Template.Spec.Template.Spec.Containers = containers
	rjs = append(rjs, *job)
	return rjs, nil

}
