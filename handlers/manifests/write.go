package manifests

import (
	"fmt"
	"path"

	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/internal/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/utils/ptr"
)

func newLokiWriteStatefulSet(opts Options) *appsv1.StatefulSet {
	l := ComponentLabels(LabelIngesterComponent, opts.Name)
	a := commonAnnotations(opts)

	replicas := int32(3)

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "loki-write",
			Labels:    l,
		},
		Spec: appsv1.StatefulSetSpec{
			PodManagementPolicy:  appsv1.OrderedReadyPodManagement,
			RevisionHistoryLimit: ptr.To(defaultRevHistoryLimit),
			Replicas:    &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels.Merge(l, GossipLabels()),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fmt.Sprintf("loki-write-%s", opts.Name),
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
							Name:  "loki-write-component",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("1Gi"),
								},
							},
							Args: []string{
								"-target=write",
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
