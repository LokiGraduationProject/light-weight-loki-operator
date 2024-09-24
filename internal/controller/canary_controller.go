package controller

import (
	"context"
	"fmt"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	"github.com/go-logr/logr"
)

// CanaryReconciler reconciles a Canary object
type CanaryReconciler struct {
	Log logr.Logger
	client.Client
	Scheme *runtime.Scheme
}

const DefaultImage = "grafana/loki-canary:latest"

// Reconcile reads that state of the cluster for a Canary object and makes changes based on the state read
func (r *CanaryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var canary lokiv1.Canary

	if err := r.Get(ctx, req.NamespacedName, &canary); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Define desired state for DaemonSet
	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      canary.Name + "-ds",
			Namespace: canary.Namespace,
			Labels:    getLabels(canary.Spec.DaemonSetLabels),
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: getLabels(canary.Spec.PodLabels),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: getLabels(canary.Spec.PodLabels),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  canary.Name,
							Image: getOrDefault(canary.Spec.Image, DefaultImage),
							Args:  buildArgs(canary.Spec),
						},
					},
				},
			},
		},
	}

	// Set Canary instance as the owner and controller of the DaemonSet
	if err := controllerutil.SetControllerReference(&canary, ds, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if the DaemonSet already exists
	existingDS := &appsv1.DaemonSet{}
	err := r.Get(ctx, client.ObjectKey{Namespace: ds.Namespace, Name: ds.Name}, existingDS)
	if err != nil && errors.IsNotFound(err) {
		// Create the DaemonSet
		if err := r.Create(ctx, ds); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else if !reflect.DeepEqual(ds.Spec, existingDS.Spec) {
		// Update the DaemonSet if necessary
		existingDS.Spec = ds.Spec
		if err := r.Update(ctx, existingDS); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Define desired state for Service
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      canary.Name + "-svc",
			Namespace: canary.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: getLabels(canary.Spec.PodLabels),
			Ports: []corev1.ServicePort{
				{
					Name:       "metrics",
					Port:       canary.Spec.Port,
					TargetPort: intstr.FromInt(3500),
				},
			},
		},
	}

	// Set Canary instance as the owner and controller of the Service
	if err := controllerutil.SetControllerReference(&canary, svc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if the Service already exists
	existingSvc := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{Namespace: svc.Namespace, Name: svc.Name}, existingSvc)
	if err != nil && errors.IsNotFound(err) {
		// Create the Service
		if err := r.Create(ctx, svc); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else if !reflect.DeepEqual(svc.Spec, existingSvc.Spec) {
		// Update the Service if necessary
		existingSvc.Spec = svc.Spec
		if err := r.Update(ctx, existingSvc); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *CanaryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1.Canary{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

// buildArgs constructs the Args slice for the DaemonSet based on the CanarySpec
func buildArgs(spec lokiv1.CanarySpec) []string {
	var args []string

	if spec.Addr != "" {
		args = append(args, fmt.Sprintf("-addr=%s", spec.Addr))
	}
	if spec.Port != 0 {
		args = append(args, fmt.Sprintf("-port=%d", spec.Port))
	}
	if spec.Buckets != 0 {
		args = append(args, fmt.Sprintf("-buckets=%d", spec.Buckets))
	}
	if spec.Interval != "" {
		args = append(args, fmt.Sprintf("-interval=%s", spec.Interval))
	}
	if spec.LabelName != "" {
		args = append(args, fmt.Sprintf("-labelname=%s", spec.LabelName))
	}
	if spec.LabelValue != "" {
		args = append(args, fmt.Sprintf("-labelvalue=%s", spec.LabelValue))
	}
	if spec.MaxWait != "" {
		args = append(args, fmt.Sprintf("-max-wait=%s", spec.MaxWait))
	}
	if spec.MetricTestInterval != "" {
		args = append(args, fmt.Sprintf("-metric-test-interval=%s", spec.MetricTestInterval))
	}
	if spec.MetricTestRange != "" {
		args = append(args, fmt.Sprintf("-metric-test-range=%s", spec.MetricTestRange))
	}
	if spec.OutOfOrderMax != "" {
		args = append(args, fmt.Sprintf("-out-of-order-max=%s", spec.OutOfOrderMax))
	}
	if spec.OutOfOrderMin != "" {
		args = append(args, fmt.Sprintf("-out-of-order-min=%s", spec.OutOfOrderMin))
	}
	if spec.OutOfOrderPercent != 0 {
		args = append(args, fmt.Sprintf("-out-of-order-percentage=%d", spec.OutOfOrderPercent))
	}
	if spec.PruneInterval != "" {
		args = append(args, fmt.Sprintf("-pruneinterval=%s", spec.PruneInterval))
	}
	if spec.Push {
		args = append(args, "-push")
	}
	if spec.QueryTimeout != "" {
		args = append(args, fmt.Sprintf("-query-timeout=%s", spec.QueryTimeout))
	}
	if spec.Size != 0 {
		args = append(args, fmt.Sprintf("-size=%d", spec.Size))
	}
	if spec.SpotCheckInitialWait != "" {
		args = append(args, fmt.Sprintf("-spot-check-initial-wait=%s", spec.SpotCheckInitialWait))
	}
	if spec.SpotCheckInterval != "" {
		args = append(args, fmt.Sprintf("-spot-check-interval=%s", spec.SpotCheckInterval))
	}
	if spec.SpotCheckMax != "" {
		args = append(args, fmt.Sprintf("-spot-check-max=%s", spec.SpotCheckMax))
	}
	if spec.SpotCheckQueryRate != "" {
		args = append(args, fmt.Sprintf("-spot-check-query-rate=%s", spec.SpotCheckQueryRate))
	}
	if spec.StreamName != "" {
		args = append(args, fmt.Sprintf("-streamname=%s", spec.StreamName))
	}
	if spec.StreamValue != "" {
		args = append(args, fmt.Sprintf("-streamvalue=%s", spec.StreamValue))
	}
	if spec.TenantID != "" {
		args = append(args, fmt.Sprintf("-tenant-id=%s", spec.TenantID))
	}
	if spec.WriteMaxBackoff != "" {
		args = append(args, fmt.Sprintf("-write-max-backoff=%s", spec.WriteMaxBackoff))
	}
	if spec.WriteMaxRetries != 0 {
		args = append(args, fmt.Sprintf("-write-max-retries=%d", spec.WriteMaxRetries))
	}
	if spec.WriteMinBackoff != "" {
		args = append(args, fmt.Sprintf("-write-min-backoff=%s", spec.WriteMinBackoff))
	}
	if spec.WriteTimeout != "" {
		args = append(args, fmt.Sprintf("-write-timeout=%s", spec.WriteTimeout))
	}

	return args
}

func getOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func getLabels(labels []lokiv1.Label) map[string]string {
	labelMap := make(map[string]string)
	for _, label := range labels {
		labelMap[label.Key] = label.Value
	}
	return labelMap
}
