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

	"github.com/converged-computing/metrics-operator/pkg/metadata"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/converged-computing/metrics-operator/pkg/specs"
)

// sysstat provides a tool "iostat" to assess a storage mount
// https://github.com/sysstat/sysstat

type IOStat struct {
	metrics.StorageGeneric
	humanReadable bool
	rate          int32
	completions   int32

	// pre and post commands
	pre  string
	post string
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
	v, ok := metric.Options["pre"]
	if ok {
		m.pre = v.StrVal
	}
	v, ok = metric.Options["post"]
	if ok {
		m.post = v.StrVal
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

func (m IOStat) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)
	command := "iostat -dxm -o JSON"
	if m.humanReadable {
		command = "iostat -dxm"
	}

	preBlock := `#!/bin/bash
# Custom pre comamand logic
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
`

	postBlock := `
%s
%s
`
	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	preBlock = fmt.Sprintf(
		preBlock,
		m.pre,
		meta,
		m.completions,
		metadata.CollectionStart,
		metadata.Separator,
		command,
		metadata.CollectionEnd,
		metadata.CollectionEnd,
		m.rate,
	)

	postBlock = fmt.Sprintf(postBlock, m.post, interactive)
	return m.StorageContainerSpec(preBlock, "", postBlock)
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
	base := metrics.BaseMetric{
		Identifier: "io-sysstat",
		Summary:    "statistics for Linux tasks (processes) : I/O, CPU, memory, etc.",
		Container:  "ghcr.io/converged-computing/metric-sysstat:latest",
	}
	storage := metrics.StorageGeneric{BaseMetric: base}
	iostat := IOStat{StorageGeneric: storage}
	metrics.Register(&iostat)
}
