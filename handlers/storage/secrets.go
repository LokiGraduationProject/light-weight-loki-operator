package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/external/k8s"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/storage"
)

func getSecrets(ctx context.Context, k k8s.Client, stack *lokiv1.LokiStack) (*corev1.Secret, error) {
	var (
		storageSecret corev1.Secret
	)

	key := client.ObjectKey{Name: stack.Spec.Storage.Secret.Name, Namespace: stack.Namespace}
	if err := k.Get(ctx, key, &storageSecret); err != nil {
		return nil, fmt.Errorf("ERROR: lokistack storage secret: %w", err)
	}

	return &storageSecret, nil
}

func extractSecrets(secretSpec lokiv1.ObjectStorageSecretSpec, objStore *corev1.Secret) (storage.Options, error) {
	storageOpts := storage.Options{
		SecretName:  objStore.Name,
		SharedStore: secretSpec.Type,
	}

	temp, err := extractS3ConfigSecret(objStore)

	storageOpts.S3 = temp

	if err != nil {
		return storage.Options{}, err
	}

	return storageOpts, nil
}

func extractS3ConfigSecret(s *corev1.Secret) (*storage.S3StorageConfig, error) {
	buckets := s.Data["bucketnames"]
	if len(buckets) == 0 {
		return nil, fmt.Errorf("%w: %s", errors.New("missing secret field"), "bucketnames")
	}

	var (
		endpoint = s.Data["endpoint"]
		id       = s.Data["access_key_id"]
		secret   = s.Data["access_key_secret"]
	)

	cfg := &storage.S3StorageConfig{
		Buckets: string(buckets),
	}

	cfg.Endpoint = string(endpoint)

	if len(id) == 0 {
		return nil, fmt.Errorf("%w: %s", errors.New("missing secret field"), "access_key_id")
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("%w: %s", errors.New("missing secret field"), "access_key_secret")
	}

	return cfg, nil
}
