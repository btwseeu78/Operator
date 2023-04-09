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
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
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
		if errors.IsNotFound(err) {
			log.Info("Controller resources Must Be deleted not found teh reuired details")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to get Operator")
		return ctrl.Result{}, err
	}
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Namespace: operator.Namespace, Name: operator.Name}, found)

	if err != nil && errors.IsNotFound(err) {
		dep := r.deploymentForOperator(operator)
		log.Info("Creating A new deployment For Operator", "Deployment.NameSpace", dep.Namespace, dep.Name)
		err := r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to Create The Deploment", dep.Name, "Namespce", dep.Namespace)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}
	deploy := r.deploymentForOperator(operator)

	if !equality.Semantic.DeepDerivative(deploy.Spec.Template, found.Spec.Template) {
		found = deploy
		log.Info("Updatng deployment Template for", "Name", found.Name, "Namespace", found.Namespace)
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Name:", found.Name, "Namespace:", found.Namespace)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}
	size := operator.Spec.Size

	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Not Possible to Scale Up The cluster")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	service := &v1.Service{}
	err = r.Get(ctx, types.NamespacedName{Namespace: operator.Namespace, Name: operator.Name}, service)
	if err != nil && errors.IsNotFound(err) {
		dep = r.serviceForOperator(operator)
		log.Info("creating a new  Service", "Name:", dep.Name, "Namespace", dep.Namespace)
		err = r.Create(ctx, service)
		if err != nil {
			log.Error(err, "Unable to find Associated Service", "Name:", dep.Name, "Namespace", dep.Namespace)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Unable to Retrive service")
		return ctrl.Result{}, err
	}

	podList := &v1.PodList{}
	listOptions := []client.ListOption{
		client.InNamespace(found.Namespace),
		client.MatchingLabels(map[string]string{"app": found.Name, "Namespace": found.Namespace}),
	}
	if err = r.List(ctx, podList, listOptions...); err != nil {
		log.Error(err, "Failed to list pods associated", "Name:", found.Name, "Namespace", found.Namespace)
		return ctrl.Result{}, err
	}
	podNames := getPodName(podList.Items)

	if !reflect.DeepEqual(operator.Status.PodList, podNames) {
		operator.Status.PodList = podNames
		log.Info("Updating operator Status Ffields")
		err = r.Status().Update(ctx, operator)
		if err != nil {
			log.Error(err, "unable to update podlist to", "Namespace", operator.Namespace, "name", operator.Name)
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func getPodName(pods []v1.Pod) []string {
	var podName []string
	for _, val := range pods {
		podName = append(podName, val.Name)
	}
	return podName
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyPythonAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frontendv1.MyPythonApp{}).
		Complete(r)
}
