package manifests

import (
	"k8s.io/apimachinery/pkg/labels"
)

const (
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

	kubernetesComponentLabel    = "app.kubernetes.io/component"
	kubernetesInstanceLabel     = "app.kubernetes.io/instance"
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