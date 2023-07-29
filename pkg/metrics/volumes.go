/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package metrics

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumeMounts returns read only volume for entrypoint scripts, etc.
func getVolumeMounts(set *api.MetricSet) []corev1.VolumeMount {
	mounts := []corev1.VolumeMount{
		{
			Name:      set.Name,
			MountPath: "/metrics_operator/",
			ReadOnly:  true,
		},
	}
	return mounts
}

// getVolumes for the Indexed Jobs
func getVolumes(set *api.MetricSet, metrics *[]Metric) []corev1.Volume {

	// Runner start scripts
	makeExecutable := int32(0777)

	// Each metric has an entrypoint script
	runnerScripts := []corev1.KeyToPath{}
	for i, _ := range *metrics {
		key := fmt.Sprintf("entrypoint-%d", i)
		runnerScript := corev1.KeyToPath{
			Key:  key,
			Path: key + ".sh",
			Mode: &makeExecutable,
		}
		runnerScripts = append(runnerScripts, runnerScript)
	}

	// TODO will need to add volumes to here for storage requests / metrics
	volumes := []corev1.Volume{
		{
			Name: set.Name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{

					// Namespace based on the cluster
					LocalObjectReference: corev1.LocalObjectReference{
						Name: set.Name,
					},
					Items: runnerScripts,
				},
			},
		},
	}
	return volumes
}
