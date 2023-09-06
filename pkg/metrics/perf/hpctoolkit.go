/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package perf

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"github.com/converged-computing/metrics-operator/pkg/jobs"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type HPCToolkit struct {
	jobs.SingleApplication
	events string
	mount  string
}

func (m HPCToolkit) Url() string {
	return "https://gitlab.com/hpctoolkit/hpctoolkit"
}

// GetVolumes to provide an empty volume for the application to share
func (m HPCToolkit) GetVolumes() map[string]api.Volume {
	return map[string]api.Volume{
		"hpctoolkit": {
			Path:     m.mount,
			EmptyVol: true,
		},
	}
}

// Validate we have an executable provided, and args and optional
func (m *HPCToolkit) Validate(ms *api.MetricSet) bool {
	if m.events == "" {
		logger.Error("One or more events for hpcrun (events) are required (e.g., -e IO).")
		return false
	}
	return true
}

// Set custom options / attributes for the metric
func (m *HPCToolkit) SetOptions(metric *api.Metric) {
	// Defaults for rate and completions
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes
	m.mount = "/opt/share"

	// UseColor set to anything means to use it
	mount, ok := metric.Options["mount"]
	if ok {
		m.mount = mount.StrVal
	}
	events, ok := metric.Options["events"]
	if ok {
		m.events = events.StrVal
	}
}

// Exported options and list options
func (m HPCToolkit) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"events": intstr.FromString(m.events),
		"mount":  intstr.FromString(m.mount),
	}
}

// Generate the replicated job for measuring the application
// TODO if the app is too fast we might miss it?
func (m HPCToolkit) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// This is the metric container entrypoint.
	// The sole purpose is just to provide the volume, meaning copying content there
	template := `#!/bin/bash

echo "Moving content from /opt/view to be in shared volume at %s"
view=$(ls /opt/views/._view/)
view="/opt/views/._view/${view}"

# Give a little extra wait time
sleep 10

viewroot="%s"
mkdir -p $viewroot/view
# We have to move both of these paths, *sigh*
cp -R ${view}/* $viewroot/view
cp -R /opt/software $viewroot/

# Sleep forever, the application needs to run and end
echo "Sleeping forever so %s can be shared and use for hpctoolkit."
sleep infinity
`
	script := fmt.Sprintf(
		template,
		m.mount,
		m.mount,
		m.mount,
	)

	// Custom logic for application entrypoint
	metadata := metrics.Metadata(spec, metric)
	custom := `

# Ensure hpcrun and software exists. This is rough, but should be OK with enough wait time
wget https://github.com/converged-computing/goshare/releases/download/2023-09-06/wait-fs
chmod +x ./wait-fs
mv ./wait-fs /usr/bin/goshare-wait-fs

# Ensure spack view is on the path, wherever it is mounted
viewbase="%s"
software="${viewbase}/software"
viewbin="${viewbase}/view/bin"
export PATH=${viewbin}:$PATH

# Wait for software directory, and give it time
goshare-wait-fs -p ${software}

# Wait for copy to finish
sleep 10

# Copy mount software to /opt/software
cp -R %s/software /opt/software

# Wait for hpcrun
goshare-wait-fs -p ${viewbin}/hpcrun

# This will work with capability SYS_ADMIN added. 
# It will only work with privileged set to true AT YOUR OWN RISK!
echo "-1" | tee /proc/sys/kernel/perf_event_paranoid

# Run hpcrun. See options with hpcrun -L
events="%s"
echo "%s"
echo "%s"
echo "%s"

# Commands to interact with output data
# hpcprof hpctoolkit-sleep-measurements
# hpcstruct hpctoolkit-sleep-measurements
# hpcviewer ./hpctoolkit-lmp-database
`

	custom = fmt.Sprintf(
		custom,
		m.mount,
		m.mount,
		m.events,
		metadata,
		metrics.CollectionStart,
		metrics.Separator,
	)

	// And the suffix (post run)
	suffix := `
echo "%s"
%s
`
	suffix = fmt.Sprintf(
		suffix,
		metrics.CollectionEnd,
		metrics.Interactive(spec.Spec.Logging.Interactive),
	)

	// NOTE: for this container the metrics entrypoint just copies and then
	// waits, and the custom application entrypoint runs the wrapped application
	// command.
	return []metrics.EntrypointScript{
		{Script: script},
		m.ApplicationEntrypoint(spec, custom, "hpcrun $events", suffix),
	}
}

func init() {
	app := jobs.SingleApplication{
		Identifier: "perf-hpctoolkit",
		Summary:    "performance tools for measurement and analysis",
		Container:  "ghcr.io/converged-computing/metric-hpctoolkit-view:latest",
	}
	HPCToolkit := HPCToolkit{SingleApplication: app}
	metrics.Register(&HPCToolkit)
}
