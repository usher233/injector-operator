/*
Copyright 2024.

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

	"k8s.io/apimachinery/pkg/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	injectorv1alpha1 "github.com/usher233/injector-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// InjectorReconciler reconciles a Injector object
type InjectorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=injector.example.com,resources=injectors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=injector.example.com,resources=injectors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=injector.example.com,resources=injectors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Injector object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *InjectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// TODO(user): your logic here
	injector := &injectorv1alpha1.Injector{}
	if err := r.Client.Get(ctx, req.NamespacedName, injector); err != nil {
		l.Info("unable to fetch Injector", "error", err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	podList := &corev1.PodList{}
	if err := r.Client.List(ctx, podList, client.InNamespace(req.Namespace)); err != nil {
		l.Info("unable to list Pods", "error", err)
		return ctrl.Result{}, err
	}

	if len(podList.Items) == 0 {
		l.Info("no Pods found")
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nginx-pod",
				Namespace: req.Namespace,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "injector",
						Image: "nginx:latest",
					},
				},
			},
		}
		if err := r.Client.Create(ctx, pod); err != nil {
			l.Info("unable to create Pod", "error", err)
			return ctrl.Result{}, err
		}
	}


	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InjectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
    For(&injectorv1alpha1.Injector{}).
    Owns(&corev1.Pod{}).
    Named("injector").
    Complete(r)
}
