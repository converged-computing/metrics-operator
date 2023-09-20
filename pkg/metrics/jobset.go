/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

// Each type of metric returns a replicated job that can be put into a common JobSet

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"github.com/converged-computing/metrics-operator/pkg/specs"

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
func GetJobSet(
	spec *api.MetricSet,
	set *MetricSet,
) (*jobset.JobSet, []*specs.ContainerSpec, error) {
	containerSpecs := []*specs.ContainerSpec{}

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
			return js, containerSpecs, err
		}

		// Generate container specs for the metric, each is associated with a replicated job
		// The containers are paired with entrypoints, and also with the replicated jobs
		// We do this so we can match addons easily. The only reason we do this outside
		// of the loop below is to allow shared logic.
		cs := m.PrepareContainers(spec, &m)

		// Prepare container and volume specs (that are changeable) e.g.,
		// 1. Create VolumeSpec across metrics and addons that can predefine volumes
		// 2. Create ContainerSpec across metrics that can predefine containers, entrypoints, volumes
		err = m.AddAddons(spec, jobs, cs)
		if err != nil {
			return js, containerSpecs, err
		}

		// Add the finalized container specs for the entire set of replicated jobs
		// We need this at the end to hand back to generate config maps
		// TODO if containers are specific to jobs, maybe need to have based on key...
		containerSpecs = append(containerSpecs, cs...)

		// Add the final set of jobs (bad decision for the pointer here, oops)
		for _, job := range jobs {
			rjs = append(rjs, (*job))
		}
	}

	// Get those replicated Jobs.
	js.Spec.ReplicatedJobs = rjs
	return js, containerSpecs, nil
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
