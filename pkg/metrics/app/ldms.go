/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package application

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/converged-computing/metrics-operator/pkg/jobs"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

type LDMS struct {
	jobs.SingleApplication

	// Custom Options
	completions int32
	command     string
	rate        int32
}

// I think this is a simulation?
func (m LDMS) Family() string {
	return metrics.PerformanceFamily
}

func (m LDMS) Url() string {
	return "https://github.com/ovis-hpc/ovis"
}

// Set custom options / attributes for the metric
func (m *LDMS) SetOptions(metric *api.Metric) {
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes
	m.rate = 10

	// Set user defined values or fall back to defaults
	m.command = "ldms_ls -h localhost -x sock -p 10444 -l -v"
	m.Workdir = "/opt"

	command, ok := metric.Options["command"]
	if ok {
		m.command = command.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		m.Workdir = workdir.StrVal
	}
	completions, ok := metric.Options["completions"]
	if ok {
		m.completions = completions.IntVal
	}
	rate, ok := metric.Options["rate"]
	if ok {
		m.rate = rate.IntVal
	}
}

// Exported options and list options
func (m LDMS) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.rate)),
		"completions": intstr.FromInt(int(m.completions)),
		"command":     intstr.FromString(m.command),
		"workdir":     intstr.FromString(m.Workdir),
	}
}
func (n LDMS) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Return lookup of entrypoint scripts
func (m LDMS) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)

	// Template for the launcher
	template := `
# Setup munge
mkdir -p /run/munge
chown -R 0 /var/log/munge /var/lib/munge /etc/munge /run/munge
# Skip munge for now, not on a cluster
# ldmsd -x sock:10444 -c /opt/sampler.conf -l /tmp/demo_ldmsd_log -v DEBUG -a munge  -r $(pwd)/ldmsd.pid
ldmsd -x sock:10444 -c /opt/sampler.conf -l /tmp/demo_ldmsd_log -v DEBUG -r $(pwd)/ldmsd.pid
echo "%s"

i=0
completions=%d
echo "%s"
while true
  do
    echo "%s"
    %s
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
echo "%s"
%s
`
	script := fmt.Sprintf(
		template,
		metadata,
		m.completions,
		metrics.CollectionStart,
		metrics.Separator,
		m.command,
		metrics.CollectionEnd,
		metrics.CollectionEnd,
		m.rate,
		metrics.CollectionEnd,
		metrics.Interactive(spec.Spec.Logging.Interactive),
	)
	return []metrics.EntrypointScript{
		{Script: script},
	}
}

func init() {
	app := jobs.SingleApplication{
		Identifier: "app-ldms",
		Summary:    "provides LDMS, a low-overhead, low-latency framework for collecting, transferring, and storing metric data on a large distributed computer system.",
		Container:  "ghcr.io/converged-computing/metric-ovis-hpc:latest",
	}
	LDMS := LDMS{SingleApplication: app}
	metrics.Register(&LDMS)
}
