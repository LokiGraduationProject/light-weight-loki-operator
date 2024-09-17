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

// LokiStackReconciler reconciles a LokiStack object
type LokiStackReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=loki.lightweight.com,resources=lokistacks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=loki.lightweight.com,resources=lokistacks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=loki.lightweight.com,resources=lokistacks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LokiStack object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *LokiStackReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	var degraded *status.DegradedError
	credentialMode, err := r.updateResources(ctx, req)
	switch {
	case errors.As(err, &degraded):
		// degraded errors are handled by status.Refresh below
	case err != nil:
		return ctrl.Result{}, err
	}

	err = status.Refresh(ctx, r.Client, req, time.Now(), credentialMode, degraded)
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

func (r *LokiStackReconciler) updateResources(ctx context.Context, req ctrl.Request) (lokiv1.CredentialMode, error) {

	credentialMode, err := handlers.CreateOrUpdateLokiStack(ctx, r.Log, req, r.Client, r.Scheme)
	if err != nil {
		return "", err
	}

	return credentialMode, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LokiStackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1.LokiStack{}).
		Complete(r)
}
