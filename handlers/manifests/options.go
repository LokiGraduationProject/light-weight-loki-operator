package manifests

import (
	"time"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/internal"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/internal/config"
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

	Timeouts TimeoutConfig

	ObjectStorage storage.Options
}

// TimeoutConfig contains the server configuration options for all Loki components
type TimeoutConfig struct {
	Loki config.HTTPTimeoutConfig
}

func calculateHTTPTimeouts(queryTimeout time.Duration) TimeoutConfig {
	idleTimeout := lokiDefaultHTTPIdleTimeout
	if queryTimeout < idleTimeout {
		idleTimeout = queryTimeout
	}

	readTimeout := queryTimeout / 10
	writeTimeout := queryTimeout + lokiQueryWriteDuration

	return TimeoutConfig{
		Loki: config.HTTPTimeoutConfig{
			IdleTimeout:  idleTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}
