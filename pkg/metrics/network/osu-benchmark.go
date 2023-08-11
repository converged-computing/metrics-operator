/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package network

import (
	"fmt"
	"path"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// ghcr.io/converged-computing/metric-osu-benchmark:latest
// https://mvapich.cse.ohio-state.edu/benchmarks/

var (
	singleSidedDir  = "/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided"
	pointToPointDir = "/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt"
	collectiveDir   = "/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/collective"
	startupDir      = "/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/startup"

	// Lookup of all OSU benchmarks available
	osuBenchmarkCommands = map[string]string{

		// Single Sided
		"osu_get_acc_latency": singleSidedDir,
		"osu_acc_latency":     singleSidedDir, // Latency Test for Accumulate
		"osu_fop_latency":     singleSidedDir,
		"osu_get_latency":     singleSidedDir, // Latency Test for Get
		"osu_put_latency":     singleSidedDir, // Latency Test for Put
		"osu_cas_latency":     singleSidedDir,
		"osu_get_bw":          singleSidedDir,
		"osu_put_bibw":        singleSidedDir,
		"osu_put_bw":          singleSidedDir,

		// Collective
		"osu_allreduce":      collectiveDir, // MPI_Allreduce Latency Test
		"osu_allgather":      collectiveDir,
		"osu_allgatherv":     collectiveDir,
		"osu_alltoall":       collectiveDir,
		"osu_alltoallv":      collectiveDir,
		"osu_barrier":        collectiveDir,
		"osu_bcast":          collectiveDir,
		"osu_gather":         collectiveDir,
		"osu_gatherv":        collectiveDir,
		"osu_iallgather":     collectiveDir,
		"osu_iallgatherv":    collectiveDir,
		"osu_iallreduce":     collectiveDir,
		"osu_ialltoall":      collectiveDir,
		"osu_ialltoallv":     collectiveDir,
		"osu_ialltoallw":     collectiveDir,
		"osu_ibarrier":       collectiveDir,
		"osu_ibcast":         collectiveDir,
		"osu_igather":        collectiveDir,
		"osu_igatherv":       collectiveDir,
		"osu_ireduce":        collectiveDir,
		"osu_iscatter":       collectiveDir,
		"osu_iscatterv":      collectiveDir,
		"osu_reduce":         collectiveDir,
		"osu_reduce_scatter": collectiveDir,
		"osu_scatter":        collectiveDir,
		"osu_scatterv":       collectiveDir,

		// Point to Point
		"osu_latency":    pointToPointDir, // Latency Test
		"osu_bibw":       pointToPointDir, // Bidirectional Bandwidth Test
		"osu_bw":         pointToPointDir, // Bandwidth Test
		"osu_latency_mp": pointToPointDir,
		"osu_latency_mt": pointToPointDir,
		"osu_mbw_mr":     pointToPointDir,
		"osu_multi_lat":  pointToPointDir,

		// Startup
		"osu_hello": startupDir,
		"osu_init":  startupDir,
	}

	// Defaults that we provide when none specified
	osuBenchmarkDefaults = []string{
		"osu_get_acc_latency",
		"osu_acc_latency",
		"osu_fop_latency",
		"osu_get_latency",
		"osu_put_latency",
		"osu_allreduce",
		"osu_latency",
		"osu_bibw",
		"osu_bw",
	}
)

type OSUBenchmark struct {
	name        string
	rate        int32
	completions int32
	description string
	container   string
	resources   *api.ContainerResources
	attributes  *api.ContainerSpec

	// Scripts
	// TODO make a function that derives this across metrics?
	workerScript   string
	launcherScript string
	commands       []string
	lookup         map[string]bool
}

// Name returns the metric name
func (m OSUBenchmark) Name() string {
	return m.name
}

// Description returns the metric description
func (m OSUBenchmark) Description() string {
	return m.description
}
func (m OSUBenchmark) Url() string {
	return "https://mvapich.cse.ohio-state.edu/benchmarks/"
}
func (m OSUBenchmark) Attributes() *api.ContainerSpec {
	return m.attributes
}

// Return container resources for the metric container
func (m OSUBenchmark) Resources() *api.ContainerResources {
	return m.resources
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
			Key:  deriveScriptKey(n.launcherScript),
			Path: path.Base(n.launcherScript),
			Mode: &makeExecutable,
		},
		{
			Key:  deriveScriptKey(n.workerScript),
			Path: path.Base(n.workerScript),
			Mode: &makeExecutable,
		},
	}
}

