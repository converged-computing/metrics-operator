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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
)

// exposeService will expose services for job networking (headless)
func (r *MetricSetReconciler) exposeServices(
	ctx context.Context,
	set *api.MetricSet,
	selector map[string]string,
) (ctrl.Result, error) {

	// This service is for the restful API
	existing := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: set.Spec.ServiceName, Namespace: set.Namespace}, existing)
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = r.createHeadlessService(ctx, set, selector)
		}
	}
	return ctrl.Result{}, err
}

// createHeadlessService creates the service
func (r *MetricSetReconciler) createHeadlessService(
	ctx context.Context,
	set *api.MetricSet,
	selector map[string]string,
) (*corev1.Service, error) {

	r.Log.Info("🤯️ Creating headless service with: ", set.Spec.ServiceName, set.Namespace)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: set.Spec.ServiceName, Namespace: set.Namespace},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector:  selector,
		},
	}
	ctrl.SetControllerReference(set, service, r.Scheme)
	err := r.Client.Create(ctx, service)
	if err != nil {
		r.Log.Error(err, "🔴 Create service", "Service", service.Name)
	}
	return service, err
}
