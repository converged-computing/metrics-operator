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
func getBaseJobSet(set *api.MetricSet) *jobset.JobSet {

	// When suspend is true we have a hard time debugging jobs, so keep false
	suspend := false
	enableDNSHostnames := false

	return &jobset.JobSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      set.Name,
			Namespace: set.Namespace,
		},
		Spec: jobset.JobSetSpec{

			// We define success when all of the jobs are done
			// This could be limited to the application in question
			FailurePolicy: &jobset.FailurePolicy{
				MaxRestarts: 0,
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
func getReplicatedJob(set *api.MetricSet, completionMode batchv1.CompletionMode) *jobset.ReplicatedJob {

	// Pod labels from the MetricSet
	podLabels := set.GetPodLabels()

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

	// We want to share the process namespace between containers
	shareProcessNamespace := true

	// Create the JobSpec for the job -> Template -> Spec
	jobspec := batchv1.JobSpec{
		BackoffLimit:          &backoffLimit,
		Completions:           &set.Spec.Application.Completions,
		Parallelism:           &set.Spec.Application.Completions,
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

	// TODO we will vary here by resources, mounts, and contianers
	// Also - should we add resources back?
	// jobspec.Template.Spec.Overhead = resources
	// Tie the jobspec to the job
	job.Template.Spec = jobspec
	return &job
}

// CreateJobSet creates a generic jobset to only run metrics.
// We typically expect to just use our own containers or test storage, and might
// extend the functions to be specific to that.
func CreateJobSet(set *api.MetricSet, metrics *[]Metric) (*jobset.JobSet, error) {
	js := getBaseJobSet(set)

	// TODO not written yet
	// This will be for storage / etc metrics that need volumes but not application logic
	return js, nil
}

// CreateApplicationJobSet creates the jobset for the metrics set given an application of interest.
// Each replicated job corresponds to one application being run, and thus one Metrics set. We use a jobset to
// store associated services alongside the job (TBA) and indexed mode to allow multiple replicas.
func CreateApplicationJobSet(set *api.MetricSet, metrics *[]Metric) (*jobset.JobSet, error) {
	js := getBaseJobSet(set)

	// We always create appliction jobsets with indexed completion
	job := getReplicatedJob(set, batchv1.NonIndexedCompletion)

	// Add volumes expecting an application (this could be general and moved up into function above)
	job.Template.Spec.Template.Spec.Volumes = getVolumes(set, metrics)

	// Derive the containers, one per metric
	// This will also include mounts for volumes
	containers, err := getContainers(set, metrics)
	if err != nil {
		return js, err
	}
	job.Template.Spec.Template.Spec.Containers = containers
	js.Spec.ReplicatedJobs = []jobset.ReplicatedJob{*job}
	return js, err
}
