package manifests

import (
	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/internal"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/storage"
)

// Options is a set of configuration values to use when building manifests such as resource sizes, etc.
// Most of this should be provided - either directly or indirectly - by the user.
type Options struct {
	Name                   string
	Namespace              string
	Image                  string
	GatewayImage           string
	GatewayBaseDomain      string
	ConfigSHA1             string
	CertRotationRequiredAt string

	Stack                lokiv1.LokiStackSpec
	ResourceRequirements internal.ComponentResources

	ObjectStorage storage.Options
}
