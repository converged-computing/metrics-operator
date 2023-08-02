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
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/cri-api/pkg/errors"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	mctrl "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/go-logr/logr"
)

// MetricReconciler reconciles a Metric object
type MetricSetReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Log        logr.Logger
	RESTClient rest.Interface
	RESTConfig *rest.Config
}

//+kubebuilder:rbac:groups=flux-framework.org,resources=metricsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=flux-framework.org,resources=metricsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=flux-framework.org,resources=metricsets/finalizers,verbs=update

//+kubebuilder:rbac:groups=jobset.x-k8s.io,resources=jobsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=jobset.x-k8s.io,resources=jobsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=jobset.x-k8s.io,resources=jobsets/finalizers,verbs=update

//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/log,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/exec,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources="",verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=batch,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources=networks,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources="services",verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources="ingresses",verbs=get;list;watch;create;update;patch;delete

//+kubebuilder:rbac:groups="",resources=events,verbs=create;watch;update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete;exec
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;list;watch;create;update;patch;delete;exec

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Metric object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MetricSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Create a new MetricSet
	var set api.MetricSet

	// Keep developer informed what is going on.
	r.Log.Info("🧀️ Event received by Metric controller!")
	r.Log.Info("Request: ", "req", req)

	// Does the metric exist yet (based on name and namespace)
	err := r.Get(ctx, req.NamespacedName, &set)
	if err != nil {

		// Create it, doesn't exist yet
		if errors.IsNotFound(err) {
			r.Log.Info("🧀️ MetricSet not found. Ignoring since object must be deleted.")

			// This should not be necessary, but the config map isn't owned by the operator
			return ctrl.Result{}, nil
		}
		r.Log.Info("🧀️ Failed to get MetricSet. Re-running reconcile.")
		return ctrl.Result{Requeue: true}, err
	}

	// Show parameters provided and validate one flux runner
	if !set.Validate() {
		r.Log.Info("🧀️ Your MetricSet config did not validate.")
		return ctrl.Result{}, nil
	}

	// Verify that all metrics are valid.
	// If the metric requires an application, the MetricSet CRD must have one!
	metrics := []mctrl.Metric{}

	for _, metric := range set.Spec.Metrics {
		m, err := mctrl.GetMetric(&metric)
		if err != nil {
			r.Log.Info(fmt.Sprintf("🧀️ We cannot find a metric named %s!", metric.Name))
			return ctrl.Result{}, nil
		}

		// We can only use the metric if it matches application or storage
		// Ensure we give verbose output if we don't intend to use something
		if m.RequiresApplication() && set.HasApplication() {
			r.Log.Info("Found application metric", metric.Name, m.Description())
			metrics = append(metrics, m)
		} else if m.RequiresStorage() && set.HasStorage() {
			r.Log.Info("Found storage metric", metric.Name, m.Description())
			metrics = append(metrics, m)
		} else if m.Standalone() && set.IsStandalone() {
			r.Log.Info("Found standalone metric", metric.Name, m.Description())
			metrics = append(metrics, m)
		} else {
			r.Log.Info("Metric %s is mismatched for expected MetricSet, skipping.", metric.Name)
		}
	}

	// Ensure the metricset is mapped to a JobSet. For design:
	// 1. If an application is provided, we pair the application at some scale with each metric as a contaienr
	// 2. If an application is not provided, we assume only storage / others that don't require an application
	result, err := r.ensureMetricSet(ctx, &set, &metrics)
	if err != nil {
		return result, err
	}

	// By the time we get here we have a Job + pods + config maps!
	// What else do we want to do?
	r.Log.Info("🧀️ MetricSet is Ready!")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MetricSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.MetricSet{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&batchv1.Job{}).
		Owns(&jobset.JobSet{}).
		Complete(r)
}
