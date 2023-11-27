/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package application

import (
	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/util/intstr"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

const (
	cabanapicIdentifier = "app-cabanapic"
	cabanapicSummary    = "structured PIC (particle in cell) proxy app"
	cabanapicContainer  = "ghcr.io/converged-computing/metric-cabanapic:latest"
)

type CabanaPIC struct {
	metrics.LauncherWorker
}

// I think this is a simulation?
func (m CabanaPIC) Family() string {
	return metrics.SimulationFamily
}

func (m CabanaPIC) Validate(set *api.MetricSet) bool {
	return true
}

func (m CabanaPIC) Url() string {
	return "https://github.com/ECP-copa/CabanaPIC"
}

// Set custom options / attributes for the metric
func (m *CabanaPIC) SetOptions(metric *api.Metric) {

	m.Identifier = cabanapicIdentifier
	m.Summary = cabanapicSummary
	m.Container = cabanapicContainer

	// Set user defined values or fall back to defaults
	m.Prefix = "/bin/bash"
	m.Command = "cbnpic"
	m.Workdir = "/opt/cabanaPIC/build"
	m.SetDefaultOptions(metric)
}

// Exported options and list options
func (m CabanaPIC) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.Command),
		"prefix":  intstr.FromString(m.Prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: cabanapicIdentifier,
		Summary:    cabanapicSummary,
		Container:  cabanapicContainer,
	}
	launcher := metrics.LauncherWorker{BaseMetric: base}
	CabanaPIC := CabanaPIC{LauncherWorker: launcher}
	metrics.Register(&CabanaPIC)
}
