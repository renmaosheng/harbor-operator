package notaryserverresources

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/goharbor/harbor-operator/pkg/factories/application"
)

const (
	PublicPort = 80
)

func (m *Manager) GetServices(ctx context.Context) ([]*corev1.Service, error) {
	operatorName := application.GetName(ctx)

	return []*corev1.Service{
		{
			// https://github.com/goharbor/harbor-helm/blob/master/templates/notary/notary-svc.yaml
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Notary.Name,
				Namespace: m.Notary.Namespace,
				Labels: map[string]string{
					"app":      m.Notary.Name,
					"operator": operatorName,
				},
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:       m.Notary.Name,
						Port:       PublicPort,
						TargetPort: intstr.FromInt(notaryServerPort),
					},
				},
				Selector: map[string]string{
					"app":      m.Notary.Name,
					"operator": operatorName,
				},
			},
		},
	}, nil
}
