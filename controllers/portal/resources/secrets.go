package portalresources

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

func (*Manager) GetSecrets(ctx context.Context) ([]*corev1.Secret, error) {
	return []*corev1.Secret{}, nil
}
