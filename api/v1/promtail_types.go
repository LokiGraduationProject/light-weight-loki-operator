package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PromtailSpec defines the desired state of Promtail
type PromtailSpec struct {
	// +optional
	Config PromtailConfig `json:"promtail_config" yaml:"promtail_config"`
}

// PromtailConfig defines the configuration of the promtail deployment
type PromtailConfig struct {
	// +optional
	Server ServerConfig `json:"server,omitempty" yaml:"server,omitempty"`

	// +optional
	Clients []ClientConfig `json:"clients,omitempty" yaml:"clients,omitempty"`

	// +optional
	Positions PositionsConfig `json:"positions,omitempty" yaml:"positions,omitempty"`

	// +optional
	TargetConfig TargetConfigConfig `json:"target_config,omitempty" yaml:"target_config,omitempty"`

	// +optional
	ScrapeConfigs []ScrapeConfigsConfig `json:"scrape_configs,omitempty" yaml:"scrape_configs,omitempty"`
}

type ServerConfig struct {
	Disable                        bool   `json:"disable,omitempty" yaml:"disable,omitempty"`
	ProfilingEnabled               bool   `json:"profiling_enabled,omitempty" yaml:"profiling_enabled,omitempty"`
	HttpListenAddress              string `json:"http_listen_address,omitempty" yaml:"http_listen_address,omitempty"`
	HttpListenPort                 int    `json:"http_listen_port,omitempty" yaml:"http_listen_port,omitempty"`
	GrpcListenAddress              string `json:"grpc_listen_address,omitempty" yaml:"grpc_listen_address,omitempty"`
	GrpcListenPort                 int    `json:"grpc_listen_port,omitempty" yaml:"grpc_listen_port,omitempty"`
	RegisterInstrumentation        bool   `json:"register_instrumentation,omitempty" yaml:"register_instrumentation,omitempty"`
	GracefulShutdownTimeout        string `json:"graceful_shutdown_timeout,omitempty" yaml:"graceful_shutdown_timeout,omitempty"`
	HttpServerReadTimeout          string `json:"http_server_read_timeout,omitempty" yaml:"http_server_read_timeout,omitempty"`
	HttpServerWriteTimeout         string `json:"http_server_write_timeout,omitempty" yaml:"http_server_write_timeout,omitempty"`
	HttpServerIdleTimeout          string `json:"http_server_idle_timeout,omitempty" yaml:"http_server_idle_timeout,omitempty"`
	GrpcServerMaxRecvMsgSize       int    `json:"grpc_server_max_recv_msg_size,omitempty" yaml:"grpc_server_max_recv_msg_size,omitempty"`
	GrpcServerMaxSendMsgSize       int    `json:"grpc_server_max_send_msg_size,omitempty" yaml:"grpc_server_max_send_msg_size,omitempty"`
	GrpcServerMaxConcurrentStreams int    `json:"grpc_server_max_concurrent_streams,omitempty" yaml:"grpc_server_max_concurrent_streams,omitempty"`
	LogLevel                       string `json:"log_level,omitempty" yaml:"log_level,omitempty"`
	HttpPathPrefix                 string `json:"http_path_prefix,omitempty" yaml:"http_path_prefix,omitempty"`
	HealthCheckTarget              bool   `json:"health_check_target,omitempty" yaml:"health_check_target,omitempty"`
	EnableRuntimeReload            bool   `json:"enable_runtime_reload,omitempty" yaml:"enable_runtime_reload,omitempty"`
}

