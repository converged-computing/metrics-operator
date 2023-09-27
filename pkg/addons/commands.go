/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package addons

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

const (
	commandsName     = "commands"
	perfCommandsName = "perf-commands"
)

// Perf addon expects the same command structure, but adds sys caps for trace and admin
type PerfAddon struct {
	CommandAddon
}

// CustomizeEntrypoint scripts
func (a *PerfAddon) CustomizeEntrypoints(
	cs []*specs.ContainerSpec,
	rjs []*jobset.ReplicatedJob,
) {
	for _, rj := range rjs {

		// Only customize if the replicated job name matches the target
		if a.target != "" && a.target != rj.Name {
			continue
		}
		a.customizeEntrypoint(cs, rj)
		a.addContainerCaps(cs, rj)
	}
}

func (a *PerfAddon) SetOptions(addon *api.MetricAddon, metric *api.MetricSet) {
	a.Identifier = perfCommandsName
	a.SetSharedCommandOptions(addon)
}

// addContainerCaps adds capabilities to a container spec
func (a *PerfAddon) addContainerCaps(
	cs []*specs.ContainerSpec,
	rj *jobset.ReplicatedJob,
) {

	// We use container names to target specific entrypoint scripts here
	for _, containerSpec := range cs {

		// Is this the right replicated job?
		if !a.isSelected(containerSpec, rj) {
			continue
		}

		// Always copy over the pre block - we need the logic to copy software
		containerSpec.Attributes.SecurityContext.AllowAdmin = true
		containerSpec.Attributes.SecurityContext.AllowPtrace = true
	}
}

// Command addons primarily edit the entrypoint commands
type CommandAddon struct {
	AddonBase

	// preBlock is run right before the command
	preBlock string

	// prefix is added as a prefix to the command
	prefix string

	// add a suffix to the command
	suffix string

	// postBlock is run after
	postBlock string

	// job name and container name targets
	target          string
	containerTarget string
}

// Doesn't make sense to have an empty command prefix / pre and post!
func (a *CommandAddon) Validate() bool {
	if a.preBlock == "" && a.prefix == "" && a.postBlock == "" && a.suffix == "" {
		logger.Error("The command addon requires one of a 'prefix', 'preBlock', 'postBlock' or 'suffix'")
		return false
	}
	return true
}

// Application family for now...
func (m CommandAddon) Family() string {
	return AddonFamilyApplication
}

func (a *CommandAddon) SetOptions(addon *api.MetricAddon, metric *api.MetricSet) {
	a.Identifier = commandsName
	a.SetSharedCommandOptions(addon)
}

// Set custom options / attributes for the metric
func (a *CommandAddon) SetSharedCommandOptions(metric *api.MetricAddon) {
	target, ok := metric.Options["target"]
	if ok {
		a.target = target.StrVal
	}
	ctarget, ok := metric.Options["containerTarget"]
	if ok {
		a.containerTarget = ctarget.StrVal
	}
	prefix, ok := metric.Options["prefix"]
	if ok {
		a.prefix = prefix.StrVal
	}
	suffix, ok := metric.Options["suffix"]
	if ok {
		a.suffix = suffix.StrVal
	}
	preBlock, ok := metric.Options["preBlock"]
	if ok {
		a.preBlock = preBlock.StrVal
	}
	postBlock, ok := metric.Options["postBlock"]
	if ok {
		a.postBlock = postBlock.StrVal
	}
}

// Underlying function that can be shared
func (a *CommandAddon) DefaultOptions() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"target":          intstr.FromString(a.target),
		"prefix":          intstr.FromString(a.prefix),
		"suffix":          intstr.FromString(a.suffix),
		"preBlock":        intstr.FromString(a.preBlock),
		"postBlock":       intstr.FromString(a.postBlock),
		"containerTarget": intstr.FromString(a.containerTarget),
	}
}

// CustomizeEntrypoint scripts
func (a *CommandAddon) CustomizeEntrypoints(
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

// isSelected determines if a container spec is targeted based on the container and rj names
func (a *CommandAddon) isSelected(
	cs *specs.ContainerSpec,
	rj *jobset.ReplicatedJob,
) bool {
	if cs.JobName != rj.Name {
		return false
	}

	// Next check if we have a target set (for the container)
	if a.containerTarget != "" && cs.Name != "" && a.containerTarget != cs.Name {
		return false
	}
	return true
}

// CustomizeEntrypoint for a single replicated job
func (a *CommandAddon) customizeEntrypoint(
	cs []*specs.ContainerSpec,
	rj *jobset.ReplicatedJob,
) {

	// Generate addon metadata
	meta := Metadata(a)

	// This should be run after the pre block of the script, and includes our preblock
	preBlock := `
echo "%s"
%s
`
	preBlock = fmt.Sprintf(preBlock, meta, a.preBlock)

	// postBlock to possibly run the hpcstruct command should come right after
	postBlock := fmt.Sprintf("\n%s", a.postBlock)

	// We use container names to target specific entrypoint scripts here
	for _, containerSpec := range cs {

		// Is this the right replicated job?
		if !a.isSelected(containerSpec, rj) {
			continue
		}

		// Always copy over the pre block - we need the logic to copy software
		containerSpec.EntrypointScript.Pre += "\n" + preBlock

		// If the post command ends with sleep infinity, tweak it
		isInteractive, updatedPost := deriveUpdatedPost(containerSpec.EntrypointScript.Post)
		containerSpec.EntrypointScript.Post = updatedPost

		// The post to run the command across nodes (when the application finishes)
		containerSpec.EntrypointScript.Post = containerSpec.EntrypointScript.Post + "\n" + postBlock
		containerSpec.EntrypointScript.Command = fmt.Sprintf("%s %s %s", a.prefix, containerSpec.EntrypointScript.Command, a.suffix)

		// If is interactive, add back sleep infinity
		if isInteractive {
			containerSpec.EntrypointScript.Post += "\nsleep infinity\n"
		}
	}
}

func init() {

	// Config map volume type
	base := AddonBase{
		Identifier: commandsName,
		Summary:    "customize a metric's entrypoints",
	}
	app := CommandAddon{AddonBase: base}
	Register(&app)

	base = AddonBase{
		Identifier: perfCommandsName,
		Summary:    "customize a metric's entrypoints expecting performance tracing (adding ptrace and admin caps)",
	}
	cmd := CommandAddon{AddonBase: base}
	perf := PerfAddon{CommandAddon: cmd}
	Register(&perf)
}
