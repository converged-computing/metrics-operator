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

type Kripke struct {
	metrics.LauncherWorker

	// Options
	command string
	prefix  string
}

func (m Kripke) Url() string {
	return "https://github.com/LLNL/Kripke"
}

// I think this is a simulation?
func (m Kripke) Family() string {
	return metrics.SolverFamily
}

// Set custom options / attributes for the metric
func (m *Kripke) SetOptions(metric *api.Metric) {

	// Set user defined values or fall back to defaults
	m.prefix = "mpirun --hostfile ./hostlist.txt"
	m.command = "kripke"
	m.Workdir = "/opt/kripke"
	m.SetDefaultOptions(metric)
}

// Validate that we can run Kripke
func (n Kripke) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m Kripke) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"prefix":  intstr.FromString(m.prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}
func (n Kripke) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

func init() {
	launcher := metrics.LauncherWorker{
		Identifier:     "app-kripke",
		Summary:        "parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids",
		Container:      "ghcr.io/converged-computing/metric-kripke:latest",
		WorkerScript:   "/metrics_operator/kripke-worker.sh",
		LauncherScript: "/metrics_operator/kripke-launcher.sh",
	}
	kripke := Kripke{LauncherWorker: launcher}
	metrics.Register(&kripke)
}
