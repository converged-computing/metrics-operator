/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

This is part of the Flux resource manager framework.
For details, see https://github.com/flux-framework.

SPDX-License-Identifier: Apache-2.0
*/

package controllers

import (
	"context"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	mctrl "github.com/converged-computing/metrics-operator/pkg/metrics"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// ensureMetricsSet creates a JobSet and associated configs
func (r *MetricSetReconciler) ensureMetricSet(
	ctx context.Context,
	spec *api.MetricSet,
	sets *map[string]mctrl.MetricSet,
) (ctrl.Result, error) {

	// First ensure config maps, typically entrypoints for custom metrics containers.
	// They are all bound to the same config map (read only volume)
	// and named by metric index or custom metric script key name
	// We could theoretically allow creating more than one JobSet here
	// and change the name to include the group type.
	_, result, err := r.ensureConfigMaps(ctx, spec, sets)
	if err != nil {
		return result, err
	}

	// Create headless service for the metrics set (which is a JobSet)
	// If we create > 1 JobSet, this should be updated
	selector := map[string]string{"metricset-name": spec.Name}
	result, err = r.exposeServices(ctx, spec, selector)
	if err != nil {
		return result, err
	}

	// Ensure we create the JobSet for the MetricSet
	// either application, storage, or standalone based
	// This could be updated to support > 1
	_, result, err = r.ensureJobSet(ctx, spec, sets)
	if err != nil {
		return result, err
	}

	// And we re-queue so the Ready condition triggers next steps!
	return ctrl.Result{Requeue: true}, nil
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

// getCluster does an actual check if we have a jobset in the namespace
func (r *MetricSetReconciler) ensureJobSet(
	ctx context.Context,
	spec *api.MetricSet,
	sets *map[string]mctrl.MetricSet,
) ([]*jobset.JobSet, ctrl.Result, error) {

	// Look for an existing job
	// We only care about the set Name/Namespace matched to one
	// This can eventually update to support > 1 if needed
	existing, err := r.getExistingJob(ctx, spec)

	// Create a new job if it does not exist
	if err != nil {

		r.Log.Info(
			"✨ Creating a new Metrics JobSet ✨",
			"Namespace:", spec.Namespace,
			"Name:", spec.Name,
		)

		// Get one JobSet to create (can eventually support > 1)
		jobsets, err := mctrl.GetJobSet(spec, sets)
		if err != nil {
			return jobsets, ctrl.Result{}, err
		}
		for _, js := range jobsets {
			err = r.createJobSet(ctx, spec, js)
			if err != nil {
				return jobsets, ctrl.Result{}, err
			}
		}
		return jobsets, ctrl.Result{}, nil

	} else {
		r.Log.Info(
			"🎉 Found existing Metrics JobSet 🎉",
			"Namespace:", existing.Namespace,
			"Name:", existing.Name,
		)
	}
	return []*jobset.JobSet{existing}, ctrl.Result{}, err
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
