/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

 SPDX-License-Identifier: MIT
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
	var spec api.MetricSet

	// Keep developer informed what is going on.
	r.Log.Info("🧀️ Event received by Metric controller!")
	r.Log.Info("Request: ", "req", req)

	// Does the metric exist yet (based on name and namespace)
	err := r.Get(ctx, req.NamespacedName, &spec)
	if err != nil {

		// Create it, doesn't exist yet
		if errors.IsNotFound(err) {
			r.Log.Info("🟥️ MetricSet not found. Ignoring since object must be deleted.")

			// This should not be necessary, but the config map isn't owned by the operator
			return ctrl.Result{}, nil
		}
		r.Log.Info("🟥️ Failed to get MetricSet. Re-running reconcile.")
		return ctrl.Result{Requeue: true}, err
	}

	// Show parameters provided and validate one flux runner
	if !spec.Validate() {
		r.Log.Info("🟥️ Your MetricSet config did not validate.")
		return ctrl.Result{}, nil
	}

	// Verify that all metrics are valid.
	// If the metric requires an application, the MetricSet CRD must have one!
	// If the metric requires storage, the MetricSet CRD must defined storage
	// Only one of application, storage, and standalone is required until
	// we see a use case that warrants this be done differently.
	metrics := []mctrl.Metric{}

	// We are allowed to create more that one MetricSet (JobSet)
	sets := map[string]mctrl.MetricSet{}
	for _, metric := range spec.Spec.Metrics {

		// Get the individual metric, the type will determine the set we add it to
		m, err := mctrl.GetMetric(&metric, &spec)
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("🟥️ We had an issue loading that metric %s!", metric.Name))
			return ctrl.Result{}, nil
		}
		metricType := m.Type()

		// Determine if we've seen the MetricSet type yet, and add either way.
		_, ok := sets[metricType]
		if !ok {
			ms, err := mctrl.GetMetricSet(metricType)
			if err != nil {
				r.Log.Info(fmt.Sprintf("🟥️ We cannot find a metricset type called %s!", metricType))
				return ctrl.Result{}, nil
			}
			sets[metricType] = ms
		}
		sets[metricType].Add(&m)
	}

	// Ensure sets all have one or more metrics
	for setName, set := range sets {
		count := len(set.Metrics())
		if count == 0 {
			r.Log.Info(fmt.Sprintf("🟥️ Metric set %s does not have any validated metrics.", setName))
			return ctrl.Result{}, nil
		}
		r.Log.Info(fmt.Sprintf("🟦️ Metric set %s has %d metrics.", setName, count))
	}
	// Currently just support one JobSet per MetricSet
	if len(sets) != 1 {
		r.Log.Info(fmt.Sprintf("🟥️ Found %d metric sets, but exactly one is allowed to correspond to a final JobSet.", len(sets)))
		return ctrl.Result{}, nil
	}

	// Currently just support one jobset for standalone
	_, ok := sets[mctrl.StandaloneMetric]
	if ok && len(metrics) > 1 {
		r.Log.Info("🟥️ The standalone type metric, by definition, must be measured on its own and not with other metrics.")
		return ctrl.Result{}, nil
	}

	// Ensure the metricset is mapped to a JobSet. For design:
	// 1. If an application is provided, we pair the application at some scale with each metric as a contaienr
	// 2. If storage is provided, we create the volumes for the metric containers
	// 3. If standalone is required, we create a JobSet with custom logic
	result, err := r.ensureMetricSet(ctx, &spec, &sets)
	if err != nil {
		r.Log.Error(err, "🟥️ Issue ensuring metric set")
		return result, err
	}

	// By the time we get here we have a Job + pods + config maps!
	// What else do we want to do?
	r.Log.Info("🧀️ MetricSet is Ready!")
	return ctrl.Result{Requeue: false}, nil
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
