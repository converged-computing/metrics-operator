package perf

import (
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// sysstat provides a tool "pidstat" that can monitor a PID (along with others)
// https://github.com/sysstat/sysstat

type PidStat struct {
	name        string
	description string
}

// Name returns the metric name
func (m PidStat) Name() string {
	return m.name
}

// Description returns the metric description
func (m PidStat) Description() string {
	return m.description
}

func init() {
	metrics.Register(PidStat{
		name:        "perf-sysstat",
		description: "statistics for Linux tasks (processes) : I/O, CPU, memory, etc.",
	})
}
