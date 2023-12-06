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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	crdv1 "rohan.com/operator/api/v1"
)

// PodFinderReconciler reconciles a PodFinder object
type PodFinderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=crd.rohan.com,resources=podfinders,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.rohan.com,resources=podfinders/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.rohan.com,resources=podfinders/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodFinder object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PodFinderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("reconcilation triggered again")

	var podFinder crdv1.PodFinder
	if err := r.Get(ctx, req.NamespacedName, &podFinder); err != nil {
		logger.Info("Unable to fetch Pod finder resource")
		return ctrl.Result{}, nil
		// return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get pods
	var podList corev1.PodList
	var podFound bool
	if err := r.List(ctx, &podList); err != nil {
		logger.Info("Unable to fetch pod list")
		return ctrl.Result{}, nil
		// return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(podList.Items) == 0 {
		logger.Info("There are no pods")
		return ctrl.Result{}, nil
	} else {
		for _, item := range podList.Items {
			if item.GetName() == podFinder.Spec.Name {
				logger.Info("PodFinder found the expected pod with name that was defined in the spec")
				podFound = true
			}
		}
	}

	// Update the status of PodFinder
	podFinder.Status.Found = podFound
	if err := r.Status().Update(ctx, &podFinder); err != nil {
		logger.Error(err, "unable to update PodFinder's found status", "status", podFound)
		return ctrl.Result{}, err
	} else {
		logger.Info("PodFinder status is updated successfully!", "status", podFound)
	}

	logger.Info("PodFinder custom resouce reconciled")

	return ctrl.Result{}, nil
}

func (r *PodFinderReconciler) mapPodReqToPodFinder(ctx context.Context, pod client.Object) []reconcile.Request {
	logger := log.FromContext(ctx)

	req := []reconcile.Request{}
	var list crdv1.PodFinderList
	if err := r.List(ctx, &list); err != nil {
		logger.Error(err, "Unable to list pod finder")
	} else {
		for _, item := range list.Items {
			if item.Spec.Name == pod.GetName() {
				req = append(req, reconcile.Request{
					NamespacedName: types.NamespacedName{Name: item.Name, Namespace: item.Namespace},
				})
				logger.Info("pod linked to a foo custom resource issued an event", "name", pod.GetName())
			}
		}
	}
	return req
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodFinderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1.PodFinder{}).
		Watches(
			&corev1.Pod{},
			handler.EnqueueRequestsFromMapFunc(r.mapPodReqToPodFinder),
		).
		Complete(r)
}
