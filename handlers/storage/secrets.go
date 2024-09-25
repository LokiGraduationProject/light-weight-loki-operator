package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"sort"

	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/external/k8s"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/status"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

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
		if apierrors.IsNotFound(err) {
			return nil, &status.DegradedError{
				Message: "Missing object storage secret",
				Reason:  lokiv1.ReasonMissingObjectStorageSecret,
				Requeue: false,
			}
		}
		return nil, fmt.Errorf("failed to lookup lokistack storage secret: %w", err)
	}

	return &storageSecret, nil
}

func extractSecrets(secretSpec lokiv1.ObjectStorageSecretSpec, objStore *corev1.Secret) (storage.Options, error) {
	hash, err := hashSecretData(objStore)
	if err != nil {
		return storage.Options{}, errors.New("error calculating hash for secret")
	}

	storageOpts := storage.Options{
		SecretName:  objStore.Name,
		SecretSHA1:  hash,
		SharedStore: secretSpec.Type,
	}

	switch secretSpec.Type {
	case lokiv1.ObjectStorageSecretS3:
		storageOpts.S3, err = extractS3ConfigSecret(objStore)
	default:
		return storage.Options{}, fmt.Errorf("%w: %s", errors.New("unknown secret type"), secretSpec.Type)
	}
	if err != nil {
		return storage.Options{}, err
	}

	return storageOpts, nil
}

func hashSecretData(s *corev1.Secret) (string, error) {
	keys := make([]string, 0, len(s.Data))
	for k := range s.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha1.New()
	for _, k := range keys {
		if _, err := h.Write([]byte(k)); err != nil {
			return "", err
		}

		if _, err := h.Write([]byte(",")); err != nil {
			return "", err
		}

		if _, err := h.Write(s.Data[k]); err != nil {
			return "", err
		}

		if _, err := h.Write([]byte(",")); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
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

	if err := validateS3Endpoint(string(endpoint)); err != nil {
		return nil, err
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("%w: %s", errors.New("missing secret field"), "access_key_id")
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("%w: %s", errors.New("missing secret field"), "access_key_secret")
	}

	return cfg, nil
}

func validateS3Endpoint(endpoint string) error {
	if len(endpoint) == 0 {
		return fmt.Errorf("%w: %s", errors.New("missing secret field"), "endpoint")
	}

	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("%w: %w", errors.New("can not parse S3 endpoint as URL"), err)
	}

	if parsedURL.Scheme == "" {
		return errors.New("endpoint for S3 must be an HTTP or HTTPS URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("%w: %s", errors.New("scheme of S3 endpoint URL is unsupported"), parsedURL.Scheme)
	}

	return nil
}
