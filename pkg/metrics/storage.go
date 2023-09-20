/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

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
