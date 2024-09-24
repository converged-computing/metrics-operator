/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"fmt"
	"log"
	"reflect"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	addons "github.com/converged-computing/metrics-operator/pkg/addons"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

var (
	Registry = map[string]Metric{}
)

// A general metric is a container added to a JobSet
type Metric interface {

	// Metadata
	Name() string
	Description() string
	Family() string
	Url() string

	// Container attributes
	Image() string
	SetContainer(string)

	// Options and exportable attributes
	SetOptions(*api.Metric)
	Options() map[string]intstr.IntOrString
	ListOptions() map[string][]intstr.IntOrString

	// Validation and append addons
	Validate(*api.MetricSet) bool
	RegisterAddon(*addons.Addon)
	AddAddons(*api.MetricSet, []*jobset.ReplicatedJob, []*specs.ContainerSpec) ([]*specs.ContainerSpec, error)
	GetAddons() []*addons.Addon

	// Attributes for JobSet, etc.
	HasSoleTenancy() bool
	ReplicatedJobs(*api.MetricSet) ([]*jobset.ReplicatedJob, error)
	SuccessJobs() []string
	Resources() *api.ContainerResources
	Attributes() *api.ContainerSpec

	// Prepare Containers. These are used to generate configmaps,
	// and populate the respective replicated jobs with containers!
	PrepareContainers(*api.MetricSet, *Metric) []*specs.ContainerSpec
}

// GetMetric returns a metric, if it is known to the metrics operator
// We also confirm that the addon exists, validate, and instantiate it.
func GetMetric(metric *api.Metric, set *api.MetricSet) (Metric, error) {

	if _, ok := Registry[metric.Name]; ok {

		// Start with the empty template, and create a copy
		// This is important so we don't preserve state to the actaul interface
		template := Registry[metric.Name]
		templateType := reflect.ValueOf(template)
		if templateType.Kind() == reflect.Ptr {
			templateType = reflect.Indirect(templateType)
		}
		m := reflect.New(templateType.Type()).Interface().(Metric)

		// Set global and custom options on the registry metric from the CRD
		m.SetOptions(metric)

		// If the metric has a custom container, set here
		if metric.Image != "" {
			m.SetContainer(metric.Image)
		}

		// Register addons, meaning adding the spec but not instantiating yet (or should we?)
		for _, a := range metric.Addons {

			logger.Infof("Attempting to add addon %s", a.Name)
			addon, err := addons.GetAddon(&a, set)
			if err != nil {
				return nil, fmt.Errorf("addon %s for metric %s did not validate", a.Name, metric.Name)
			}
			logger.Infof("Registering addon %s", a.Name)
			m.RegisterAddon(&addon)
		}

		// After options are set, final validation
		if !m.Validate(set) {
			return nil, fmt.Errorf("%s did not validate", metric.Name)
		}
		return m, nil
	}
	return nil, fmt.Errorf("%s is not a registered Metric type", metric.Name)
}

// Register a new Metric type, adding it to the Registry
func Register(m Metric) {
	name := m.Name()
	if _, ok := Registry[name]; ok {
		log.Fatalf("Metric: %s has already been added to the registry\n", m)
	}
	Registry[name] = m
}
