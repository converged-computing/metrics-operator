/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package metrics

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
)

// getContainers gets containers for a Hyperqueue node
func getContainers(set *api.MetricSet, metrics *[]Metric) ([]corev1.Container, error) {

	// Assume we can pull once for now, this could be changed to allow
	// corev2.PullAlways
	pullPolicy := corev1.PullIfNotPresent
	containers := []corev1.Container{}

	// Currently we share the same mounts across containers, makes life easier!
	mounts := getVolumeMounts(set)

	// Create one container per metric!
	// Each needs to have the sys trace capability to see the application pids
	for i, m := range *metrics {
		script := fmt.Sprintf("/metrics_operator/entrypoint-%d.sh", i)
		command := []string{"/bin/bash", script}

		// TODO specify container resources here?

		// Assemble the container for the node
		// Name the container by the metric for now
		newContainer := corev1.Container{
			Name:            m.Name(),
			Image:           m.Image(),
			ImagePullPolicy: pullPolicy,
			VolumeMounts:    mounts,
			Stdin:           true,
			TTY:             true,
			Command:         command,
			SecurityContext: &corev1.SecurityContext{
				Capabilities: &corev1.Capabilities{
					Add: []corev1.Capability{"SYS_PTRACE"},
				},
			},
		}
		// Only add the working directory if it's defined
		workdir := m.WorkingDir()
		if workdir != "" {
			newContainer.WorkingDir = workdir
		}

		// Ports and environment
		// TODO this should be added when needed
		ports := []corev1.ContainerPort{}
		envars := []corev1.EnvVar{}
		newContainer.Ports = ports
		newContainer.Env = envars
		containers = append(containers, newContainer)
	}

	// If our metric set has an application, add it last
	if set.HasApplication() {
		command := []string{"/bin/bash", "-c", set.Spec.Application.Entrypoint}
		appContainer := corev1.Container{
			Name:            "app",
			Image:           set.Spec.Application.Image,
			ImagePullPolicy: pullPolicy,
			Stdin:           true,
			TTY:             true,
			Command:         command,
		}
		containers = append(containers, appContainer)
	}
	return containers, nil
}
