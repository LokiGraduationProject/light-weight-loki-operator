package storage

import (
	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
)

// Options is used to configure Loki to integrate with
// supported object storages.
type Options struct {
	Schemas                 []lokiv1.ObjectStorageSchema
	SharedStore             lokiv1.ObjectStorageSecretType
	CredentialMode          lokiv1.CredentialMode
	AllowStructuredMetadata bool

	S3 *S3StorageConfig

	SecretName string
}

// S3StorageConfig for S3 storage config
type S3StorageConfig struct {
	Endpoint       string
	Region         string
	Buckets        string
	Audience       string
	STS            bool
	SSE            S3SSEConfig
	ForcePathStyle bool
}

type S3SSEType string

const (
	SSEKMSType S3SSEType = "SSE-KMS"
	SSES3Type  S3SSEType = "SSE-S3"
)

type S3SSEConfig struct {
	Type                 S3SSEType
	KMSKeyID             string
	KMSEncryptionContext string
}
