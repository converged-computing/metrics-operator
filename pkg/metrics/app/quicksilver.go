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

type Quicksilver struct {
	metrics.LauncherWorker

	// Custom Options
	command string
	prefix  string
}

// I think this is a simulation?
func (m Quicksilver) Family() string {
	return metrics.SimulationFamily
}

func (m Quicksilver) Url() string {
	return "https://github.com/LLNL/Quicksilver"
}

// Set custom options / attributes for the metric
func (m *Quicksilver) SetOptions(metric *api.Metric) {
	// Set user defined values or fall back to defaults
	m.prefix = "mpirun --hostfile ./hostlist.txt"
	m.command = "qs /opt/quicksilver/Examples/CORAL2_Benchmark/Problem1/Coral2_P1.inp"
	m.Workdir = "/opt/quicksilver/Examples"
	m.SetDefaultOptions(metric)
}

// Exported options and list options
func (m Quicksilver) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"prefix":  intstr.FromString(m.prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func init() {
	launcher := metrics.LauncherWorker{
		Identifier:     "app-quicksilver",
		Summary:        "A proxy app for the Monte Carlo Transport Code",
		Container:      "ghcr.io/converged-computing/metric-quicksilver:latest",
		WorkerScript:   "/metrics_operator/quicksilver-worker.sh",
		LauncherScript: "/metrics_operator/quicksilver-launcher.sh",
	}
	Quicksilver := Quicksilver{LauncherWorker: launcher}
	metrics.Register(&Quicksilver)
}
