package notaryserverresources

import (
	"context"
	"encoding/json"
	"path"
	"time"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
	"github.com/goharbor/harbor-operator/pkg/factories/logger"
)

const (
	initImage       = "hairyhenderson/gomplate"
	apiPort         = 6060 // https://github.com/quay/clair/blob/c39101e9b8206401d8b9cb631f3aee47a24ab889/cmd/clair/config.go#L64
	healthPort      = 6061 // https://github.com/quay/clair/blob/c39101e9b8206401d8b9cb631f3aee47a24ab889/cmd/clair/config.go#L63
	adapterPort     = 8080
	clairConfigPath = "/etc/clair"

	livenessProbeInitialDelay = 300 * time.Second
)

var (
	revisionHistoryLimit int32 = 0 // nolint:golint
	varFalse                   = false
)

func (m *Manager) GetDeployments(ctx context.Context) ([]*appsv1.Deployment, error) { // nolint:funlen
	operatorName := application.GetName(ctx)
	harborName := m.Clair.GetName()

	image, err := m.Clair.Spec.GetImage()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get image")
	}

	adapterImage, err := m.Clair.Spec.GetAdapterImage()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get adapter image")
	}

	vulnsrc, err := json.Marshal(m.Clair.Spec.VulnerabilitySources)
	if err != nil {
		logger.Get(ctx).Error(err, "invalid vulnerability sources")
	}

	return []*appsv1.Deployment{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Clair.Name,
				Namespace: m.Clair.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.ClairName,
					"harbor":   harborName,
					"operator": operatorName,
				},
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app":      goharborv1alpha2.ClairName,
						"harbor":   harborName,
						"operator": operatorName,
					},
				},
				Replicas: m.Clair.Spec.Replicas,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"configuration/checksum": m.GetConfigMapsCheckSum(),
							"secret/checksum":        m.GetSecretsCheckSum(),
							"operator/version":       application.GetVersion(ctx),
						},
						Labels: map[string]string{
							"app":      goharborv1alpha2.ClairName,
							"harbor":   harborName,
							"operator": operatorName,
						},
					},
					Spec: corev1.PodSpec{
						NodeSelector:                 m.Clair.Spec.NodeSelector,
						AutomountServiceAccountToken: &varFalse,
						Volumes: []corev1.Volume{
							{
								Name: "config-template",
								VolumeSource: corev1.VolumeSource{
									ConfigMap: &corev1.ConfigMapVolumeSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: m.Clair.Name,
										},
										Items: []corev1.KeyToPath{
											{
												Key:  configKey,
												Path: configKey,
											},
										},
									},
								},
							}, {
								Name:         "config",
								VolumeSource: corev1.VolumeSource{},
							},
						},
						InitContainers: []corev1.Container{
							{
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
								Env: []corev1.EnvVar{
									{
										Name:  "vulnsrc",
										Value: string(vulnsrc),
									},
								},
								EnvFrom: []corev1.EnvFromSource{
									{
										SecretRef: &corev1.SecretEnvSource{
											Optional: &varFalse,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: m.Clair.Spec.DatabaseSecret,
											},
										},
									},
								},
							},
						},
						Containers: []corev1.Container{
							{
								Name:  "clair",
								Image: image,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: apiPort,
									}, {
										ContainerPort: healthPort,
									},
								},

								Env: []corev1.EnvVar{
									{ // // https://github.com/goharbor/harbor/blob/master/make/photon/prepare/templates/clair/clair_env.jinja
										Name:  "HTTP_PROXY",
										Value: "",
									}, {
										Name:  "HTTPS_PROXY",
										Value: "",
									}, {
										Name:  "NO_PROXY",
										Value: "",
									}, { // https://github.com/goharbor/harbor/blob/master/make/photon/prepare/templates/clair/postgres_env.jinja
										Name: "POSTGRES_PASSWORD",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												Key:      goharborv1alpha2.HarborClairDatabasePasswordKey,
												Optional: &varFalse,
												LocalObjectReference: corev1.LocalObjectReference{
													Name: m.Clair.Spec.DatabaseSecret,
												},
											},
										},
									},
								},
								Command:         []string{"/home/clair/clair"},
								Args:            []string{"-config", path.Join(clairConfigPath, configKey)},
								ImagePullPolicy: corev1.PullAlways,
								LivenessProbe: &corev1.Probe{
									Handler: corev1.Handler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: "/health",
											Port: intstr.FromInt(healthPort),
										},
									},
								},
								ReadinessProbe: &corev1.Probe{
									Handler: corev1.Handler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: "/health",
											Port: intstr.FromInt(healthPort),
										},
									},
								},
								VolumeMounts: []corev1.VolumeMount{
									{
										MountPath: path.Join(clairConfigPath, configKey),
										Name:      "config",
										SubPath:   configKey,
									},
								},
							}, {
								Name:  "clair-adapter",
								Image: adapterImage,
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: adapterPort,
									},
								},

								Env: []corev1.EnvVar{
									{
										Name: "SCANNER_STORE_REDIS_URL",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												Key:      goharborv1alpha2.HarborClairAdapterBrokerURLKey,
												Optional: &varFalse,
												LocalObjectReference: corev1.LocalObjectReference{
													Name: m.Clair.Spec.Adapter.RedisSecret,
												},
											},
										},
									}, {
										Name: "SCANNER_STORE_REDIS_NAMESPACE",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												Key:      goharborv1alpha2.HarborClairAdapterBrokerNamespaceKey,
												Optional: &varFalse,
												LocalObjectReference: corev1.LocalObjectReference{
													Name: m.Clair.Spec.Adapter.RedisSecret,
												},
											},
										},
									},
								},
								EnvFrom: []corev1.EnvFromSource{
									{
										Prefix: "clair_db_",
										SecretRef: &corev1.SecretEnvSource{
											Optional: &varFalse,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: m.Clair.Spec.DatabaseSecret,
											},
										},
									}, {
										ConfigMapRef: &corev1.ConfigMapEnvSource{
											Optional: &varFalse,
											LocalObjectReference: corev1.LocalObjectReference{
												Name: m.Clair.Name,
											},
										},
									},
								},

								ImagePullPolicy: corev1.PullAlways,
								LivenessProbe: &corev1.Probe{
									Handler: corev1.Handler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: "/probe/healthy",
											Port: intstr.FromInt(adapterPort),
										},
									},
									InitialDelaySeconds: int32(livenessProbeInitialDelay.Seconds()),
								},
								ReadinessProbe: &corev1.Probe{
									Handler: corev1.Handler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: "/probe/healthy",
											Port: intstr.FromInt(adapterPort),
										},
									},
								},
								VolumeMounts: []corev1.VolumeMount{
									{
										MountPath: path.Join(clairConfigPath, configKey),
										Name:      "config",
										SubPath:   configKey,
									},
								},
							},
						},
						Priority: m.Clair.Spec.Priority,
					},
				},
				RevisionHistoryLimit: &revisionHistoryLimit,
			},
		},
	}, nil
}
