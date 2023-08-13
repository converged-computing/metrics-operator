/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package network

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

// These are common templates for networking standalone apps.

// LauncherWorkerNetwork is a launcher + worker setup for networking apps.
// it is almost functionally equivalent to the LauncherWorker app setup
// but I'm keeping them separate for now
type LauncherWorkerNetwork struct {
	name        string
	rate        int32
	completions int32
	description string
	container   string
	workdir     string
	resources   *api.ContainerResources
	attributes  *api.ContainerSpec

	// Scripts
	workerScript   string
	launcherScript string
	launcherLetter string
}

// Name returns the metric name
func (m LauncherWorkerNetwork) Name() string {
	return m.name
}

// Description returns the metric description
func (m LauncherWorkerNetwork) Description() string {
	return m.description
}

// Jobs required for success condition (n is the LauncherWorkerNetwork run)
func (m *LauncherWorkerNetwork) SuccessJobs() []string {
	return []string{m.launcherLetter}
}

// Container variables
func (n LauncherWorkerNetwork) Type() string {
	return metrics.StandaloneMetric
}
func (n LauncherWorkerNetwork) Image() string {
	return n.container
}
func (n LauncherWorkerNetwork) WorkingDir() string {
	return n.workdir
}

// Return container resources for the metric container
func (n LauncherWorkerNetwork) Resources() *api.ContainerResources {
	return n.resources
}
func (n LauncherWorkerNetwork) Attributes() *api.ContainerSpec {
	return n.attributes
}

func (n LauncherWorkerNetwork) getMetricsKeyToPath() []corev1.KeyToPath {
	// Runner start scripts
	makeExecutable := int32(0777)

	// Each metric has an entrypoint script
	return []corev1.KeyToPath{
		{
			Key:  deriveScriptKey(n.launcherScript),
			Path: path.Base(n.launcherScript),
			Mode: &makeExecutable,
		},
		{
			Key:  deriveScriptKey(n.workerScript),
			Path: path.Base(n.workerScript),
			Mode: &makeExecutable,
		},
	}
}

// Replicated Jobs are custom for this standalone metric
func (m LauncherWorkerNetwork) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {

	js := []jobset.ReplicatedJob{}

	// Generate a replicated job for the launcher (LauncherWorkerNetwork) and workers
	launcher, err := metrics.GetReplicatedJob(spec, false, 1, 1, m.launcherLetter, false)
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

	// runnerScripts are custom for a LauncherWorkerNetwork jobset
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

func (n Netmark) ListOptions() map[string][]intstr.IntOrString {
	return map[string][]intstr.IntOrString{}
}

// Validate that we can run a network. At least one launcher and worker is required
func (n LauncherWorkerNetwork) Validate(spec *api.MetricSet) bool {
	return spec.Spec.Pods >= 2
}

// Given a full path, derive the key from the script name minus the extension
func deriveScriptKey(path string) string {

	// Basename
	path = filepath.Base(path)

	// Remove the extension, and this assumes we don't have double .
	return strings.Split(path, ".")[0]
}

func (n LauncherWorkerNetwork) finalizeEntrypoints(launcherTemplate string, workerTemplate string) []metrics.EntrypointScript {
	return []metrics.EntrypointScript{
		{
			Name:   deriveScriptKey(n.launcherScript),
			Path:   n.launcherScript,
			Script: launcherTemplate,
		},
		{
			Name:   deriveScriptKey(n.workerScript),
			Path:   n.workerScript,
			Script: workerTemplate,
		},
	}
}

// Get common hostlist for launcher/worker app
func (n *LauncherWorkerNetwork) getHostlist(spec *api.MetricSet) string {

	// The launcher has a different hostname, n for netmark
	hosts := fmt.Sprintf("%s-%s-0-0.%s.%s.svc.cluster.local\n",
		spec.Name, n.launcherLetter, spec.Spec.ServiceName, spec.Namespace,
	)
	// Add number of workers
	for i := 0; i < int(spec.Spec.Pods-1); i++ {
		hosts += fmt.Sprintf("%s-w-0-%d.%s.%s.svc.cluster.local\n",
			spec.Name, i, spec.Spec.ServiceName, spec.Namespace)
	}
	return hosts
}
