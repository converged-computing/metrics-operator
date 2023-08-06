/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package network

import (
	"fmt"
	"path"
	"strconv"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// This library is currently private

type Netmark struct {
	name        string
	rate        int32
	completions int32
	description string
	container   string

	// Scripts
	workerScript      string
	launcherScript    string
	workerScriptKey   string
	launcherScriptKey string

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

// Name returns the metric name
func (m Netmark) Name() string {
	return m.name
}
func (m Netmark) Url() string {
	return ""
}

// Description returns the metric description
func (m Netmark) Description() string {
	return m.description
}

// Jobs required for success condition (n is the netmark run)
func (m Netmark) SuccessJobs() []string {
	return []string{"n"}
}

// Container variables
func (n Netmark) Type() string {
	return metrics.StandaloneMetric
}
func (n Netmark) Image() string {
	return n.container
}
func (n Netmark) WorkingDir() string {
	return ""
}

func (n Netmark) getMetricsKeyToPath() []corev1.KeyToPath {
	// Runner start scripts
	makeExecutable := int32(0777)

	// Each metric has an entrypoint script
	return []corev1.KeyToPath{
		{
			Key:  n.launcherScriptKey,
			Path: path.Base(n.launcherScript),
			Mode: &makeExecutable,
		},
		{
			Key:  n.workerScriptKey,
			Path: path.Base(n.workerScript),
			Mode: &makeExecutable,
		},
	}
}

// Replicated Jobs are custom for this standalone metric
func (m Netmark) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {

	// Generate a replicated job for the launcher (netmark) and workers
	launcher := metrics.GetReplicatedJob(spec, false, 1, 1, "n")
	workers := metrics.GetReplicatedJob(spec, false, spec.Spec.Pods-1, spec.Spec.Pods-1, "w")

	// Add volumes defined under storage.
	v := map[string]api.Volume{}
	if spec.HasStorage() {
		v["storage"] = spec.Spec.Storage.Volume
	}

	// runnerScripts are custom for a netmark jobset
	runnerScripts := m.getMetricsKeyToPath()

	volumes := metrics.GetVolumes(spec, runnerScripts, v)
	launcher.Template.Spec.Template.Spec.Volumes = volumes
	workers.Template.Spec.Template.Spec.Volumes = volumes

	// Prepare container specs, one for launcher and one for workers
	launcherSpec := []metrics.ContainerSpec{
		{
			Image:   m.container,
			Name:    "launcher",
			Command: []string{"/bin/bash", m.launcherScript},
		},
	}
	workerSpec := []metrics.ContainerSpec{
		{
			Image:   m.container,
			Name:    "workers",
			Command: []string{"/bin/bash", m.workerScript},
		},
	}
	js := []jobset.ReplicatedJob{}

	// Derive the containers, one per metric
	// This will also include mounts for volumes
	launcherContainers, err := metrics.GetContainers(spec, launcherSpec, v, false)
	if err != nil {
		fmt.Printf("issue creating launcher containers %s", err)
		return js, err
	}
	workerContainers, err := metrics.GetContainers(spec, workerSpec, v, false)
	if err != nil {
		fmt.Printf("issue creating worker containers %s", err)
		return js, err
	}
	launcher.Template.Spec.Template.Spec.Containers = launcherContainers
	workers.Template.Spec.Template.Spec.Containers = workerContainers
	js = []jobset.ReplicatedJob{*launcher, *workers}
	return js, nil
}

// Set custom options / attributes for the metric
func (m *Netmark) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions

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

// Validate that we can run Netmark
func (n Netmark) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (n Netmark) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":              intstr.FromInt(int(n.rate)),
		"completions":       intstr.FromInt(int(n.completions)),
		"tasks":             intstr.FromInt(int(n.tasks)),
		"warmups":           intstr.FromInt(int(n.warmups)),
		"trials":            intstr.FromInt(int(n.trials)),
		"sendReceiveCycles": intstr.FromInt(int(n.sendReceiveCycles)),
		"messageSize":       intstr.FromInt(int(n.messageSize)),
		"storeEachTrial":    intstr.FromString(strconv.FormatBool(n.storeEachTrial)),
	}
}
func (n Netmark) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Return lookup of entrypoint scripts
func (m Netmark) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Generate hostlists
	// The launcher has a different hostname, n for netmark
	hosts := fmt.Sprintf("%s-n-0-0.%s.%s.svc.cluster.local\n",
		spec.Name, spec.Spec.ServiceName, spec.Namespace,
	)
	// Add number of workers
	for i := 0; i < int(spec.Spec.Pods-1); i++ {
		hosts += fmt.Sprintf("%s-w-0-%d.%s.%s.svc.cluster.local\n",
			spec.Name, i, spec.Spec.ServiceName, spec.Namespace)
	}
	storeTrial := ""
	if m.storeEachTrial {
		storeTrial = "-s"
	}

	prefixTemplate := `#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
whoami
# Show ourselves!
cat ${0}

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
	`
	prefix := fmt.Sprintf(
		prefixTemplate,
		m.tasks,
		spec.Spec.Pods,
		hosts,
	)

	// Template for the launcher
	template := `
mpirun -f ./hostlist.txt -np $np /usr/local/bin/netmark.x -w %d -t %d -c %d -b %d %s
`
	launcherTemplate := prefix + fmt.Sprintf(
		template,
		m.warmups,
		m.trials,
		m.sendReceiveCycles,
		m.messageSize,
		storeTrial,
	)

	// The worker just has sleep infinity added
	workerTemplate := prefix + "\nsleep infinity"

	// Return the script templates for each of launcher and worker
	return []metrics.EntrypointScript{
		{
			Name:   m.launcherScriptKey,
			Path:   m.launcherScript,
			Script: launcherTemplate,
		},
		{
			Name:   m.workerScriptKey,
			Path:   m.workerScript,
			Script: workerTemplate,
		},
	}
}

func init() {
	metrics.Register(
		&Netmark{
			name:              "network-netmark",
			description:       "point to point networking tool",
			container:         "vanessa/netmark:latest",
			workerScript:      "/metrics_operator/netmark-worker.sh",
			launcherScript:    "/metrics_operator/netmark-launcher.sh",
			workerScriptKey:   "netmark-worker",
			launcherScriptKey: "netmark-launcher",
		})
}
