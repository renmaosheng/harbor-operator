package notaryserverresources

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
)

const (
	PublicPort        = 80
	AdapterPublicPort = 8080
)

func (m *Manager) GetServices(ctx context.Context) ([]*corev1.Service, error) {
	operatorName := application.GetName(ctx)

	return []*corev1.Service{
		{
			// https://github.com/goharbor/harbor-helm/blob/master/templates/clair/clair-svc.yaml
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Clair.Name,
				Namespace: m.Clair.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.ClairName,
					"operator": operatorName,
				},
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:       "api",
						Port:       PublicPort,
						TargetPort: intstr.FromInt(apiPort),
					}, {
						Name: "healthcheck",
						Port: healthPort,
					}, {
						Name:       "adapter",
						Port:       AdapterPublicPort,
						TargetPort: intstr.FromInt(adapterPort),
					},
				},
				Selector: map[string]string{
					"app":      goharborv1alpha2.ClairName,
					"operator": operatorName,
				},
			},
		},
	}, nil
}
