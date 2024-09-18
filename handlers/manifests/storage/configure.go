package storage

import (
	lokiv1 "github.com/LokiGraduationProject/light-weight-loki-operator/api/v1"
	"github.com/ViaQ/logerr/v2/kverrors"
	"github.com/imdario/mergo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// ConfigureDeployment appends additional pod volumes and container env vars, args, volume mounts
// based on the object storage type. Currently supported amendments:
// - All: Ensure object storage secret mounted and auth projected as env vars.
// - GCS: Ensure env var GOOGLE_APPLICATION_CREDENTIALS in container
// - S3: Ensure mounting custom CA configmap if any TLSConfig given
func ConfigureDeployment(d *appsv1.Deployment, opts Options) error {
	switch opts.SharedStore {
	case lokiv1.ObjectStorageSecretS3:
		err := configureDeployment(d, opts)
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

func ConfigureStatefulSet(d *appsv1.StatefulSet, opts Options) error {
	switch opts.SharedStore {
	case lokiv1.ObjectStorageSecretS3:
		if err := configureStatefulSet(d, opts); err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

// ConfigureDeployment merges the object storage secret volume into the deployment spec.
// With this, the deployment will expose credentials specific environment variables.
func configureDeployment(d *appsv1.Deployment, opts Options) error {
	p := ensureObjectStoreCredentials(&d.Spec.Template.Spec, opts)
	if err := mergo.Merge(&d.Spec.Template.Spec, p, mergo.WithOverride); err != nil {
		return kverrors.Wrap(err, "failed to merge gcs object storage spec ")
	}

	return nil
}

// ConfigureStatefulSet merges a the object storage secrect volume into the statefulset spec.
// With this, the statefulset will expose credentials specific environment variable.
func configureStatefulSet(s *appsv1.StatefulSet, opts Options) error {
	p := ensureObjectStoreCredentials(&s.Spec.Template.Spec, opts)
	if err := mergo.Merge(&s.Spec.Template.Spec, p, mergo.WithOverride); err != nil {
		return kverrors.Wrap(err, "failed to merge gcs object storage spec ")
	}

	return nil
}

func ensureObjectStoreCredentials(p *corev1.PodSpec, opts Options) corev1.PodSpec {
	container := p.Containers[0].DeepCopy()
	volumes := p.Volumes
	secretName := opts.SecretName

	volumes = append(volumes, corev1.Volume{
		Name: secretName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
			},
		},
	})

	container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
		Name:      secretName,
		ReadOnly:  false,
		MountPath: secretDirectory,
	})

	return corev1.PodSpec{
		Containers: []corev1.Container{
			*container,
		},
		Volumes: volumes,
	}
}
