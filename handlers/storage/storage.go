package storage

import (
	"context"
	"fmt"
	"time"

	// configv1 "github.com/grafana/loki/operator/apis/config/v1"
	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/external/k8s"

	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/storage"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/status"
)

// BuildOptions returns the object storage options to generate Kubernetes resource manifests
// which require access to object storage buckets.
// The returned error can be a status.DegradedError in the following cases:
//   - The user-provided object storage secret is missing.
//   - The object storage Secret data is invalid.
//   - The object storage schema config is invalid.
//   - The object storage CA ConfigMap is missing if one referenced.
//   - The object storage CA ConfigMap data is invalid.
//   - The object storage token cco auth secret is missing (Only on OpenShift STS-clusters)
func BuildOptions(ctx context.Context, k k8s.Client, stack *lokiv1.LokiStack) (storage.Options, error) {
	storageSecret, err := getSecrets(ctx, k, stack)
	if err != nil {
		return storage.Options{}, err
	}

	objStore, err := extractSecrets(stack.Spec.Storage.Secret, storageSecret)
	if err != nil {
		return storage.Options{}, &status.DegradedError{
			Message: fmt.Sprintf("Invalid object storage secret contents: %s", err),
			Reason:  lokiv1.ReasonInvalidObjectStorageSecret,
			Requeue: false,
		}
	}

	now := time.Now().UTC()
	storageSchemas, err := storage.BuildSchemaConfig(
		now,
		stack.Spec.Storage,
		stack.Status.Storage,
	)
	if err != nil {
		return storage.Options{}, &status.DegradedError{
			Message: fmt.Sprintf("Invalid object storage schema contents: %s", err),
			Reason:  lokiv1.ReasonInvalidObjectStorageSchema,
			Requeue: false,
		}
	}

	objStore.Schemas = storageSchemas

	return objStore, nil
}
