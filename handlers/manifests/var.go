package manifests

import (
	"fmt"
	"time"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	httpPort                 = 3100
	internalHTTPPort         = 3101
	grpcPort                 = 9095
	lokiHTTPPortName         = "metrics"
	lokiInternalHTTPPortName = "healthchecks"
	lokiGRPCPortName         = "grpclb"
	protocolTCP              = "TCP"
	lokiLivenessPath         = "/loki/api/v1/status/buildinfo"
	lokiReadinessPath        = "/ready"
	configVolumeName         = "config"

	gossipPort                       = 7946
	gossipInstanceAddrEnvVarName     = "HASH_RING_INSTANCE_ADDR"
	gossipInstanceAddrEnvVarTemplate = "${" + gossipInstanceAddrEnvVarName + "}"
	lokiGossipPortName               = "gossip-ring"

	// AnnotationLokiConfigHash stores the last SHA1 hash of the loki configuration
	AnnotationLokiConfigHash string = "loki.grafana.com/config-hash"
	// AnnotationLokiObjectStoreHash stores the last SHA1 hash of the loki object storage credetials.
	AnnotationLokiObjectStoreHash string = "loki.grafana.com/object-store-hash"

	// LabelCompactorComponent is the label value for the compactor component
	LabelCompactorComponent string = "compactor"
	// LabelDistributorComponent is the label value for the distributor component
	LabelDistributorComponent string = "distributor"
	// LabelIngesterComponent is the label value for the ingester component
	LabelIngesterComponent string = "ingester"
	// LabelQuerierComponent is the label value for the querier component
	LabelQuerierComponent string = "querier"
	// LabelQueryFrontendComponent is the label value for the query frontend component
	LabelQueryFrontendComponent string = "query-frontend"
	// LabelIndexGatewayComponent is the label value for the lokiStack-index-gateway component
	LabelIndexGatewayComponent string = "index-gateway"
	// LabelRulerComponent is the label value for the lokiStack-ruler component
	LabelRulerComponent string = "ruler"
	// LabelGatewayComponent is the label value for the lokiStack-gateway component
	LabelGatewayComponent string = "lokistack-gateway"

	// EnvRelatedImageLoki is the environment variable to fetch the Loki image pullspec.
	EnvRelatedImageLoki = "RELATED_IMAGE_LOKI"

	// DefaultContainerImage declares the default fallback for loki image.
	DefaultContainerImage = "docker.io/grafana/loki:3.1.1"

	kubernetesComponentLabel = "app.kubernetes.io/component"
	kubernetesInstanceLabel  = "app.kubernetes.io/instance"

	lokiFrontendContainerName = "loki-query-frontend"

	kubernetesNodeOSLabel       = "kubernetes.io/os"
	kubernetesNodeOSLinux       = "linux"
	kubernetesNodeHostnameLabel = "kubernetes.io/hostname"

	dataDirectory     = "/tmp/loki"
	storageVolumeName = "storage"

	saTokenVolumeName            = "bound-sa-token"
	saTokenExpiration      int64 = 3600
	saTokenVolumeMountPath       = "/var/run/secrets/storage/serviceaccount"

	ServiceAccountTokenFilePath = saTokenVolumeMountPath + "/token"

	secretDirectory  = "/etc/storage/secrets"
	storageTLSVolume = "storage-tls"
	caDirectory      = "/etc/storage/ca"

	tokenAuthConfigVolumeName = "token-auth-config"
	tokenAuthConfigDirectory  = "/etc/storage/token-auth"
)

// GossipLabels is the list of labels that should be assigned to components using the gossip ring
func GossipLabels() map[string]string {
	return map[string]string{
		"loki.grafana.com/gossip": "true",
	}
}

const (
	// lokiDefaultQueryTimeout contains the default query timeout. It should match the value mentioned in the CRD
	// definition and also the default in the `sizes.go`.
	lokiDefaultQueryTimeout    = 3 * time.Minute
	lokiDefaultHTTPIdleTimeout = 30 * time.Second
	lokiQueryWriteDuration     = 1 * time.Minute

	gatewayReadDuration  = 30 * time.Second
	gatewayWriteDuration = 2 * time.Minute
)

var (
	defaultTimeoutConfig = calculateHTTPTimeouts(lokiDefaultQueryTimeout)

	defaultRevHistoryLimit int32 = 10
	defaultConfigMapMode   int32 = 420
	volumeFileSystemMode         = corev1.PersistentVolumeFilesystem
)

// ComponentLabels is a list of all commonLabels including the app.kubernetes.io/component:<component> label
func ComponentLabels(component, stackName string) labels.Set {
	return labels.Merge(commonLabels(stackName), map[string]string{
		kubernetesComponentLabel: component,
	})
}

func commonLabels(stackName string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "lokistack",
		kubernetesInstanceLabel:        stackName,
		"app.kubernetes.io/managed-by": "lokistack-controller",
		"app.kubernetes.io/created-by": "lokistack-controller",
	}
}

