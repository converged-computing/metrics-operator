/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package network

import (
	"fmt"
	"strconv"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/converged-computing/metrics-operator/pkg/metadata"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/converged-computing/metrics-operator/pkg/specs"
)

// This library is currently private
const (
	netmarkIdentifier = "network-netmark"
	netmarkSummary    = "point to point networking tool"
	netmarkContainer  = "vanessa/netmark:latest"
)

type Netmark struct {
	metrics.LauncherWorker

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

	m.Identifier = netmarkIdentifier
	m.Summary = netmarkSummary
	m.Container = netmarkContainer

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
	st, ok := metric.Options["soleTenancy"]
	if ok {
		if st.StrVal == "false" || st.StrVal == "no" {
			m.SoleTenancy = false
		}
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

func (m Netmark) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)

	// The launcher has a different hostname, n for netmark
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
		meta,
		m.tasks,
		spec.Spec.Pods,
		hosts,
		metadata.CollectionStart,
	)

	// Netmark main command
	command := "mpirun -f ./hostlist.txt -np $np /usr/local/bin/netmark.x -w %d -t %d -c %d -b %d %s"
	command = fmt.Sprintf(
		command,
		m.warmups,
		m.trials,
		m.sendReceiveCycles,
		m.messageSize,
		storeTrial,
	)
	// The preBlock is also the prefix
	preBlock := prefix

	postBlock := `
ls
echo "NETMARK RTT.CSV START"
cat RTT.csv
echo "NETMARK RTT.CSV END"
echo "%s"
%s
`
	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	postBlock = fmt.Sprintf(
		postBlock,
		metadata.CollectionEnd,
		interactive,
	)

	// The worker just has a preBlock with the prefix and the command is to sleep
	launcherEntrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(m.LauncherScript),
		Path:    m.LauncherScript,
		Pre:     preBlock,
		Command: command,
		Post:    postBlock,
	}

	// Entrypoint for the worker
	workerEntrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(m.WorkerScript),
		Path:    m.WorkerScript,
		Pre:     prefix,
		Command: "sleep infinity",
	}

	// Container spec for the launcher
	launcherContainer := m.GetLauncherContainerSpec(launcherEntrypoint)
	workerContainer := m.GetWorkerContainerSpec(workerEntrypoint)

	// Return the script templates for each of launcher and worker
	return []*specs.ContainerSpec{&launcherContainer, &workerContainer}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: netmarkIdentifier,
		Summary:    netmarkSummary,
		Container:  netmarkContainer,
	}
	launcher := metrics.LauncherWorker{BaseMetric: base}
	netmark := Netmark{LauncherWorker: launcher}
	metrics.Register(&netmark)
}
