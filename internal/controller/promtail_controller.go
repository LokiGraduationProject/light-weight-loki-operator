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

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

// PromtailReconciler reconciles a Promtail object
type PromtailReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=loki.grafana.com,resources=promtails,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=loki.grafana.com,resources=promtails/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=loki.grafana.com,resources=promtails/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Promtail object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *PromtailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Promtail 객체 가져오기
	promtail := &lokiv1.Promtail{}
	err := r.Get(ctx, req.NamespacedName, promtail)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// ConfigMap 생성 또는 업데이트
	if err := r.createOrUpdateConfigMap(ctx, promtail); err != nil {
		return ctrl.Result{}, err
	}

	// DaemonSet 생성 또는 업데이트
	if err := r.createOrUpdateDaemonSet(ctx, promtail); err != nil {
		return ctrl.Result{}, err
	}

	// RBAC 리소스 생성 또는 업데이트
	if err := r.createOrUpdateClusterRole(ctx); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.createOrUpdateServiceAccount(ctx, promtail.Namespace); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.createOrUpdateClusterRoleBinding(ctx, promtail.Namespace); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PromtailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1.Promtail{}).
		Complete(r)
}

func (r *PromtailReconciler) createOrUpdateConfigMap(ctx context.Context, promtail *lokiv1.Promtail) error {
	// PromtailSpec을 YAML로 변환
	promtailConfigYAML, err := generatePromtailConfigYAML(promtail.Spec.Config)
	if err != nil {
		return err
	}

	// ConfigMap 생성 또는 업데이트
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "promtail-config",
			Namespace: promtail.Namespace,
		},
		Data: map[string]string{
			"promtail.yaml": promtailConfigYAML,
		},
	}

	// ConfigMap이 이미 존재하는지 확인
	existingConfigMap := &corev1.ConfigMap{}
	err = r.Get(ctx, client.ObjectKey{Name: configMap.Name, Namespace: configMap.Namespace}, existingConfigMap)
	if err == nil {
		// 존재하면 업데이트
		configMap.ResourceVersion = existingConfigMap.ResourceVersion
		return r.Update(ctx, configMap)
	} else {
		// 존재하지 않으면 생성
		return r.Create(ctx, configMap)
	}
}

func generatePromtailConfigYAML(config lokiv1.PromtailConfig) (string, error) {
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

func (r *PromtailReconciler) createOrUpdateDaemonSet(ctx context.Context, promtail *lokiv1.Promtail) error {
	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "promtail-daemonset",
			Namespace: promtail.Namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "promtail",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": "promtail",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "promtail-serviceaccount",
					Containers: []corev1.Container{
						{
							Name:  "promtail-container",
							Image: "grafana/promtail",
							Args: []string{
								"-config.file=/etc/promtail/promtail.yaml",
							},
							Env: []corev1.EnvVar{
								{
									Name: "HOSTNAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "logs",
									MountPath: "/var/log",
								},
								{
									Name:      "promtail-config",
									MountPath: "/etc/promtail",
								},
								{
									Name:      "varlibdockercontainers",
									MountPath: "/var/lib/docker/containers",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "logs",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/log",
								},
							},
						},
						{
							Name: "varlibdockercontainers",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/docker/containers",
								},
							},
						},
						{
							Name: "promtail-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "promtail-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// DaemonSet이 이미 존재하는지 확인
	existingDaemonSet := &appsv1.DaemonSet{}
	err := r.Get(ctx, client.ObjectKey{Name: daemonSet.Name, Namespace: daemonSet.Namespace}, existingDaemonSet)
	if err == nil {
		// 존재하면 업데이트
		daemonSet.ResourceVersion = existingDaemonSet.ResourceVersion
		return r.Update(ctx, daemonSet)
	} else {
		// 존재하지 않으면 생성
		return r.Create(ctx, daemonSet)
	}
}

// ClusterRole 생성
func (r *PromtailReconciler) createOrUpdateClusterRole(ctx context.Context) error {
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "promtail-clusterrole",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes", "services", "pods"},
				Verbs:     []string{"get", "watch", "list"},
			},
		},
	}

	existingClusterRole := &rbacv1.ClusterRole{}
	err := r.Get(ctx, client.ObjectKey{Name: clusterRole.Name}, existingClusterRole)
	if err == nil {
		// 존재하면 업데이트
		clusterRole.ResourceVersion = existingClusterRole.ResourceVersion
		return r.Update(ctx, clusterRole)
	} else {
		// 존재하지 않으면 생성
		return r.Create(ctx, clusterRole)
	}
}

// ServiceAccount 생성
func (r *PromtailReconciler) createOrUpdateServiceAccount(ctx context.Context, namespace string) error {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "promtail-serviceaccount",
			Namespace: namespace,
		},
	}

	existingSA := &corev1.ServiceAccount{}
	err := r.Get(ctx, client.ObjectKey{Name: serviceAccount.Name, Namespace: namespace}, existingSA)
	if err == nil {
		// 존재하면 업데이트 필요 없음 (변경사항이 없는 경우)
		return nil
	} else {
		// 존재하지 않으면 생성
		return r.Create(ctx, serviceAccount)
	}
}

// ClusterRoleBinding 생성
func (r *PromtailReconciler) createOrUpdateClusterRoleBinding(ctx context.Context, namespace string) error {
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "promtail-clusterrolebinding",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "promtail-serviceaccount",
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "promtail-clusterrole",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	existingCRB := &rbacv1.ClusterRoleBinding{}
	err := r.Get(ctx, client.ObjectKey{Name: clusterRoleBinding.Name}, existingCRB)
	if err == nil {
		// 존재하면 업데이트
		clusterRoleBinding.ResourceVersion = existingCRB.ResourceVersion
		return r.Update(ctx, clusterRoleBinding)
	} else {
		// 존재하지 않으면 생성
		return r.Create(ctx, clusterRoleBinding)
	}
}
