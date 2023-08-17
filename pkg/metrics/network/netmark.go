/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package network

import (
	"fmt"
	"strconv"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"

	jobs "github.com/converged-computing/metrics-operator/pkg/jobs"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// This library is currently private

type Netmark struct {
	jobs.LauncherWorker

	// Options
	tasks int32

	// number of warmups
	warmups int32

	// number of trials
	trials int32

	// number of send-recv cycles
	sendReceiveCycles int32

	// message size in bytes
	messageSize int32

	// storage each trial flag
	storeEachTrial bool
}

// Family returns the network family
func (n Netmark) Family() string {
	return metrics.NetworkFamily
}

func (m Netmark) Url() string {
	return ""
}

// Set custom options / attributes for the metric
func (m *Netmark) SetOptions(metric *api.Metric) {
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes
	m.LauncherLetter = "n"

	// One pod per hostname
	m.SoleTenancy = true

	// Set user defined values or fall back to defaults
	// If we have tasks defined, use it! Otherwise fall back to 2 (likely demo)
	m.tasks = 0
	m.warmups = 10
	m.trials = 20
	m.sendReceiveCycles = 20
	m.messageSize = 0
	m.storeEachTrial = true

	// This could be improved :)
	tasks, ok := metric.Options["tasks"]
	if ok {
		m.tasks = tasks.IntVal
	}
	warmups, ok := metric.Options["warmups"]
	if ok {
		m.warmups = warmups.IntVal
	}
	trials, ok := metric.Options["trials"]
	if ok {
		m.trials = trials.IntVal
	}
	messageSize, ok := metric.Options["messageSize"]
	if ok {
		m.messageSize = messageSize.IntVal
	}
	sendReceiveCycle, ok := metric.Options["sendReceiveCycles"]
	if ok {
		m.sendReceiveCycles = sendReceiveCycle.IntVal
	}
	storeEachTrial, ok := metric.Options["storeEachTrial"]
	if ok {
		if storeEachTrial.StrVal == "true" || storeEachTrial.StrVal == "yes" {
			m.storeEachTrial = true
		}
		if storeEachTrial.StrVal == "false" || storeEachTrial.StrVal == "no" {
			m.storeEachTrial = false
		}
	}
}

// Exported options and list options
func (n Netmark) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"tasks":             intstr.FromInt(int(n.tasks)),
		"warmups":           intstr.FromInt(int(n.warmups)),
		"trials":            intstr.FromInt(int(n.trials)),
		"sendReceiveCycles": intstr.FromInt(int(n.sendReceiveCycles)),
		"messageSize":       intstr.FromInt(int(n.messageSize)),
		"storeEachTrial":    intstr.FromString(strconv.FormatBool(n.storeEachTrial)),
	}
}

// Return lookup of entrypoint scripts
func (m Netmark) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)
	hosts := m.GetHostlist(spec)

	// Add boolean flag to store the trial?
	storeTrial := ""
	if m.storeEachTrial {
		storeTrial = "-s"
	}

	prefixTemplate := `#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
echo "%s"

# If we have zero tasks, default to workers * nproc
np=%d
pods=%d
if [[ $np -eq 0 ]]; then
	np=$(nproc)
	np=$(( $pods*$np ))
fi

# Write the hosts file
cat <<EOF > ./hostlist.txt
%s
EOF

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10
echo "%s"
`
	prefix := fmt.Sprintf(
		prefixTemplate,
		metadata,
		m.tasks,
		spec.Spec.Pods,
		hosts,
		metrics.CollectionStart,
	)

	// Template for the launcher
	template := `
mpirun -f ./hostlist.txt -np $np /usr/local/bin/netmark.x -w %d -t %d -c %d -b %d %s
ls
echo "NETMARK RTT.CSV START"
cat RTT.csv
echo "NETMARK RTT.CSV END"
echo "%s"
%s
`
	launcherTemplate := prefix + fmt.Sprintf(
		template,
		m.warmups,
		m.trials,
		m.sendReceiveCycles,
		m.messageSize,
		storeTrial,
		metrics.CollectionEnd,
		metrics.Interactive(spec.Spec.Logging.Interactive),
	)

	// The worker just has sleep infinity added
	workerTemplate := prefix + "\nsleep infinity"
	return m.FinalizeEntrypoints(launcherTemplate, workerTemplate)
}

func init() {
	launcher := jobs.LauncherWorker{
		Identifier:     "network-netmark",
		Summary:        "point to point networking tool",
		Container:      "vanessa/netmark:latest",
		WorkerScript:   "/metrics_operator/netmark-worker.sh",
		LauncherScript: "/metrics_operator/netmark-launcher.sh",
	}
	netmark := Netmark{LauncherWorker: launcher}
	metrics.Register(&netmark)
}
