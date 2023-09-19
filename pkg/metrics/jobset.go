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

const podLabelAppName = "app.kubernetes.io/name"

// GetJobSet is called by the controller to return a JobSet for the MetricSet
// Although we currently just support 1 (and return an array of 1)
// this could eventually support > 1 for cases that warrant it.
func GetJobSet(
	spec *api.MetricSet,
	set *MetricSet,
) ([]*jobset.JobSet, error) {

	// Assume we can eventually support >1 jobset
	jobsets := []*jobset.JobSet{}

	// TODO each metric needs to provide some listing of success jobs...
	// Success Set we expect some subset of the replicated job names
	successJobs := getSuccessJobs(set.Metrics())

	// A base JobSet can hold one or more replicated jobs
	js := getBaseJobSet(spec, successJobs)

	// Get one or more replicated jobs, some number from each metric
	rjs := []jobset.ReplicatedJob{}

	// Get one replicated job per metric, and for each, extend with addons
	for _, metric := range set.Metrics() {

		// The metric exposes it's own replicated jobs
		// Since these are custom functions, we add addons / containers / volumes consistently after
		m := (*metric)
		jobs, err := m.ReplicatedJobs(spec)
		if err != nil {
			return jobsets, err
		}

		// Prepare container and volume specs (that are changeable) e.g.,
		// 1. Create VolumeSpec across metrics and addons that can predefine volumes
		// 2. Create ContainerSpec across metrics that can predefine containers, entrypoints, volumes

		// Add addons!
		for _, rj := range jobs {
			m.AddAddons(&rj)
		}

		// Add the final set of jobs
		rjs = append(rjs, jobs...)
	}

	// Get those replicated Jobs.
	js.Spec.ReplicatedJobs = rjs
	jobsets = append(jobsets, js)
	return jobsets, nil
}

// Get list of strings that define successful for a jobset.
// Since these are from replicatedJobs in metrics, we collect from there
func getSuccessJobs(metrics []*Metric) []string {

	// Each metric can define if it's jobs are required for success
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
func getBaseJobSet(set *api.MetricSet, successSet []string) *jobset.JobSet {

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
	return &js
}

// getReplicatedJob returns the base of the replicated job
func GetReplicatedJob(spec *api.MetricSet, set *MetricSet) (*jobset.ReplicatedJob, error) {

	// TODO add way to customize replicated job name, is it needed?
	// Pod labels from the MetricSet
	podLabels := spec.GetPodLabels()

	// Always indexed completion mode to have predictable hostnames
	completionMode := batchv1.IndexedCompletion

	// We only expect one replicated job (for now) so give it a short name for DNS
	job := jobset.ReplicatedJob{
		Name: ReplicatedJobName,
		Template: batchv1.JobTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:      spec.Name,
				Namespace: spec.Namespace,
			},
		},
		// This is the default, but let's be explicit
		Replicas: 1,
	}

	// This should default to true
	setAsFDQN := !spec.Spec.DontSetFQDN

	// Is there an application that might warrant sharing the namespace?
	// TODO we could add logic from metrics here too, if use case arises
	shareProcessNamespace := spec.HasApplication()

	// Create the JobSpec for the job -> Template -> Spec
	jobspec := batchv1.JobSpec{
		BackoffLimit: &backoffLimit,
		Parallelism:  &spec.Spec.Pods,

		// For now assume completions == pods
		Completions:           &spec.Spec.Pods,
		CompletionMode:        &completionMode,
		ActiveDeadlineSeconds: &spec.Spec.DeadlineSeconds,

		// Note there is parameter to limit runtime
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:      spec.Name,
				Namespace: spec.Namespace,
				Labels:    podLabels,
			},
			Spec: corev1.PodSpec{
				// matches the service
				Subdomain:     spec.Spec.ServiceName,
				RestartPolicy: corev1.RestartPolicyOnFailure,

				// This is important to share the process namespace!
				SetHostnameAsFQDN:     &setAsFDQN,
				ShareProcessNamespace: &shareProcessNamespace,
				ServiceAccountName:    spec.Spec.Pod.ServiceAccountName,
				NodeSelector:          spec.Spec.Pod.NodeSelector,
			},
		},
	}

	// Do we want sole tenancy?
	if set.HasSoleTenancy() {
		jobspec.Template.Spec.Affinity = getAffinity(spec)
	}

	// Assemble pull secrets for the application.
	// Metric containers are required to be public.
	jobspec.Template.Spec.ImagePullSecrets = GetPullSecrets(spec, set)

	// Tie the jobspec to the job
	job.Template.Spec = jobspec
	return &job, nil
}

// GetPullSecrets for a metric set (and optionally an application container)
func GetPullSecrets(spec *api.MetricSet, set *MetricSet) []corev1.LocalObjectReference {

	secrets := []corev1.LocalObjectReference{}
	if spec.Spec.Application.PullSecret != "" {
		secrets = append(secrets, corev1.LocalObjectReference{Name: spec.Spec.Application.PullSecret})
	}

	// For now we require metric containers to be public.
	return secrets
}

// getAffinity returns to pod affinity to ensure 1 address / node
func getAffinity(set *api.MetricSet) *corev1.Affinity {
	return &corev1.Affinity{
		// Prefer to schedule pods on the same zone
		PodAffinity: &corev1.PodAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: 100,
					PodAffinityTerm: corev1.PodAffinityTerm{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									// added in getPodLabels
									Key:      podLabelAppName,
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{set.Name},
								},
							},
						},
						TopologyKey: "topology.kubernetes.io/zone",
					},
				},
			},
		},
		// Prefer to schedule pods on different nodes
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
				{
					Weight: 100,
					PodAffinityTerm: corev1.PodAffinityTerm{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									// added in getPodLabels
									Key:      podLabelAppName,
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{set.Name},
								},
							},
						},
						TopologyKey: "kubernetes.io/hostname",
					},
				},
			},
		},
	}
}
