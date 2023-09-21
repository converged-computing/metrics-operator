/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"github.com/converged-computing/metrics-operator/pkg/metadata"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// These are common templates for standalone apps.
// They define the interface of a Metric.

// These are used for network and job names, etc.
var (
	defaultLauncherLetter = "l"
	defaultWorkerLetter   = "w"
)

// LauncherWorker is a launcher + worker setup for apps. These need to
// be accessible by other packages (and not conflict with function names)
type LauncherWorker struct {
	BaseMetric

	Identifier    string
	Summary       string
	Container     string
	Workdir       string
	ResourceSpec  *api.ContainerResources
	AttributeSpec *api.ContainerSpec

	// A metric can have one or more addons
	Addons []*api.MetricAddon

	// Most laucher workers have a command
	Command string
	Prefix  string

	// Scripts
	WorkerScript      string
	LauncherScript    string
	LauncherLetter    string
	WorkerContainer   string
	LauncherContainer string
	WorkerLetter      string
}

// Family returns a generic performance family
func (m LauncherWorker) Family() string {
	return PerformanceFamily
}

// Jobs required for success condition (n is the LauncherWorker run)
func (m *LauncherWorker) SuccessJobs() []string {
	m.ensureDefaultNames()
	return []string{m.LauncherLetter}
}

// Set default options / attributes for the launcher metric
func (m *LauncherWorker) SetDefaultOptions(metric *api.Metric) {
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	command, ok := metric.Options["command"]
	if ok {
		m.Command = command.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		m.Workdir = workdir.StrVal
	}
	prefix, ok := metric.Options["prefix"]
	if ok {
		m.Prefix = prefix.StrVal
	}
}

// Ensure the worker and launcher default names are set
func (m *LauncherWorker) ensureDefaultNames() {
	// Ensure we set the default launcher letter, if not set
	if m.LauncherLetter == "" {
		m.LauncherLetter = defaultLauncherLetter
	}
	if m.WorkerLetter == "" {
		m.WorkerLetter = defaultWorkerLetter
	}
	if m.LauncherScript == "" {
		m.LauncherScript = "/metrics_operator/launcher.sh"
	}
	if m.WorkerScript == "" {
		m.WorkerScript = "/metrics_operator/worker.sh"
	}
	if m.LauncherContainer == "" {
		m.LauncherContainer = "launcher"
	}
	if m.WorkerContainer == "" {
		m.WorkerContainer = "workers"
	}
}

func (m *LauncherWorker) PrepareContainers(
	spec *api.MetricSet,
	metric *Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := Metadata(spec, metric)
	hosts := m.GetHostlist(spec)
	prefix := m.GetCommonPrefix(meta, m.Command, hosts)
	logger.Infof("COMMAND %s", m.Command)

	preBlock := `
echo "%s"
`

	postBlock := `
echo "%s"
%s
`
	command := fmt.Sprintf("%s ./problem.sh", m.Prefix)
	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	preBlock = prefix + fmt.Sprintf(preBlock, metadata.Separator)
	postBlock = fmt.Sprintf(postBlock, metadata.CollectionEnd, interactive)

	// Entrypoint for the launcher
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

// GetCommonPrefix returns a common prefix for the worker/ launcher script, setting up hosts, etc.
func (m *LauncherWorker) GetCommonPrefix(
	meta string,
	command string,
	hosts string,
) string {

	// Generate problem.sh with command only if we have one!
	if command != "" {
		command = fmt.Sprintf(`# Write the command file
cat <<EOF > ./problem.sh
#!/bin/bash
%s
EOF
chmod +x ./problem.sh`, command)
	}

	prefixTemplate := `#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
echo "%s"
# Write the hosts file
cat <<EOF > ./hostlist.txt
%s
EOF

%s

# Allow network to ready (this could be a variable)
echo "Sleeping for 10 seconds waiting for network..."
sleep 10
echo "%s"
`
	return fmt.Sprintf(
		prefixTemplate,
		meta,
		hosts,
		command,
		metadata.CollectionStart,
	)
}

// AddWorkers generates worker jobs, only if we have them
func (m *LauncherWorker) AddWorkers(spec *api.MetricSet) (*jobset.ReplicatedJob, error) {

	numWorkers := spec.Spec.Pods - 1
	workers, err := AssembleReplicatedJob(spec, false, numWorkers, numWorkers, m.WorkerLetter, m.SoleTenancy)
	if err != nil {
		return workers, err
	}
	return workers, nil
}

func (m *LauncherWorker) GetLauncherContainerSpec(
	entrypoint specs.EntrypointScript,
) specs.ContainerSpec {
	spec := specs.ContainerSpec{
		JobName:          m.LauncherLetter,
		Image:            m.Image(),
		Name:             m.LauncherContainer,
		EntrypointScript: entrypoint,
		Resources:        m.ResourceSpec,
		Attributes:       m.AttributeSpec,
	}
	if m.Workdir != "" {
		spec.WorkingDir = m.Workdir
	}
	return spec
}
func (m *LauncherWorker) GetWorkerContainerSpec(
	entrypoint specs.EntrypointScript,
) specs.ContainerSpec {

	// Container spec for the launcher
	spec := specs.ContainerSpec{
		JobName:          m.WorkerLetter,
		Image:            m.Image(),
		Name:             m.WorkerContainer,
		EntrypointScript: entrypoint,
		Resources:        m.ResourceSpec,
		Attributes:       m.AttributeSpec,
	}
	if m.Workdir != "" {
		spec.WorkingDir = m.Workdir
	}
	return spec
}

// Replicated Jobs are custom for a launcher worker
func (m *LauncherWorker) ReplicatedJobs(spec *api.MetricSet) ([]*jobset.ReplicatedJob, error) {

	js := []*jobset.ReplicatedJob{}
	m.ensureDefaultNames()

	// Generate a replicated job for the launcher (LauncherWorker) and workers
	launcher, err := AssembleReplicatedJob(spec, false, 1, 1, m.LauncherLetter, m.SoleTenancy)
	if err != nil {
		return js, err
	}

	numWorkers := spec.Spec.Pods - 1
	var workers *jobset.ReplicatedJob

	// Generate the replicated job with just a launcher, or launcher and workers
	if numWorkers > 0 {
		workers, err = m.AddWorkers(spec)
		if err != nil {
			return js, err
		}
		js = []*jobset.ReplicatedJob{launcher, workers}
	} else {
		js = []*jobset.ReplicatedJob{launcher}
	}
	return js, nil
}

// Validate that we can run a network. At least one launcher and worker is required
func (m LauncherWorker) Validate(spec *api.MetricSet) bool {
	isValid := spec.Spec.Pods >= 2
	if !isValid {
		logger.Errorf("Pods for a Launcher Worker app must be >=2. This app is invalid.")
	}
	return isValid
}

// Get common hostlist for launcher/worker app
func (m *LauncherWorker) GetHostlist(spec *api.MetricSet) string {
	m.ensureDefaultNames()

	// The launcher has a different hostname, n for netmark
	hosts := fmt.Sprintf("%s-%s-0-0.%s.%s.svc.cluster.local\n",
		spec.Name, m.LauncherLetter, spec.Spec.ServiceName, spec.Namespace,
	)
	// Add number of workers
	for i := 0; i < int(spec.Spec.Pods-1); i++ {
		hosts += fmt.Sprintf("%s-%s-0-%d.%s.%s.svc.cluster.local\n",
			spec.Name, m.WorkerLetter, i, spec.Spec.ServiceName, spec.Namespace)
	}
	return hosts
}
