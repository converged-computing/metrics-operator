/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package addons

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"github.com/converged-computing/metrics-operator/pkg/metadata"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// https://github.com/IBM/mpitrace
const (
	mpitraceIdentifier = "perf-mpitrace"
)

type MPITrace struct {
	SpackView

	// Target is the name of the replicated job to customize entrypoint logic for
	target string

	// ContainerTarget is the name of the container to add the entrypoint logic to
	containerTarget string
}

func (m MPITrace) Family() string {
	return AddonFamilyPerformance
}

// AssembleVolumes to provide an empty volume for the application to share
// We also need to provide a config map volume for our container spec
func (m MPITrace) AssembleVolumes() []specs.VolumeSpec {
	return m.GetSpackViewVolumes()
}

// Validate we have an executable provided, and args and optional
func (a *MPITrace) Validate() bool {
	return true
}

// Set custom options / attributes for the metric
func (a *MPITrace) SetOptions(metric *api.MetricAddon, m *api.MetricSet) {

	a.EntrypointPath = "/metrics_operator/mpitrace-entrypoint.sh"
	a.image = "ghcr.io/converged-computing/metric-mpitrace:rocky"
	a.SetDefaultOptions(metric)
	a.Mount = "/opt/share"
	a.VolumeName = "mpitrace"
	a.Identifier = mpitraceIdentifier
	a.SpackViewContainer = "mpitrace"
	a.InitContainer = true

	mount, ok := metric.Options["mount"]
	if ok {
		a.Mount = mount.StrVal
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
	image, ok := metric.Options["image"]
	if ok {
		a.image = image.StrVal
	}
}

// Exported options and list options
func (a *MPITrace) Options() map[string]intstr.IntOrString {
	options := a.DefaultOptions()
	options["mount"] = intstr.FromString(a.Mount)
	return options
}

// CustomizeEntrypoint scripts
func (a *MPITrace) CustomizeEntrypoints(
	cs []*specs.ContainerSpec,
	rjs []*jobset.ReplicatedJob,
) {
	logger.Infof("🟧️ Customizing entrypoints for %s\n", rjs)

	for _, rj := range rjs {
		logger.Infof("🟧️ Comparing job target %s vs job name %s\n", a.target, rj.Name)

		// Only customize if the replicated job name matches the target
		if a.target != "" && a.target != rj.Name {
			continue
		}
		a.customizeEntrypoint(cs, rj)
	}
}

// CustomizeEntrypoint for a single replicated job
func (a *MPITrace) customizeEntrypoint(
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
libmpitraceso=${viewbase}/view/lib/libmpitrace.so

# Important to add AFTER in case software in container duplicated
export PATH=$PATH:${viewbin}
	
# Wait for software directory, and give it time
goshare-wait-fs -p ${software}
	
# Wait for copy to finish
sleep 10
	
# Copy mount software to /opt/software
cp -R %s/software /opt/software
	
# Wait for file indicator that copy is done
goshare-wait-fs -p ${viewbase}/metrics-operator-done.txt

# A small extra wait time to be conservative
sleep 5
echo "%s"
echo "%s"
`
	preBlock = fmt.Sprintf(
		preBlock,
		meta,
		a.Mount,
		a.Mount,
		metadata.CollectionStart,
		metadata.Separator,
	)

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
		containerSpec.EntrypointScript.Command = fmt.Sprintf(
			"export LD_PRELOAD=${libmpitraceso}\n%s\nunset LD_PRELOAD",
			containerSpec.EntrypointScript.Command,
		)

		// If is interactive, add back sleep infinity
		if isInteractive {
			containerSpec.EntrypointScript.Post += "\nsleep infinity\n"
		}
	}
}

func init() {
	base := AddonBase{
		Identifier: mpitraceIdentifier,
		Summary:    "library for measuring communication in distributed-memory parallel applications that use MPI",
	}
	app := ApplicationAddon{AddonBase: base}
	spack := SpackView{ApplicationAddon: app}
	tracer := MPITrace{SpackView: spack}
	Register(&tracer)
}
