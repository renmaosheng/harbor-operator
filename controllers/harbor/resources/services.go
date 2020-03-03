package harborresources

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

func (m *Manager) GetServices(ctx context.Context) ([]*corev1.Service, error) {
	return []*corev1.Service{}, nil
}
