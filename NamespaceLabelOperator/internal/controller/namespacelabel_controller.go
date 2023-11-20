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
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	multitenancyv1 "github.com/Dolevkle/NamespaceLabelOperator/api/v1"
)

// NamespaceLabelReconciler reconciles a NamespaceLabel object
type NamespaceLabelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=multitenancy.example.org,resources=namespacelabels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multitenancy.example.org,resources=namespacelabels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=multitenancy.example.org,resources=namespacelabels/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NamespaceLabel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *NamespaceLabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	nsLabel := &multitenancyv1.NamespaceLabel{}

	log.Info("Reconciling namespacelabel")

	// Fetch the namespacelabel instance
	if err := r.Get(ctx, req.NamespacedName, nsLabel); err != nil {
		// namespacelabel object not found, it might have been deleted
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	ns := &corev1.Namespace{}

	log.Info("Ensuring Namespace", "namespace", nsLabel.Name)

	// Fetch the corresponding Namespace
	// Define a ns object
	if err := r.ensureNamespace(ctx, nsLabel, ns); err != nil {
		log.Error(err, "unable to ensure Namespace", "namespace", ns)
		return ctrl.Result{}, err
	}

	// Reconcile Namespace labels
	lbls := nsLabel.Spec.Labels

	if lbls == nil {
		lbls = make(map[string]string)
	}

	log.Info("synchronizing namespace labels with NamespaceLabels")

	if !reflect.DeepEqual(ns.Labels, lbls) {
		ns.Labels = lbls
		if err := r.Update(ctx, ns); err != nil {
			log.Error(err, "unable to update Namespace labels")
			return ctrl.Result{}, err
		}

		log.Info("Namespace labels reconciled")
	}

	return ctrl.Result{}, nil
}

// ensureNamespace is a function that ensures that namespace exists,
// else it creates one.
func (r *NamespaceLabelReconciler) ensureNamespace(ctx context.Context, nsLabel *multitenancyv1.NamespaceLabel, ns *corev1.Namespace) error {
	log := log.FromContext(ctx)

	namespaceName := nsLabel.Name

	// Attempt to get the namespace with the provided name
	err := r.Get(ctx, client.ObjectKey{Name: namespaceName}, ns)
	if err != nil {
		// If the namespace doesn't exist, create it
		if errors.IsNotFound(err) {
			log.Info("Creating Namespace", "namespace", namespaceName)
			namespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:   namespaceName,
					Labels: nsLabel.Spec.Labels,
				},
			}
			// Attempt to create the namespace
			if err = r.Create(ctx, namespace); err != nil {
				log.Error(err, "unable to create Namespace")
				return err
			}
			return nil
		} else {
			return err
		}
	} else {
		// If the namespace already exists, check for required annotations
		log.Info("Namespace already exists", "namespace", namespaceName)
		return nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multitenancyv1.NamespaceLabel{}).
		Complete(r)
}
