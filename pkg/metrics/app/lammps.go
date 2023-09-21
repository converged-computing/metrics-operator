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

type Lammps struct {
	metrics.LauncherWorker

	// Options
	command string
}

func (m Lammps) Url() string {
	return "https://www.lammps.org/"
}

// I think this is a simulation?
func (m Lammps) Family() string {
	return metrics.SimulationFamily
}

// Set custom options / attributes for the metric
func (m *Lammps) SetOptions(metric *api.Metric) {
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Set user defined values or fall back to defaults
	// This is a more manual approach that puts the user in charge of determining the entire command
	// This more closely matches what we might do on HPC :)
	m.command = "mpirun --hostfile ./hostlist.txt -np 2 --map-by socket lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite"
	m.Workdir = "/opt/lammps/examples/reaxff/HNS"

	// This could be improved :)
	command, ok := metric.Options["command"]
	if ok {
		m.command = command.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		m.Workdir = workdir.StrVal
	}
}

// Validate that we can run Lammps
func (n Lammps) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m Lammps) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: "app-lammps",
		Summary:    "LAMMPS molecular dynamic simulation",
		Container:  "ghcr.io/converged-computing/metric-lammps:latest",
	}
	launcher := metrics.LauncherWorker{
		BaseMetric:     base,
		WorkerScript:   "/metrics_operator/lammps-worker.sh",
		LauncherScript: "/metrics_operator/lammps-launcher.sh",
	}
	lammps := Lammps{LauncherWorker: launcher}
	metrics.Register(&lammps)
}
