/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package network

import (
	"fmt"
	"path"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// ghcr.io/converged-computing/metric-osu-benchmark:latest
// https://mvapich.cse.ohio-state.edu/benchmarks/

var (
	// I chose this set that can run on my small local machine
	// likely others work on more powerful machines
	osuBenchmarkCommands = map[string]bool{
		"osu_acc_latency":     true,
		"osu_fop_latency":     true,
		"osu_get_acc_latency": true,
		"osu_get_latency":     true,
		"osu_put_latency":     true,
	}
)

type OSUBenchmark struct {
	name        string
	rate        int32
	completions int32
	description string
	container   string

	// Scripts
	// TODO make a function that derives this across metrics?
	workerScript      string
	launcherScript    string
	workerScriptKey   string
	launcherScriptKey string
	commands          []string
	lookup            map[string]bool
}

// Name returns the metric name
func (m OSUBenchmark) Name() string {
	return m.name
}

// Description returns the metric description
func (m OSUBenchmark) Description() string {
	return m.description
}

// Jobs required for success condition (l is the osu benchmark launcher)
func (m OSUBenchmark) SuccessJobs() []string {
	return []string{"l"}
}

// Container variables
func (n OSUBenchmark) Type() string {
	return metrics.StandaloneMetric
}
func (n OSUBenchmark) Image() string {
	return n.container
}
func (n OSUBenchmark) WorkingDir() string {
	return "/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided"
}

func (n OSUBenchmark) getMetricsKeyToPath() []corev1.KeyToPath {
	// Runner start scripts
	makeExecutable := int32(0777)

	// Each metric has an entrypoint script
	return []corev1.KeyToPath{
		{
			Key:  n.launcherScriptKey,
			Path: path.Base(n.launcherScript),
			Mode: &makeExecutable,
		},
		{
			Key:  n.workerScriptKey,
			Path: path.Base(n.workerScript),
			Mode: &makeExecutable,
		},
	}
}

// Replicated Jobs are custom for this standalone metric
func (m OSUBenchmark) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {

	// Generate a replicated job for the launcher (netmark) and workers
	launcher := metrics.GetReplicatedJob(spec, false, 1, 1, "l")
	workers := metrics.GetReplicatedJob(spec, false, 1, 1, "w")

	// Add volumes defined under storage.
	v := map[string]api.Volume{}
	if spec.HasStorage() {
		v["storage"] = spec.Spec.Storage.Volume
	}

	runnerScripts := m.getMetricsKeyToPath()
	volumes := metrics.GetVolumes(spec, runnerScripts, v)
	launcher.Template.Spec.Template.Spec.Volumes = volumes
	workers.Template.Spec.Template.Spec.Volumes = volumes

	// Prepare container specs, one for launcher and one for workers
	launcherSpec := []metrics.ContainerSpec{
		{
			Image:      m.container,
			Name:       "launcher",
			Command:    []string{"/bin/bash", m.launcherScript},
			WorkingDir: m.WorkingDir(),
		},
	}
	workerSpec := []metrics.ContainerSpec{
		{
			Image:      m.container,
			Name:       "workers",
			Command:    []string{"/bin/bash", m.workerScript},
			WorkingDir: m.WorkingDir(),
		},
	}
	js := []jobset.ReplicatedJob{}

	// Derive the containers, one per metric
	// This will also include mounts for volumes
	// TODO allow ptrace and getting root filesystem?
	launcherContainers, err := metrics.GetContainers(spec, launcherSpec, v, false)
	if err != nil {
		fmt.Printf("issue creating launcher containers %s", err)
		return js, err
	}
	workerContainers, err := metrics.GetContainers(spec, workerSpec, v, false)
	if err != nil {
		fmt.Printf("issue creating worker containers %s", err)
		return js, err
	}
	launcher.Template.Spec.Template.Spec.Containers = launcherContainers
	workers.Template.Spec.Template.Spec.Containers = workerContainers
	js = []jobset.ReplicatedJob{*launcher, *workers}
	return js, nil
}

// Determine if the command is poised to run
func (m *OSUBenchmark) hasCommand(command string) bool {
	_, exists := m.lookup[command]
	return exists
}

func (m *OSUBenchmark) addCommand(command string) {
	m.commands = append(m.commands, command)
	m.lookup[command] = true
}

// Set custom options / attributes for the metric
func (m *OSUBenchmark) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions
	m.lookup = map[string]bool{}
	m.commands = []string{}

	// We are allowed to specify just one command
	opts, ok := metric.ListOptions["commands"]
	if ok {
		// Parse list options that are valid
		for _, opt := range opts {
			_, ok := osuBenchmarkCommands[opt.StrVal]
			if ok && !m.hasCommand(opt.StrVal) {
				m.addCommand(opt.StrVal)
			}
		}
	}
	if len(m.commands) == 0 {
		for command := range osuBenchmarkCommands {
			if !m.hasCommand(command) {
				m.addCommand(command)
			}
		}
	}
}

