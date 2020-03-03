package registryctlresources

import (
	"context"

	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
	"github.com/ovh/configstore"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
	"github.com/goharbor/harbor-operator/pkg/factories/logger"
)

const (
	defaultKeyAlgorithm = certv1.RSAKeyAlgorithm
	defaultKeySize      = 4096
)

type certificateEncryption struct {
	KeySize      int
	KeyAlgorithm certv1.KeyAlgorithm
}

func (m *Manager) GetCertificates(ctx context.Context) ([]*certv1.Certificate, error) {
	operatorName := application.GetName(ctx)
	harborName := m.RegistryController.Name

	encryption := &certificateEncryption{
		KeySize:      defaultKeySize,
		KeyAlgorithm: defaultKeyAlgorithm,
	}

	item, err := configstore.Filter().Slice("certificate-encryption").Unmarshal(func() interface{} { return &certificateEncryption{} }).GetFirstItem()
	if err == nil {
		l := logger.Get(ctx)

		// todo
		encryptionConfig, err := item.Unmarshaled()
		if err != nil {
			l.Error(err, "Invalid encryption certificate config: use default value")
		} else {
			encryption = encryptionConfig.(*certificateEncryption)
		}
	}

	return []*certv1.Certificate{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.RegistryController.Name,
				Namespace: m.RegistryController.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.RegistryName,
					"harbor":   harborName,
					"operator": operatorName,
				},
			},
			Spec: certv1.CertificateSpec{
				CommonName:   m.RegistryController.Spec.PublicURL,
				Organization: []string{"Harbor Operator"},
				SecretName:   m.RegistryController.Name,
				KeySize:      encryption.KeySize,
				KeyAlgorithm: encryption.KeyAlgorithm,
				// https://github.com/goharbor/harbor/blob/ba4764c61d7da76f584f808f7d16b017db576fb4/src/jobservice/generateCerts.sh#L24-L26
				KeyEncoding: certv1.PKCS1,
				DNSNames:    []string{m.RegistryController.Spec.PublicURL},
				IssuerRef:   m.RegistryController.Spec.CertificateIssuerRef,
			},
		},
	}, nil
}
