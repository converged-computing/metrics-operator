/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package jobs

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
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
	Identifier    string
	Summary       string
	Container     string
	Workdir       string
	ResourceSpec  *api.ContainerResources
	AttributeSpec *api.ContainerSpec

	// Scripts
	WorkerScript   string
	LauncherScript string
	LauncherLetter string
	WorkerLetter   string
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
	return metrics.PerformanceFamily
}

// Jobs required for success condition (n is the LauncherWorker run)
func (m *LauncherWorker) SuccessJobs() []string {
	m.ensureDefaultNames()
	return []string{m.LauncherLetter}
}

// Container variables
func (n LauncherWorker) Type() string {
	return metrics.StandaloneMetric
}
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
}

// GetCommonPrefix returns a common prefix for the worker/ launcher script, setting up hosts, etc.
func (m LauncherWorker) GetCommonPrefix(
	metadata string,
	command string,
	hosts string,
) string {

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

# Write the command file
cat <<EOF > ./problem.sh
#!/bin/bash
%s
EOF
chmod +x ./problem.sh

# Allow network to ready (this could be a variable)
echo "Sleeping for 10 seconds waiting for network..."
sleep 10
echo "%s"
`
	return fmt.Sprintf(
		prefixTemplate,
		metadata,
		m.WorkingDir(),
		hosts,
		command,
		metrics.CollectionStart,
	)
}

// Replicated Jobs are custom for this standalone metric
func (m *LauncherWorker) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {

	js := []jobset.ReplicatedJob{}
	m.ensureDefaultNames()

	// Generate a replicated job for the launcher (LauncherWorker) and workers
	launcher, err := metrics.GetReplicatedJob(spec, false, 1, 1, m.LauncherLetter, false)
	if err != nil {
		return js, err
	}

	workers, err := metrics.GetReplicatedJob(spec, false, spec.Spec.Pods-1, spec.Spec.Pods-1, m.WorkerLetter, false)
	if err != nil {
		return js, err
	}

	// Add volumes defined under storage.
	v := map[string]api.Volume{}
	if spec.HasStorage() {
		v["storage"] = spec.Spec.Storage.Volume
	}

	// runnerScripts are custom for a LauncherWorker jobset
	runnerScripts := m.getMetricsKeyToPath()

	volumes := metrics.GetVolumes(spec, runnerScripts, v)
	launcher.Template.Spec.Template.Spec.Volumes = volumes
	workers.Template.Spec.Template.Spec.Volumes = volumes

	// Prepare container specs, one for launcher and one for workers
	launcherSpec := []metrics.ContainerSpec{
		{
			Image:      m.Container,
			Name:       "launcher",
			Command:    []string{"/bin/bash", m.LauncherScript},
			Resources:  m.ResourceSpec,
			Attributes: m.AttributeSpec,
		},
	}
	workerSpec := []metrics.ContainerSpec{
		{
			Image:      m.Container,
			Name:       "workers",
			Command:    []string{"/bin/bash", m.WorkerScript},
			Resources:  m.ResourceSpec,
			Attributes: m.AttributeSpec,
		},
	}

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

func (m LauncherWorker) FinalizeEntrypoints(launcherTemplate string, workerTemplate string) []metrics.EntrypointScript {
	return []metrics.EntrypointScript{
		{
			Name:   deriveScriptKey(m.LauncherScript),
			Path:   m.LauncherScript,
			Script: launcherTemplate,
		},
		{
			Name:   deriveScriptKey(m.WorkerScript),
			Path:   m.WorkerScript,
			Script: workerTemplate,
		},
	}
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
