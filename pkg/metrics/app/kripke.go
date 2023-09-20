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

	"github.com/converged-computing/metrics-operator/pkg/metadata"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/converged-computing/metrics-operator/pkg/specs"
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
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Set user defined values or fall back to defaults
	m.prefix = "mpirun --hostfile ./hostlist.txt"
	m.command = "kripke"
	m.Workdir = "/opt/kripke"

	// This could be improved :)
	command, ok := metric.Options["command"]
	if ok {
		m.command = command.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		m.Workdir = workdir.StrVal
	}
	mpirun, ok := metric.Options["mpirun"]
	if ok {
		m.prefix = mpirun.StrVal
	}
}

// Validate that we can run Kripke
func (n Kripke) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m Kripke) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"mpirun":  intstr.FromString(m.prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}
func (n Kripke) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

func (m Kripke) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)
	hosts := m.GetHostlist(spec)
	prefix := m.GetCommonPrefix(meta, m.command, hosts)

	preBlock := `
echo "%s"
`

	postBlock := `
echo "%s"
%s
`
	command := fmt.Sprintf("%s ./problem.sh", m.prefix)
	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	preBlock = prefix + fmt.Sprintf(preBlock, metadata.Separator)
	postBlock = fmt.Sprintf(postBlock, metadata.CollectionEnd, interactive)

	// Entrypoint for the launcher
	launcherEntrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(m.LauncherScript),
		Path:    m.LauncherScript,
		Pre:     preBlock,
		Command: command,
		Post:    postBlock,
	}

	// Entrypoint for the worker
	workerEntrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(m.WorkerScript),
		Path:    m.WorkerScript,
		Pre:     prefix,
		Command: "sleep infinity",
	}

	// Container spec for the launcher
	launcherContainer := m.GetLauncherContainerSpec(launcherEntrypoint)
	workerContainer := m.GetWorkerContainerSpec(workerEntrypoint)

	// Return the script templates for each of launcher and worker
	return []*specs.ContainerSpec{&launcherContainer, &workerContainer}

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
