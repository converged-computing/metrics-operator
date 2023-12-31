/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	corev1 "k8s.io/api/core/v1"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"github.com/converged-computing/metrics-operator/pkg/specs"
)

// Security context defaults
var (
	capAdmin  = corev1.Capability("SYS_ADMIN")
	capPtrace = corev1.Capability("SYS_PTRACE")
)

// getReplicatedJobContainers gets containers (sidecar and init)
// for the replicated job, also generating needed mounts, etc.
func getReplicatedJobContainers(
	set *api.MetricSet,
	rj *jobset.ReplicatedJob,
	containerSpecs []specs.ContainerSpec,
	volumes []specs.VolumeSpec,
) ([]corev1.Container, []corev1.Container, error) {

	// We only generate containers from specs that match the replicated job name
	containers := []corev1.Container{}
	initContainers := []corev1.Container{}

	// Assume we can pull once for now, this could be changed to allow pull always
	pullPolicy := corev1.PullIfNotPresent

	// Currently we share the same mounts across containers, makes life easier!
	mounts := getVolumeMounts(set, volumes)

	// Keep track of any specs that have privileged, then the app needs it
	hasPrivileged := false

	// Each needs to have the sys trace capability to see the application pids
	for _, cs := range containerSpecs {

		logger.Infof("Checking container spec %s", cs)

		// Skip containers not intended for the replicated job
		if cs.JobName != "" && cs.JobName != rj.Name {
			continue
		}
		hasPrivileged = hasPrivileged || cs.Attributes.SecurityContext.Privileged
		resources, err := getContainerResources(cs.Resources)
		if err != nil {
			return containers, initContainers, err
		}

		// If a command is provided, use it first
		command := []string{"/bin/bash", cs.EntrypointScript.Path}
		if len(cs.Command) > 0 {
			command = cs.Command
		}
		// Create the actual container from the spec
		newContainer := corev1.Container{
			Name:            cs.Name,
			Image:           cs.Image,
			ImagePullPolicy: pullPolicy,
			VolumeMounts:    mounts,
			Stdin:           true,
			TTY:             true,
			Command:         command,
			SecurityContext: &corev1.SecurityContext{
				Privileged: &cs.Attributes.SecurityContext.Privileged,
			},
		}

		// Add capabilities to the security context
		caps := []corev1.Capability{}

		// Should we allow sharing the process namespace?
		if cs.Attributes.SecurityContext.AllowPtrace {
			caps = append(caps, capPtrace)
		}
		if cs.Attributes.SecurityContext.AllowAdmin {
			caps = append(caps, capAdmin)
		}
		newContainer.SecurityContext.Capabilities = &corev1.Capabilities{Add: caps}

		// Only add the working directory if it's defined
		if cs.WorkingDir != "" {
			newContainer.WorkingDir = cs.WorkingDir
		}

		// Ports and environment (add when needed)
		ports := []corev1.ContainerPort{}
		envars := []corev1.EnvVar{}
		newContainer.Ports = ports
		newContainer.Env = envars
		newContainer.Resources = resources

		// Add as an init container, or a sidecar container
		if cs.InitContainer {
			initContainers = append(initContainers, newContainer)

		} else {
			containers = append(containers, newContainer)
		}
	}
	logger.Infof("🟪️ Adding %d init containers\n", len(initContainers))
	logger.Infof("🟪️ Adding %d containers\n", len(containers))
	return containers, initContainers, nil
}
