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

const (
	lammpsIdentifier = "app-lammps"
	lammpsSummary    = "LAMMPS molecular dynamic simulation"
	lammpsContainer  = "ghcr.io/converged-computing/metric-lammps:latest"
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

	// Default metric options, these are overridden when we reflect
	m.Identifier = lammpsIdentifier
	m.Summary = lammpsSummary
	m.Container = lammpsContainer

	// Ensure we set sole tenancy if desired
	st, ok := metric.Options["soleTenancy"]
	if ok && st.StrVal == "true" || st.StrVal == "yes" {
		m.SoleTenancy = false
	}

	// Set user defined values or fall back to defaults
	// This is a more manual approach that puts the user in charge of determining the entire command
	// This more closely matches what we might do on HPC :)
	m.Command = "mpirun --hostfile ./hostlist.txt -np 2 --map-by socket lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite"
	m.Workdir = "/opt/lammps/examples/reaxff/HNS"
	m.SetDefaultOptions(metric)
}

// LAMMPS can be run on one node
func (m Lammps) Validate(spec *api.MetricSet) bool {
	return true
}

// Exported options and list options
func (m Lammps) Options() map[string]intstr.IntOrString {
	values := map[string]intstr.IntOrString{
		"command":    intstr.FromString(m.Command),
		"workdir":    intstr.FromString(m.Workdir),
		"soleTenancy": intstr.FromString("false"),
	}
	if m.SoleTenancy {
		values["soleTenancy"] = intstr.FromString("true")
	}
	return values
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

// TODO can we have a new function instead?
func init() {
	base := metrics.BaseMetric{
		Identifier: lammpsIdentifier,
		Summary:    lammpsSummary,
		Container:  lammpsContainer,
	}
	launcher := metrics.LauncherWorker{BaseMetric: base}
	lammps := Lammps{LauncherWorker: launcher}
	metrics.Register(&lammps)
}
