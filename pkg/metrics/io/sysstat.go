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

	"github.com/converged-computing/metrics-operator/pkg/jobs"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// sysstat provides a tool "iostat" to assess a storage mount
// https://github.com/sysstat/sysstat

type IOStat struct {
	jobs.StorageGeneric
	humanReadable bool
	rate          int32
	completions   int32
}

func (m IOStat) Url() string {
	return "https://github.com/sysstat/sysstat"
}

// Set custom options / attributes for the metric
func (m *IOStat) SetOptions(metric *api.Metric) {
	m.rate = 10
	m.completions = 0 // infinite
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Does the person want human readable instead of table?
	value, ok := metric.Options["human"]
	if ok {
		if value.StrVal == "true" {
			m.humanReadable = true
		}
	}
	rate, ok := metric.Options["rate"]
	if ok {
		m.rate = rate.IntVal
	}
	completions, ok := metric.Options["completions"]
	if ok {
		m.completions = completions.IntVal
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
# Custom pre command
%s
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
        %s
		exit 0
    fi
	sleep %d
	let i=i+1
done
# Custom post command after done, if we get here
%s
%s
`
	script := fmt.Sprintf(
		template,
		spec.Spec.Storage.Commands.Pre,
		metadata,
		m.completions,
		metrics.CollectionStart,
		metrics.Separator,
		command,
		metrics.CollectionEnd,
		spec.Spec.Storage.Commands.Post,
		m.rate,
		spec.Spec.Storage.Commands.Post,
		metrics.Interactive(spec.Spec.Logging.Interactive),
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

func init() {
	storage := jobs.StorageGeneric{
		Identifier: "io-sysstat",
		Summary:    "statistics for Linux tasks (processes) : I/O, CPU, memory, etc.",
		Container:  "ghcr.io/converged-computing/metric-sysstat:latest",
	}
	iostat := IOStat{StorageGeneric: storage}
	metrics.Register(&iostat)
}
