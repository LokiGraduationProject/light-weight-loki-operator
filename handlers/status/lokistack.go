package status

import (
	"fmt"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
)

// DegradedError contains information about why the managed LokiStack has an invalid configuration.
type DegradedError struct {
	Message string
	Reason  lokiv1.LokiStackConditionReason
	Requeue bool
}

func (e *DegradedError) Error() string {
	return fmt.Sprintf("cluster degraded: %s", e.Message)
}
