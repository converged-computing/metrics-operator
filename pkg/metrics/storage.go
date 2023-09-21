/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"github.com/converged-computing/metrics-operator/pkg/specs"
)

// These are common templates for storage apps.
// They define the interface of a Metric.

type StorageGeneric struct {
	BaseMetric
}

// Family returns the storage family
func (m StorageGeneric) Family() string {
	return StorageFamily
}

// By default assume storage does not have sole tenancy
func (m StorageGeneric) HasSoleTenancy() bool {
	return false
}

// StorageContainerSpec gets the storage container spec
// This is identical to the application spec and could be combined
func (m *StorageGeneric) StorageContainerSpec(
	preBlock string,
	command string,
	postBlock string,
) []*specs.ContainerSpec {

	entrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(DefaultEntrypointScript),
		Path:    DefaultEntrypointScript,
		Pre:     preBlock,
		Command: command,
		Post:    postBlock,
	}

	return []*specs.ContainerSpec{{
		JobName:          ReplicatedJobName,
		Image:            m.Image(),
		Name:             "storage",
		WorkingDir:       m.Workdir,
		EntrypointScript: entrypoint,
		Resources:        m.ResourceSpec,
		Attributes:       m.AttributeSpec,
	}}
}
