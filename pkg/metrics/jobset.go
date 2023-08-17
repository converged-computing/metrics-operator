/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

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
	// This is the default Replicated Job Name optional for use
	ReplicatedJobName = "m"
	backoffLimit      = int32(100)
	tenancyLabel      = "metrics-operator-tenancy"
	soleTenancyValue  = "sole-tenancy"
)

// GetJobSet is called by the controller to return some JobSet based
// on the type: application, storage, or standalone
func GetJobSet(
	spec *api.MetricSet,
	sets *map[string]MetricSet,
) ([]*jobset.JobSet, error) {

	// Assume we can eventually support >1 jobset
	jobsets := []*jobset.JobSet{}

	// Assume we have one jobset type
	for _, set := range *sets {
		// For a standalone, we expect one JobSet with 1+ replicatedJobs, and a custom
		// Success Set we expect some subset of the replicated job names
		successJobs := getSuccessJobs(set.Metrics())

		// A base JobSet can hold one or more replicated jobs
		js := getBaseJobSet(spec, successJobs, set.HasSoleTenancy())

		// Get one or more replicated jobs, depending on the type
		rjs, err := set.ReplicatedJobs(spec)
		if err != nil {
			return jobsets, err
		}

		// Get those replicated Jobs.
		js.Spec.ReplicatedJobs = rjs
		jobsets = append(jobsets, js)
	}
	return jobsets, nil
}

// Get list of strings that define successful for a jobset.
// Since these are from replicatedJobs in metrics, we collect from there
func getSuccessJobs(metrics []*Metric) []string {

	// Success jobs are always the default replicatedJobName for storage and application
	// Use a map akin to a set
	successJobs := map[string]bool{}
	for _, m := range metrics {
		for _, sj := range (*m).SuccessJobs() {
			successJobs[sj] = true
		}
	}
	onSuccess := []string{}
	for sj, _ := range successJobs {
		onSuccess = append(onSuccess, sj)
	}
	return onSuccess
}

// getBaseJobSet shared for either an application or isolated jobset
func getBaseJobSet(set *api.MetricSet, successSet []string, soleTenancy bool) *jobset.JobSet {

	// When suspend is true we have a hard time debugging jobs, so keep false
	suspend := false
	enableDNSHostnames := false

	js := jobset.JobSet{
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

	// Do we want to assign 1 node: 1 pod? We can use Pod Anti-affinity for that
	if soleTenancy {
		js.ObjectMeta.Annotations = map[string]string{jobset.ExclusiveKey: "kubernetes.io/hostname"}
	}
	return &js
}

// getReplicatedJob returns the base of the replicated job
func GetReplicatedJob(
	set *api.MetricSet,
	shareProcessNamespace bool,
	pods int32,
	completions int32,
	jobname string,
) (*jobset.ReplicatedJob, error) {

	// Default replicated job name, if not set
	if jobname == "" {
		jobname = ReplicatedJobName
	}

	// Pod labels from the MetricSet
	podLabels := set.GetPodLabels()

	// Indexed mode if >=2 pods
	completionMode := batchv1.NonIndexedCompletion
	if set.Spec.Pods > 1 {
		completionMode = batchv1.IndexedCompletion
	}

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

	// Do we have a pull secret for the application image?
	if set.Spec.Application.PullSecret != "" {
		jobspec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{Name: set.Spec.Application.PullSecret},
		}
	}
	// Tie the jobspec to the job
	job.Template.Spec = jobspec
	return &job, nil
}
