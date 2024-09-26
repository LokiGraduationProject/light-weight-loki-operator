package manifests

import (
	"crypto/sha1"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/internal/config"
)

// LokiConfigMap creates the single configmap containing the loki configuration for the whole cluster
func LokiConfigMap(opt Options) (*corev1.ConfigMap, string, error) {
	cfg := ConfigOptions(opt)

	c, err := config.Build(cfg)
	if err != nil {
		return nil, "", err
	}

	s := sha1.New()
	_, err = s.Write(c)
	if err != nil {
		return nil, "", err
	}

	sha1C := fmt.Sprintf("%x", s.Sum(nil))

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   fmt.Sprintf("%s-config", opt.Name),
			Labels: commonLabels(opt.Name),
		},
		Data: map[string]string{
			config.LokiConfigFileName: string(c),
		},
	}, sha1C, nil
}

func ConfigOptions(opt Options) config.Options {

	protocol := "http"

	shippers := []string{}
	boltdb := false
	tsdb := false
	for _, schema := range opt.Stack.Storage.Schemas {
		if !boltdb && (schema.Version == lokiv1.ObjectStorageSchemaV11 || schema.Version == lokiv1.ObjectStorageSchemaV12) {
			shippers = append(shippers, "boltdb")
			boltdb = true
		} else if !tsdb {
			shippers = append(shippers, "tsdb")
			tsdb = true
		}
	}

	return config.Options{
		Stack:     opt.Stack,
		Namespace: opt.Namespace,
		Name:      opt.Name,
		Compactor: config.Address{
			FQDN: fqdn(NewCompactorGRPCService(opt).GetName(), opt.Namespace),
			Port: grpcPort,
		},
		FrontendWorker: config.Address{
			FQDN: fqdn(NewQueryFrontendGRPCService(opt).GetName(), opt.Namespace),
			Port: grpcPort,
		},
		GossipRing: gossipRingConfig(opt.Name, opt.Namespace, opt.Stack.HashRing, opt.Stack.Replication),
		Querier: config.Address{
			Protocol: protocol,
			FQDN:     fqdn(NewQuerierHTTPService(opt).GetName(), opt.Namespace),
			Port:     httpPort,
		},
		IndexGateway: config.Address{
			FQDN: fqdn(NewIndexGatewayGRPCService(opt).GetName(), opt.Namespace),
			Port: grpcPort,
		},
		StorageDirectory: dataDirectory,
		MaxConcurrent: config.MaxConcurrent{
			AvailableQuerierCPUCores: int32(opt.ResourceRequirements.Querier.Requests.Cpu().Value()),
		},
		Shippers:      shippers,
		ObjectStorage: opt.ObjectStorage,
		HTTPTimeouts:  opt.Timeouts.Loki,
	}
}

var deleteWorkerCountMap = map[lokiv1.LokiStackSizeType]uint{
	lokiv1.SizeOneXDemo:       10,
	lokiv1.SizeOneXExtraSmall: 10,
	lokiv1.SizeOneXSmall:      150,
	lokiv1.SizeOneXMedium:     150,
}

func gossipRingConfig(stackName, stackNs string, spec *lokiv1.HashRingSpec, replication *lokiv1.ReplicationSpec) config.GossipRing {
	var (
		instanceAddr string
		enableIPv6   bool
	)
	if spec != nil && spec.Type == lokiv1.HashRingMemberList && spec.MemberList != nil {
		switch spec.MemberList.InstanceAddrType {
		case lokiv1.InstanceAddrPodIP:
			instanceAddr = gossipInstanceAddrEnvVarTemplate
		case lokiv1.InstanceAddrDefault:
			// Do nothing use loki defaults
		default:
			// Do nothing use loki defaults
		}

		if spec.MemberList.EnableIPv6 {
			enableIPv6 = true
			instanceAddr = gossipInstanceAddrEnvVarTemplate
		}
	}

	return config.GossipRing{
		EnableIPv6:                     enableIPv6,
		InstanceAddr:                   instanceAddr,
		InstancePort:                   grpcPort,
		BindPort:                       gossipPort,
		MembersDiscoveryAddr:           fqdn(BuildLokiGossipRingService(stackName).GetName(), stackNs),
		EnableInstanceAvailabilityZone: replication != nil && len(replication.Zones) > 0,
	}
}
