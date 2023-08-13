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
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	corev1 "k8s.io/api/core/v1"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// These are common (often standalone) templates that can be shared between standalone apps.
// As an example, it's a common pattern to have a launcher with one or more workers.
// We can reduce the redundancy of the code (and make it more streamlined to create
// custom applications) with an ability to share redundant logic here.

// LauncherWorkerApp expects a launcher and worker setup
type LauncherWorkerApp struct {
	name        string
	rate        int32
	completions int32
	description string
	container   string
	workdir     string
	resources   *api.ContainerResources
	attributes  *api.ContainerSpec

	// Common scripts for a launcher/worker design
	workerScript   string
	launcherScript string
}

// Name returns the metric name
func (m LauncherWorkerApp) Name() string {
	return m.name
}

// Description returns the metric description
func (m LauncherWorkerApp) Description() string {
	return m.description
}

// Jobs required for success condition (n is the LauncherWorkerApp run)
func (m LauncherWorkerApp) SuccessJobs() []string {
	return []string{"l"}
}

// Container variables
func (m LauncherWorkerApp) Type() string {
	return metrics.StandaloneMetric
}
func (m LauncherWorkerApp) Image() string {
	return m.container
}
func (m LauncherWorkerApp) WorkingDir() string {
	return m.workdir
}

// Return container resources for the metric container
func (m LauncherWorkerApp) Resources() *api.ContainerResources {
	return m.resources
}
func (m LauncherWorkerApp) Attributes() *api.ContainerSpec {
	return m.attributes
}

func (m LauncherWorkerApp) getMetricsKeyToPath() []corev1.KeyToPath {
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

// Given a full path, derive the key from the script name minus the extension
func deriveScriptKey(path string) string {

	// Basename
	path = filepath.Base(path)

	// Remove the extension, and this assumes we don't have double .
	return strings.Split(path, ".")[0]
}

// Replicated Jobs are custom for this standalone metric
func (m LauncherWorkerApp) ReplicatedJobs(spec *api.MetricSet) ([]jobset.ReplicatedJob, error) {

	js := []jobset.ReplicatedJob{}

	// Generate a replicated job for the launcher (LauncherWorkerApp) and workers
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

	// runnerScripts are custom for a LauncherWorkerApp jobset
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

func (m LauncherWorkerApp) finalizeEntrypoints(launcherTemplate string, workerTemplate string) []metrics.EntrypointScript {
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

// Get common hostlist for launcher/worker app
func (m LauncherWorkerApp) getHostlist(spec *api.MetricSet) string {

	// The launcher has a different hostname, n for LauncherWorkerApp
	hosts := fmt.Sprintf("%s-l-0-0.%s.%s.svc.cluster.local\n",
		spec.Name, spec.Spec.ServiceName, spec.Namespace,
	)
	// Add number of workers
	for i := 0; i < int(spec.Spec.Pods-1); i++ {
		hosts += fmt.Sprintf("%s-w-0-%d.%s.%s.svc.cluster.local\n",
			spec.Name, i, spec.Spec.ServiceName, spec.Namespace)
	}
	return hosts
}
