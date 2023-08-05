/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package perf

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// sysstat provides a tool "pidstat" that can monitor a PID (along with others)
// https://github.com/sysstat/sysstat

type PidStat struct {
	name        string
	rate        int32
	completions int32
	description string
	container   string
}

// Name returns the metric name
func (m PidStat) Name() string {
	return m.name
}

// Description returns the metric description
func (m PidStat) Description() string {
	return m.description
}

// Validation
func (m PidStat) Validate(spec *api.MetricSet) bool {
	return true
}

// Container variables
func (m PidStat) Image() string {
	return m.container
}
func (m PidStat) WorkingDir() string {
	return ""
}
func (m PidStat) Url() string {
	return "https://github.com/sysstat/sysstat"
}

// Set custom options / attributes for the metric
func (m *PidStat) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions
}

func (m PidStat) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {
	return []jobset.ReplicatedJob{}, nil
}

// Generate the replicated job for measuring the application
// We provide the entire Metrics Set (including the application) if we need
// to extract metadata from elsewhere
// TODO need to think of more clever way to export the values?
// Save to somewhere?
// TODO if the app is too fast we might miss it?
func (m PidStat) EntrypointScripts(spec *api.MetricSet) []metrics.EntrypointScript {

	template := `#!/bin/bash

# Download the wait binary
wget https://github.com/converged-computing/goshare/releases/download/2023-07-27/wait
chmod +x ./wait
mv ./wait /usr/bin/goshare-wait
echo "Waiting for application PID..."
pid=$(goshare-wait -c "%s" -q)

i=0
completions=%d
while true
  do
    echo "CPU STATISTICS TIMEPOINT ${i}
    pidstat -p ${pid} -u -h
    echo "KERNEL STATISTICS TIMEPOINT ${i}
    pidstat -p ${pid} -d -h
    echo "POLICY TIMEPOINT ${i}
    pidstat -p ${pid} -R -h
    echo "PAGEFAULTS and MEMORY ${i}
	pidstat -p ${pid} -r -h
    echo "STACK UTILIZATION ${i}
	pidstat -p ${pid} -s -h
    echo "THREADS ${i}	
	pidstat -p ${pid} -t -h
    echo "KERNEL TABLES ${i}	
	pidstat -p ${pid} -v -h
    echo "TASK SWITCHING ${i}	
	pidstat -p ${pid} -w -h
	# Check if still running
	ps -p ${pid} > /dev/null
    retval=$?
	if [[ $retval -ne 0 ]]; then
	    exit 0
    fi
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
	return []metrics.EntrypointScript{
		{Script: fmt.Sprintf(template, spec.Spec.Application.Command, m.completions, m.rate)},
	}
}

func (m PidStat) SuccessJobs() []string {
	return []string{}
}

func (m PidStat) Type() string {
	return metrics.ApplicationMetric
}

func init() {
	metrics.Register(
		&PidStat{
			name:        "perf-sysstat",
			description: "statistics for Linux tasks (processes) : I/O, CPU, memory, etc.",
			container:   "ghcr.io/converged-computing/metric-sysstat:latest",
		})
}