// Replicated Jobs are custom for this standalone metric
func (m OSUBenchmark) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {

	js := []jobset.ReplicatedJob{}

	// Generate a replicated job for the launcher (netmark) and workers, both need sole tenancy
	launcher, err := metrics.GetReplicatedJob(spec, false, 1, 1, "l", false)
	if err != nil {
		return js, err
	}
	workers, err := metrics.GetReplicatedJob(spec, false, 1, 1, "w", false)
	if err != nil {
		return js, err
	}

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
			Resources:  m.resources,
			Attributes: m.attributes,
		},
	}
	workerSpec := []metrics.ContainerSpec{
		{
			Image:      m.container,
			Name:       "workers",
			Command:    []string{"/bin/bash", m.workerScript},
			WorkingDir: m.WorkingDir(),
			Resources:  m.resources,
			Attributes: m.attributes,
		},
	}

	// Derive the containers, one per metric
	// This will also include mounts for volumes
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
	// Get the fullpath from our lookup
	fullpath := path.Join(osuBenchmarkCommands[command], command)
	m.commands = append(m.commands, fullpath)
	m.lookup[command] = true
}

// Set custom options / attributes for the metric
func (m *OSUBenchmark) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions
	m.lookup = map[string]bool{}
	m.commands = []string{}
	m.resources = &metric.Resources
	m.attributes = &metric.Attributes

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

	// If not selected or found, fall back to default list
	if len(m.commands) == 0 {
		for _, command := range osuBenchmarkDefaults {
			if !m.hasCommand(command) {
				m.addCommand(command)
			}
		}
	}
}

// Exported options and list options
func (m OSUBenchmark) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.rate)),
		"completions": intstr.FromInt(int(m.completions)),
	}
}
func (m OSUBenchmark) ListOptions() map[string][]intstr.IntOrString {
	commands := []intstr.IntOrString{}
	for _, command := range m.commands {
		commands = append(commands, intstr.FromString(command))
	}
	return map[string][]intstr.IntOrString{
		"commands": commands,
	}
}

// OSU Benchmarks MUST be run with two nodes
func (m OSUBenchmark) Validate(spec *api.MetricSet) bool {
	if len(m.commands) == 0 {
		fmt.Printf("🟥️ OSUBenchmark not valid, requires 1+ commands.")
		return false
	}
	return spec.Spec.Pods == 2
}

// Return lookup of entrypoint scripts
func (m OSUBenchmark) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)

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

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10

# Write the hosts file
launcher=$(getent hosts %s | awk '{ print $1 }')
worker=$(getent hosts %s | awk '{ print $1 }')
echo "${launcher}" >> ./hostfile.txt
echo "${worker}" >> ./hostfile.txt

# Show metadata for run
echo "%s"
`
	prefix := fmt.Sprintf(prefixTemplate, launcherHost, workerHost, metadata)

	// Prepare list of commands, e.g.,
	// mpirun -f ./hostlist.txt -np 2 ./osu_acc_latency (mpich)
	// mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_fop_latency (openmpi)
	// Sleep a little more to allow worker to write launcher hostname
	commands := fmt.Sprintf("\nsleep 5\necho %s\n", metrics.CollectionStart)
	for _, command := range m.commands {

		// This starts the line with a separator for the new section
		line := fmt.Sprintf("mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 %s", command)
		commands += fmt.Sprintf("echo %s\necho \"%s\"\n%s\n", metrics.Separator, line, line)
	}

	// Close the commands block
	commands += fmt.Sprintf("echo %s\n", metrics.CollectionEnd)

	// Template for the launcher
	launcherTemplate := fmt.Sprintf("%s\n%s", prefix, commands)

	// The worker just has sleep infinity added, and getting the ip address of the launcher
	workerTemplate := prefix + "\nsleep infinity"

	// Return the script templates for each of launcher and worker
	return []metrics.EntrypointScript{
		{
			Name:   deriveScriptKey(m.launcherScript),
			Path:   m.launcherScript,
			Script: launcherTemplate,
		},
		{
			Name:   deriveScriptKey(m.workerScript),
			Path:   m.workerScript,
			Script: workerTemplate,
		},
	}
}

func init() {
	metrics.Register(
		&OSUBenchmark{
			name:           "network-osu-benchmark",
			description:    "point to point MPI benchmarks",
			container:      "ghcr.io/converged-computing/metric-osu-benchmark:latest",
			workerScript:   "/metrics_operator/osu-worker.sh",
			launcherScript: "/metrics_operator/osu-launcher.sh",
		})
}
