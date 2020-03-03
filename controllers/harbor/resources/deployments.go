package harborresources

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
)

func (m *Manager) GetDeployments(ctx context.Context) ([]*appsv1.Deployment, error) {
	return []*appsv1.Deployment{}, nil
}
