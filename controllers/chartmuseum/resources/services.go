package chartmuseumresources

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

	return []*corev1.Service{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.ChartMuseum.Name,
				Namespace: m.ChartMuseum.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.ChartMuseumName,
					"operator": operatorName,
				},
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Port:       PublicPort,
						TargetPort: intstr.FromInt(port),
					},
				},
				Selector: map[string]string{
					"app":      goharborv1alpha2.ChartMuseumName,
					"operator": operatorName,
				},
			},
		},
	}, nil
}
