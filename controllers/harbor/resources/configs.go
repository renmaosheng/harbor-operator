package harborresources

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

func (m *Manager) GetConfigMaps(ctx context.Context) ([]*corev1.ConfigMap, error) { // nolint:funlen
	return []*corev1.ConfigMap{}, nil
}
