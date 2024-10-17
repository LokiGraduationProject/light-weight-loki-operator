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
	Read_http             Address
	Read_grpc             Address
	Write                 Address
	Backend               Address
	GossipRing            GossipRing
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

// GossipRing defines the memberlist configuration
type GossipRing struct {
	// EnableIPv6 is optional, memberlist IPv6 support
	EnableIPv6 bool
	// InstanceAddr is optional, defaults to private networks
	InstanceAddr string
	// InstancePort is required
	InstancePort int
	// BindPort is the port for listening to gossip messages
	BindPort int
	// MembersDiscoveryAddr is required
	MembersDiscoveryAddr           string
	EnableInstanceAvailabilityZone bool
}
