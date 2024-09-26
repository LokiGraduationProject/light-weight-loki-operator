package storage

import (
	"sort"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
)

func BuildSchemas(schemas []lokiv1.ObjectStorageSchema) []lokiv1.ObjectStorageSchema {
	sortedSchemas := make([]lokiv1.ObjectStorageSchema, len(schemas))
	copy(sortedSchemas, schemas)

	sort.SliceStable(sortedSchemas, func(i, j int) bool {
		iDate, _ := sortedSchemas[i].EffectiveDate.UTCTime()
		jDate, _ := sortedSchemas[j].EffectiveDate.UTCTime()

		return iDate.Before(jDate)
	})

	return reduceSortedSchemas(sortedSchemas)
}

func reduceSortedSchemas(schemas []lokiv1.ObjectStorageSchema) []lokiv1.ObjectStorageSchema {
	version := ""
	reduced := []lokiv1.ObjectStorageSchema{}

	for _, schema := range schemas {
		strSchemaVersion := string(schema.Version)

		if version != strSchemaVersion {
			version = strSchemaVersion
			reduced = append(reduced, schema)
		}
	}

	return reduced
}
