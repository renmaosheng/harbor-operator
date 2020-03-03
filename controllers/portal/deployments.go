package portal

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
	"github.com/pkg/errors"
)

const (
	port = 8080
)

var (
	revisionHistoryLimit int32 = 0 // nolint:golint
	varFalse                   = false
)

func (r *Reconciler) GetDeployments(ctx context.Context) (*appsv1.Deployment, error) { // nolint:funlen
	operatorName := application.GetName(ctx)

	image, err := m.Portal.Spec.GetImage()
	if err != nil {
		return nil, errors.Wrap(err, "cannot get image")
	}

	return *appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Portal.Name,
			Namespace: m.Portal.Namespace,
			Labels: map[string]string{
				"app":      goharborv1alpha2.PortalName,
				"operator": operatorName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":      goharborv1alpha2.PortalName,
					"operator": operatorName,
				},
			},
			Replicas: m.Portal.Spec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"configuration/checksum": "",
						"secret/checksum":        "",
						"operator/version":       application.GetVersion(ctx),
					},
					Labels: map[string]string{
						"app":      goharborv1alpha2.PortalName,
						"operator": operatorName,
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector:                 m.Portal.Spec.NodeSelector,
					AutomountServiceAccountToken: &varFalse,
					Containers: []corev1.Container{
						{
							Name:  "portal",
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: port,
								},
							},

							ImagePullPolicy: corev1.PullAlways,
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(port),
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(port),
									},
								},
							},
						},
					},
					Priority: m.Portal.Spec.Priority,
				},
			},
			RevisionHistoryLimit: &revisionHistoryLimit,
			Paused:               false,
		},
	}, nil
}
