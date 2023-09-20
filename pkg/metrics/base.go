/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"github.com/converged-computing/metrics-operator/pkg/addons"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// BaseMetric provides shared attributes across Metric types
type BaseMetric struct {
	Identifier string
	Summary    string
	Container  string
	Workdir    string

	ResourceSpec  *api.ContainerResources
	AttributeSpec *api.ContainerSpec

	// A metric can have one or more addons
	Addons map[string]*addons.Addon
}

// RegisterAddon adds an addon to the set, assuming it's already validated
func (m BaseMetric) RegisterAddon(addon *addons.Addon) {
	a := (*addon)
	m.Addons[a.Name()] = addon
}

// Name returns the metric name
func (m BaseMetric) Name() string {
	return m.Identifier
}

// Description returns the metric description
func (m BaseMetric) Description() string {
	return m.Summary
}

// Container
func (m BaseMetric) Image() string {
	return m.Container
}

// WorkingDir does not matter
func (m BaseMetric) WorkingDir() string {
	return m.Workdir
}

// Return container resources for the metric container
func (m BaseMetric) Resources() *api.ContainerResources {
	return m.ResourceSpec
}
func (m BaseMetric) Attributes() *api.ContainerSpec {
	return m.AttributeSpec
}

// Validation
func (m BaseMetric) Validate(set *api.MetricSet) bool {
	return true
}

func (m BaseMetric) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Jobs required for success condition (n is the netmark run)
func (m BaseMetric) SuccessJobs() []string {
	return []string{}
}

func (m BaseMetric) ReplicatedJobs(set *api.MetricSet) ([]*jobset.ReplicatedJob, error) {
	return []*jobset.ReplicatedJob{}, nil
}

// Add registered addons to replicated jobs
func (m BaseMetric) AddAddons(
	spec *api.MetricSet,
	rjs []*jobset.ReplicatedJob,

	// These container specs include all replicated jobs
	containerSpecs []*specs.ContainerSpec,
) error {

	// VolumeMounts can be generated from container specs
	// For each addon, do custom logic depending on the type
	// These are the main set of volumes, containers we are going to add
	volumes := []specs.VolumeSpec{}

	// These are addon container specs
	addonContainers := []specs.ContainerSpec{}
	for _, addon := range m.Addons {

		a := (*addon)

		// Assemble volume specs that addons provide
		// These are assumed to exist, and we create mounts for them only
		volumes = append(volumes, a.AssembleVolumes()...)

		// Assemble containers that addons provide, also as specs
		addonContainers = append(addonContainers, a.AssembleContainers()...)

		// Allow the addons to customize the container entrypoints, specific to the job name
		// It's important that this set does not include other addon container specs
		a.CustomizeEntrypoints(containerSpecs, rjs)
	}

	// Add containers to the replicated job (filtered based on matching names)
	containers := addonContainers
	for _, cs := range containerSpecs {
		containers = append(containers, (*cs))
	}

	// Generate actual containers and volumes for each replicated job
	for _, rj := range rjs {

		// We also include the addon volumes, which generally need mount points
		rjContainers, err := getReplicatedJobContainers(spec, rj, containers, volumes)
		if err != nil {
			return err
		}
		rj.Template.Spec.Template.Spec.Containers = rjContainers

		// And volumes!
		// containerSpecs are used to generate our metric entrypoint volumes
		// volumes indicate existing volumes
		rj.Template.Spec.Template.Spec.Volumes = getReplicatedJobVolumes(spec, containerSpecs, volumes)
	}
	return nil
}

// Addons returns a list of addons, removing them from the key value lookup
func (m BaseMetric) GetAddons() []*addons.Addon {
	addons := []*addons.Addon{}
	for _, addon := range m.Addons {
		addons = append(addons, addon)
	}
	return addons
}
