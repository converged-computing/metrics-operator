/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
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

	// A custom container can be used to replace the application
	// (typically advanced users only)
	CustomContainer string
	ResourceSpec    *api.ContainerResources
	AttributeSpec   *api.ContainerSpec

	// If we ask for sole tenancy, we assign 1 pod / hostname
	SoleTenancy bool

	// A metric can have one or more addons
	Addons map[string]*addons.Addon
}

// RegisterAddon adds an addon to the set, assuming it's already validated
func (m *BaseMetric) RegisterAddon(addon *addons.Addon) {
	a := (*addon)
	m.InitAddons()
	logger.Infof("🟧️ Registering addon %s", a.Name())
	m.Addons[a.Name()] = addon
}

// InitAddons ensures we don't have an empty map
func (m *BaseMetric) InitAddons() {
	if m.Addons == nil {
		logger.Infof("🟧️ Resetting addons - they are unset.")
		m.Addons = map[string]*addons.Addon{}
	}
}

// Name returns the metric name
func (m BaseMetric) Name() string {
	return m.Identifier
}

// Set a custom container
func (m *BaseMetric) SetContainer(container string) {
	m.Container = container
}

// Description returns the metric description
func (m BaseMetric) Description() string {
	return m.Summary
}

// Container
func (m *BaseMetric) Image() string {
	return m.Container
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
	if m.Identifier == "" {
		logger.Errorf("Metric %s is missing an identifier.\n", m)
		return false
	}
	return true
}

func (m BaseMetric) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Jobs required for success condition (n is the netmark run)
func (m BaseMetric) SuccessJobs() []string {
	return []string{}
}

func (m BaseMetric) HasSoleTenancy() bool {
	return m.SoleTenancy
}

// Default replicated jobs will generate for N pods, with no shared process namespace (e.g., storage)
func (m *BaseMetric) ReplicatedJobs(spec *api.MetricSet) ([]*jobset.ReplicatedJob, error) {

	js := []*jobset.ReplicatedJob{}

	// An empty jobname will default to "m" the ReplicatedJobName provided by the operator
	rj, err := AssembleReplicatedJob(spec, false, spec.Spec.Pods, spec.Spec.Pods, "", m.SoleTenancy)
	if err != nil {
		return js, err
	}
	js = []*jobset.ReplicatedJob{rj}
	return js, nil
}

// SetDefaultOptions that are shared (possibly)
// TODO this doesn't do anything given an interface and needs
// a different placement.
func (m BaseMetric) SetDefaultOptions(metric *api.Metric) {
	st, ok := metric.Options["soleTenancy"]
	if ok && st.StrVal == "false" || st.StrVal == "no" {
		m.SoleTenancy = false
	}
	if ok && st.StrVal == "true" || st.StrVal == "yes" {
		m.SoleTenancy = true
	}
}

// Add registered addons to replicated jobs
// Container specs returned are assumed to be config maps that need to be written
func (m BaseMetric) AddAddons(
	spec *api.MetricSet,
	rjs []*jobset.ReplicatedJob,

	// These container specs include all replicated jobs
	containerSpecs []*specs.ContainerSpec,
) ([]*specs.ContainerSpec, error) {

	// Ensure we have created the map!
	m.InitAddons()

	// VolumeMounts can be generated from container specs
	// For each addon, do custom logic depending on the type
	// These are the main set of volumes, containers we are going to add
	// Organize volumes by unique name
	volumes := []specs.VolumeSpec{}

	// These are addon container specs
	addonContainers := []specs.ContainerSpec{}

	// These are container specs that need to be written to configmaps
	cms := []*specs.ContainerSpec{}

	for _, addon := range m.Addons {
		a := (*addon)

		logger.Infof("🟧️ Including Addon", a.Name())
		volumes = append(volumes, a.AssembleVolumes()...)

		// Assemble containers that addons provide, also as specs
		assembleContainers := a.AssembleContainers()

		// Sidecar containers
		for _, assembleContainer := range assembleContainers {

			// Any container specs that need to be created later as config maps are kept in cms
			if assembleContainer.NeedsWrite {
				cms = append(cms, &assembleContainer)
			}
			addonContainers = append(addonContainers, assembleContainer)
		}

		// Allow the addons to customize the container entrypoints, specific to the job name
		// It's important that this set does not include other addon container specs
		a.CustomizeEntrypoints(containerSpecs, rjs)
	}

	// There is a bug here showing lots of nil but I don't know why
	logger.Infof("🟧️ Volumes that are going to be added %s\n", volumes)

	// Add containers to the replicated job (filtered based on matching names)
	containers := addonContainers
	for _, cs := range containerSpecs {
		containers = append(containers, (*cs))
	}

	// Generate actual containers and volumes for each replicated job
	for _, rj := range rjs {

		// We also include the addon volumes, which generally need mount points
		rjContainers, initContainers, err := getReplicatedJobContainers(spec, rj, containers, volumes)
		if err != nil {
			return cms, err
		}
		rj.Template.Spec.Template.Spec.Containers = rjContainers
		rj.Template.Spec.Template.Spec.InitContainers = initContainers

		// And volumes!
		// containerSpecs are used to generate our metric entrypoint volumes
		// volumes indicate existing volumes
		rj.Template.Spec.Template.Spec.Volumes = getReplicatedJobVolumes(spec, containerSpecs, volumes)
	}
	return cms, nil
}

// Addons returns a list of addons, removing them from the key value lookup
func (m BaseMetric) GetAddons() []*addons.Addon {
	addons := []*addons.Addon{}
	for _, addon := range m.Addons {
		addons = append(addons, addon)
	}
	return addons
}
