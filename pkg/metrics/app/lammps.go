/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package application

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

type Lammps struct {
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
}

// Name returns the metric name
func (m Lammps) Name() string {
	return m.name
}
func (m Lammps) Url() string {
	return "https://www.lammps.org/"
}

// Description returns the metric description
func (m Lammps) Description() string {
	return m.description
}

// Jobs required for success condition (n is the Lammps run)
func (m Lammps) SuccessJobs() []string {
	return []string{"l"}
}

// Container variables
func (m Lammps) Type() string {
	return metrics.StandaloneMetric
}
func (m Lammps) Image() string {
	return m.container
}
func (m Lammps) WorkingDir() string {
	return ""
}

// Return container resources for the metric container
func (m Lammps) Resources() *api.ContainerResources {
	return m.resources
}
func (m Lammps) Attributes() *api.ContainerSpec {
	return m.attributes
}

func (m Lammps) getMetricsKeyToPath() []corev1.KeyToPath {
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
func (m Lammps) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {

	js := []jobset.ReplicatedJob{}

	// Generate a replicated job for the launcher (Lammps) and workers
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

	// runnerScripts are custom for a Lammps jobset
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
func (m *Lammps) SetOptions(metric *api.Metric) {
	m.rate = metric.Rate
	m.completions = metric.Completions
	m.resources = &metric.Resources
	m.attributes = &metric.Attributes

	// Set user defined values or fall back to defaults
	// This is a more manual approach that puts the user in charge of determining the entire command
	// This more closely matches what we might do on HPC :)
	m.command = "mpirun --hostfile ./hostlist.txt -np 2 --map-by socket lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite"
	m.workdir = "/opt/lammps/examples/reaxff/HNS"

	// This could be improved :)
	command, ok := metric.Options["command"]
	if ok {
		m.command = command.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		m.workdir = workdir.StrVal
	}
}

// Validate that we can run Lammps
func (n Lammps) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Exported options and list options
func (m Lammps) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.rate)),
		"completions": intstr.FromInt(int(m.completions)),
		"command":     intstr.FromString(m.command),
		"workdir":     intstr.FromString(m.workdir),
	}
}
func (n Lammps) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Given a full path, derive the key from the script name minus the extension
func deriveScriptKey(path string) string {

	// Basename
	path = filepath.Base(path)

	// Remove the extension, and this assumes we don't have double .
	return strings.Split(path, ".")[0]
}

// Return lookup of entrypoint scripts
func (m Lammps) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Metadata to add to beginning of run
	metadata := metrics.Metadata(spec, metric)

	// Generate hostlists
	// The launcher has a different hostname, n for Lammps
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
		metrics.CollectionStart,
	)

	// Template for the launcher
	template := `
echo "%s"
%s
echo "%s"
%s
`
	launcherTemplate := prefix + fmt.Sprintf(
		template,
		metrics.Separator,
		m.command,
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
		&Lammps{
			name:           "app-lammps",
			description:    "LAMMPS molecular dynamic simulation",
			container:      "ghcr.io/converged-computing/metric-lammps:latest",
			workerScript:   "/metrics_operator/lammps-worker.sh",
			launcherScript: "/metrics_operator/lammps-launcher.sh",
		})
}