// OSU Benchmarks MUST be run with two nodes
func (n OSUBenchmark) Validate(spec *api.MetricSet) bool {
	if len(n.commands) == 0 {
		fmt.Printf("🟥️ OSUBenchmark not valid, requires 1+ commands.")
		return false
	}
	return spec.Spec.Pods == 2
}

// Return lookup of entrypoint scripts
func (m OSUBenchmark) EntrypointScripts(spec *api.MetricSet) []metrics.EntrypointScript {

	// Generate hostlists
	// The launcher has a different hostname, n for netmark
	launcherHost := fmt.Sprintf("%s-l-0-0.%s.%s.svc.cluster.local",
		spec.Name, spec.Spec.ServiceName, spec.Namespace,
	)
	workerHost := fmt.Sprintf("%s-w-0-0.%s.%s.svc.cluster.local",
		spec.Name, spec.Spec.ServiceName, spec.Namespace,
	)
	prefixTemplate := `#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
whoami
# Show ourselves!
cat ${0}
	
# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10

# Write the hosts file
launcher=$(getent hosts %s | awk '{ print $1 }')
worker=$(getent hosts %s | awk '{ print $1 }')
echo "${launcher}" >> ./hostfile.txt
echo "${worker}" >> ./hostfile.txt
`
	prefix := fmt.Sprintf(prefixTemplate, launcherHost, workerHost)

	// Prepare list of commands, e.g.,
	// mpirun -f ./hostlist.txt -np 2 ./osu_acc_latency (mpich)
	// mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_fop_latency (openmpi)
	// Sleep a little more to allow worker to write launcher hostname
	commands := "sleep 5\n"
	for _, command := range m.commands {
		line := fmt.Sprintf("mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./%s", command)
		commands += fmt.Sprintf("echo \"%s\"\n%s\n", line, line)
	}

	// Template for the launcher
	//
	launcherTemplate := fmt.Sprintf("%s\n%s", prefix, commands)

	// The worker just has sleep infinity added, and getting the ip address of the launcher
	workerTemplate := prefix + "\nsleep infinity"

	// Return the script templates for each of launcher and worker
	return []metrics.EntrypointScript{
		{
			Name:   m.launcherScriptKey,
			Path:   m.launcherScript,
			Script: launcherTemplate,
		},
		{
			Name:   m.workerScriptKey,
			Path:   m.workerScript,
			Script: workerTemplate,
		},
	}
}

func init() {
	metrics.Register(
		&OSUBenchmark{
			name:              "network-osu-benchmark",
			description:       "point to point MPI benchmarks, see https://mvapich.cse.ohio-state.edu/benchmarks/",
			container:         "ghcr.io/converged-computing/metric-osu-benchmark:latest",
			workerScript:      "/metrics_operator/osu-worker.sh",
			launcherScript:    "/metrics_operator/osu-launcher.sh",
			workerScriptKey:   "osu-worker",
			launcherScriptKey: "osu-launcher",
		})
}
