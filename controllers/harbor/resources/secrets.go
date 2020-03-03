package harborresources

import (
	"context"

	"github.com/sethvargo/go-password/password"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
)

const (
	keyLength = 16
	secretKey = "secretKey"
)

func (m *Manager) GetSecrets(ctx context.Context) ([]*corev1.Secret, error) {
	operatorName := application.GetName(ctx)

	return []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Harbor.Name,
				Namespace: m.Harbor.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.CoreName,
					"operator": operatorName,
				},
			},
			StringData: map[string]string{
				"secret":  password.MustGenerate(keyLength, 5, 0, false, true),
				secretKey: password.MustGenerate(keyLength, 5, 0, false, true),
			},
		},
	}, nil
}
