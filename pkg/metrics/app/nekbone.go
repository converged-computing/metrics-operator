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

type Nekbone struct {
	metrics.LauncherWorker

	// Custom Options
	command string
	prefix  string
}

// I think this is a simulation?
func (m Nekbone) Family() string {
	return metrics.SolverFamily
}

func (m Nekbone) Url() string {
	return "https://github.com/Nek5000/Nekbone"
}

// Set custom options / attributes for the metric
func (m *Nekbone) SetOptions(metric *api.Metric) {
	// Set user defined values or fall back to defaults
	m.prefix = "/bin/bash"
	m.command = "mpiexec --hostfile ./hostlist.txt -np 2 ./nekbone"
	m.Workdir = "/root/nekbone-3.0/test/example2"
	m.SetDefaultOptions(metric)
}

// Exported options and list options
func (m Nekbone) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"prefix":  intstr.FromString(m.prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: "app-nekbone",
		Summary:    "A mini-app derived from the Nek5000 CFD code which is a high order, incompressible Navier-Stokes CFD solver based on the spectral element method. The conjugate gradiant solve is compute intense, contains small messages and frequent allreduces.",
		Container:  "ghcr.io/converged-computing/metric-nekbone:latest",
	}
	launcher := metrics.LauncherWorker{BaseMetric: base}
	Nekbone := Nekbone{LauncherWorker: launcher}
	metrics.Register(&Nekbone)
}
