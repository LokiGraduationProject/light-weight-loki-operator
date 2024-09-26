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
	"errors"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/status"
)

type LokiStackReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=pods;nodes;services;endpoints;configmaps;secrets;serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=pods;nodes;services;endpoints;configmaps;secrets;serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=loki.lightweight.com,resources=lokistacks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=loki.lightweight.com,resources=lokistacks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=loki.lightweight.com,resources=lokistacks/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=apps,resources=deployments;daemonsets;statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=get;list;watch;create;update
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings;clusterroles;roles;rolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=loki.lightweight.com,resources=canaries,verbs=get;list;watch
//+kubebuilder:rbac:groups=loki.lightweight.com,resources=promtails,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=persistentvolumes;persistentvolumeclaims,verbs=get;list;watch;create;update;delete




func (r *LokiStackReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var degraded *status.DegradedError
	err := r.updateResources(ctx, req)
	switch {
	case errors.As(err, &degraded):
	case err != nil:
		return ctrl.Result{}, err
	}

	err = status.Refresh(ctx, r.Client, req, time.Now(), degraded)
	if err != nil {
		return ctrl.Result{}, err
	}

	if degraded != nil {
		return ctrl.Result{
			Requeue: degraded.Requeue,
		}, nil
	}

	return ctrl.Result{}, nil
}

func (r *LokiStackReconciler) updateResources(ctx context.Context, req ctrl.Request) error {

	err := handlers.CreateOrUpdateLokiStack(ctx, r.Log, req, r.Client, r.Scheme)
	if err != nil {
		return err
	}

	return nil
}

func (r *LokiStackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1.LokiStack{}).
		Complete(r)
}
