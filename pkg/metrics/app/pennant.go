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

type Pennant struct {
	metrics.LauncherWorker

	// Custom Options
	command string
	prefix  string
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
	// Set user defined values or fall back to defaults
	m.prefix = "mpirun --hostfile ./hostlist.txt"
	m.command = "pennant /opt/pennant/test/sedovsmall/sedovsmall.pnt"
	m.Workdir = "/opt/pennant/test"
	m.SetDefaultOptions(metric)
}

// Exported options and list options
func (m Pennant) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"prefix":  intstr.FromString(m.prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: "app-pennant",
		Summary:    "Unstructured mesh hydrodynamics for advanced architectures ",
		Container:  "ghcr.io/converged-computing/metric-pennant:latest",
	}
	launcher := metrics.LauncherWorker{
		BaseMetric:     base,
		WorkerScript:   "/metrics_operator/pennant-worker.sh",
		LauncherScript: "/metrics_operator/pennant-launcher.sh",
	}
	Pennant := Pennant{LauncherWorker: launcher}
	metrics.Register(&Pennant)
}
