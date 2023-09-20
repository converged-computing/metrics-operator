/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"github.com/converged-computing/metrics-operator/pkg/metadata"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

	// If we ask for sole tenancy, we assign 1 pod / hostname
	SoleTenancy bool

	// Scripts
	WorkerScript      string
	LauncherScript    string
	LauncherLetter    string
	WorkerContainer   string
	LauncherContainer string
	WorkerLetter      string
}

func (m LauncherWorker) HasSoleTenancy() bool {
	return m.SoleTenancy
}

// Name returns the metric name
func (m LauncherWorker) Name() string {
	return m.Identifier
}

// Description returns the metric description
func (m LauncherWorker) Description() string {
	return m.Summary
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

// Container variables
func (n LauncherWorker) Image() string {
	return n.Container
}
func (m LauncherWorker) WorkingDir() string {
	return m.Workdir
}

// Return container resources for the metric container
func (m LauncherWorker) Resources() *api.ContainerResources {
	return m.ResourceSpec
}
func (m LauncherWorker) Attributes() *api.ContainerSpec {
	return m.AttributeSpec
}

func (m LauncherWorker) getMetricsKeyToPath() []corev1.KeyToPath {
	// Runner start scripts
	makeExecutable := int32(0777)

	// Each metric has an entrypoint script
	return []corev1.KeyToPath{
		{
			Key:  deriveScriptKey(m.LauncherScript),
			Path: path.Base(m.LauncherScript),
			Mode: &makeExecutable,
		},
		{
			Key:  deriveScriptKey(m.WorkerScript),
			Path: path.Base(m.WorkerScript),
			Mode: &makeExecutable,
		},
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

// GetCommonPrefix returns a common prefix for the worker/ launcher script, setting up hosts, etc.
func (m LauncherWorker) GetCommonPrefix(
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
# Change directory to where we will run (and write hostfile)
cd %s
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
		m.WorkingDir(),
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
	return specs.ContainerSpec{
		JobName:          m.LauncherLetter,
		Image:            m.Image(),
		Name:             m.LauncherContainer,
		WorkingDir:       m.Workdir,
		EntrypointScript: entrypoint,
		Resources:        m.ResourceSpec,
		Attributes:       m.AttributeSpec,
	}
}
func (m *LauncherWorker) GetWorkerContainerSpec(
	entrypoint specs.EntrypointScript,
) specs.ContainerSpec {

	// Container spec for the launcher
	return specs.ContainerSpec{
		JobName:          m.WorkerLetter,
		Image:            m.Image(),
		Name:             m.WorkerContainer,
		WorkingDir:       m.Workdir,
		EntrypointScript: entrypoint,
		Resources:        m.ResourceSpec,
		Attributes:       m.AttributeSpec,
	}
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

func (m LauncherWorker) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Validate that we can run a network. At least one launcher and worker is required
func (m LauncherWorker) Validate(spec *api.MetricSet) bool {
	isValid := spec.Spec.Pods >= 2
	if !isValid {
		logger.Errorf("Pods for a Launcher Worker app must be >=2. This app is invalid.")
	}
	return isValid
}

// Given a full path, derive the key from the script name minus the extension
func deriveScriptKey(path string) string {

	// Basename
	path = filepath.Base(path)

	// Remove the extension, and this assumes we don't have double .
	return strings.Split(path, ".")[0]
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
