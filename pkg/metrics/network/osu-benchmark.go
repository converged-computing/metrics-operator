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

type BenchmarkConfig struct {
	Workdir string
	Flags   string
}

var (
	singleSidedDir  = "/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided"
	pointToPointDir = "/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt"
	collectiveDir   = "/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/collective"
	startupDir      = "/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/startup"

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

	// Lookup of all OSU benchmarks available
	osuBenchmarkCommands = map[string]BenchmarkConfig{

		// Single Sided
		"osu_get_acc_latency": {Workdir: singleSidedDir, Flags: ""},
		"osu_acc_latency":     {Workdir: singleSidedDir, Flags: ""}, // Latency Test for Accumulate
		"osu_fop_latency":     {Workdir: singleSidedDir, Flags: ""},
		"osu_get_latency":     {Workdir: singleSidedDir, Flags: ""}, // Latency Test for Get
		"osu_put_latency":     {Workdir: singleSidedDir, Flags: ""}, // Latency Test for Put
		"osu_cas_latency":     {Workdir: singleSidedDir, Flags: ""},
		"osu_get_bw":          {Workdir: singleSidedDir, Flags: ""},
		"osu_put_bibw":        {Workdir: singleSidedDir, Flags: ""},
		"osu_put_bw":          {Workdir: singleSidedDir, Flags: ""},

		// Collective
		// For allreduce this should work, need to test -np $np -map-by ppr:1:node -rank-by core
		"osu_allreduce":      {Workdir: collectiveDir, Flags: "-np 2 -map-by ppr:1:node -rank-by core"}, // MPI_Allreduce Latency Test
		"osu_allgather":      {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:1:node -rank-by core"},
		"osu_allgatherv":     {Workdir: collectiveDir, Flags: ""},
		"osu_alltoall":       {Workdir: collectiveDir, Flags: ""},
		"osu_alltoallv":      {Workdir: collectiveDir, Flags: ""},
		"osu_barrier":        {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:1:node -rank-by core"},
		"osu_bcast":          {Workdir: collectiveDir, Flags: ""},
		"osu_gather":         {Workdir: collectiveDir, Flags: ""},
		"osu_gatherv":        {Workdir: collectiveDir, Flags: ""},
		"osu_iallgather":     {Workdir: collectiveDir, Flags: ""},
		"osu_iallgatherv":    {Workdir: collectiveDir, Flags: ""},
		"osu_iallreduce":     {Workdir: collectiveDir, Flags: ""},
		"osu_ialltoall":      {Workdir: collectiveDir, Flags: ""},
		"osu_ialltoallv":     {Workdir: collectiveDir, Flags: ""},
		"osu_ialltoallw":     {Workdir: collectiveDir, Flags: ""},
		"osu_ibarrier":       {Workdir: collectiveDir, Flags: ""},
		"osu_ibcast":         {Workdir: collectiveDir, Flags: ""},
		"osu_igather":        {Workdir: collectiveDir, Flags: ""},
		"osu_igatherv":       {Workdir: collectiveDir, Flags: ""},
		"osu_ireduce":        {Workdir: collectiveDir, Flags: ""},
		"osu_iscatter":       {Workdir: collectiveDir, Flags: ""},
		"osu_iscatterv":      {Workdir: collectiveDir, Flags: ""},
		"osu_reduce":         {Workdir: collectiveDir, Flags: ""},
		"osu_reduce_scatter": {Workdir: collectiveDir, Flags: ""},
		"osu_scatter":        {Workdir: collectiveDir, Flags: ""},
		"osu_scatterv":       {Workdir: collectiveDir, Flags: ""},

		// Point to Point
		"osu_latency":    {Workdir: pointToPointDir, Flags: "-np 2 -map-by ppr:1:node"}, // Latency Test
		"osu_bibw":       {Workdir: pointToPointDir, Flags: "-np 2 -map-by ppr:1:node"}, // Bidirectional Bandwidth Test
		"osu_bw":         {Workdir: pointToPointDir, Flags: "-np 2 -map-by ppr:1:node"}, // Bandwidth Test
		"osu_latency_mp": {Workdir: pointToPointDir, Flags: ""},
		"osu_latency_mt": {Workdir: pointToPointDir, Flags: ""},
		"osu_mbw_mr":     {Workdir: pointToPointDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core"},
		"osu_multi_lat":  {Workdir: pointToPointDir, Flags: ""},

		// Startup
		"osu_hello": {Workdir: startupDir, Flags: ""},
		"osu_init":  {Workdir: startupDir, Flags: ""},
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
	workerScript   string
	launcherScript string
	commands       []string
	tasks          int32
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
	m.commands = append(m.commands, command)
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

	// Don't use default tasks
	tasks, ok := metric.Options["tasks"]
	if ok {
		m.tasks = tasks.IntVal
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
		"tasks":       intstr.FromInt(int(m.tasks)),
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

# If we have zero tasks, default to workers * nproc for total tasks
# This is only for non point to point benchmarks
np=%d
pods=%d
# Tasks per node, not total
tasks=$(nproc)
if [[ $np -eq 0 ]]; then
	np=$(( $pods*$tasks ))
fi

echo "Number of tasks (nproc on one node) is $tasks"
echo "Number of tasks total (across $pods nodes) is $np"

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
	prefix := fmt.Sprintf(
		prefixTemplate,
		m.tasks,
		spec.Spec.Pods,
		launcherHost,
		workerHost,
		metadata,
	)

	// Prepare list of commands, e.g.,
	// mpirun -f ./hostlist.txt -np 2 ./osu_acc_latency (mpich)
	// mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 ./osu_fop_latency (openmpi)
	// Sleep a little more to allow worker to write launcher hostname
	commands := fmt.Sprintf("\nsleep 5\necho %s\n", metrics.CollectionStart)
	for _, executable := range m.commands {

		command := path.Join(osuBenchmarkCommands[executable].Workdir, executable)

		// Flags can vary by command
		flags := osuBenchmarkCommands[executable].Flags
		if flags == "" {
			flags = "-np 2 -map-by ppr:1:node"
		}

		// Assume always 2 nodes for now
		line := fmt.Sprintf("mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 %s %s", flags, command)
		commands += fmt.Sprintf("echo %s\necho \"%s\"\n%s\n", metrics.Separator, line, line)
	}

	// Close the commands block
	commands += fmt.Sprintf("echo %s\n", metrics.CollectionEnd)

	// Template for the launcher with interactive mode, if desired
	launcherTemplate := fmt.Sprintf("%s\n%s\n%s", prefix, commands, metrics.Interactive(spec.Spec.Logging.Interactive))

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
