package io

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// sysstat provides a tool "iostat" to assess a storage mount
// https://github.com/sysstat/sysstat

type IOStat struct {
	name                string
	rate                int32
	completions         int32
	description         string
	container           string
	requiresApplication bool
	requiresStorage     bool
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
func (m IOStat) EntrypointScript(set *api.MetricSet) string {

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
	// NOTE: the entrypoint is the entrypoint for the container, while
	// the command is expected to be what we are monitoring. Often
	// they are the same thing.
	return fmt.Sprintf(template, m.completions, m.rate)
}

// Does the metric require an application container?
func (m IOStat) RequiresApplication() bool {
	return m.requiresApplication
}
func (m IOStat) RequiresStorage() bool {
	return m.requiresStorage
}

func init() {
	metrics.Register(
		&IOStat{
			name:                "io-sysstat",
			description:         "statistics for Linux tasks (processes) : I/O, CPU, memory, etc.",
			requiresApplication: false,
			requiresStorage:     true,
			container:           "ghcr.io/converged-computing/metric-sysstat:latest",
		})
}
