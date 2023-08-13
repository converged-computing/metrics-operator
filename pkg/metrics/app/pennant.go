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

type Pennant struct {
	jobs.LauncherWorker

	// Custom Options
	workdir string
	command string
	mpirun  string
}

func (m Pennant) Url() string {
	return "https://github.com/LLNL/pennant"
}

// Set custom options / attributes for the metric
func (m *Pennant) SetOptions(metric *api.Metric) {
	m.Rate = metric.Rate
	m.Completions = metric.Completions
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Set user defined values or fall back to defaults
	m.mpirun = "mpirun --hostfile ./hostlist.txt"
	m.command = "pennant /opt/pennant/test/sedovsmall/sedovsmall.pnt"
	m.workdir = "/opt/pennant/test"

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

// Exported options and list options
func (m Pennant) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.Rate)),
		"completions": intstr.FromInt(int(m.Completions)),
		"command":     intstr.FromString(m.command),
		"mpirun":      intstr.FromString(m.mpirun),
		"workdir":     intstr.FromString(m.workdir),
	}
}
func (n Pennant) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Return lookup of entrypoint scripts
func (m Pennant) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)
	hosts := m.GetHostlist(spec)

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
	return m.FinalizeEntrypoints(launcherTemplate, workerTemplate)
}

func init() {
	launcher := jobs.LauncherWorker{
		Identifier:     "app-pennant",
		Summary:        "Unstructured mesh hydrodynamics for advanced architectures ",
		Container:      "ghcr.io/converged-computing/metric-pennant:latest",
		WorkerScript:   "/metrics_operator/pennant-worker.sh",
		LauncherScript: "/metrics_operator/pennant-launcher.sh",
	}
	Pennant := Pennant{LauncherWorker: launcher}
	metrics.Register(&Pennant)
}
