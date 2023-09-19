/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package specs

import (
	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
)

// Specs are used to generate configurations for containers and volumes of
// a jobset before we finalize their creation

// A ContainerSpec is used by a metric to define a container
type ContainerSpec struct {
	Command    []string
	Image      string
	Name       string
	WorkingDir string
	Resources  *api.ContainerResources
	Attributes *api.ContainerSpec
}

// DELETE ME
//type VolumeSpec struct {
//	EntrypointScript string
//	EntrypointName   string
//}

// Named entrypoint script for a container
type EntrypointScript struct {
	Name   string
	Path   string
	Script string
}
