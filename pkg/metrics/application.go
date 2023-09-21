/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// These are common templates for application metrics
var (
	DefaultEntrypointScript = "/metrics_operator/entrypoint-0.sh"
)

// SingleApplication is a Metric base for a simple application metric
// be accessible by other packages (and not conflict with function names)
type SingleApplication struct {
	BaseMetric
}

func (m SingleApplication) HasSoleTenancy() bool {
	return false
}

// Default SingleApplication is generic performance family
func (m SingleApplication) Family() string {
	return PerformanceFamily
}

func (m *SingleApplication) ApplicationContainerSpec(
	preBlock string,
	command string,
	postBlock string,
) []*specs.ContainerSpec {

	entrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(DefaultEntrypointScript),
		Path:    DefaultEntrypointScript,
		Pre:     preBlock,
		Command: command,
		Post:    postBlock,
	}

	return []*specs.ContainerSpec{{
		JobName:          ReplicatedJobName,
		Image:            m.Image(),
		Name:             "app",
		WorkingDir:       m.WorkingDir,
		EntrypointScript: entrypoint,
		Resources:        m.ResourceSpec,
		Attributes:       m.AttributeSpec,
	}}

}

// Replicated Jobs are custom for a launcher worker
func (m *SingleApplication) ReplicatedJobs(spec *api.MetricSet) ([]*jobset.ReplicatedJob, error) {

	js := []*jobset.ReplicatedJob{}

	// Generate a replicated job for the applicatino
	// An empty jobname will default to "m" the ReplicatedJobName provided by the operator
	rj, err := AssembleReplicatedJob(spec, true, spec.Spec.Pods, spec.Spec.Pods, "", m.SoleTenancy)
	if err != nil {
		return js, err
	}
	js = []*jobset.ReplicatedJob{rj}
	return js, nil
}
