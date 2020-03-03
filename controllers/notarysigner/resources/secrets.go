package notaryserverresources

import (
	"context"
	"crypto/sha256"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func (*Manager) GetSecrets(ctx context.Context) ([]*corev1.Secret, error) {
	return []*corev1.Secret{}, nil
}

func (m *Manager) GetSecretsCheckSum() string {
	// TODO get generation of the secret
	value := fmt.Sprintf("%s\n%s", m.Notary.Spec.CertificateSecret, m.Notary.Spec.DatabaseSecret)
	sum := sha256.New().Sum([]byte(value))

	return fmt.Sprintf("%x", sum)
}
