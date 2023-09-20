/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package application

import (
	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

type Laghos struct {
	metrics.LauncherWorker

	// Custom Options
	command string
	prefix  string
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
	// Set user defined values or fall back to defaults
	m.prefix = "/bin/bash"
	m.command = "mpirun -np 4 --hostfile ./hostlist.txt ./laghos"
	m.Workdir = "/workflow/laghos"
	m.SetDefaultOptions(metric)
}

// Exported options and list options
func (m Laghos) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"prefix":  intstr.FromString(m.prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func init() {
	launcher := metrics.LauncherWorker{
		Identifier: "app-laghos",
		Summary:    "LAGrangian High-Order Solver",
		Container:  "ghcr.io/converged-computing/metric-laghos:latest",
	}

	Laghos := Laghos{LauncherWorker: launcher}
	metrics.Register(&Laghos)
}
