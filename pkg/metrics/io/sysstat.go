/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package io

import (
	"fmt"
	"strconv"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// sysstat provides a tool "iostat" to assess a storage mount
// https://github.com/sysstat/sysstat

type IOStat struct {
	name          string
	rate          int32
	completions   int32
	description   string
	container     string
	humanReadable bool
}

// Name returns the metric name
func (m IOStat) Name() string {
	return m.name
}
func (m IOStat) Url() string {
	return "https://github.com/sysstat/sysstat"
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

	// Does the person want human readable instead of table?
	value, ok := metric.Options["human"]
	if ok {
		if value.StrVal == "true" {
			m.humanReadable = true
		}
	}
}

// Generate the entrypoint for measuring the storage
func (m IOStat) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Prepare metadata for set and separator
	metadata := metrics.Metadata(spec, metric)
	command := "iostat -dxm -o JSON"
	if m.humanReadable {
		command = "iostat -dxm"
	}
	template := `#!/bin/bash
i=0
echo "%s"
completions=%d
echo "%s"
while true
  do
    echo "%s"
	%s
	# Note we can do iostat -o JSON
	if [[ $completions -ne 0 ]] && [[ $i -eq $completions ]]; then
    	echo "%s"
    	exit 0
    fi
	sleep %d
	let i=i+1
done
`
	script := fmt.Sprintf(
		template,
		metadata,
		m.completions,
		metrics.CollectionStart,
		metrics.Separator,
		command,
		metrics.CollectionEnd,
		m.rate,
	)
	// The entrypoint is the entrypoint for the container, while
	// the command is expected to be what we are monitoring. Often
	// they are the same thing. We return an empty Name so it's automatically
	// assigned
	return []metrics.EntrypointScript{
		{Script: script},
	}

}

// Exported options and list options
func (m IOStat) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.rate)),
		"completions": intstr.FromInt(int(m.completions)),
		"human":       intstr.FromString(strconv.FormatBool(m.humanReadable)),
	}
}
func (m IOStat) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
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