type ClientConfig struct {
	// Loki의 HTTP URL
	// +optional
	URL string `json:"url" yaml:"url"`

	// HTTP 헤더 설정
	// +optional
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`

	// 테넌트 ID
	// +optional
	TenantID string `json:"tenant_id,omitempty" yaml:"tenant_id,omitempty"`

	// 최대 대기 시간
	// +optional
	BatchWait string `json:"batchwait,omitempty" yaml:"batchwait,omitempty"`

	// 최대 배치 크기
	// +optional
	BatchSize int `json:"batchsize,omitempty" yaml:"batchsize,omitempty"`

	// Basic Auth 설정
	// +optional
	BasicAuth *BasicAuthConfig `json:"basic_auth,omitempty" yaml:"basic_auth,omitempty"`

	// OAuth2 설정
	// +optional
	OAuth2 *OAuth2Config `json:"oauth2,omitempty" yaml:"oauth2,omitempty"`

	// Bearer Token 설정
	// +optional
	BearerToken     string `json:"bearer_token,omitempty" yaml:"bearer_token,omitempty"`
	BearerTokenFile string `json:"bearer_token_file,omitempty" yaml:"bearer_token_file,omitempty"`

	// 프록시 설정
	// +optional
	ProxyURL string `json:"proxy_url,omitempty" yaml:"proxy_url,omitempty"`

	// TLS 설정
	// +optional
	TLSConfig *TLSConfig `json:"tls_config,omitempty" yaml:"tls_config,omitempty"`

	// 백오프 설정
	// +optional
	BackoffConfig *BackoffConfig `json:"backoff_config,omitempty" yaml:"backoff_config,omitempty"`

	// Rate Limit 관련 설정
	// +optional
	DropRateLimitedBatches bool `json:"drop_rate_limited_batches,omitempty" yaml:"drop_rate_limited_batches,omitempty"`

	// 외부 레이블 설정
	// +optional
	ExternalLabels map[string]string `json:"external_labels,omitempty" yaml:"external_labels,omitempty"`

	// 요청 타임아웃
	// +optional
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

type BasicAuthConfig struct {
	// +optional
	Username string `json:"username,omitempty" yaml:"username,omitempty"`

	// +optional
	Password string `json:"password,omitempty" yaml:"password,omitempty"`

	// +optional
	PasswordFile string `json:"password_file,omitempty" yaml:"password_file,omitempty"`
}

type OAuth2Config struct {
	// +optional
	ClientID string `json:"client_id,omitempty" yaml:"client_id,omitempty"`

	// +optional
	ClientSecret string `json:"client_secret,omitempty" yaml:"client_secret,omitempty"`

	// +optional
	ClientSecretFile string `json:"client_secret_file,omitempty" yaml:"client_secret_file,omitempty"`

	// +optional
	Scopes []string `json:"scopes,omitempty" yaml:"scopes,omitempty"`

	// +optional
	TokenURL string `json:"token_url,omitempty" yaml:"token_url,omitempty"`

	// +optional
	EndpointParams map[string]string `json:"endpoint_params,omitempty" yaml:"endpoint_params,omitempty"`
}

type TLSConfig struct {
	// +optional
	CAFile string `json:"ca_file,omitempty" yaml:"ca_file,omitempty"`

	// +optional
	CertFile string `json:"cert_file,omitempty" yaml:"cert_file,omitempty"`

	// +optional
	KeyFile string `json:"key_file,omitempty" yaml:"key_file,omitempty"`

	// +optional
	ServerName string `json:"server_name,omitempty" yaml:"server_name,omitempty"`

	// +optional
	InsecureSkipVerify bool `json:"insecure_skip_verify,omitempty" yaml:"insecure_skip_verify,omitempty"`
}

type BackoffConfig struct {
	// +optional
	MinPeriod string `json:"min_period,omitempty" yaml:"min_period,omitempty"`

	// +optional
	MaxPeriod string `json:"max_period,omitempty" yaml:"max_period,omitempty"`

	// +optional
	MaxRetries int `json:"max_retries,omitempty" yaml:"max_retries,omitempty"`
}

type PositionsConfig struct {
	// +optional
	Filename string `json:"filename,omitempty" yaml:"filename,omitempty"`

	// +optional
	SyncPeriod string `json:"sync_period,omitempty" yaml:"sync_period,omitempty"`

	// +optional
	IgnoreInvalidYAML bool `json:"ignore_invalid_yaml,omitempty" yaml:"ignore_invalid_yaml,omitempty"`
}

type TargetConfigConfig struct {
	// +optional
	SyncPeriod string `json:"sync_period,omitempty" yaml:"sync_period,omitempty"`
}

type ScrapeConfigsConfig struct {
	// +optional
	JobName string `json:"job_name" yaml:"job_name"`

	KubernetesSDConfigs []KubernetesSDConfig `json:"kubernetes_sd_configs,omitempty" yaml:"kubernetes_sd_configs,omitempty"`

	PipelineStages []PipelineStage `json:"pipeline_stages,omitempty" yaml:"pipeline_stages,omitempty"`

	RelabelConfigs []RelabelConfigsConfig `json:"relabel_configs,omitempty" yaml:"relabel_configs,omitempty"`
}

type KubernetesSDConfig struct {
	// +optional
	Role string `json:"role" yaml:"role"`
}

type EmptyStruct struct{}

type PipelineStage struct {
	Docker *EmptyStruct `json:"docker,omitempty" yaml:"docker,omitempty"`
	Cri    *EmptyStruct `json:"cri,omitempty" yaml:"cri,omitempty"`
}

type RelabelConfigsConfig struct {
	// Their content is concatenated using the configured separator.
	SourceLabels []string `json:"source_labels,omitempty" yaml:"source_labels,omitempty"`

	// Separator placed between concatenated source label values.
	Separator string `json:"separator,omitempty" yaml:"separator,omitempty"`

	// TargetLabel is the label to which the resulting value is written in a replace action.
	TargetLabel string `json:"target_label,omitempty" yaml:"target_label,omitempty"`

	// Regex is the regular expression against which the extracted value is matched.
	Regex string `json:"regex,omitempty" yaml:"regex,omitempty"`

	// Modulus to take of the hash of the source label values.
	Modulus uint64 `json:"modulus,omitempty" yaml:"modulus,omitempty"`

	// Replacement is the value against which a regex replace is performed if the regular expression matches.
	Replacement string `json:"replacement,omitempty" yaml:"replacement,omitempty"`

	// Action to perform based on regex matching.
	Action string `json:"action,omitempty" yaml:"action,omitempty"`
}

// PromtailStatus defines the observed state of Promtail
type PromtailStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Promtail is the Schema for the promtails API
type Promtail struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PromtailSpec   `json:"spec,omitempty"`
	Status PromtailStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PromtailList contains a list of Promtail
type PromtailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Promtail `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Promtail{}, &PromtailList{})
}
