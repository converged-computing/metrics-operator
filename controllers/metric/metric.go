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
	set *api.MetricSet,
	metrics *[]mctrl.Metric,
) (ctrl.Result, error) {

	// First ensure config maps, typically entrypoints for custom metrics containers.
	// They are all bound to the same config map (read only volume) and named by metric index
	_, result, err := r.ensureConfigMaps(ctx, set, metrics)
	if err != nil {
		return result, err
	}

	// Create headless service for the metrics set (which is a JobSet)
	selector := map[string]string{"metricset-name": set.Name}
	result, err = r.exposeServices(ctx, set, selector)
	if err != nil {
		return result, err
	}

	// Ensure we create the JobSet for the MetricSet
	// either application or generic based
	_, result, err = r.ensureJobSet(ctx, set, metrics)
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
	set *api.MetricSet,
	metrics *[]mctrl.Metric,
) (*jobset.JobSet, ctrl.Result, error) {

	// Look for an existing job
	existing, err := r.getExistingJob(ctx, set)

	// Create a new job if it does not exist
	if err != nil {

		r.Log.Info(
			"✨ Creating a new Metrics JobSet ✨",
			"Namespace:", set.Namespace,
			"Name:", set.Name,
		)

		var js *jobset.JobSet
		if set.HasApplication() {
			r.Log.Info("Creating application JobSet for MetricSet")
			js, err = mctrl.GetApplicationJobSet(set, metrics)

		} else if set.HasStorage() {
			r.Log.Info("Creating storage JobSet for MetricSet")
			js, err = mctrl.GetStorageJobSet(set, metrics)

		} else {

			// We shouldn't get here
			r.Log.Info("A MetricSet must be for an application or storage.")
			return js, ctrl.Result{}, err
		}
		ctrl.SetControllerReference(set, js, r.Scheme)
		if err != nil {
			return js, ctrl.Result{}, err
		}
		err = r.Client.Create(ctx, js)
		if err != nil {
			r.Log.Error(
				err,
				"Failed to create new Metrics JobSet",
				"Namespace:", js.Namespace,
				"Name:", js.Name,
			)
			return existing, ctrl.Result{}, err
		}
		return js, ctrl.Result{}, err

	} else {
		r.Log.Info(
			"🎉 Found existing Metrics JobSet 🎉",
			"Namespace:", existing.Namespace,
			"Name:", existing.Name,
		)
	}
	ctrl.SetControllerReference(set, existing, r.Scheme)
	return existing, ctrl.Result{}, err
}
