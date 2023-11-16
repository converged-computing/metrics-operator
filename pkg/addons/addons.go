/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package addons

import (
	"fmt"
	"log"
	"reflect"

	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// An addon can support adding volumes, containers, or otherwise customizing the jobset.

var (
	Registry               = make(map[string]Addon)
	AddonFamilyPerformance = "performance"
	AddonFamilyVolume      = "volume"
	AddonFamilyApplication = "application"
	AddonFamilyWorkload    = "workload"
)

// A general metric is a container added to a JobSet
type Addon interface {

	// Metadata
	Name() string
	Family() string
	Description() string

	// Options and exportable attributes
	SetOptions(*api.MetricAddon, *api.MetricSet)
	Options() map[string]intstr.IntOrString
	ListOptions() map[string][]intstr.IntOrString
	MapOptions() map[string]map[string]intstr.IntOrString

	// What addons can control:
	AssembleVolumes() []specs.VolumeSpec
	AssembleContainers() []specs.ContainerSpec

	CustomizeEntrypoints([]*specs.ContainerSpec, []*jobset.ReplicatedJob)

	// Instead of exposing individual pieces (volumes, settings, etc)
	// We simply allow it to modify the job
	// Attributes for JobSet, etc.
	Validate() bool
}

// Shared based of metadata and functions
type AddonBase struct {
	Identifier string
	Url        string
	Summary    string
	Family     string

	options     map[string]intstr.IntOrString
	listOptions map[string][]intstr.IntOrString
	mapOptions  map[string]map[string]intstr.IntOrString
}

func (b *AddonBase) SetOptions(addon *api.MetricAddon, metric *api.MetricSet)             {}
func (b *AddonBase) CustomizeEntrypoints([]*specs.ContainerSpec, []*jobset.ReplicatedJob) {}

func (b *AddonBase) Validate() bool {
	return true
}
func (b *AddonBase) AssembleContainers() []specs.ContainerSpec {
	return []specs.ContainerSpec{}
}

// Assemble Volumes (for now) just generates one
func (b *AddonBase) AssembleVolumes() []specs.VolumeSpec {
	return []specs.VolumeSpec{}
}

func (b *AddonBase) Description() string {
	return b.Summary
}
func (b *AddonBase) Name() string {
	return b.Identifier
}
func (b *AddonBase) Options() map[string]intstr.IntOrString {
	return b.options
}
func (b *AddonBase) ListOptions() map[string][]intstr.IntOrString {
	return b.listOptions
}
func (b *AddonBase) MapOptions() map[string]map[string]intstr.IntOrString {
	return b.mapOptions
}

// GetAddon looks up and validates an addon
func GetAddon(a *api.MetricAddon, set *api.MetricSet) (Addon, error) {

	// We don't want to change the addon interface/struct itself
	template, ok := Registry[a.Name]
	if !ok {
		return nil, fmt.Errorf("%s is not a known addon", a.Name)
	}
	templateType := reflect.ValueOf(template)
	if templateType.Kind() == reflect.Ptr {
		templateType = reflect.Indirect(templateType)
	}
	addon := reflect.New(templateType.Type()).Interface().(Addon)

	// Set options before validation
	addon.SetOptions(a, set)

	// Validate the addon
	if !addon.Validate() {
		return nil, fmt.Errorf("Addon %s did not validate", a.Name)
	}
	return addon, nil
}

// Register a new addon!
func Register(a Addon) {
	name := a.Name()
	if _, ok := Registry[name]; ok {
		log.Fatalf("Addon: %s has already been added to the addon registry", name)
	}
	Registry[name] = a
}
