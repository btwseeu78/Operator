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
	"fmt"

	frontendv1 "frontendapp/api/v1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

)

// MyPythonAppReconciler reconciles a MyPythonApp object
type MyPythonAppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=frontend.stickers.com,resources=mypythonapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=frontend.stickers.com,resources=mypythonapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=frontend.stickers.com,resources=mypythonapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyPythonApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MyPythonAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("MyPythonApps", req.NamespacedName)
	fmt.Println(log) //placeholder for now
	operator := &frontendv1.MyPythonApp{}
	err := r.Get(ctx, req.NamespacedName, operator)
	if err != nil {
		if errors.
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyPythonAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frontendv1.MyPythonApp{}).
		Complete(r)
}
