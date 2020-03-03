package harborresources

import (
	"context"
	"fmt"

	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
)

const (
	notaryCertificateName = "notary-certificate"
)

func (m *Manager) GetCertificates(ctx context.Context) ([]*certv1.Certificate, error) {
	operatorName := application.GetName(ctx)

	return []*certv1.Certificate{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Harbor.NormalizeComponentName(notaryCertificateName),
				Namespace: m.Harbor.Namespace,
				Labels: map[string]string{
					"app":      notaryCertificateName,
					"operator": operatorName,
				},
			},
			Spec: certv1.CertificateSpec{
				CommonName:   m.Harbor.Spec.Components.NotarySigner.CommonName,
				Organization: m.Harbor.Spec.Components.NotarySigner.Organization,
				SecretName:   fmt.Sprintf("%s-%s", m.Harbor.NormalizeComponentName(goharborv1alpha2.NotarySignerName), "certificate"),
				KeySize:      m.Harbor.Spec.Components.NotarySigner.KeySize,
				KeyAlgorithm: certv1.RSAKeyAlgorithm,
				KeyEncoding:  certv1.PKCS1,
				DNSNames: []string{
					m.Harbor.NormalizeComponentName(goharborv1alpha2.NotaryName),
					m.Harbor.Spec.Components.NotaryServer.PublicURL,
				},
				IssuerRef: m.Harbor.Spec.CertificateIssuerRef,
			},
		},
	}, nil
}
