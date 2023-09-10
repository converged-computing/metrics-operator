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
	"k8s.io/apimachinery/pkg/util/intstr"

	jobs "github.com/converged-computing/metrics-operator/pkg/jobs"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// ghcr.io/converged-computing/metric-osu-benchmark:latest
// https://mvapich.cse.ohio-state.edu/benchmarks/

type BenchmarkConfig struct {
	Workdir  string
	Flags    string
	HostFile string
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

		// Single Sided (all require exactly 2 processes)
		"osu_get_acc_latency": {Workdir: singleSidedDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"},
		"osu_acc_latency":     {Workdir: singleSidedDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"}, // Latency Test for Accumulate
		"osu_fop_latency":     {Workdir: singleSidedDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"},
		"osu_get_latency":     {Workdir: singleSidedDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"}, // Latency Test for Get
		"osu_put_latency":     {Workdir: singleSidedDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"}, // Latency Test for Put
		"osu_cas_latency":     {Workdir: singleSidedDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"},
		"osu_get_bw":          {Workdir: singleSidedDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"},
		"osu_put_bibw":        {Workdir: singleSidedDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"},
		"osu_put_bw":          {Workdir: singleSidedDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"},

		// Collective
		// For allreduce this should work, need to test -np $np -map-by ppr:1:node -rank-by core
		"osu_allreduce":      {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"}, // MPI_Allreduce Latency Test
		"osu_allgather":      {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_allgatherv":     {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_alltoall":       {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_alltoallv":      {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_barrier":        {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_bcast":          {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_gather":         {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_gatherv":        {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_iallgather":     {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_iallgatherv":    {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_iallreduce":     {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_ialltoall":      {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_ialltoallv":     {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_ialltoallw":     {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_ibarrier":       {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_ibcast":         {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_igather":        {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_igatherv":       {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_ireduce":        {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_iscatter":       {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_iscatterv":      {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_reduce":         {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_reduce_scatter": {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_scatter":        {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_scatterv":       {Workdir: collectiveDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},

		// Point to Point (commented if requires 2 processes)
		"osu_latency":    {Workdir: pointToPointDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"}, // Latency Test (requires 2)
		"osu_bibw":       {Workdir: pointToPointDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"}, // Bidirectional Bandwidth Test (requires 2)
		"osu_bw":         {Workdir: pointToPointDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"}, // Bandwidth Test (requires 2)
		"osu_latency_mp": {Workdir: pointToPointDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"}, // requires 2
		"osu_latency_mt": {Workdir: pointToPointDir, Flags: "-np 2 -map-by ppr:1:node", HostFile: "./hostlist-pairs.txt"}, // requires 2
		"osu_mbw_mr":     {Workdir: pointToPointDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_multi_lat":  {Workdir: pointToPointDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},

		// Startup
		"osu_hello": {Workdir: startupDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
		"osu_init":  {Workdir: startupDir, Flags: "-np $np -map-by ppr:$tasks:node -rank-by core", HostFile: "./hostlist.txt"},
	}
)

type OSUBenchmark struct {
	jobs.LauncherWorker

	// Custom options
	commands []string
	tasks    int32
	lookup   map[string]bool
	runAll   bool
	flags    string
	timed    bool
}

func (m OSUBenchmark) Url() string {
	return "https://mvapich.cse.ohio-state.edu/benchmarks/"
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
	m.lookup = map[string]bool{}
	m.commands = []string{}
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// One pod per hostname
	m.SoleTenancy = true

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
	st, ok := metric.Options["sole-tenancy"]
	if ok && st.StrVal == "false" || st.StrVal == "no" {
		m.SoleTenancy = false
	}
	runAll, ok := metric.Options["all"]
	if ok && runAll.StrVal == "true" || runAll.StrVal == "yes" {
		m.runAll = true
	}
	timed, ok := metric.Options["timed"]
	if ok && timed.StrVal == "true" || timed.StrVal == "yes" {
		m.timed = true
	}
	flags, ok := metric.Options["flags"]
	if ok {
		m.flags = flags.StrVal
	}

	// If not selected or found, fall back to default list
	if len(m.commands) == 0 {
		for _, command := range osuBenchmarkDefaults {
			if !m.hasCommand(command) {
				m.addCommand(command)
			}
		}
	}

	// Run ALL the benchmarks (living dangerously)!
	if m.runAll {
		for command := range osuBenchmarkCommands {
			if !m.hasCommand(command) {
				m.addCommand(command)
			}
		}
	}
}

// Exported options and list options
func (m OSUBenchmark) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"sole-tenancy": intstr.FromString(fmt.Sprintf("%v", m.SoleTenancy)),
		"tasks":        intstr.FromInt(int(m.tasks)),
		"flags":        intstr.FromString(m.flags),
		"timed":        intstr.FromString(fmt.Sprintf("%v", m.timed)),
		"all":          intstr.FromString(fmt.Sprintf("%v", m.runAll)),
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

// OSU Benchmarks pair to pair must be run with only two nodes
func (m OSUBenchmark) Validate(spec *api.MetricSet) bool {
	if len(m.commands) == 0 {
		fmt.Printf("🟥️ OSUBenchmark not valid, requires 1+ commands.")
		return false
	}
	return true
}

// Family returns the network family
func (n OSUBenchmark) Family() string {
	return metrics.NetworkFamily
}

// Return lookup of entrypoint scripts
func (m OSUBenchmark) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)

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

echo "Number of tasks (nproc on one node) is $tasks"
echo "Number of tasks total (across $pods nodes) is $np"

# Allow network to ready (we need the hostnames / ip addresses to be there)
echo "Sleeping for 60 seconds waiting for network..."
sleep 60

# Write the hosts file.
cat <<EOF > ./hostnames.txt
%s
EOF

%s

# prepare hostlist for pair to pair
cat hostlist.txt | head -2 > ./hostlist-pairs.txt

echo "Hostlist"
cat ./hostlist.txt

echo "Hostlist for Pair to Pair"
cat ./hostlist-pairs.txt

# Show metadata for run
echo "%s"
`
	prefix := fmt.Sprintf(
		prefixTemplate,
		m.tasks,
		spec.Spec.Pods,
		hosts,
		metrics.TemplateConvertHostnames,
		metadata,
	)

	// Do we want timed?
	mpirun := "mpirun"
	if m.timed {
		mpirun = "time mpirun"
	}
	// Prepare list of commands, e.g.,
	// mpirun -f ./hostlist.txt -np 2 ./osu_acc_latency (mpich)
	// mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 ./osu_fop_latency (openmpi)
	// Sleep a little more to allow worker to write launcher hostname
	commands := fmt.Sprintf("\nsleep 5\necho %s\n", metrics.CollectionStart)
	for _, executable := range m.commands {

		workDir := osuBenchmarkCommands[executable].Workdir
		command := path.Join(workDir, executable)

		// Flags can vary by command, or be overridden by the user
		flags := osuBenchmarkCommands[executable].Flags
		if m.flags != "" {
			flags = m.flags
		}
		hostfile := osuBenchmarkCommands[executable].HostFile

		// Some pair to pair is for 2 nodes
		var line string
		if workDir == pointToPointDir || workDir == singleSidedDir {
			line = fmt.Sprintf("%s --hostfile %s --allow-run-as-root %s %s", mpirun, hostfile, flags, command)
		} else {
			line = fmt.Sprintf("%s --hostfile %s --allow-run-as-root %s %s", mpirun, hostfile, flags, command)
		}
		commands += fmt.Sprintf("echo %s\necho \"%s\"\n%s\n", metrics.Separator, line, line)
	}

	// Close the commands block
	commands += fmt.Sprintf("echo %s\n", metrics.CollectionEnd)

	// Template for the launcher with interactive mode, if desired
	launcherTemplate := fmt.Sprintf("%s\n%s\n%s", prefix, commands, metrics.Interactive(spec.Spec.Logging.Interactive))

	// The worker just has sleep infinity added, and getting the ip address of the launcher
	workerTemplate := prefix + "\nsleep infinity"
	return m.FinalizeEntrypoints(launcherTemplate, workerTemplate)
}

func init() {
	launcher := jobs.LauncherWorker{
		Identifier:     "network-osu-benchmark",
		Summary:        "point to point MPI benchmarks",
		Container:      "ghcr.io/converged-computing/metric-osu-benchmark:latest",
		WorkerScript:   "/metrics_operator/osu-worker.sh",
		LauncherScript: "/metrics_operator/osu-launcher.sh",
	}
	osu := OSUBenchmark{LauncherWorker: launcher}
	metrics.Register(&osu)
}
