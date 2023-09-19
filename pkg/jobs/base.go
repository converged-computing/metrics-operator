/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package jobs

import (
	"github.com/converged-computing/metrics-operator/pkg/addons"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// BaseMetric provides shared attributes across Metric types
type BaseMetric struct {

	// A metric can have one or more addons
	Addons map[string]*addons.Addon

}

// RegisterAddon adds an addon to the set, assuming it's already validated
func (m BaseMetric) RegisterAddon(addon *addons.Addon) {
	a := (*addon)
	m.Addons[a.Name()] = addon
}


// Get the entrypoint name
func (m BaseMetric) GetEntrypointName() string {
	return fmt.Sprintf("entrypoint-%d", m.Identifier)
}

// Get the entrypoint script
func (m BaseMetric) GetEntrypointScript() string {
	return fmt.Sprintf("%s.sh", m.GetEntrypointName())
}

// Add registered addons to a replicated job
// Since we have the metric (and also know containers here) we do that
func (m BaseMetric) AddAddons(rj *jobset.ReplicatedJob) *jobset.ReplicatedJob {

	// VolumeMounts for containers, which we only need to know names for
	mounts := []corev1.KeyToPath{
		{
			Key:  m.GetEntrypointName(),
			Path: m.GetEntrypointScript(),
			Mode: &makeExecutable,
		},
	}

	// For each addon, do custom logic depending on the type
	// These are the main set of volumes, containers we are going to add
	volumes := []corev1.Volume{}
	for _, a := range m.Addons {

		// Add volumes that addons provide
		volumes = append(volumes, a.GetVolumes())

		// Add containers that addons provide, but as specs

	}

	return rj
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

	if err != nil {
		return jobsets, err
	}



// Addons returns a list of addons, removing them from the key value lookup
func (m BaseMetric) GetAddons() []*addons.Addon {
	addons := []*addons.Addon{}
	for _, addon := range m.Addons {
		addons = append(addons, addon)
	}
	return addons
}

// TODO move shared logic into this class..
