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
	laghosIdentifier = "app-laghos"
	laghosSummary    = "LAGrangian High-Order Solver"
	laghosContainer  = "ghcr.io/converged-computing/metric-laghos:latest"
)

type Laghos struct {
	metrics.LauncherWorker
}

// I think this is a simulation?
func (m Laghos) Family() string {
	return metrics.SolverFamily
}

func (m Laghos) Url() string {
	return "https://github.com/CEED/Laghos"
}

// Set custom options / attributes for the metric
func (m *Laghos) SetOptions(metric *api.Metric) {

	m.Identifier = laghosIdentifier
	m.Summary = laghosSummary
	m.Container = laghosSummary

	// Set user defined values or fall back to defaults
	m.Prefix = "/bin/bash"
	m.Command = "mpirun -np 4 --hostfile ./hostlist.txt ./laghos"
	m.Workdir = "/workflow/laghos"
	m.SetDefaultOptions(metric)
}

// Exported options and list options
func (m Laghos) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.Command),
		"prefix":  intstr.FromString(m.Prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: laghosIdentifier,
		Summary:    laghosSummary,
		Container:  laghosContainer,
	}
	launcher := metrics.LauncherWorker{BaseMetric: base}
	Laghos := Laghos{LauncherWorker: launcher}
	metrics.Register(&Laghos)
}
