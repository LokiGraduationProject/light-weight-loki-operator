package manifests

import (
	"fmt"
	"math"
	"path"

	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/internal/config"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/storage"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func BuildBackendComponent(opts Options) ([]client.Object, error) {
	statefulSet := newBackendStatefulSet(opts)

	if err := storage.ConfigureStatefulSet(statefulSet, opts.ObjectStorage); err != nil {
		return nil, err
	}

	if err := configureHashRingEnv(&statefulSet.Spec.Template.Spec, opts); err != nil {
		return nil, err
	}

	if err := configureProxyEnv(&statefulSet.Spec.Template.Spec, opts); err != nil {
		return nil, err
	}

	if err := configureReplication(&statefulSet.Spec.Template, opts.Stack.Replication, LabelBackendComponent, opts.Name); err != nil {
		return nil, err
	}

	return []client.Object{
		statefulSet,
		NewBackendGRPCService(opts),
		NewBackendHTTPService(opts),
		NewBackendPodDisruptionBudget(opts),
	}, nil
}

func newBackendStatefulSet(opts Options) *appsv1.StatefulSet {
	l := ComponentLabels(LabelBackendComponent, opts.Name)
	a := commonAnnotations(opts)

	replicas := int32(1)

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "loki-backend",
			Labels: l,
		},
		Spec: appsv1.StatefulSetSpec{
			PodManagementPolicy:  appsv1.OrderedReadyPodManagement,
			RevisionHistoryLimit: ptr.To(defaultRevHistoryLimit),
			Replicas:             &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels.Merge(l, GossipLabels()),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fmt.Sprintf("loki-backend-%s", opts.Name),
					Labels:      labels.Merge(l, GossipLabels()),
					Annotations: a,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: opts.Name,
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									DefaultMode: &defaultConfigMapMode,
									LocalObjectReference: corev1.LocalObjectReference{
										Name: lokiConfigMapName(opts.Name),
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Image: opts.Image,
							Name:  "loki-backend-component",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("1Gi"),
								},
							},
							Args: []string{
								"-target=backend",
								fmt.Sprintf("-config.file=%s", path.Join(config.LokiConfigMountDir, config.LokiConfigFileName)),
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
								{
									Name:          lokiGRPCPortName,
									ContainerPort: grpcPort,
									Protocol:      protocolTCP,
								},
								{
									Name:          lokiGossipPortName,
									ContainerPort: gossipPort,
									Protocol:      protocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      configVolumeName,
									ReadOnly:  false,
									MountPath: config.LokiConfigMountDir,
								},
								{
									Name:      storageVolumeName,
									ReadOnly:  false,
									MountPath: dataDirectory,
								},
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: "File",
							ImagePullPolicy:          "IfNotPresent",
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Labels: l,
						Name:   storageVolumeName,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("10Gi"),
							},
						},
						StorageClassName: ptr.To(opts.Stack.StorageClassName),
						VolumeMode:       &volumeFileSystemMode,
					},
				},
			},
		},
	}
}

func NewBackendGRPCService(opts Options) *corev1.Service {
	labels := ComponentLabels(LabelBackendComponent, opts.Name)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceNameBackendGRPC(opts.Name),
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name:       lokiGRPCPortName,
					Port:       grpcPort,
					Protocol:   protocolTCP,
					TargetPort: intstr.IntOrString{IntVal: grpcPort},
				},
			},
			Selector: labels,
		},
	}
}

func NewBackendHTTPService(opts Options) *corev1.Service {
	serviceName := serviceNameBackendHTTP(opts.Name)
	labels := ComponentLabels(LabelBackendComponent, opts.Name)

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

func NewBackendPodDisruptionBudget(opts Options) *policyv1.PodDisruptionBudget {
	l := ComponentLabels(LabelBackendComponent, opts.Name)

	// Have at least N-1 replicas available, unless N==1 in which case the minimum available is 1.
	replicas := int32(1)
	ma := intstr.FromInt(int(math.Max(1, float64(replicas-1))))

	return &policyv1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodDisruptionBudget",
			APIVersion: policyv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    l,
			Name:      BackendName(opts.Name),
			Namespace: opts.Namespace,
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: l,
			},
			MinAvailable: &ma,
		},
	}
}
