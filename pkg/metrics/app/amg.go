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

// AMG is a launcher + workers metric application
type AMG struct {
	metrics.LauncherWorker

	// Custom Options
	workdir string
	command string
	prefix  string
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
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Set user defined values or fall back to defaults
	m.prefix = "mpirun --hostfile ./hostlist.txt"
	m.command = "amg"
	m.workdir = "/opt/AMG"

	// This could be improved :)
	command, ok := metric.Options["command"]
	if ok {
		m.command = command.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		m.workdir = workdir.StrVal
	}
	mpirun, ok := metric.Options["mpirun"]
	if ok {
		m.prefix = mpirun.StrVal
	}
}

// Validate that we can run AMG
func (n AMG) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m AMG) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"mpirun":  intstr.FromString(m.prefix),
		"workdir": intstr.FromString(m.workdir),
	}
}

func (m AMG) PrepareContainers(
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
		Identifier:     "app-amg",
		Summary:        "parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids",
		Container:      "ghcr.io/converged-computing/metric-amg:latest",
		WorkerScript:   "/metrics_operator/amg-worker.sh",
		LauncherScript: "/metrics_operator/amg-launcher.sh",
	}
	amg := AMG{LauncherWorker: launcher}
	metrics.Register(&amg)
}
