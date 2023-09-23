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
	bdasIdentifier = "app-bdas"
	bdasSummary    = "The big data analytic suite contains the K-Means observation label, PCA, and SVM benchmarks."
	bdasContainer  = "ghcr.io/converged-computing/metric-bdas:latest"
)

type BDAS struct {
	metrics.LauncherWorker
}

// I think this is a simulation?
func (m BDAS) Family() string {
	return metrics.MachineLearningFamily
}

func (m BDAS) Url() string {
	return "https://asc.llnl.gov/sites/asc/files/2020-09/BDAS_Summary_b4bcf27_0.pdf"
}

// Set custom options / attributes for the metric
func (m *BDAS) SetOptions(metric *api.Metric) {

	// Metadatqa
	m.Identifier = bdasIdentifier
	m.Summary = bdasSummary
	m.Container = bdasContainer

	// Set user defined values or fall back to defaults
	m.Prefix = "/bin/bash"
	m.Command = "mpirun --allow-run-as-root -np 4 --hostfile ./hostlist.txt Rscript /opt/bdas/benchmarks/r/princomp.r 250 50"
	m.Workdir = "/opt/bdas/benchmarks/r"

	// Examples from guide
	// mpirun -np num_ranks Rscript princomp.r num_local_rows num_global_cols
	// mpirun -np 16 Rscript princomp.r 1000 250
	m.SetDefaultOptions(metric)
}

// Exported options and list options
func (m BDAS) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.Command),
		"prefix":  intstr.FromString(m.Prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

func (m BDAS) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)
	hosts := m.GetHostlist(spec)
	prefix := m.GetCommonPrefix(meta, m.Command, hosts)

	preBlock := `
echo "%s"

# We need ip addresses for openmpi
mv ./hostlist.txt ./hostnames.txt
for h in $(cat ./hostnames.txt); do
  if [[ "${h}" != "" ]]; then
	if [[ "${h}" == "$(hostname)" ]]; then
		hostname -I | awk '{print $1}' >> hostlist.txt
	else
		host $h | cut -d ' ' -f 4 >> hostlist.txt
	fi
  fi  
done
echo "Hostlist"
cat ./hostlist.txt
`

	postBlock := `
echo "%s"
%s
`
	command := fmt.Sprintf("%s ./problem.sh", m.Prefix)
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
	base := metrics.BaseMetric{
		Identifier: bdasIdentifier,
		Summary:    bdasSummary,
		Container:  bdasContainer,
	}
	launcher := metrics.LauncherWorker{BaseMetric: base}
	BDAS := BDAS{LauncherWorker: launcher}
	metrics.Register(&BDAS)
}
