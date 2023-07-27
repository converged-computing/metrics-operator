/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

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
	create "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/go-logr/logr"
)

// MetricReconciler reconciles a Metric object
type MetricReconciler struct {
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
func (r *MetricReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

	// Prepare a group of metrics. For now each will be creating a jobset
	// But I'd like to test having more than one grouped together via
	// a replicated job.
	metrics := []api.Metric{}
	for _, metric := range set.Spec.Metrics {
		m, err := create.GetMetric(metric.Name)
		if err != nil {
			r.Log.Info("🧀️ We cannot find a metric named %s!", metric.Name)
			return ctrl.Result{}, nil
		}
		r.Log.Info("Found metric %s: %s", metric.Name, m.Description())
		metrics = append(metrics, metric)
	}

	// By the time we get here we have a Job + pods + config maps!
	// What else do we want to do?
	r.Log.Info("🧀️ MetricSet is Ready!")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MetricReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.MetricSet{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&jobset.JobSet{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
