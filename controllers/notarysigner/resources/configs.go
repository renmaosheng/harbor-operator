package notaryserverresources

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/markbates/pkger"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/goharbor/harbor-operator/pkg/factories/application"
)

const (
	signerConfigKey = "signer.json"
)

var (
	once         sync.Once
	signerConfig []byte
)

func InitConfigMaps() {
	// We can't use a constant containing file path. Pkger don't understant if it's not the value passed as parameter.
	// const templatePath = "/my/Path"
	// pkger.Open(templatePath) --> Doesn't work
	signerFile, signerErr := pkger.Open("/assets/templates/notary/signer.json")
	if signerErr != nil {
		panic(errors.Wrapf(signerErr, "cannot open Notary Signer configuration template %s", "/assets/templates/notary/signer.json"))
	}
	defer signerFile.Close()

	signerConfig, signerErr = ioutil.ReadAll(signerFile)
	if signerErr != nil {
		panic(errors.Wrapf(signerErr, "cannot read Notary Signer configuration template %s", "/assets/templates/notary/signer.json"))
	}
}

func (m *Manager) GetConfigMaps(ctx context.Context) ([]*corev1.ConfigMap, error) {
	once.Do(InitConfigMaps)

	operatorName := application.GetName(ctx)

	return []*corev1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Notary.Name,
				Namespace: m.Notary.Namespace,
				Labels: map[string]string{
					"app":      m.Notary.Name,
					"operator": operatorName,
				},
			},
			BinaryData: map[string][]byte{
				signerConfigKey: signerConfig,
			},
		},
	}, nil
}

func (m *Manager) GetConfigMapsCheckSum() string {
	sum := sha256.New().Sum(signerConfig)

	return fmt.Sprintf("%x", sum)
}
