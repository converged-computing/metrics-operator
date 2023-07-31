package metrics

/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

// Each type of metric returns a replicated job that can be put into a common JobSet

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"

	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

var (
	// Keep this short so DNS doesn't risk overflow
	replicatedJobName = "m"
	backoffLimit      = int32(100)
)

// getBaseJobSet shared for either an application or isolated jobset
func getBaseJobSet(set *api.MetricSet, successSet []string) *jobset.JobSet {

	// When suspend is true we have a hard time debugging jobs, so keep false
	suspend := false
	enableDNSHostnames := false

	return &jobset.JobSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      set.Name,
			Namespace: set.Namespace,
		},
		Spec: jobset.JobSetSpec{
			FailurePolicy: &jobset.FailurePolicy{
				MaxRestarts: 0,
			},
			SuccessPolicy: &jobset.SuccessPolicy{
				Operator:             "All",
				TargetReplicatedJobs: successSet,
			},

			Network: &jobset.Network{
				EnableDNSHostnames: &enableDNSHostnames,
				Subdomain:          set.Spec.ServiceName,
			},

			// This might be the control for child jobs (worker)
			// But I don't think we need this anymore.
			Suspend: &suspend,
		},
	}
}

// getReplicatedJob returns the base of the replicated job
func getReplicatedJob(
	set *api.MetricSet,
	shareProcessNamespace bool,
) *jobset.ReplicatedJob {

	// Pod labels from the MetricSet
	podLabels := set.GetPodLabels()

	completionMode := batchv1.NonIndexedCompletion
	if set.Spec.Completions > 1 {
		completionMode = batchv1.IndexedCompletion
	}

	// We only expect one replicated job (for now) so give it a short name for DNS
	job := jobset.ReplicatedJob{
		Name: replicatedJobName,
		Template: batchv1.JobTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:      set.Name,
				Namespace: set.Namespace,
			},
		},
		// This is the default, but let's be explicit
		Replicas: 1,
	}

	// Create the JobSpec for the job -> Template -> Spec
	jobspec := batchv1.JobSpec{
		BackoffLimit:          &backoffLimit,
		Completions:           &set.Spec.Completions,
		Parallelism:           &set.Spec.Completions,
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
				ShareProcessNamespace: &shareProcessNamespace,
			},
		},
	}

	// Do we have a pull secret for the application image?
	if set.Spec.Application.PullSecret != "" {
		jobspec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{Name: set.Spec.Application.PullSecret},
		}
	}

	// Should we add resources back?
	// jobspec.Template.Spec.Overhead = resources
	// Tie the jobspec to the job
	job.Template.Spec = jobspec
	return &job
}

// GetStorageJobSet creates a jobset intending to mount storage.
// For a storage metric, the metrics are the main containers, at some
// replica level (completions) with shared volumes
func GetStorageJobSet(set *api.MetricSet, metrics *[]Metric) (*jobset.JobSet, error) {
	js := getBaseJobSet(set, []string{replicatedJobName})

	// We don't need a shared process namespace here
	job := getReplicatedJob(set, false)

	// Volumes for job template and containers
	volumes := map[string]api.Volume{"storage": set.Spec.Storage.Volume}

	// Add volumes for storage
	job.Template.Spec.Template.Spec.Volumes = getVolumes(set, metrics, volumes)

	// Derive the containers, one per metric, and volumes are included
	containers, err := getContainers(set, metrics, volumes)
	if err != nil {
		return js, err
	}

	job.Template.Spec.Template.Spec.Containers = containers
	js.Spec.ReplicatedJobs = []jobset.ReplicatedJob{*job}
	return js, nil
}

// CreateApplicationJobSet creates the jobset for the metrics set given an application of interest.
// Each replicated job corresponds to one application being run, and thus one Metrics set.
// For an application, the metrics are sidecar containers to the application
func GetApplicationJobSet(set *api.MetricSet, metrics *[]Metric) (*jobset.JobSet, error) {

	// Done/successful when main application completed
	// TODO when the jobset has customization for indexed completions, update here
	js := getBaseJobSet(set, []string{replicatedJobName})

	// true indicates a shareProcessNamespace
	job := getReplicatedJob(set, true)

	// Add volumes expecting an application.
	job.Template.Spec.Template.Spec.Volumes = getVolumes(
		set, metrics, set.Spec.Application.Volumes,
	)

	// Derive the containers, one per metric
	// This will also include mounts for volumes
	containers, err := getContainers(set, metrics, set.Spec.Application.Volumes)
	if err != nil {
		return js, err
	}
	job.Template.Spec.Template.Spec.Containers = containers
	js.Spec.ReplicatedJobs = []jobset.ReplicatedJob{*job}
	return js, err
}
