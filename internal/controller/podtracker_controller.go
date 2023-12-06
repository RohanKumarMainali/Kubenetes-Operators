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

package controller

import (
	"context"

	crdv1 "devops.toolbox/controller/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// PodTrackerReconciler reconciles a PodTracker object
type PodTrackerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=crd.devops.toolbox,resources=podtrackers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.devops.toolbox,resources=podtrackers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.devops.toolbox,resources=podtrackers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodTracker object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PodTrackerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciler is getting triggered")

	// var podTracker crdv1.PodTracker
	// if err := r.Get(ctx, req.NamespacedName, &podTracker); err != nil {
	// 	logger.Info("unable to fetch pod tracker")
	// 	return ctrl.Result{}, client.IgnoreNotFound(err)
	// }

	var podTrackerList crdv1.PodTrackerList
	if err := r.List(ctx, &podTrackerList); err != nil {
		logger.Info("Unable to fetch pod tracker list")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(podTrackerList.Items) == 0 {
		logger.Info("No pod tracker found")
		return ctrl.Result{}, nil
	} else {
		var podObject corev1.Pod
		if err := r.Get(context.Background(), req.NamespacedName, &podObject); err != nil {
			logger.Info("unable to fetch pod tracker")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		logger.Info("Send email notifications")
	}

	return ctrl.Result{}, nil
}

func (r *PodTrackerReconciler) HandlePodEvents(ctx context.Context, pod client.Object) []reconcile.Request {
	if pod.GetNamespace() != "default" {
		return []reconcile.Request{}
	}
	namespacedName := types.NamespacedName{
		Namespace: pod.GetNamespace(),
		Name:      pod.GetName(),
	}

	var podObject corev1.Pod
	err := r.Get(context.Background(), namespacedName, &podObject)
	if err != nil {
		return []reconcile.Request{}
	}

	if len(podObject.Annotations) == 0 {
		log.Log.V(1).Info("No annotations set, so this pod is becoming the tracked one", "pod", podObject.Name)
	} else if podObject.GetAnnotations()["exampleAnnotations"] == "crd.devops.toolbox" {
		log.Log.V(1).Info("Found a managed pod, Let's report it", "pod", podObject.Name)
	} else {
		return []reconcile.Request{}
	}

	podObject.SetAnnotations(map[string]string{
		"exampleAnnotations": "crd.devops.toolbox",
	})

	if err := r.Update(context.TODO(), &podObject); err != nil {
		log.Log.V(1).Info("Error trying to update pod", "err", err)
	}

	requests := []reconcile.Request{
		{NamespacedName: namespacedName},
	}
	return requests
	// return []reconcile.Request{}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodTrackerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1.PodTracker{}).
		Watches(
			&corev1.Pod{},
			handler.EnqueueRequestsFromMapFunc(r.HandlePodEvents),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}