// configureAffinity returns an Affinity struture that can be used directly
// in a Deployment/StatefulSet. Parameters will affected configuration of the
// different fields in Affinity (NodeAffinity, PodAffinity, PodAntiAffinity).
func configureAffinity(componentLabel, stackName string, enableNodeAffinity bool, cSpec *lokiv1.LokiComponentSpec) *corev1.Affinity {
	affinity := &corev1.Affinity{
		NodeAffinity:    defaultNodeAffinity(enableNodeAffinity),
		PodAntiAffinity: defaultPodAntiAffinity(componentLabel, stackName),
	}
	if cSpec.PodAntiAffinity != nil {
		affinity.PodAntiAffinity = cSpec.PodAntiAffinity
	}
	return affinity
}

// defaultNodeAffinity if enabled will require pods to run on Linux nodes
func defaultNodeAffinity(enableNodeAffinity bool) *corev1.NodeAffinity {
	if !enableNodeAffinity {
		return nil
	}

	return &corev1.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
			NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Key:      kubernetesNodeOSLabel,
							Operator: corev1.NodeSelectorOpIn,
							Values: []string{
								kubernetesNodeOSLinux,
							},
						},
					},
				},
			},
		},
	}
}

// defaultPodAntiAffinity for components in podAntiAffinityComponents will
// configure pods, of a LokiStack, to preferably not run on the same node
func defaultPodAntiAffinity(componentLabel, stackName string) *corev1.PodAntiAffinity {
	return &corev1.PodAntiAffinity{
		PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
			{
				Weight: 100,
				PodAffinityTerm: corev1.PodAffinityTerm{
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: componentInstanceLabels(componentLabel, stackName),
					},
					TopologyKey: kubernetesNodeHostnameLabel,
				},
			},
		},
	}
}

func lokiLivenessProbe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   lokiLivenessPath,
				Port:   intstr.FromInt(3100),
				Scheme: corev1.URISchemeHTTP,
			},
		},
		TimeoutSeconds:   2,
		PeriodSeconds:    30,
		FailureThreshold: 10,
		SuccessThreshold: 1,
	}
}

func lokiReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   lokiReadinessPath,
				Port:   intstr.FromInt(3100),
				Scheme: corev1.URISchemeHTTP,
			},
		},
		PeriodSeconds:       10,
		InitialDelaySeconds: 15,
		TimeoutSeconds:      1,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
}

func componentInstanceLabels(component string, stackName string) map[string]string {
	return map[string]string{
		kubernetesInstanceLabel:  stackName,
		kubernetesComponentLabel: component,
	}
}

func commonAnnotations(opts Options) map[string]string {
	a := map[string]string{
		AnnotationLokiConfigHash: opts.ConfigSHA1,
	}

	return a
}

func lokiConfigMapName(stackName string) string {
	return fmt.Sprintf("%s-config", stackName)
}

// CompactorName is the name of the compactor statefulset
func CompactorName(stackName string) string {
	return fmt.Sprintf("%s-compactor", stackName)
}

// DistributorName is the name of the distributor deployment
func DistributorName(stackName string) string {
	return fmt.Sprintf("%s-distributor", stackName)
}

// IngesterName is the name of the compactor statefulset
func IngesterName(stackName string) string {
	return fmt.Sprintf("%s-ingester", stackName)
}

// QuerierName is the name of the querier deployment
func QuerierName(stackName string) string {
	return fmt.Sprintf("%s-querier", stackName)
}

// QueryFrontendName is the name of the query-frontend statefulset
func QueryFrontendName(stackName string) string {
	return fmt.Sprintf("%s-query-frontend", stackName)
}

// IndexGatewayName is the name of the index-gateway statefulset
func IndexGatewayName(stackName string) string {
	return fmt.Sprintf("%s-index-gateway", stackName)
}

func serviceNameQuerierHTTP(stackName string) string {
	return fmt.Sprintf("%s-querier-http", stackName)
}

func serviceNameQuerierGRPC(stackName string) string {
	return fmt.Sprintf("%s-querier-grpc", stackName)
}

func serviceNameIngesterGRPC(stackName string) string {
	return fmt.Sprintf("%s-ingester-grpc", stackName)
}

func serviceNameIngesterHTTP(stackName string) string {
	return fmt.Sprintf("%s-ingester-http", stackName)
}

func serviceNameDistributorGRPC(stackName string) string {
	return fmt.Sprintf("%s-distributor-grpc", stackName)
}

func serviceNameDistributorHTTP(stackName string) string {
	return fmt.Sprintf("%s-distributor-http", stackName)
}

func serviceNameCompactorGRPC(stackName string) string {
	return fmt.Sprintf("%s-compactor-grpc", stackName)
}

func serviceNameCompactorHTTP(stackName string) string {
	return fmt.Sprintf("%s-compactor-http", stackName)
}

func serviceNameQueryFrontendGRPC(stackName string) string {
	return fmt.Sprintf("%s-query-frontend-grpc", stackName)
}

func serviceNameQueryFrontendHTTP(stackName string) string {
	return fmt.Sprintf("%s-query-frontend-http", stackName)
}

func serviceNameIndexGatewayHTTP(stackName string) string {
	return fmt.Sprintf("%s-index-gateway-http", stackName)
}

func serviceNameIndexGatewayGRPC(stackName string) string {
	return fmt.Sprintf("%s-index-gateway-grpc", stackName)
}

func fqdn(serviceName, namespace string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace)
}
