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

// AMG is a launcher + workers metric application
type AMG struct {
	LauncherWorkerApp

	// Custom Options
	workdir string
	command string
	mpirun  string
}

func (m AMG) Url() string {
	return "https://github.com/LLNL/AMG"
}

// Set custom options / attributes for the metric
func (m *AMG) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions
	m.resources = &metric.Resources
	m.attributes = &metric.Attributes

	// Set user defined values or fall back to defaults
	m.mpirun = "mpirun --hostfile ./hostlist.txt"
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
		m.mpirun = mpirun.StrVal
	}
}

// Validate that we can run AMG
func (n AMG) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m AMG) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.rate)),
		"completions": intstr.FromInt(int(m.completions)),
		"command":     intstr.FromString(m.command),
		"mpirun":      intstr.FromString(m.mpirun),
		"workdir":     intstr.FromString(m.workdir),
	}
}
func (n AMG) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Return lookup of entrypoint scripts
func (m AMG) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)

	// Generate hostlists
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

# Write the command file for mpirun
cat <<EOF > ./problem.sh
#!/bin/bash
%s
EOF
chmod +x ./problem.sh

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
		m.command,
		metrics.CollectionStart,
	)

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
		m.mpirun,
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
		name:           "app-amg",
		description:    "parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids",
		container:      "ghcr.io/converged-computing/metric-amg:latest",
		workerScript:   "/metrics_operator/amg-worker.sh",
		launcherScript: "/metrics_operator/amg-launcher.sh",
	}
	amg := AMG{LauncherWorkerApp: launcher}
	metrics.Register(&amg)
}
