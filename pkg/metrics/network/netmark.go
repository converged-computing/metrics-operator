package network

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// This library is currently private

type Netmark struct {
	name                string
	rate                int32
	completions         int32
	description         string
	container           string
	standalone          bool
	requiresApplication bool
	requiresStorage     bool

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

// Description returns the metric description
func (m Netmark) Description() string {
	return m.description
}

// Container
func (m Netmark) Image() string {
	return m.container
}

// WorkingDir does not matter
func (m Netmark) WorkingDir() string {
	return ""
}

// Set custom options / attributes for the metric
func (m *Netmark) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions

	// Set user defined values or fall back to defaults
	// If we have tasks defined, use it! Otherwise fall back to 2 (likely demo)
	m.tasks = 2
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

// Setup Netmark with ssh and a hostlist.txt for hostnames
// TODO we need a way to specify done with index 0
func (m Netmark) EntrypointScript(set *api.MetricSet) string {

	// Generate hostlists
	hosts := ""
	for i := 0; i < int(set.Spec.Pods); i++ {
		hosts += fmt.Sprintf("%s-m-0-%d.%s.%s.svc.cluster.local\n",
			set.Name, i, set.Spec.ServiceName, set.Namespace)
	}

	storeTrial := ""
	if m.storeEachTrial {
		storeTrial = "-s"
	}

	template := `#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
whoami
# Show ourselves!
cat ${0}

# Write the hosts file
cat <<EOF > ./hostlist.txt
%s
EOF

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10

if [ $JOB_COMPLETION_INDEX = 0 ]; then
   mpirun -f ./hostlist.txt -np %d /usr/local/bin/netmark.x -w %d -t %d -c %d -b %d %s     
else
   sleep infinity
fi
`
	return fmt.Sprintf(
		template,
		hosts,
		m.tasks,
		m.warmups,
		m.trials,
		m.sendReceiveCycles,
		m.messageSize,
		storeTrial,
	)
}

// Does the metric require an application container?
func (m Netmark) RequiresApplication() bool {
	return m.requiresApplication
}
func (m Netmark) RequiresStorage() bool {
	return m.requiresStorage
}
func (m Netmark) Standalone() bool {
	return m.standalone
}

func init() {
	metrics.Register(
		&Netmark{
			name:                "network-netmark",
			description:         "point to point networking tool",
			requiresApplication: false,
			requiresStorage:     false,
			standalone:          true,
			container:           "vanessa/netmark:latest",
		})
}
