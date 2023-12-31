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
	qsIdentifier = "app-quicksilver"
	qsSummary    = "A proxy app for the Monte Carlo Transport Code"
	qsContainer  = "ghcr.io/converged-computing/metric-quicksilver:latest"
)

type Quicksilver struct {
	metrics.LauncherWorker
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

	m.Identifier = qsIdentifier
	m.Summary = qsSummary
	m.Container = qsContainer

	// Set user defined values or fall back to defaults
	m.Prefix = "mpirun --hostfile ./hostlist.txt"
	m.Command = "qs /opt/quicksilver/Examples/CORAL2_Benchmark/Problem1/Coral2_P1.inp"
	m.Workdir = "/opt/quicksilver/Examples"
	m.SetDefaultOptions(metric)
}

// Exported options and list options
func (m Quicksilver) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.Command),
		"prefix":  intstr.FromString(m.Prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: qsIdentifier,
		Summary:    qsSummary,
		Container:  qsContainer,
	}
	launcher := metrics.LauncherWorker{BaseMetric: base}
	Quicksilver := Quicksilver{LauncherWorker: launcher}
	metrics.Register(&Quicksilver)
}
