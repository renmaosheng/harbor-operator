package registryctlresources

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
)

const (
	PublicPort = 80
)

func (m *Manager) GetServices(ctx context.Context) ([]*corev1.Service, error) {
	operatorName := application.GetName(ctx)
	harborName := m.RegistryController.Name

	return []*corev1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.RegistryController.Name,
				Namespace: m.RegistryController.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.RegistryName,
					"harbor":   harborName,
					"operator": operatorName,
				},
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:       "registry",
						TargetPort: intstr.FromInt(apiPort),
						Port:       PublicPort,
					}, {
						Name: "registry-debug",
						Port: metricsPort,
					}, {
						Name: "controller",
						Port: ctlAPIPort,
					},
				},
				Selector: map[string]string{
					"app":      goharborv1alpha2.RegistryName,
					"harbor":   harborName,
					"operator": operatorName,
				},
			},
		},
	}, nil
}
