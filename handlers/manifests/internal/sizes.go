package internal

import (
	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ComponentResources struct {
	Write ResourceRequirements
	Read  corev1.ResourceRequirements
}

type ResourceRequirements struct {
	Limits          corev1.ResourceList
	Requests        corev1.ResourceList
	PVCSize         resource.Quantity
	PDBMinAvailable int
}

var StackSizeTable = map[lokiv1.LokiStackSizeType]lokiv1.LokiStackSpec{
	lokiv1.SizeOneXDemo: {
		Size: lokiv1.SizeOneXDemo,
		Limits: &lokiv1.LimitsSpec{
			Global: &lokiv1.LimitsTemplateSpec{
				IngestionLimits: &lokiv1.IngestionLimitSpec{
					IngestionRate:           4,
					IngestionBurstSize:      6,
					MaxLabelNameLength:      1024,
					MaxLabelValueLength:     2048,
					MaxLabelNamesPerSeries:  30,
					MaxLineSize:             256000,
					PerStreamDesiredRate:    3,
					PerStreamRateLimit:      5,
					PerStreamRateLimitBurst: 15,
				},
				QueryLimits: &lokiv1.QueryLimitSpec{
					MaxEntriesLimitPerQuery: 5000,
					MaxChunksPerQuery:       2000000,
					MaxQuerySeries:          500,
					QueryTimeout:            "3m",
					CardinalityLimit:        100000,
					MaxVolumeSeries:         1000,
				},
			},
		},
	},
}

var ResourceRequirementsTable = map[lokiv1.LokiStackSizeType]ComponentResources{
	lokiv1.SizeOneXDemo: {
		Write: ResourceRequirements{
			PVCSize: resource.MustParse("10Gi"),
		},
	},
}
