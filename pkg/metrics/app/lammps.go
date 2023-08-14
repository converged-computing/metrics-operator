/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package application

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/converged-computing/metrics-operator/pkg/jobs"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

type Lammps struct {
	jobs.LauncherWorker

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
	m.Rate = metric.Rate
	m.Completions = metric.Completions
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
		"rate":        intstr.FromInt(int(m.Rate)),
		"completions": intstr.FromInt(int(m.Completions)),
		"command":     intstr.FromString(m.command),
		"workdir":     intstr.FromString(m.Workdir),
	}
}
func (n Lammps) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Return lookup of entrypoint scripts
func (m Lammps) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)
	hosts := m.GetHostlist(spec)
	prefix := m.GetCommonPrefix(metadata, m.command, hosts)

	// Template for the launcher
	template := `
echo "%s"
%s
echo "%s"
%s
`
	launcherTemplate := prefix + fmt.Sprintf(
		template,
		metrics.Separator,
		m.command,
		metrics.CollectionEnd,
		metrics.Interactive(spec.Spec.Logging.Interactive),
	)

	// The worker just has sleep infinity added
	workerTemplate := prefix + "\nsleep infinity"

	// Return the script templates for each of launcher and worker
	return m.FinalizeEntrypoints(launcherTemplate, workerTemplate)
}

func init() {
	launcher := jobs.LauncherWorker{
		Identifier:     "app-lammps",
		Summary:        "LAMMPS molecular dynamic simulation",
		Container:      "ghcr.io/converged-computing/metric-lammps:latest",
		WorkerScript:   "/metrics_operator/lammps-worker.sh",
		LauncherScript: "/metrics_operator/lammps-launcher.sh",
	}
	lammps := Lammps{LauncherWorker: launcher}
	metrics.Register(&lammps)
}
