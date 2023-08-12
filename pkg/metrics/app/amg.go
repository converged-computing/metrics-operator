/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package application

import (
	"fmt"
	"path"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

type AMG struct {
	name        string
	rate        int32
	completions int32
	description string
	container   string
	resources   *api.ContainerResources
	attributes  *api.ContainerSpec

	// Scripts
	workerScript   string
	launcherScript string

	// Options
	workdir string
	command string
	mpirun  string
}

// Name returns the metric name
func (m AMG) Name() string {
	return m.name
}
func (m AMG) Url() string {
	return "https://github.com/LLNL/AMG"
}

// Description returns the metric description
func (m AMG) Description() string {
	return m.description
}

// Jobs required for success condition (n is the AMG run)
func (m AMG) SuccessJobs() []string {
	return []string{"l"}
}

// Container variables
func (m AMG) Type() string {
	return metrics.StandaloneMetric
}
func (m AMG) Image() string {
	return m.container
}
func (m AMG) WorkingDir() string {
	return ""
}

// Return container resources for the metric container
func (m AMG) Resources() *api.ContainerResources {
	return m.resources
}
func (m AMG) Attributes() *api.ContainerSpec {
	return m.attributes
}

func (m AMG) getMetricsKeyToPath() []corev1.KeyToPath {
	// Runner start scripts
	makeExecutable := int32(0777)

	// Each metric has an entrypoint script
	return []corev1.KeyToPath{
		{
			Key:  deriveScriptKey(m.launcherScript),
			Path: path.Base(m.launcherScript),
			Mode: &makeExecutable,
		},
		{
			Key:  deriveScriptKey(m.workerScript),
			Path: path.Base(m.workerScript),
			Mode: &makeExecutable,
		},
	}
}

// Replicated Jobs are custom for this standalone metric
func (m AMG) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {

	js := []jobset.ReplicatedJob{}

	// Generate a replicated job for the launcher (AMG) and workers
	launcher, err := metrics.GetReplicatedJob(spec, false, 1, 1, "l", false)
	if err != nil {
		return js, err
	}

	workers, err := metrics.GetReplicatedJob(spec, false, spec.Spec.Pods-1, spec.Spec.Pods-1, "w", false)
	if err != nil {
		return js, err
	}

	// Add volumes defined under storage.
	v := map[string]api.Volume{}
	if spec.HasStorage() {
		v["storage"] = spec.Spec.Storage.Volume
	}

	// runnerScripts are custom for a AMG jobset
	runnerScripts := m.getMetricsKeyToPath()

	volumes := metrics.GetVolumes(spec, runnerScripts, v)
	launcher.Template.Spec.Template.Spec.Volumes = volumes
	workers.Template.Spec.Template.Spec.Volumes = volumes

	// Prepare container specs, one for launcher and one for workers
	launcherSpec := []metrics.ContainerSpec{
		{
			Image:      m.container,
			Name:       "launcher",
			Command:    []string{"/bin/bash", m.launcherScript},
			Resources:  m.resources,
			Attributes: m.attributes,
		},
	}
	workerSpec := []metrics.ContainerSpec{
		{
			Image:      m.container,
			Name:       "workers",
			Command:    []string{"/bin/bash", m.workerScript},
			Resources:  m.resources,
			Attributes: m.attributes,
		},
	}
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
func (m *AMG) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions
	m.resources = &metric.Resources
	m.attributes = &metric.Attributes

	// Set user defined values or fall back to defaults
	m.mpirun = "mpirun --hostfile ./hostlist.txt"
	m.command = "amg"
	m.workdir = "/opt/AMG"

	// This could be improved :)
	command, ok := metric.Options["command"]
	if ok {
		m.command = command.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		m.workdir = workdir.StrVal
	}
	mpirun, ok := metric.Options["mpirun"]
	if ok {
		m.mpirun = mpirun.StrVal
	}
}

// Validate that we can run AMG
func (n AMG) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m AMG) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.rate)),
		"completions": intstr.FromInt(int(m.completions)),
		"command":     intstr.FromString(m.command),
		"mpirun":      intstr.FromString(m.mpirun),
		"workdir":     intstr.FromString(m.workdir),
	}
}
func (n AMG) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Return lookup of entrypoint scripts
func (m AMG) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)

	// Generate hostlists
	// The launcher has a different hostname, n for AMG
	hosts := fmt.Sprintf("%s-l-0-0.%s.%s.svc.cluster.local\n",
		spec.Name, spec.Spec.ServiceName, spec.Namespace,
	)
	// Add number of workers
	for i := 0; i < int(spec.Spec.Pods-1); i++ {
		hosts += fmt.Sprintf("%s-w-0-%d.%s.%s.svc.cluster.local\n",
			spec.Name, i, spec.Spec.ServiceName, spec.Namespace)
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

# Write the command file for mpirun
cat <<EOF > ./problem.sh
#!/bin/bash
%s
EOF
chmod +x ./problem.sh

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10
echo "%s"
`
	prefix := fmt.Sprintf(
		prefixTemplate,
		metadata,
		m.workdir,
		hosts,
		m.command,
		metrics.CollectionStart,
	)

	// Template for the launcher
	template := `
echo "%s"
%s ./problem.sh
echo "%s"
%s
`
	launcherTemplate := prefix + fmt.Sprintf(
		template,
		metrics.Separator,
		m.mpirun,
		metrics.CollectionEnd,
		metrics.Interactive(spec.Spec.Logging.Interactive),
	)

	// The worker just has sleep infinity added
	workerTemplate := prefix + "\nsleep infinity"

	// Return the script templates for each of launcher and worker
	return []metrics.EntrypointScript{
		{
			Name:   deriveScriptKey(m.launcherScript),
			Path:   m.launcherScript,
			Script: launcherTemplate,
		},
		{
			Name:   deriveScriptKey(m.workerScript),
			Path:   m.workerScript,
			Script: workerTemplate,
		},
	}
}

func init() {
	metrics.Register(
		&AMG{
			name:           "app-amg",
			description:    "parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids",
			container:      "ghcr.io/converged-computing/metric-amg:latest",
			workerScript:   "/metrics_operator/amg-worker.sh",
			launcherScript: "/metrics_operator/amg-launcher.sh",
		})
}
