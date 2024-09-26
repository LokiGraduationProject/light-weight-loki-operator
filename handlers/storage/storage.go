package storage

import (
	"context"
	"fmt"

	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/external/k8s"

	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/manifests/storage"
	"github.com/LokiGraduationProject/light-weight-loki-operator/handlers/status"
)

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

	storageSchemas := storage.BuildSchemas(stack.Spec.Storage.Schemas)

	objStore.Schemas = storageSchemas

	return objStore, nil
}
