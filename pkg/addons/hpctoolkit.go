/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package addons

import (
	"fmt"
	"strings"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"github.com/converged-computing/metrics-operator/pkg/metadata"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// HPCToolkit is an addon that provides a container to collect performance metrics
// Commands to interact with output data
// hpcstruct hpctoolkit-sleep-measurements
// hpcprof hpctoolkit-sleep-measurements
// hpcviewer ./hpctoolkit-lmp-database
const (
	hpctoolkitIdentifier = "perf-hpctoolkit"
)

type HPCToolkit struct {
	SpackView

	// Target is the name of the replicated job to customize entrypoint logic for
	target string

	// Output files
	// This is the main output file, and then the database is this + -database
	output string

	// Run a post analysis with hpcstruct and hpcprof to generate a database
	postAnalysis bool

	// ContainerTarget is the name of the container to add the entrypoint logic to
	containerTarget string
	events          string

	// For mpirun and similar, mpirun needs to wrap hpcrun and the command, e.g.,
	// mpirun <MPI args> hpcrun <hpcrun args> <app> <app args>
	prefix string
}

func (m HPCToolkit) Family() string {
	return AddonFamilyPerformance
}

// AssembleVolumes to provide an empty volume for the application to share
// We also need to provide a config map volume for our container spec
func (m HPCToolkit) AssembleVolumes() []specs.VolumeSpec {
	return m.GetSpackViewVolumes()
}

// Validate we have an executable provided, and args and optional
func (a *HPCToolkit) Validate() bool {
	if a.events == "" {
		logger.Error("The HPCtoolkit application addon requires one or more 'events' for hpcrun (e.g., -e IO).")
		return false
	}
	return true
}

// Set custom options / attributes for the metric
func (a *HPCToolkit) SetOptions(metric *api.MetricAddon, m *api.MetricSet) {

	a.EntrypointPath = "/metrics_operator/hpctoolkit-entrypoint.sh"
	a.image = "ghcr.io/converged-computing/metric-hpctoolkit-view:ubuntu"
	a.SetDefaultOptions(metric)
	a.Mount = "/opt/share"
	a.VolumeName = "hpctoolkit"
	a.output = "hpctoolkit-result"
	a.postAnalysis = true
	a.Identifier = hpctoolkitIdentifier
	a.SpackViewContainer = "hpctoolkit"

	// UseColor set to anything means to use it
	output, ok := metric.Options["output"]
	if ok {
		a.output = output.StrVal
	}
	mount, ok := metric.Options["mount"]
	if ok {
		a.Mount = mount.StrVal
	}
	prefix, ok := metric.Options["prefix"]
	if ok {
		a.prefix = prefix.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		a.workdir = workdir.StrVal
	}
	target, ok := metric.Options["target"]
	if ok {
		a.target = target.StrVal
	}
	ctarget, ok := metric.Options["containerTarget"]
	if ok {
		a.containerTarget = ctarget.StrVal
	}
	events, ok := metric.Options["events"]
	if ok {
		a.events = events.StrVal
	}
	image, ok := metric.Options["image"]
	if ok {
		a.image = image.StrVal
	}
	// This will work via a ssh command
	postAnalysis, ok := metric.Options["postAnalysis"]
	if ok {
		if postAnalysis.StrVal == "no" || postAnalysis.StrVal == "false" {
			a.postAnalysis = false
		}
	}
}

// Exported options and list options
func (a *HPCToolkit) Options() map[string]intstr.IntOrString {
	options := a.DefaultOptions()
	options["events"] = intstr.FromString(a.events)
	options["mount"] = intstr.FromString(a.Mount)
	options["prefix"] = intstr.FromString(a.prefix)
	return options
}

// CustomizeEntrypoint scripts
func (a *HPCToolkit) CustomizeEntrypoints(
	cs []*specs.ContainerSpec,
	rjs []*jobset.ReplicatedJob,
) {
	for _, rj := range rjs {

		// Only customize if the replicated job name matches the target
		if a.target != "" && a.target != rj.Name {
			continue
		}
		a.customizeEntrypoint(cs, rj)
	}

}

// CustomizeEntrypoint for a single replicated job
func (a *HPCToolkit) customizeEntrypoint(
	cs []*specs.ContainerSpec,
	rj *jobset.ReplicatedJob,
) {

	// Generate addon metadata
	meta := Metadata(a)

	// This should be run after the pre block of the script
	preBlock := `
echo "%s"
# Ensure hpcrun and software exists. This is rough, but should be OK with enough wait time
wget -q https://github.com/converged-computing/goshare/releases/download/2023-09-06/wait-fs
chmod +x ./wait-fs
mv ./wait-fs /usr/bin/goshare-wait-fs
	
# Ensure spack view is on the path, wherever it is mounted
viewbase="%s"
software="${viewbase}/software"
viewbin="${viewbase}/view/bin"
hpcrunpath=${viewbin}/hpcrun

# Important to add AFTER in case software in container duplicated
export PATH=$PATH:${viewbin}
	
# Wait for software directory, and give it time
goshare-wait-fs -p ${software}
	
# Wait for copy to finish
sleep 10
	
# Copy mount software to /opt/software
cp -R %s/software /opt/software
	
# Wait for hpcrun and marker to indicate copy is done
goshare-wait-fs -p ${viewbin}/hpcrun
goshare-wait-fs -p ${viewbase}/metrics-operator-done.txt

# A small extra wait time to be conservative
sleep 5

# This will work with capability SYS_ADMIN added.
# It will only work with privileged set to true AT YOUR OWN RISK!
echo "-1" | tee /proc/sys/kernel/perf_event_paranoid

# The output path for the analysis
output="%s"

# Run hpcrun. See options with hpcrun -L
events="%s"

# Write a script to run for the post block analysis
here=$(pwd)
cat <<EOF > ./post-run.sh
#!/bin/bash
# Input path should be consistent between nodes
cd ${here}
${viewbin}/hpcstruct ${output}
${viewbin}/hpcprof -o ${output}-database ${output}
EOF
chmod +x ./post-run.sh

echo "%s"
echo "%s"
`
	preBlock = fmt.Sprintf(
		preBlock,
		meta,
		a.Mount,
		a.Mount,
		a.output,
		a.events,
		metadata.CollectionStart,
		metadata.Separator,
	)

	// postBlock to possibly run the hpcstruct command should come right after
	postBlock := ""
	if a.postAnalysis {
		postBlock = `
for host in $(cat ./hostlist.txt); do
    echo "Running post analysis for host ${host}"
    if [[ "$host" == "$(hostname)" ]]; then
    	bash ./post-run.sh
    else
        ssh ${host} ${workdir}/post-run.sh
    fi
done
echo "METRICS-OPERATOR HPCTOOLKIT Post analysis done."
`
	}

	// Add the working directory, if defined
	if a.workdir != "" {
		preBlock += fmt.Sprintf(`
workdir="%s"
echo "Changing directory to ${workdir}"
cd ${workdir}			
`, a.workdir)
	}

	// We use container names to target specific entrypoint scripts here
	for _, containerSpec := range cs {

		// First check - is this the right replicated job?
		if containerSpec.JobName != rj.Name {
			continue
		}

		// Always copy over the pre block - we need the logic to copy software
		containerSpec.EntrypointScript.Pre += "\n" + preBlock

		// Next check if we have a target set (for the container)
		if a.containerTarget != "" && containerSpec.Name != "" && a.containerTarget != containerSpec.Name {
			continue
		}

		// If the post command ends with sleep infinity, tweak it
		isInteractive, updatedPost := deriveUpdatedPost(containerSpec.EntrypointScript.Post)
		containerSpec.EntrypointScript.Post = updatedPost

		// The post to run the command across nodes (when the application finishes)
		containerSpec.EntrypointScript.Post = containerSpec.EntrypointScript.Post + "\n" + postBlock
		containerSpec.EntrypointScript.Command = fmt.Sprintf("%s $hpcrunpath -o $output $events %s", a.prefix, containerSpec.EntrypointScript.Command)

		// If is interactive, add back sleep infinity
		if isInteractive {
			containerSpec.EntrypointScript.Post += "\nsleep infinity\n"
		}
	}
}

// update a post command to not end in sleep
func deriveUpdatedPost(post string) (bool, string) {
	if strings.HasSuffix(post, "sleep infinity\n") {
		updated := strings.Split(post, "\n")
		// This is actually two lines
		updated = updated[:len(updated)-2]
		return true, strings.Join(updated, "\n")
	}
	return false, post
}

func init() {
	base := AddonBase{
		Identifier: hpctoolkitIdentifier,
		Summary:    "performance tools for measurement and analysis",
	}
	app := ApplicationAddon{AddonBase: base}
	spack := SpackView{ApplicationAddon: app}
	toolkit := HPCToolkit{SpackView: spack}
	Register(&toolkit)
}
