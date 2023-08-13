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

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

type Lammps struct {
	LauncherWorkerApp

	// Options
	workdir string
	command string
}

func (m Lammps) Url() string {
	return "https://www.lammps.org/"
}

// Set custom options / attributes for the metric
func (m *Lammps) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions
	m.resources = &metric.Resources
	m.attributes = &metric.Attributes

	// Set user defined values or fall back to defaults
	// This is a more manual approach that puts the user in charge of determining the entire command
	// This more closely matches what we might do on HPC :)
	m.command = "mpirun --hostfile ./hostlist.txt -np 2 --map-by socket lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite"
	m.workdir = "/opt/lammps/examples/reaxff/HNS"

	// This could be improved :)
	command, ok := metric.Options["command"]
	if ok {
		m.command = command.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		m.workdir = workdir.StrVal
	}
}

// Validate that we can run Lammps
func (n Lammps) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m Lammps) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.rate)),
		"completions": intstr.FromInt(int(m.completions)),
		"command":     intstr.FromString(m.command),
		"workdir":     intstr.FromString(m.workdir),
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
	hosts := m.getHostlist(spec)

	prefixTemplate := `#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
echo "%s"
# Change directory to where we will run (and write hostfile)
cd %s
# Write the hosts file
cat <<EOF > ./hostlist.txt
%s
EOF

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10
echo "%s"
`
	prefix := fmt.Sprintf(
		prefixTemplate,
		metadata,
		m.workdir,
		hosts,
		metrics.CollectionStart,
	)

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
	return m.finalizeEntrypoints(launcherTemplate, workerTemplate)
}

func init() {
	launcher := LauncherWorkerApp{
		name:           "app-lammps",
		description:    "LAMMPS molecular dynamic simulation",
		container:      "ghcr.io/converged-computing/metric-lammps:latest",
		workerScript:   "/metrics_operator/lammps-worker.sh",
		launcherScript: "/metrics_operator/lammps-launcher.sh",
	}
	lammps := Lammps{LauncherWorkerApp: launcher}
	metrics.Register(&lammps)
}
