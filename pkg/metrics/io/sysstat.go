/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package io

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// sysstat provides a tool "iostat" to assess a storage mount
// https://github.com/sysstat/sysstat

type IOStat struct {
	name        string
	rate        int32
	completions int32
	description string
	container   string
}

// Name returns the metric name
func (m IOStat) Name() string {
	return m.name
}

// Description returns the metric description
func (m IOStat) Description() string {
	return m.description
}

// Container
func (m IOStat) Image() string {
	return m.container
}

// WorkingDir does not matter
func (m IOStat) WorkingDir() string {
	return ""
}

// Validation
func (m IOStat) Validate(set *api.MetricSet) bool {
	return true
}

// Set custom options / attributes for the metric
func (m *IOStat) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions
}

// Generate the replicated job for measuring the application
// We provide the entire Metrics Set (including the application) if we need
// to extract metadata from elsewhere
// TODO need to think of more clever way to export the values?
// Save to somewhere?
// TODO if the app is too fast we might miss it?
func (m IOStat) EntrypointScripts(set *api.MetricSet) []metrics.EntrypointScript {

	template := `#!/bin/bash
i=0
completions=%d
while true
  do
    echo "IOSTAT TIMEPOINT ${i}"
    iostat 
	# Note we can do iostat -o JSON
	if [[ $completions -ne 0 ]] && [[ $i -eq $completions ]]; then
    	exit 0
    fi
	sleep %d
	let i=i+1 
done
`
	// The entrypoint is the entrypoint for the container, while
	// the command is expected to be what we are monitoring. Often
	// they are the same thing. We return an empty Name so it's automatically
	// assigned
	return []metrics.EntrypointScript{
		{Script: fmt.Sprintf(template, m.completions, m.rate)},
	}

}

// Jobs required for success condition (n is the netmark run)
func (m IOStat) SuccessJobs() []string {
	return []string{}
}

func (m IOStat) Type() string {
	return metrics.StorageMetric
}
func (m IOStat) ReplicatedJobs(set *api.MetricSet) ([]jobset.ReplicatedJob, error) {
	return []jobset.ReplicatedJob{}, nil
}

func init() {
	metrics.Register(
		&IOStat{
			name:        "io-sysstat",
			description: "statistics for Linux tasks (processes) : I/O, CPU, memory, etc.",
			container:   "ghcr.io/converged-computing/metric-sysstat:latest",
		})
}
