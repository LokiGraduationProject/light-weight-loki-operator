package config

import (
	"time"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/storage"
)

// Options is used to render the loki-config.yaml file template
type Options struct {
	Stack lokiv1.LokiStackSpec

	Namespace             string
	Name                  string
	Compactor             Address
	FrontendWorker        Address
	Querier               Address
	IndexGateway          Address
	StorageDirectory      string
	MaxConcurrent         MaxConcurrent
	EnableRemoteReporting bool
	Shippers              []string

	ObjectStorage storage.Options

	HTTPTimeouts HTTPTimeoutConfig
}

// Address FQDN and port for a k8s service.
type Address struct {
	// Protocol is optional
	Protocol string
	// FQDN is required
	FQDN string
	// Port is required
	Port int
}

// HTTPTimeoutConfig defines the HTTP server config options.
type HTTPTimeoutConfig struct {
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// MaxConcurrent for concurrent query processing.
type MaxConcurrent struct {
	AvailableQuerierCPUCores int32
}
