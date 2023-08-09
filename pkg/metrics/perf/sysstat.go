/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package perf

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"k8s.io/apimachinery/pkg/util/intstr"
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

	// Options
	useColor bool
	showPIDS bool
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

	// UseColor set to anything means to use it
	_, ok := metric.Options["color"]
	if ok {
		m.useColor = true
	}
	_, ok = metric.Options["pids"]
	if ok {
		m.showPIDS = true
	}

}

func (m PidStat) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {
	return []jobset.ReplicatedJob{}, nil
}

// Exported options and list options
func (m PidStat) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.rate)),
		"completions": intstr.FromInt(int(m.completions)),
	}
}
func (m PidStat) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Generate the replicated job for measuring the application
// TODO if the app is too fast we might miss it?
func (m PidStat) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)

	useColor := ""
	if !m.useColor {
		useColor = "export NO_COLOR=true"
	}

	showPIDS := ""
	if m.showPIDS {
		showPIDS = "ps aux"
	}
	template := `#!/bin/bash

echo "%s"
# Download the wait binary
wget https://github.com/converged-computing/goshare/releases/download/2023-07-27/wait
chmod +x ./wait
mv ./wait /usr/bin/goshare-wait
echo "PIDSTAT COMMAND START"
echo "%s"
echo "PIDSTAT COMMAND END"
echo "Waiting for application PID..."
pid=$(goshare-wait -c "%s" -q)

# Set color or not
%s

# See https://kellyjonbrazil.github.io/jc/docs/parsers/pidstat
# for how we get lovely json
i=0
completions=%d
echo "%s"
while true
  do
    echo "%s"
    %s
    echo "CPU STATISTICS"
    pidstat -p ${pid} -u -h | jc --pidstat
    echo "KERNEL STATISTICS"
    pidstat -p ${pid} -d -h | jc --pidstat
    echo "POLICY"
    pidstat -p ${pid} -R -h | jc --pidstat
    echo "PAGEFAULTS"
    pidstat -p ${pid} -r -h | jc --pidstat
    echo "STACK UTILIZATION"
    pidstat -p ${pid} -s -h | jc --pidstat
    echo "THREADS"
    pidstat -p ${pid} -t -h | jc --pidstat
    echo "KERNEL TABLES"
    pidstat -p ${pid} -v -h | jc --pidstat
    echo "TASK SWITCHING"
    pidstat -p ${pid} -w -h | jc --pidstat
    # Check if still running
    ps -p ${pid} > /dev/null
    retval=$?
    if [[ $retval -ne 0 ]]; then
        echo "%s"
        exit 0
    fi
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
		spec.Spec.Application.Command,
		spec.Spec.Application.Command,
		useColor,
		m.completions,
		metrics.CollectionStart,
		metrics.Separator,
		showPIDS,
		metrics.CollectionEnd,
		metrics.CollectionEnd,
		m.rate,
	)

	// NOTE: the entrypoint is the entrypoint for the container, while
	// the command is expected to be what we are monitoring. Often
	// they are the same thing.
	return []metrics.EntrypointScript{
		{Script: script},
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
