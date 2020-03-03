package notaryserverresources

import (
	"context"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/goharbor/harbor-operator/pkg/factories/application"
	"github.com/pkg/errors"
)

const (
	migrationDatabaseURL = "postgresql://$(username):$(password)@$(host):$(port)/$(database)?sslmode=$(ssl)"
	initImage            = "hairyhenderson/gomplate"
	notaryServerPort     = 4443
	notarySignerPort     = 7899
)

var (
	revisionHistoryLimit     int32 = 0 // nolint:golint
	varFalse                       = false
	notarySignerKeyAlgorithm       = "ecdsa"
)

func (m *Manager) GetDeployments(ctx context.Context) ([]*appsv1.Deployment, error) { // nolint:funlen
	operatorName := application.GetName(ctx)

	image, err := m.Notary.Spec.GetImage()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get image")
	}

	dbMigratorImage, err := m.Notary.Spec.GetDBMigratorImage()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get dbMigrator image")
	}

	return []*appsv1.Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Notary.Name,
				Namespace: m.Notary.Namespace,
				Labels: map[string]string{
					"app":      m.Notary.Name,
					"operator": operatorName,
				},
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app":      m.Notary.Name,
						"operator": operatorName,
					},
				},
				Replicas: m.Notary.Spec.Replicas,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"configuration/checksum": m.GetConfigMapsCheckSum(),
							"secret/checksum":        m.GetSecretsCheckSum(),
							"operator/version":       application.GetVersion(ctx),
						},
						Labels: map[string]string{
							"app":      m.Notary.Name,
							"operator": operatorName,
						},
					},
					Spec: corev1.PodSpec{
						NodeSelector:                 m.Notary.Spec.NodeSelector,
						AutomountServiceAccountToken: &varFalse,
						Volumes: []corev1.Volume{
							{
								Name: "config-template",
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: m.Notary.Name,
										},
									},
								},
							}, {
								Name:         "config",
								VolumeSource: corev1.VolumeSource{},
							}, {
								Name: "notary-certificate",
								VolumeSource: corev1.VolumeSource{
									Secret: &corev1.SecretVolumeSource{
										SecretName: m.Notary.Spec.CertificateSecret,
									},
								},
							},
						},
						InitContainers: []corev1.Container{
							{
								Name:  "init-db",
								Image: dbMigratorImage,
								Args:  []string{"-c", "signer", "-d", migrationDatabaseURL},
								EnvFrom: []corev1.EnvFromSource{
									{
										SecretRef: &corev1.SecretEnvSource{
											Optional: &varFalse,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: m.Notary.Spec.DatabaseSecret,
											},
										},
									},
								},
							}, {
								Name:       "configuration",
								Image:      initImage,
								WorkingDir: "/workdir",
								Args:       []string{"--input-dir", "/workdir", "--output-dir", "/processed"},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "config-template",
										MountPath: "/workdir",
										ReadOnly:  true,
									}, {
										Name:      "config",
										MountPath: "/processed",
										ReadOnly:  false,
									},
								},
								EnvFrom: []corev1.EnvFromSource{
									{
										SecretRef: &corev1.SecretEnvSource{
											Optional: &varFalse,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: m.Notary.Spec.DatabaseSecret,
											},
										},
									},
								},
								Env: []corev1.EnvVar{
									{
										Name:  "notary_signer_port",
										Value: strconv.Itoa(notarySignerPort),
									},
								},
							},
						},
						Containers: []corev1.Container{
							{
								Name:            "notary-signer",
								Image:           image,
								Args:            []string{"notary-signer", "-config", "/etc/notary/signer-config.json"},
								ImagePullPolicy: corev1.PullAlways,
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "notary-certificate",
										MountPath: "/etc/ssl/notary/tls.crt",
										SubPath:   "tls.crt",
									}, {
										Name:      "notary-certificate",
										MountPath: "/etc/ssl/notary/tls.key",
										SubPath:   "tls.key",
									}, {
										Name:      "config",
										MountPath: "/etc/notary/signer-config.json",
										SubPath:   signerConfigKey,
									},
								},
								Env: []corev1.EnvVar{
									{
										Name:  "NOTARY_SIGNER_DEFAULTALIAS",
										Value: "defaultalias",
									},
								},
							},
						},
						Priority: m.Notary.Spec.Priority,
					},
				},
				RevisionHistoryLimit: &revisionHistoryLimit,
			},
		},
	}, nil
}
