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

	"github.com/converged-computing/metrics-operator/pkg/metadata"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/converged-computing/metrics-operator/pkg/specs"
)

// https://www.netlib.org/benchmark/hpl/
// https://ulhpc-tutorials.readthedocs.io/en/production/parallel/mpi/HPL/

// Default input file hpl.dat
// The output of this is Ns, memory is in GiB
// -m 128 -NB 192 -r 0.3 -N 2: translates to --mem 128 -NB ${blocksize} -r 0.3 -N ${pods}
var (
	inputData = `HPLinpack benchmark input file
Innovative Computing Laboratory, University of Tennessee
HPL.out      output file name (if any)
6           device out (6=stdout,7=stderr,file)
1            # of problems sizes (N)
${size}        Ns
1            # of NBs
${blocksize}              NBs
${row_or_colmajor_pmapping}            PMAP process mapping (0=Row-,1=Column-major)
1            # of process grids (P x Q)
${tasks}    Ps  PxQ must equal nprocs
${pods}      Qs
16.0         threshold
1            # of panel fact
${pfact}            PFACTs (0=left, 1=Crout, 2=Right)
1            # of recursive stopping criterium
${nbmin}            NBMINs (>= 1)
1            # of panels in recursion
${ndiv}            NDIVs
1            # of recursive panel fact.
${rfact}            RFACTs (0=left, 1=Crout, 2=Right)
1            # of broadcast
${bcast}            BCASTs (0=1rg,1=1rM,2=2rg,3=2rM,4=Lng,5=LnM)
1            # of lookahead depth
${depth}            DEPTHs (>=0)
${swap}            SWAP (0=bin-exch,1=long,2=mix)
${swapping_threshold}           swapping threshold (default had 64)
${L1_transposed}            L1 in (0=transposed,1=no-transposed) form
${U_transposed}            U  in (0=transposed,1=no-transposed) form
1            Equilibration (0=no,1=yes)
${mem_alignment}            memory alignment in double (> 0) (4,8,16)
`
)

type HPL struct {
	metrics.LauncherWorker

	// Custom Options
	mpiargs string
	tasks   int32

	// r, alpha or beta is a target memory occupation.
	// to get the highest performance we want this to be close to 100% (because large matrices have a lower communication/compute ratio)
	// But that can take very long, so setting a lower ratio is often useful
	ratio  string
	memory int32

	// data inputs
	// note that size is calculated from compute_N
	// blocksize is the NBs "number blocks" value
	blocksize int32
	pfact     int32
	nbmin     int32
	ndiv      int32

	// PMAP process mapping (0=Row-,1=Column-major)
	row_or_colmajor_pmapping int32

	// (0=left, 1=Crout, 2=Right)
	rfact int32

	// (0=1rg,1=1rM,2=2rg,3=2rM,4=Lng,5=LnM)
	bcast int32

	// number of lookahead depth
	depth int32

	// (0=bin-exch,1=long,2=mix)
	swap int32

	// swapping_threshold
	swappingThreshold int32

	// (0=transposed,1=no-transposed)
	l1tranposed int32

	// (0=transposed,1=no-transposed)
	utransposed int32

	// memory alignment in double (> 0) (4,8,16)
	memAlignment int32
}

// I think this is a simulation?
func (m HPL) Family() string {
	return metrics.SolverFamily
}

func (m HPL) Url() string {
	return "https://www.netlib.org/benchmark/hpl/"
}

// Set custom options / attributes for the metric
func (m *HPL) SetOptions(metric *api.Metric) {
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Defaults for hpl.dat values.
	// memory and pods (nodes) calculated on the fly, unless otherwise provided
	m.ratio = "0.3"

	// data inputs - note we don't validate these. We trust the user (dangerous...)
	// note that size is calculated from compute_N
	// blocksize is the NBs "number blocks" value, should be between 1 and 64?
	m.blocksize = 1

	//  PMAP process mapping (0=Row-,1=Column-major)
	m.row_or_colmajor_pmapping = 0

	// PFACTs (0=left, 1=Crout, 2=Right)
	m.pfact = 0

	// NBMINs (>= 1), maybe 1 to 4?
	m.nbmin = 1

	// NDIVs
	m.ndiv = 2

	// (0=left, 1=Crout, 2=Right)
	m.rfact = 0

	// BCASTs (0=1rg,1=1rM,2=2rg,3=2rM,4=Lng,5=LnM)
	m.bcast = 0

	// DEPTHs (>=0)
	m.depth = 0

	// SWAP (0=bin-exch,1=long,2=mix)
	m.swap = 0

	// swapping threshold (e.g., 64, 128)
	m.swappingThreshold = 64

	// (0=transposed,1=no-transposed)
	m.l1tranposed = 0
	m.utransposed = 0

	// memory alignment in double (> 0) (4,8,16)
	m.memAlignment = 4

	args, ok := metric.Options["mpiargs"]
	if ok {
		m.mpiargs = args.StrVal
	}
	tasks, ok := metric.Options["tasks"]
	if ok {
		m.tasks = tasks.IntVal
	}
	// paramters for compute_N
	value, ok := metric.Options["ratio"]
	if ok {
		m.ratio = value.StrVal
	}

	// parameters for hpl.dat
	value, ok = metric.Options["blocksize"]
	if ok {
		m.blocksize = value.IntVal
	}
	value, ok = metric.Options["row_or_colmajor_pmapping"]
	if ok {
		m.row_or_colmajor_pmapping = value.IntVal
	}
	value, ok = metric.Options["pfact"]
	if ok {
		m.pfact = value.IntVal
	}
	value, ok = metric.Options["nbmin"]
	if ok {
		m.nbmin = value.IntVal
	}
	value, ok = metric.Options["ndiv"]
	if ok {
		m.ndiv = value.IntVal
	}
	value, ok = metric.Options["rfact"]
	if ok {
		m.rfact = value.IntVal
	}
	value, ok = metric.Options["bcast"]
	if ok {
		m.bcast = value.IntVal
	}
	value, ok = metric.Options["depth"]
	if ok {
		m.depth = value.IntVal
	}
	value, ok = metric.Options["swap"]
	if ok {
		m.swap = value.IntVal
	}
	value, ok = metric.Options["swappingThreshold"]
	if ok {
		m.swappingThreshold = value.IntVal
	}
	value, ok = metric.Options["l1transposed"]
	if ok {
		m.l1tranposed = value.IntVal
	}
	value, ok = metric.Options["utransposed"]
	if ok {
		m.utransposed = value.IntVal
	}
	value, ok = metric.Options["memAlignment"]
	if ok {
		m.memAlignment = value.IntVal
	}
}

