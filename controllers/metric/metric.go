/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package controllers

import (
	"context"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	mctrl "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// ensureMetricsSet creates a JobSet and associated configs
func (r *MetricSetReconciler) ensureMetricSet(
	ctx context.Context,
	spec *api.MetricSet,
	set *mctrl.MetricSet,
) (ctrl.Result, error) {

	// Ensure we create the JobSet for the MetricSet
	// We get back container specs to use for generating configmaps
	// This doesn't actually create the jobset
	js, cs, result, exists, err := r.getJobSet(ctx, spec, set)
	if err != nil {
		return result, err
	}

	// Now create config maps...
	// The config maps need to exist before the jobsets, etc.
	_, result, err = r.ensureConfigMaps(ctx, spec, set, cs)
	if err != nil {
		return result, err
	}

	// And finally, the jobset
	if !exists {
		err = r.createJobSet(ctx, spec, js)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// Create headless service for the metrics set (which is a JobSet)
	// If we create > 1 JobSet, this should be updated
	selector := map[string]string{"metricset-name": spec.Name}
	result, err = r.exposeServices(ctx, spec, selector)
	if err != nil {
		return result, err
	}

	return ctrl.Result{}, nil
}

// getExistingJob gets an existing job that matches our CRD
func (r *MetricSetReconciler) getExistingJob(
	ctx context.Context,
	set *api.MetricSet,
) (*jobset.JobSet, error) {

	existing := &jobset.JobSet{}
	err := r.Client.Get(
		ctx,
		types.NamespacedName{
			Name:      set.Name,
			Namespace: set.Namespace,
		},
		existing,
	)
	return existing, err
}

// getJobset retrieves the existing jobset (or generates the spec for a new one)
func (r *MetricSetReconciler) getJobSet(
	ctx context.Context,
	spec *api.MetricSet,
	set *mctrl.MetricSet,
) (*jobset.JobSet, []*specs.ContainerSpec, ctrl.Result, bool, error) {

	// Look for an existing job
	js, err := r.getExistingJob(ctx, spec)
	cs := []*specs.ContainerSpec{}

	// Create a new job if it does not exist
	if err != nil {

		// TODO test checking for is not found error
		r.Log.Info(
			"✨ Creating a new Metrics JobSet ✨",
			"Namespace:", spec.Namespace,
			"Name:", spec.Name,
		)

		// Get one JobSet and container specs to create config maps
		js, cs, err := mctrl.GetJobSet(spec, set)

		// We don't create it here, we need configmaps first
		return js, cs, ctrl.Result{}, false, err

	}
	r.Log.Info(
		"🎉 Found existing Metrics JobSet 🎉",
		"Namespace:", js.Namespace,
		"Name:", js.Name,
	)
	return js, cs, ctrl.Result{}, true, err
}

// createJobSet handles the creation operator
func (r *MetricSetReconciler) createJobSet(
	ctx context.Context,
	spec *api.MetricSet,
	js *jobset.JobSet,
) error {
	r.Log.Info(
		"🎉 Creating Metrics JobSet 🎉",
		"Namespace:", js.Namespace,
		"Name:", js.Name,
	)

	// Controller reference always needs to be set before creation
	ctrl.SetControllerReference(spec, js, r.Scheme)
	err := r.Client.Create(ctx, js)
	if err != nil {
		r.Log.Error(
			err,
			"Failed to create new Metrics JobSet",
			"Namespace:", js.Namespace,
			"Name:", js.Name,
		)
		return err
	}
	return nil
}
