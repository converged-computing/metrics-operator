/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
)

// Security context defaults
var (
	capAdmin  = corev1.Capability("SYS_ADMIN")
	capPtrace = corev1.Capability("SYS_PTRACE")
)

// A ContainerSpec is used by a metric to define a container
type ContainerSpec struct {
	Command    []string
	Image      string
	Name       string
	WorkingDir string
	Resources  *api.ContainerResources
	Attributes *api.ContainerSpec
}

// Named entrypoint script for a container
type EntrypointScript struct {
	Name   string
	Path   string
	Script string
}

// getContainers gets containers for a set of metrics
func getContainers(
	set *api.MetricSet,
	metrics []*Metric,
	volumes map[string]api.Volume,
) ([]corev1.Container, error) {

	containers := []ContainerSpec{}

	// Create one container per metric!
	// Each needs to have the sys trace capability to see the application pids
	for i, m := range metrics {

		metric := (*m)
		script := fmt.Sprintf("/metrics_operator/entrypoint-%d.sh", i)
		command := []string{"/bin/bash", script}

		newContainer := ContainerSpec{
			Command:    command,
			Image:      metric.Image(),
			WorkingDir: metric.WorkingDir(),
			Name:       metric.Name(),
			Resources:  metric.Resources(),
			Attributes: metric.Attributes(),
		}
		containers = append(containers, newContainer)
	}
	return GetContainers(set, containers, volumes, false, false)
}

// GetContainers based on one or more container specs
func GetContainers(
	set *api.MetricSet,
	specs []ContainerSpec,
	volumes map[string]api.Volume,
	allowPtrace bool,
	allowAdmin bool,
) ([]corev1.Container, error) {

	// Assume we can pull once for now, this could be changed to allow
	// corev2.PullAlways
	pullPolicy := corev1.PullIfNotPresent
	containers := []corev1.Container{}

	// Currently we share the same mounts across containers, makes life easier!
	mounts := getVolumeMounts(set, volumes)

	// Keep track of any specs that have privileged, then the app needs it
	hasPrivileged := false

	// Each needs to have the sys trace capability to see the application pids
	for _, s := range specs {

		hasPrivileged = hasPrivileged || s.Attributes.SecurityContext.Privileged

		// Get resources for container
		resources, err := getContainerResources(s.Resources)
		logger.Info("🌀 Metric", "Container.Resources", resources)
		if err != nil {
			return containers, err
		}

		// Create one container per metric!
		// Name the container by the metric for now
		newContainer := corev1.Container{
			Name:            s.Name,
			Image:           s.Image,
			ImagePullPolicy: pullPolicy,
			VolumeMounts:    mounts,
			Stdin:           true,
			TTY:             true,
			Command:         s.Command,
			SecurityContext: &corev1.SecurityContext{
				Privileged: &s.Attributes.SecurityContext.Privileged,
			},
		}

		// Add capabilities to the security context
		caps := []corev1.Capability{}

		// Should we allow sharing the process namespace?
		if allowPtrace {
			caps = append(caps, capPtrace)
		}
		if allowAdmin {
			caps = append(caps, capAdmin)
		}
		newContainer.SecurityContext.Capabilities = &corev1.Capabilities{Add: caps}

		// Only add the working directory if it's defined
		if s.WorkingDir != "" {
			newContainer.WorkingDir = s.WorkingDir
		}

		// Ports and environment
		// TODO this should be added when needed
		ports := []corev1.ContainerPort{}
		envars := []corev1.EnvVar{}
		newContainer.Ports = ports
		newContainer.Env = envars
		newContainer.Resources = resources
		containers = append(containers, newContainer)
	}

	// If our metric set has an application, add it last
	// We currently accept resources for an application (but not metrics yet)
	if set.HasApplication() {

		// Prepare container resources
		resources, err := getContainerResources(&set.Spec.Application.Resources)
		logger.Info("🌀 Application", "Container.Resources", resources)
		if err != nil {
			return containers, err
		}

		// The application security context can have admin (but should not have the same process sharing)
		securityContext := &corev1.SecurityContext{}
		if allowAdmin {
			securityContext.Capabilities = &corev1.Capabilities{
				Add: []corev1.Capability{capAdmin},
			}
			securityContext.Privileged = &hasPrivileged
		}

		// Minimally this is set.Spec.Application.Entrypoint executed in a bash script
		// But for an application metric with a volume, there can be custom logic
		command := []string{"/bin/bash", DefaultApplicationEntrypoint}
		appContainer := corev1.Container{
			Name:            "app",
			Image:           set.Spec.Application.Image,
			ImagePullPolicy: pullPolicy,
			VolumeMounts:    mounts,
			Stdin:           true,
			TTY:             true,
			Command:         command,
			SecurityContext: securityContext,
		}
		if set.Spec.Application.WorkingDir != "" {
			appContainer.WorkingDir = set.Spec.Application.WorkingDir
		}
		containers = append(containers, appContainer)
	}
	logger.Infof("🟪️ Adding %d containers\n", len(containers))
	return containers, nil
}
