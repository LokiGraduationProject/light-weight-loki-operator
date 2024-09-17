package v1

import (
	"time"
)

// StorageSchemaEffectiveDate defines the type for the Storage Schema Effect Date
//
// +kubebuilder:validation:Pattern:="^([0-9]{4,})([-]([0-9]{2})){2}$"
type StorageSchemaEffectiveDate string

// UTCTime returns the date as a time object in the UTC time zone
func (d StorageSchemaEffectiveDate) UTCTime() (time.Time, error) {
	return time.Parse("2006-01-02", string(d))
}
