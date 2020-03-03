package notaryserverresources

import (
	"context"

	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
)

func (m *Manager) GetCertificates(ctx context.Context) ([]*certv1.Certificate, error) {
	return []*certv1.Certificate{}, nil
}
