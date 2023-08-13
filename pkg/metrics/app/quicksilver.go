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

type Quicksilver struct {
	jobs.LauncherWorker

	// Custom Options
	workdir string
	command string
	mpirun  string
}

func (m Quicksilver) Url() string {
	return "https://github.com/LLNL/Quicksilver"
}

// Set custom options / attributes for the metric
func (m *Quicksilver) SetOptions(metric *api.Metric) {
	m.Rate = metric.Rate
	m.Completions = metric.Completions
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Set user defined values or fall back to defaults
	m.mpirun = "mpirun --hostfile ./hostlist.txt"
	m.command = "qs /opt/quicksilver/Examples/CORAL2_Benchmark/Problem1/Coral2_P1.inp"
	m.workdir = "/opt/quicksilver/Examples"

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
func (m Quicksilver) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.Rate)),
		"completions": intstr.FromInt(int(m.Completions)),
		"command":     intstr.FromString(m.command),
		"mpirun":      intstr.FromString(m.mpirun),
		"workdir":     intstr.FromString(m.workdir),
	}
}
func (n Quicksilver) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Return lookup of entrypoint scripts
func (m Quicksilver) EntrypointScripts(
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
		Identifier:     "app-quicksilver",
		Summary:        "A proxy app for the Monte Carlo Transport Code",
		Container:      "ghcr.io/converged-computing/metric-quicksilver:latest",
		WorkerScript:   "/metrics_operator/quicksilver-worker.sh",
		LauncherScript: "/metrics_operator/quicksilver-launcher.sh",
	}
	Quicksilver := Quicksilver{LauncherWorker: launcher}
	metrics.Register(&Quicksilver)
}
