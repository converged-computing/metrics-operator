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
	pennantIdentifier = "app-pennant"
	pennantSummary    = "Unstructured mesh hydrodynamics for advanced architectures "
	pennantContainer  = "ghcr.io/converged-computing/metric-pennant:latest"
)

type Pennant struct {
	metrics.LauncherWorker
}

// I think this is a simulation?
func (m Pennant) Family() string {
	return metrics.SimulationFamily
}

func (m Pennant) Url() string {
	return "https://github.com/LLNL/pennant"
}

// Set custom options / attributes for the metric
func (m *Pennant) SetOptions(metric *api.Metric) {

	m.Container = pennantContainer
	m.Identifier = pennantIdentifier
	m.Summary = pennantSummary

	// Set user defined values or fall back to defaults
	m.Prefix = "mpirun --hostfile ./hostlist.txt"
	m.Command = "pennant /opt/pennant/test/sedovsmall/sedovsmall.pnt"
	m.Workdir = "/opt/pennant/test"
	m.SetDefaultOptions(metric)
}

// Exported options and list options
func (m Pennant) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.Command),
		"prefix":  intstr.FromString(m.Prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: pennantIdentifier,
		Summary:    pennantSummary,
		Container:  pennantContainer,
	}
	launcher := metrics.LauncherWorker{BaseMetric: base}
	Pennant := Pennant{LauncherWorker: launcher}
	metrics.Register(&Pennant)
}
