/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package specs

import (
	"path/filepath"
	"strings"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// Specs are used to generate configurations for containers and volumes of
// a jobset before we finalize their creation

// A ContainerSpec is used by a metric to define a container
// Job name and container name allow us to associate a script with a replicated job
type ContainerSpec struct {
	JobName          string
	Image            string
	Name             string
	WorkingDir       string
	EntrypointScript EntrypointScript
	Resources        *api.ContainerResources
	Attributes       *api.ContainerSpec
}

// VolumeSpec includes one or more volumes and mount, etc. location
type VolumeSpec struct {
	Volume   corev1.Volume
	ReadOnly bool
	Path     string
	Mount    bool
}

// Named entrypoint script for a container
type EntrypointScript struct {
	Name   string
	Path   string
	Script string

	// Pre block typically includes metadata
	Pre string

	// This is the main command, provided in case we need to wrap it in something
	Command string

	// Anything after the command!
	Post string
}

// Given a full path, derive the key from the script name minus the extension
func DeriveScriptKey(path string) string {

	// Basename
	path = filepath.Base(path)

	// Remove the extension, and this assumes we don't have double .
	return strings.Split(path, ".")[0]
}
