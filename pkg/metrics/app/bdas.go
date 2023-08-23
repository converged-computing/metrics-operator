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

type BDAS struct {
	jobs.LauncherWorker

	// Custom Options
	command string
	prefix  string
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
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Set user defined values or fall back to defaults
	m.prefix = "/bin/bash"
	m.command = "mpirun --allow-run-as-root -np 4 --hostfile ./hostlist.txt Rscript /opt/bdas/benchmarks/r/princomp.r 250 50"
	m.Workdir = "/opt/bdas/benchmarks/r"

	// Examples from guide
	// mpirun -np num_ranks Rscript princomp.r num_local_rows num_global_cols
	// mpirun -np 16 Rscript princomp.r 1000 250

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
func (m BDAS) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"prefix":  intstr.FromString(m.prefix),
		"workdir": intstr.FromString(m.Workdir),
	}
}

// Return lookup of entrypoint scripts
func (m BDAS) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)
	hosts := m.GetHostlist(spec)
	prefix := m.GetCommonPrefix(metadata, m.command, hosts)

	// Template for the launcher
	// TODO need to finish adding here when BDAS rebuild done
	template := `
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
echo "%s"
echo "%s"
%s ./problem.sh
echo "%s"
%s
`
	launcherTemplate := prefix + fmt.Sprintf(
		template,
		metadata,
		metrics.CollectionStart,
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
		Identifier: "app-bdas",
		Summary:    "The big data analytic suite contains the K-Means observation label, PCA, and SVM benchmarks.",
		Container:  "ghcr.io/converged-computing/metric-bdas:latest",
	}

	BDAS := BDAS{LauncherWorker: launcher}
	metrics.Register(&BDAS)
}
