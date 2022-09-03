/*
Copyright 2022.

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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// SecretReconciler reconciles a Secret object
type SecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=secrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Secret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *SecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log.Log.Info("Reconciling for Request", "req", req.NamespacedName)
	var currentSecrets corev1.Secret
	if err := r.Get(ctx, req.NamespacedName, &currentSecrets); err != nil {
		return ctrl.Result{}, err
	}

	if currentSecrets.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(&currentSecrets, FINALIZER) {
			// TODO: We get this situation when :-
			// 1. New Config is added to the system
			// 2. Existing config with this labels have been added to the system
			controllerutil.AddFinalizer(&currentSecrets, FINALIZER)
			fmt.Println(currentSecrets.GetObjectMeta().GetFinalizers())
			if err := r.Update(ctx, &currentSecrets); err != nil {
				return ctrl.Result{}, err
			}

		}
	} else {
		if controllerutil.ContainsFinalizer(&currentSecrets, FINALIZER) {
			// This handles the case where we have deleted config map and we need to update something into external system
			// Remove Your things

			controllerutil.RemoveFinalizer(&currentSecrets, FINALIZER)
			if err := r.Update(ctx, &currentSecrets); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil

	// TODO(user): your logic here

}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Secret{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(ce event.CreateEvent) bool {
				labels := ce.Object.GetLabels()
				if _, found := labels[SEARCH_LABEL]; found && controllerutil.ContainsFinalizer(ce.Object, FINALIZER) {
					log.Log.Info("New Secret Created", "secret", ce.Object.GetName())
					return true
				}
				return false
			},
			UpdateFunc: func(ue event.UpdateEvent) bool {
				fmt.Println(ue.ObjectOld.GetGeneration(), ue.ObjectNew.GetGeneration())
				labels := ue.ObjectNew.GetLabels()
				if _, found := labels[SEARCH_LABEL]; found {
					return true
				}
				return false
			},
			DeleteFunc: func(de event.DeleteEvent) bool {
				return true
			},
		}).
		Complete(r)
}
