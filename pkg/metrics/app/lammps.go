/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package application

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/converged-computing/metrics-operator/pkg/metadata"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/converged-computing/metrics-operator/pkg/specs"
)

type Lammps struct {
	metrics.LauncherWorker
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
	// Set user defined values or fall back to defaults
	// This is a more manual approach that puts the user in charge of determining the entire command
	// This more closely matches what we might do on HPC :)
	m.Command = "mpirun --hostfile ./hostlist.txt -np 2 --map-by socket lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite"
	m.Workdir = "/opt/lammps/examples/reaxff/HNS"
	m.SetDefaultOptions(metric)
}

// Validate that we can run Lammps
func (n Lammps) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m Lammps) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.Command),
		"workdir": intstr.FromString(m.Workdir),
	}
}
func (n Lammps) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Prepare containers with jobs and entrypoint scripts
func (m Lammps) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)
	hosts := m.GetHostlist(spec)
	prefix := m.GetCommonPrefix(meta, m.Command, hosts)

	// Template blocks for launcher script
	preBlock := `
echo "%s"
`

	postBlock := `
echo "%s"
%s
`
	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	preBlock = prefix + fmt.Sprintf(preBlock, metadata.Separator)
	postBlock = fmt.Sprintf(postBlock, metadata.CollectionEnd, interactive)

	// Entrypoint for the launcher
	launcherEntrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(m.LauncherScript),
		Path:    m.LauncherScript,
		Pre:     preBlock,
		Command: m.Command,
		Post:    postBlock,
	}

	// Entrypoint for the worker
	// Just has a sleep infinity added to the prefix
	workerEntrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(m.WorkerScript),
		Path:    m.WorkerScript,
		Pre:     prefix,
		Command: "sleep infinity",
	}

	// These are associated with replicated jobs via JobName
	launcherContainer := m.GetLauncherContainerSpec(launcherEntrypoint)
	workerContainer := m.GetWorkerContainerSpec(workerEntrypoint)

	// Return the script templates for each of launcher and worker
	return []*specs.ContainerSpec{&launcherContainer, &workerContainer}
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
