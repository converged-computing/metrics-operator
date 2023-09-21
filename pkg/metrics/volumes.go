/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"path/filepath"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	corev1 "k8s.io/api/core/v1"
)

var (
	makeExecutable = int32(0777)
)

// GetVolumeMounts returns read only volume for entrypoint scripts, etc.
func getVolumeMounts(
	set *api.MetricSet,
	volumes []specs.VolumeSpec,
) []corev1.VolumeMount {

	// This is for the core entrypoints (that are generated as config maps here)
	mounts := []corev1.VolumeMount{
		{
			Name:      set.Name,
			MountPath: "/metrics_operator/",
			ReadOnly:  true,
		},
	}

	// This is for any extra or special entrypoints
	for _, vs := range volumes {

		// Is this volume indicated for mount?
		if vs.Mount {
			mount := corev1.VolumeMount{
				Name:      vs.Volume.Name,
				MountPath: vs.Path,
				ReadOnly:  vs.ReadOnly,
			}
			mounts = append(mounts, mount)
		}
	}
	return mounts
}

// Get MetricsKeyToPath assumes we have a predictible listing of metrics
// scripts. This is applicable for storage and application metrics
func generateOperatorItems(containerSpecs []*specs.ContainerSpec) []corev1.KeyToPath {
	// Each metric has an entrypoint script
	runnerScripts := []corev1.KeyToPath{}
	for _, cs := range containerSpecs {

		// This is relative to the directory
		path := filepath.Base(cs.EntrypointScript.Path)
		runnerScript := corev1.KeyToPath{
			Key:  cs.EntrypointScript.Name,
			Path: path,
			Mode: &makeExecutable,
		}
		runnerScripts = append(runnerScripts, runnerScript)
	}
	return runnerScripts
}

// Add extra config maps to the metrics_operator set from addons
// These are distinct because the operator needs to create them too
func getExtraConfigmaps(volumes []specs.VolumeSpec) []corev1.KeyToPath {

	// Each metric has an entrypoint script
	runnerScripts := []corev1.KeyToPath{}

	for _, addedVolume := range volumes {

		// Check that the typs is config map
		if addedVolume.Volume.ConfigMap == nil {
			continue
		}
		// This will error if it's not a config map :)
		if addedVolume.Volume.Name == "" {
			for _, item := range addedVolume.Volume.ConfigMap.Items {
				runnerScripts = append(runnerScripts, item)
			}
		}
	}
	return runnerScripts
}

// getVolumes adds expected entrypoints along with added volumes (storage or applications)
// This function is intended for a set with a listing of metrics
func getReplicatedJobVolumes(
	set *api.MetricSet,
	cs []*specs.ContainerSpec,
	addedVolumes []specs.VolumeSpec,
) []corev1.Volume {

	// These are for the main entrypoints in /metrics_operator
	runnerScripts := generateOperatorItems(cs)

	// Any volumes that don't have a Name in added we need to generate under the operator
	extraCMs := getExtraConfigmaps(addedVolumes)
	runnerScripts = append(runnerScripts, extraCMs...)

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
	existingVolumes := getAddonVolumes(addedVolumes)
	volumes = append(volumes, existingVolumes...)
	return volumes
}

// Get Addon Volumes for the cluster. This can include:
func getAddonVolumes(vs []specs.VolumeSpec) []corev1.Volume {
	volumes := []corev1.Volume{}
	for _, volume := range vs {
		// If the volume doesn't have a name, it was added to the metrics_operator namespace
		if volume.Volume.Name == "" {
			continue
		}
		logger.Infof("Adding volume %s\n", &volume.Volume)
		volumes = append(volumes, volume.Volume)
	}
	logger.Infof("Volumes %s\n", volumes)
	return volumes
}
