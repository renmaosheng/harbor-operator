package portalresources

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

func (*Manager) GetConfigMaps(ctx context.Context) ([]*corev1.ConfigMap, error) {
	return []*corev1.ConfigMap{}, nil
}
