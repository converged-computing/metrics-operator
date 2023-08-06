/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// GetVolumeMounts returns read only volume for entrypoint scripts, etc.
func getVolumeMounts(
	set *api.MetricSet,
	volumes map[string]api.Volume,
) []corev1.VolumeMount {
	mounts := []corev1.VolumeMount{
		{
			Name:      set.Name,
			MountPath: "/metrics_operator/",
			ReadOnly:  true,
		},
	}

	// Add on application volumes/claims
	for volumeName, volume := range volumes {
		mount := corev1.VolumeMount{
			Name:      volumeName,
			MountPath: volume.Path,
			ReadOnly:  volume.ReadOnly,
		}
		mounts = append(mounts, mount)
	}
	return mounts
}

// Get MetricsKeyToPath assumes we have a predictible listing of metrics
// scripts. This is applicable for storage and application metrics
func GetMetricsKeyToPath(metrics []*Metric) []corev1.KeyToPath {
	// Runner start scripts
	makeExecutable := int32(0777)

	// Each metric has an entrypoint script
	runnerScripts := []corev1.KeyToPath{}
	for i, _ := range metrics {
		key := fmt.Sprintf("entrypoint-%d", i)
		runnerScript := corev1.KeyToPath{
			Key:  key,
			Path: key + ".sh",
			Mode: &makeExecutable,
		}
		runnerScripts = append(runnerScripts, runnerScript)
	}
	return runnerScripts
}

// getVolumes adds expected entrypoints along with addedvolumes (storage or applications)
// This function is intended for a set with a listing of metrics
func GetVolumes(
	set *api.MetricSet,
	runnerScripts []corev1.KeyToPath,
	addedVolumes map[string]api.Volume,
) []corev1.Volume {

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
	existingVolumes := getExistingVolumes(addedVolumes)
	volumes = append(volumes, existingVolumes...)
	return volumes
}

// GetStandaloneVolumes is intended for a single metric, where the volumes
// are provided as custom EntrypointScripts
func GetStandaloneVolumes(
	set *api.MetricSet,
	scripts []EntrypointScript,
	addedVolumes map[string]api.Volume,
) []corev1.Volume {

	// Runner start scripts
	makeExecutable := int32(0777)

	// Each metric has an entrypoint script
	runnerScripts := []corev1.KeyToPath{}
	for i, script := range scripts {
		key := script.Name
		if key == "" {
			key = fmt.Sprintf("entrypoint-%d", i)
		}
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

	existingVolumes := getExistingVolumes(addedVolumes)
	volumes = append(volumes, existingVolumes...)
	return volumes
}

// Get Existing volumes for the cluster. This can include:
// config maps
// secrets
// persistent volumes / claims
func getExistingVolumes(existing map[string]api.Volume) []corev1.Volume {
	volumes := []corev1.Volume{}
	for volumeName, volumeMeta := range existing {

		var newVolume corev1.Volume
		if volumeMeta.HostPath != "" {

			newVolume = corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: volumeMeta.HostPath,
					},
				},
			}

		} else if volumeMeta.SecretName != "" {
			newVolume = corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: volumeMeta.SecretName,
					},
				},
			}

		} else if volumeMeta.ConfigMapName != "" {

			// Prepare items as key to path
			items := []corev1.KeyToPath{}
			for key, path := range volumeMeta.Items {
				newItem := corev1.KeyToPath{
					Key:  key,
					Path: path,
				}
				items = append(items, newItem)
			}

			// This is a config map volume with items
			newVolume = corev1.Volume{
				Name: volumeMeta.ConfigMapName,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: volumeMeta.ConfigMapName,
						},
						Items: items,
					},
				},
			}

		} else {

			// Fall back to persistent volume claim
			newVolume = corev1.Volume{
				Name: volumeName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: volumeMeta.ClaimName,
					},
				},
			}
		}
		volumes = append(volumes, newVolume)
	}
	return volumes
}
