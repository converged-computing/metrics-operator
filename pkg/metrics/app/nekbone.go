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

type Nekbone struct {
	jobs.LauncherWorker

	// Custom Options
	command string
	prefix  string
}

// I think this is a simulation?
func (m Nekbone) Family() string {
	return metrics.SolverFamily
}

func (m Nekbone) Url() string {
	return "https://github.com/Nek5000/Nekbone"
}

// Set custom options / attributes for the metric
func (m *Nekbone) SetOptions(metric *api.Metric) {
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Set user defined values or fall back to defaults
	m.prefix = "/bin/bash"
	m.command = "mpiexec --hostfile ./hostlist.txt -np 2 ./nekbone"
	m.Workdir = "/root/nekbone-3.0/test/example2"

	// This could be improved :)
	command, ok := metric.Options["command"]
	if ok {
		m.command = command.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		m.Workdir = workdir.StrVal
	}
	prefix, ok := metric.Options["prefix"]
	if ok {
		m.prefix = prefix.StrVal
	}
}

// Exported options and list options
func (m Nekbone) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"prefix":  intstr.FromString(m.prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}
func (n Nekbone) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Return lookup of entrypoint scripts
func (m Nekbone) EntrypointScripts(
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
%s ./problem.sh
echo "%s"
%s
`
	launcherTemplate := prefix + fmt.Sprintf(
		template,
		metrics.Separator,
		m.prefix,
		metrics.CollectionEnd,
		metrics.Interactive(spec.Spec.Logging.Interactive),
	)

	// The worker just has sleep infinity added
	workerTemplate := prefix + "\nsleep infinity"
	return m.FinalizeEntrypoints(launcherTemplate, workerTemplate)
}

func init() {
	launcher := jobs.LauncherWorker{
		Identifier: "app-nekbone",
		Summary:    "A mini-app derived from the Nek5000 CFD code which is a high order, incompressible Navier-Stokes CFD solver based on the spectral element method. The conjugate gradiant solve is compute intense, contains small messages and frequent allreduces.",
		Container:  "ghcr.io/converged-computing/metric-nekbone:latest",
	}
	Nekbone := Nekbone{LauncherWorker: launcher}
	metrics.Register(&Nekbone)
}
