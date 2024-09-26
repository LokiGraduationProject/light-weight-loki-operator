package manifests

import (
	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/internal"
	"github.com/ViaQ/logerr/kverrors"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BuildAll builds all manifests required to run a Loki Stack
func BuildAll(opts Options, log logr.Logger) ([]client.Object, error) {
	ll := log.WithValues("lokistack", "default", "event", "BuildAll")

	ll.Info("B")
	res := make([]client.Object, 0)

	ll.Info("C")
	sa := BuildServiceAccount(opts)

	ll.Info("D")
	cm, sha1C, mapErr := LokiConfigMap(opts, log)
	if mapErr != nil {
		return nil, mapErr
	}
	opts.ConfigSHA1 = sha1C

	ll.Info("E")
	distributorObjs, err := BuildDistributor(opts, log)
	if err != nil {
		return nil, err
	}

	ll.Info("F")
	ingesterObjs, err := BuildIngester(opts)
	if err != nil {
		return nil, err
	}

	ll.Info("G")
	querierObjs, err := BuildQuerier(opts)
	if err != nil {
		return nil, err
	}

	ll.Info("H")
	compactorObjs, err := BuildCompactor(opts)
	if err != nil {
		return nil, err
	}

	ll.Info("I")
	queryFrontendObjs, err := BuildQueryFrontend(opts, log)
	if err != nil {
		return nil, err
	}

	ll.Info("J")
	indexGatewayObjs, err := BuildIndexGateway(opts)
	if err != nil {
		return nil, err
	}

	res = append(res, cm)
	res = append(res, sa)
	res = append(res, distributorObjs...)
	res = append(res, ingesterObjs...)
	res = append(res, querierObjs...)
	res = append(res, compactorObjs...)
	res = append(res, queryFrontendObjs...)
	res = append(res, indexGatewayObjs...)
	res = append(res, BuildLokiGossipRingService(opts.Name))

	return res, nil
}

func DefaultLokiStackSpec(size lokiv1.LokiStackSizeType) *lokiv1.LokiStackSpec {
	defaults := internal.StackSizeTable[size]
	return (&defaults).DeepCopy()
}

func ApplyDefaultSettings(opts *Options) error {
	spec := DefaultLokiStackSpec(opts.Stack.Size)

	if err := mergo.Merge(spec, opts.Stack, mergo.WithOverride); err != nil {
		return kverrors.Wrap(err, "failed merging stack user options", "name", opts.Name)
	}

	opts.ResourceRequirements = internal.ResourceRequirementsTable[opts.Stack.Size]
	opts.Stack = *spec

	return nil
}
