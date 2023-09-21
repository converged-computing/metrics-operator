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

// AMG is a launcher + workers metric application
type AMG struct {
	metrics.LauncherWorker
}

func (m AMG) Url() string {
	return "https://github.com/LLNL/AMG"
}

// I think this is a simulation?
func (m AMG) Family() string {
	return metrics.SolverFamily
}

// Set custom options / attributes for the metric
func (m *AMG) SetOptions(metric *api.Metric) {
	// Set user defined values or fall back to defaults
	m.Prefix = "mpirun --hostfile ./hostlist.txt"
	m.Command = "amg"
	m.Workdir = "/opt/AMG"
	m.SetDefaultOptions(metric)
}

// Validate that we can run AMG
func (n AMG) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m AMG) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.Command),
		"prefix":  intstr.FromString(m.Prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: "app-amg",
		Summary:    "parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids",
		Container:  "ghcr.io/converged-computing/metric-amg:latest",
	}
	launcher := metrics.LauncherWorker{
		BaseMetric:     base,
		WorkerScript:   "/metrics_operator/amg-worker.sh",
		LauncherScript: "/metrics_operator/amg-launcher.sh",
	}
	amg := AMG{LauncherWorker: launcher}
	metrics.Register(&amg)
}
