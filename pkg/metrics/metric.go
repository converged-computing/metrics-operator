/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package metric

import (
	"fmt"
	"log"
)

var (
	Registry = make(map[string]Metric)
)

// A Metric defines a generic interface for the operator to interact with
// The functionality of different metric types might vary based on the type
// All metrics return a JobSet of some type (and potentially a replicated job)
type Metric interface {
	Description() string
	Name() string
}

// GetMetric returns the Component specified by name from `Registry`.
func GetMetric(name string) (Metric, error) {
	if _, ok := Registry[name]; ok {
		return Registry[name], nil
	}
	return nil, fmt.Errorf("%s is not a registered Metric type", name)
}

// Register a new Metric type, adding it to the Registry
func Register(m Metric) {
	name := m.Name()
	if _, ok := Registry[name]; ok {
		log.Fatalf("Metric: %s has already been added to the registry", name)
	}
	Registry[name] = m
}
