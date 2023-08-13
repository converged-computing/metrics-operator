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

type Kripke struct {
	LauncherWorkerApp

	// Options
	workdir string
	command string
	mpirun  string
}

func (m Kripke) Url() string {
	return "https://github.com/LLNL/Kripke"
}

// Set custom options / attributes for the metric
func (m *Kripke) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions
	m.resources = &metric.Resources
	m.attributes = &metric.Attributes

	// Set user defined values or fall back to defaults
	m.mpirun = "mpirun --hostfile ./hostlist.txt"
	m.command = "kripke"
	m.workdir = "/opt/kripke"

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

// Validate that we can run Kripke
func (n Kripke) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m Kripke) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.rate)),
		"completions": intstr.FromInt(int(m.completions)),
		"command":     intstr.FromString(m.command),
		"mpirun":      intstr.FromString(m.mpirun),
		"workdir":     intstr.FromString(m.workdir),
	}
}
func (n Kripke) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Return lookup of entrypoint scripts
func (m Kripke) EntrypointScripts(
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
	return m.finalizeEntrypoints(launcherTemplate, workerTemplate)
}

func init() {
	launcher := LauncherWorkerApp{
		name:           "app-kripke",
		description:    "parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids",
		container:      "ghcr.io/converged-computing/metric-kripke:latest",
		workerScript:   "/metrics_operator/kripke-worker.sh",
		launcherScript: "/metrics_operator/kripke-launcher.sh",
	}
	kripke := Kripke{LauncherWorkerApp: launcher}
	metrics.Register(&kripke)
}
