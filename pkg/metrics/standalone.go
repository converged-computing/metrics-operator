/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// A Standalone metric is typically going to provide its own logic for one or more replicated jobss
func (m *StandaloneMetricSet) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {
	rjs := []jobset.ReplicatedJob{}
	for _, metric := range m.Metrics() {
		jobs, err := GetStandaloneReplicatedJobs(spec, metric, spec.Spec.Application.Volumes)
		if err != nil {
			return rjs, err
		}
		rjs = append(rjs, jobs...)
	}
	return rjs, nil
}

// Create a standalone JobSet, one without volumes or application
// This will be definition be a JobSet for only one metric
func GetStandaloneReplicatedJobs(
	spec *api.MetricSet,
	metric *Metric,
	volumes map[string]api.Volume,
) ([]jobset.ReplicatedJob, error) {

	m := (*metric)

	// Does the metric provide its own logic?
	rjs, err := m.ReplicatedJobs(spec)
	if err != nil {
		return rjs, err
	}

	if len(rjs) == 0 {
		return rjs, fmt.Errorf("custom standalone metrics require a replicated job set")
	}
	return rjs, nil
}
