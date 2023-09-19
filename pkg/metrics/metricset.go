/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

var (
	RegistrySet = make(map[string]MetricSet)
)

const (

	// Metric Family Types (these likely can be changed)
	StorageFamily         = "storage"
	MachineLearningFamily = "machine-learning"
	NetworkFamily         = "network"
	SimulationFamily      = "simulation"
	SolverFamily          = "solver"

	// Generic (more than one type, CPU/io, etc)
	ProxyAppFamily    = "proxyapp"
	PerformanceFamily = "performance"
)

// A MetricSet includes one or more metrics that are assembled into a JobSet
type MetricSet struct {
	name        string
	metrics     []*Metric
	metricNames map[string]bool
}

func (m MetricSet) Metrics() []*Metric {
	return m.metrics
}
func (m MetricSet) Exists(metric *Metric) bool {
	_, ok := m.metricNames[(*metric).Name()]
	return ok
}

// Determine if any metrics in the set need sole tenancy
// This is defined on the level of the jobset for now
func (m MetricSet) HasSoleTenancy() bool {
	for _, m := range m.metrics {
		if (*m).HasSoleTenancy() {
			return true
		}
	}
	return false
}

func (m *MetricSet) Add(metric *Metric) {
	if !m.Exists(metric) {
		m.metrics = append(m.metrics, metric)
		m.metricNames[(*metric).Name()] = true
	}
}
func (m *MetricSet) EntrypointScripts(set *api.MetricSet) []EntrypointScript {
	return consolidateEntrypointScripts(m.metrics, set)
}

// AssembleReplicatedJob is used by metrics to assemble a custom, replicated job.
func AssembleReplicatedJob(
	set *api.MetricSet,
	shareProcessNamespace bool,
	pods int32,
	completions int32,
	jobname string,
	soleTenancy bool,
) (*jobset.ReplicatedJob, error) {

	// Default replicated job name, if not set
	if jobname == "" {
		jobname = ReplicatedJobName
	}

	// Pod labels from the MetricSet
	podLabels := set.GetPodLabels()

	// Always indexed completion mode to have predictable hostnames
	completionMode := batchv1.IndexedCompletion

	// We only expect one replicated job (for now) so give it a short name for DNS
	job := jobset.ReplicatedJob{
		Name: jobname,
		Template: batchv1.JobTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:      set.Name,
				Namespace: set.Namespace,
			},
		},
		// This is the default, but let's be explicit
		Replicas: 1,
	}

	// This should default to true
	setAsFDQN := !set.Spec.DontSetFQDN

	// Create the JobSpec for the job -> Template -> Spec
	jobspec := batchv1.JobSpec{
		BackoffLimit:          &backoffLimit,
		Parallelism:           &pods,
		Completions:           &completions,
		CompletionMode:        &completionMode,
		ActiveDeadlineSeconds: &set.Spec.DeadlineSeconds,

		// Note there is parameter to limit runtime
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:      set.Name,
				Namespace: set.Namespace,
				Labels:    podLabels,
			},
			Spec: corev1.PodSpec{
				// matches the service
				Subdomain:     set.Spec.ServiceName,
				RestartPolicy: corev1.RestartPolicyOnFailure,

				// This is important to share the process namespace!
				SetHostnameAsFQDN:     &setAsFDQN,
				ShareProcessNamespace: &shareProcessNamespace,
				ServiceAccountName:    set.Spec.Pod.ServiceAccountName,
				NodeSelector:          set.Spec.Pod.NodeSelector,
			},
		},
	}

	// Do we want sole tenancy?
	if soleTenancy {
		jobspec.Template.Spec.Affinity = getAffinity(set)
	}

	// Do we have a pull secret for the application image?
	// Metric containers are currently required to be public
	if set.Spec.Application.PullSecret != "" {
		jobspec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{Name: set.Spec.Application.PullSecret},
		}
	}
	// Tie the jobspec to the job
	job.Template.Spec = jobspec
	return &job, nil
}

// get an application default entrypoint, if not determined by metric
// NOTE: if the default is not used, we currently just support one metric
// that requires a volume or custom logic. This could be changed
// but my brain is too goobley right now.
func getApplicationDefaultEntrypoint(set *api.MetricSet) string {
	template := `#!/bin/bash
	exec %s	
`
	return fmt.Sprintf(template, set.Spec.Application.Entrypoint)
}

// consolidateEntrypointScripts from a metric set into one list
func consolidateEntrypointScripts(metrics []*Metric, set *api.MetricSet) []EntrypointScript {
	scripts := []EntrypointScript{}
	seenApplicationEntry := false
	for _, metric := range metrics {
		for _, script := range (*metric).EntrypointScripts(set, metric) {
			if script.Path == DefaultApplicationEntrypoint {
				seenApplicationEntry = true
			}
			scripts = append(scripts, script)
		}
	}

	// If we have an application and we haven't seen the application-0.sh, add it
	if set.HasApplication() && !seenApplicationEntry {
		script := getApplicationDefaultEntrypoint(set)
		scripts = append(scripts, EntrypointScript{
			Script: script,
			Path:   DefaultApplicationEntrypoint,
			Name:   DefaultApplicationName,
		})
	}
	return scripts
}
