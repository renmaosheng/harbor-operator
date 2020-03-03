package jobserviceresources

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/sethvargo/go-password/password"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
)

const (
	keyLength = 32
	secretKey = "secret"
)

func (m *Manager) GetSecrets(ctx context.Context) ([]*corev1.Secret, error) {
	operatorName := application.GetName(ctx)
	harborName := m.JobService.Name

	return []*corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.JobService.Name,
				Namespace: m.JobService.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.JobServiceName,
					"harbor":   harborName,
					"operator": operatorName,
				},
			},
			Type: corev1.SecretTypeOpaque,
			StringData: map[string]string{
				secretKey: password.MustGenerate(keyLength, 10, 10, false, true),
			},
		},
	}, nil
}

func (m *Manager) GetSecretsCheckSum() string {
	// TODO get generation of the secrets
	value := m.JobService.Spec.RedisSecret
	sum := sha256.New().Sum([]byte(value))

	return fmt.Sprintf("%x", sum)
}
