/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package network

import (
	"fmt"
	"path"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/converged-computing/metrics-operator/pkg/metadata"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/converged-computing/metrics-operator/pkg/specs"
)

// ghcr.io/converged-computing/metric-osu-benchmark:latest
// https://mvapich.cse.ohio-state.edu/benchmarks/

var (

	// Directory (app) name and executable in /root/chatterbug
	ChatterbugApps = map[string]string{
		"pairs":         "pairs.x",
		"ping-pong":     "ping-ping.x",
		"spread":        "spread.x",
		"stencil3d":     "stencil3d.x",
		"stencil4d":     "stencil4d.x",
		"subcom2d-coll": "subcom2d-coll.x",
		"subcom2d-a2a":  "subcom2d-a2a.x",
		"unstr-mesh":    "unstr-mesh.x",
	}
)

type Chatterbug struct {
	metrics.LauncherWorker

	// Custom options
	command string
	tasks   int32
	lookup  map[string]bool
	// mpirun options (e.g., -N 2)
	mpirun string
	// args for executable
	args string
}

func (m Chatterbug) Url() string {
	return "https://github.com/hpcgroup/chatterbug"
}

// Determine if the command is poised to run
func (m *Chatterbug) hasCommand(command string) bool {
	_, exists := m.lookup[command]
	return exists
}

// Set custom options / attributes for the metric
func (m *Chatterbug) SetOptions(metric *api.Metric) {
	m.lookup = map[string]bool{}

	// Default command and args (for a demo)
	m.command = "stencil3d"
	m.args = "./stencil3d.x 2 2 2 10 10 10 4 1"
	m.mpirun = "-N 8"
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// One pod per hostname?
	m.SoleTenancy = true

	// We are allowed to specify just one command
	command, ok := metric.Options["command"]
	if ok {
		_, ok := ChatterbugApps[command.StrVal]
		if !ok {
			fmt.Printf("🟥️ Chatterbug command %s is not known", command.StrVal)
		} else {
			m.command = command.StrVal
		}
	}

	// Don't use default tasks
	tasks, ok := metric.Options["tasks"]
	if ok {
		m.tasks = tasks.IntVal
	}
	st, ok := metric.Options["sole-tenancy"]
	if ok && st.StrVal == "false" || st.StrVal == "no" {
		m.SoleTenancy = false
	}
	mpirun, ok := metric.Options["mpirun"]
	if ok {
		m.mpirun = mpirun.StrVal
	}
	args, ok := metric.Options["args"]
	if ok {
		m.args = args.StrVal
	}
}

// Exported options and list options
func (m Chatterbug) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"command": intstr.FromString(m.command),
		"tasks":   intstr.FromInt(int(m.tasks)),
		"args":    intstr.FromString(m.args),
		"mpirun":  intstr.FromString(m.mpirun),
	}
}

// Family returns the network family
func (n Chatterbug) Family() string {
	return metrics.NetworkFamily
}

func (m Chatterbug) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)

	// The launcher has a different hostname, n for netmark
	hosts := m.GetHostlist(spec)
	prefixTemplate := `#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &

# If we have zero tasks, default to workers * nproc for total tasks
# This is only for non point to point benchmarks
np=%d
pods=%d
# Tasks per node, not total
tasks=$(nproc)
if [[ $np -eq 0 ]]; then
	np=$(( $pods*$tasks ))
fi

# note this isn't used by the job run, it is for the user FYI
echo "Number of tasks (nproc on one node) is $tasks"
echo "Number of tasks total (across $pods nodes) is $np"

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10

# Write the hosts file.
cat <<EOF > ./hostnames.txt
%s
EOF

# openmpi is evil and we need the ip addresses
for h in $(cat ./hostnames.txt); do
   if [[ "$h" == "" ]]; then
      continue
   fi
   address=$(getent hosts $h | awk '{ print $1 }')
   echo "${address}" >> ./hostlist.txt
done   

cat ./hostlist.txt
# Show metadata for run
echo "%s"
`
	prefix := fmt.Sprintf(
		prefixTemplate,
		m.tasks,
		spec.Spec.Pods,
		hosts,
		meta,
	)

	// Prepare command for chatterbug
	commands := fmt.Sprintf("\nsleep 5\necho %s\n", metadata.CollectionStart)

	// Full path to, e.g., /root/chatterbug/stencil3d/stencil3d.x
	command := path.Join("/root/chatterbug", m.command, ChatterbugApps[m.command])
	line := fmt.Sprintf("mpirun --hostfile ./hostlist.txt --allow-run-as-root %s %s %s", m.mpirun, command, m.args)

	commands += fmt.Sprintf("echo %s\necho \"%s\"\n", metadata.Separator, line)

	// The pre block has the prefix and commands, up to the echo of the command (line)
	preBlock := fmt.Sprintf("%s\n%s", prefix, commands)

	// The post block has the collection end and interactive option
	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	postBlock := fmt.Sprintf("echo %s\n%s\n", metadata.CollectionEnd, interactive)

	// The worker just has a preBlock with the prefix and the command is to sleep
	launcherEntrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(m.LauncherScript),
		Path:    m.LauncherScript,
		Pre:     preBlock,
		Command: line,
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
		Identifier: "network-chatterbug",
		Summary:    "A suite of communication proxies for HPC applications",
		Container:  "ghcr.io/converged-computing/metric-chatterbug:latest",
	}
	launcher := metrics.LauncherWorker{
		BaseMetric:     base,
		WorkerScript:   "/metrics_operator/chatterbug-worker.sh",
		LauncherScript: "/metrics_operator/chatterbug-launcher.sh",
	}
	bug := Chatterbug{LauncherWorker: launcher}
	metrics.Register(&bug)
}
