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
	var js *jobset.JobSet
	if set.HasApplication() {
		js, err = mctrl.CreateApplicationJobSet(set, metrics)
	} else {
		js, err = mctrl.CreateJobSet(set, metrics)
	}
	if err != nil {
		return result, err
	}
	ctrl.SetControllerReference(set, js, r.Scheme)

	// And we re-queue so the Ready condition triggers next steps!
	return ctrl.Result{Requeue: true}, nil
}
