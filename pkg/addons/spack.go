/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package addons

import (
	"fmt"
	"path/filepath"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	corev1 "k8s.io/api/core/v1"
)

// A spack view expects to copy a view from /opt/view into a mount
// This is a virtual struct in that it just provides shared functions for others
type SpackView struct {
	ApplicationAddon

	Setup              string
	SpackViewContainer string
	VolumeName         string
	EntrypointPath     string
	Mount              string
}

// Generate a container spec that will map to a listing of containers for the replicated job
func (a *SpackView) AssembleContainers() []specs.ContainerSpec {

	// The entrypoint script
	// This is the addon container entrypoint, we don't care about metadata here
	// The sole purpose is just to provide the volume, meaning copying content there
	template := `#!/bin/bash

# Extra setup (optional) for a spack view
%s

echo "Moving content from /opt/view to be in shared volume at %s"
view=$(ls /opt/views/._view/)
view="/opt/views/._view/${view}"

# Give a little extra wait time
sleep 10

viewroot="%s"
mkdir -p $viewroot/view
# We have to move both of these paths, *sigh*
cp -R ${view}/* $viewroot/view
cp -R /opt/software $viewroot/

# This is a marker to indicate the copy is done
touch $viewroot/metrics-operator-done.txt

# Sleep forever, the application needs to run and end
echo "Sleeping forever so %s can be shared and use for %s."
sleep infinity
`
	script := fmt.Sprintf(
		template,
		a.Setup,
		a.Mount,
		a.Mount,
		a.Mount,
		a.Identifier,
	)

	// Leave the name empty to generate in the namespace of the metric set (e.g., set.Name)
	entrypoint := specs.EntrypointScript{
		Name:   a.VolumeName,
		Path:   a.EntrypointPath,
		Script: filepath.Base(a.EntrypointPath),
		Pre:    script,
	}

	// The resource spec and attributes for now are empty (might redo this design)
	// We assume they inherit the resources / attributes of the pod for now
	// We don't use JobName here because we don't associate addon containers
	// with other addon entrypoints
	return []specs.ContainerSpec{
		{
			Image:            a.image,
			Name:             a.SpackViewContainer,
			EntrypointScript: entrypoint,
			Resources:        &api.ContainerResources{},
			Attributes: &api.ContainerSpec{
				SecurityContext: api.SecurityContext{
					Privileged: a.privileged,
				},
			},
			// We need to write this config map!
			NeedsWrite: true,
		},
	}
}

// AssembleVolumes to provide an empty volume for the application to share
// We also need to provide a config map volume for our container spec
func (m *SpackView) GetSpackViewVolumes() []specs.VolumeSpec {

	volume := corev1.Volume{
		Name: m.VolumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}

	// Prepare items as key to path
	items := []corev1.KeyToPath{
		{
			Key:  m.VolumeName,
			Path: filepath.Base(m.EntrypointPath),
		},
	}

	// This is a config map volume with items
	// It needs to be created in the same metrics operator namespace
	// Thus we only need the items!
	configVolume := corev1.Volume{
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				Items: items,
			},
		},
	}

	// EmptyDir should be ReadOnly False, and we don't need a mount for it
	return []specs.VolumeSpec{
		{
			Volume: volume,
			Mount:  true,
			Path:   m.Mount,
		},

		// Mount is set to false here because we mount via metrics_operator
		{
			Volume:   configVolume,
			ReadOnly: true,
			Mount:    false,
			Path:     filepath.Dir(m.EntrypointPath),
		},
	}
}
