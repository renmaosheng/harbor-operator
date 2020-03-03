package registryresources

import (
	"context"

	netv1 "k8s.io/api/networking/v1beta1"
)

func (m *Manager) GetIngresses(ctx context.Context) ([]*netv1.Ingress, error) {
	return []*netv1.Ingress{}, nil
}
