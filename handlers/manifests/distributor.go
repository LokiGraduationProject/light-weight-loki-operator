package manifests

import (
	"fmt"
	"path"

	"github.com/go-logr/logr"

	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/internal/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BuildDistributor returns a list of k8s objects for Loki Distributor
func BuildDistributor(opts Options, log logr.Logger) ([]client.Object, error) {
	deployment := NewDistributorDeployment(opts, log)

	return []client.Object{
		deployment,
		NewDistributorHTTPService(opts),
		newDistributorPodDisruptionBudget(opts),
	}, nil
}

// NewDistributorDeployment creates a deployment object for a distributor
func NewDistributorDeployment(opts Options, log logr.Logger) *appsv1.Deployment {
	ll := log.WithValues("lokistack", "default", "event", "distributor")

	l := ComponentLabels(LabelDistributorComponent, opts.Name)
	a := commonAnnotations(opts)
	ll.Info(LabelDistributorComponent)
	ll.Info(opts.Name)
	// ll.Info(opts.Stack.DefaultNodeAffinity)
	podSpec := corev1.PodSpec{
		ServiceAccountName: opts.Name,
		Volumes: []corev1.Volume{
			{
				Name: "config",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						DefaultMode: &defaultConfigMapMode,
						LocalObjectReference: corev1.LocalObjectReference{
							Name: fmt.Sprintf("%s-config", opts.Name),
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Image: opts.Image,
				Name:  "loki-distributor",
				Resources: corev1.ResourceRequirements{
					Limits:   opts.ResourceRequirements.Distributor.Limits,
					Requests: opts.ResourceRequirements.Distributor.Requests,
				},
				Args: []string{
					"-target=distributor",
					fmt.Sprintf("-config.file=%s", path.Join(config.LokiConfigMountDir, config.LokiConfigFileName)),
					fmt.Sprintf("-runtime-config.file=%s", path.Join(config.LokiConfigMountDir, config.LokiRuntimeConfigFileName)),
					"-config.expand-env=true",
				},
				ReadinessProbe: lokiReadinessProbe(),
				LivenessProbe:  lokiLivenessProbe(),
				Ports: []corev1.ContainerPort{
					{
						Name:          lokiHTTPPortName,
						ContainerPort: httpPort,
						Protocol:      protocolTCP,
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      configVolumeName,
						ReadOnly:  false,
						MountPath: config.LokiConfigMountDir,
					},
				},
				TerminationMessagePath:   "/dev/termination-log",
				TerminationMessagePolicy: "File",
				ImagePullPolicy:          "IfNotPresent",
			},
		},
	}

	if opts.Stack.Template != nil && opts.Stack.Template.Distributor != nil {
		podSpec.Tolerations = opts.Stack.Template.Distributor.Tolerations
		podSpec.NodeSelector = opts.Stack.Template.Distributor.NodeSelector
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   DistributorName(opts.Name),
			Labels: l,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To(opts.Stack.Template.Distributor.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: l,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fmt.Sprintf("loki-distributor-%s", opts.Name),
					Labels:      l,
					Annotations: a,
				},
				Spec: podSpec,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
		},
	}
}

// NewDistributorHTTPService creates a k8s service for the distributor HTTP endpoint
func NewDistributorHTTPService(opts Options) *corev1.Service {
	serviceName := serviceNameDistributorHTTP(opts.Name)
	labels := ComponentLabels(LabelDistributorComponent, opts.Name)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceName,
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       lokiHTTPPortName,
					Port:       httpPort,
					Protocol:   protocolTCP,
					TargetPort: intstr.IntOrString{IntVal: httpPort},
				},
			},
			Selector: labels,
		},
	}
}

// newDistributorPodDisruptionBudget returns a PodDisruptionBudget for the LokiStack
// Distributor pods.
func newDistributorPodDisruptionBudget(opts Options) *policyv1.PodDisruptionBudget {
	l := ComponentLabels(LabelDistributorComponent, opts.Name)
	mu := intstr.FromInt(1)
	return &policyv1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodDisruptionBudget",
			APIVersion: policyv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    l,
			Name:      DistributorName(opts.Name),
			Namespace: opts.Namespace,
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: l,
			},
			MinAvailable: &mu,
		},
	}
}