// Exported options and list options
func (m HPL) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"ratio":              intstr.FromString(m.mpiargs),
		"mpiargs":            intstr.FromString(m.mpiargs),
		"tasks":              intstr.FromInt(int(m.tasks)),
		"workdir":            intstr.FromString(m.Workdir),
		"blocksize":          intstr.FromInt(int(m.blocksize)),
		"pfact":              intstr.FromInt(int(m.pfact)),
		"nbmin":              intstr.FromInt(int(m.nbmin)),
		"ndiv":               intstr.FromInt(int(m.ndiv)),
		"rfact":              intstr.FromInt(int(m.rfact)),
		"bcast":              intstr.FromInt(int(m.bcast)),
		"depth":              intstr.FromInt(int(m.depth)),
		"swap":               intstr.FromInt(int(m.swap)),
		"memory":             intstr.FromInt(int(m.memory)),
		"swappableThreshold": intstr.FromInt(int(m.swappingThreshold)),
		"l1transposed":       intstr.FromInt(int(m.l1tranposed)),
		"utransposed":        intstr.FromInt(int(m.utransposed)),
		"memAlignment":       intstr.FromInt(int(m.memAlignment)),
	}
}

func (m HPL) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)
	hosts := m.GetHostlist(spec)
	prefix := m.GetCommonPrefix(meta, "", hosts)

	// Memory command since could mess up templating
	memoryCmd := `awk '/MemFree/ { printf "%.3f \n", $2/1024/1024 }' /proc/meminfo`

	preBlock := `
# Source spack environment
. /opt/spack-environment/activate.sh
		
# Calculate memory, if not defined
memory=%d
if [[ $memory -eq 0 ]]; then
	memory=$(%s)
fi
		
echo "Memory is ${memory}"
		
np=%d
pods=%d
# Tasks per node, not total
tasks=$(nproc)
if [[ $np -eq 0 ]]; then
	np=$(( $pods*$tasks ))
fi
		
echo "Number of tasks (nproc on one node) is $tasks"
echo "Number of tasks total (across $pods nodes) is $np"
		
blocksize=%d
ratio=%s
		
# This calculates the compute value - retrieved from tutorials in /opt/view/bin
compute_script="compute_N -m ${memory} -NB ${blocksize} -r ${ratio} -N ${pods}"
echo $compute_script
# This is the size, variable "N" in the hpl.dat (not confusing or anything)
size=$(${compute_script})
echo "Compute size is ${size}"
		
# Define rest of envars we need for template
row_or_colmajor_pmapping=%d
pfact=%d
nbmin=%d
ndiv=%d
rfact=%d
bcast=%d
depth=%d
swap=%d
swapping_threshold=%d
L1_transposed=%d
U_transposed=%d
mem_alignment=%d
		
# Write the input file (this parses environment variables too)
cat <<EOF > ./hpl.dat
%s
EOF
		
cp ./hostlist.txt ./hostnames.txt
rm ./hostlist.txt
%s
		
echo "%s"
# This is in /root/hpl/bin/linux/xhpl
`

	postBlock := `
echo "%s"
%s
`
	command := fmt.Sprintf("mpirun --allow-run-as-root --hostfile ./hostlist.txt -np $np %s xhpl", m.mpiargs)
	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	preBlock = prefix + fmt.Sprintf(
		preBlock,
		m.memory,
		memoryCmd,
		m.tasks,
		spec.Spec.Pods,
		m.blocksize,
		m.ratio,
		m.row_or_colmajor_pmapping,
		m.pfact,
		m.nbmin,
		m.ndiv,
		m.rfact,
		m.bcast,
		m.depth,
		m.swap,
		m.swappingThreshold,
		m.l1tranposed,
		m.utransposed,
		m.memAlignment,
		inputData,
		metrics.TemplateConvertHostnames,
		metadata.Separator,
	)
	postBlock = fmt.Sprintf(postBlock, metadata.CollectionEnd, interactive)

	// Entrypoint for the launcher
	launcherEntrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(m.LauncherScript),
		Path:    m.LauncherScript,
		Pre:     preBlock,
		Command: command,
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
		Identifier: "app-hpl",
		Summary:    "High-Performance Linpack (HPL)",
		Container:  "ghcr.io/converged-computing/metric-hpl-spack:latest",
	}
	launcher := metrics.LauncherWorker{BaseMetric: base}
	HPL := HPL{LauncherWorker: launcher}
	metrics.Register(&HPL)
}
