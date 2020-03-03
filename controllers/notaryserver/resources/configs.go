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
	serverConfigKey = "server.json"
)

var (
	once         sync.Once
	serverConfig []byte
)

func InitConfigMaps() {
	// We can't use a constant containing file path. Pkger don't understant if it's not the value passed as parameter.
	// const templatePath = "/my/Path"
	// pkger.Open(templatePath) --> Doesn't work
	serverFile, serverErr := pkger.Open("/assets/templates/notary/server.json")
	if serverErr != nil {
		panic(errors.Wrapf(serverErr, "cannot open Notary Server configuration template %s", "/assets/templates/notary/server.json"))
	}
	defer serverFile.Close()

	serverConfig, serverErr = ioutil.ReadAll(serverFile)
	if serverErr != nil {
		panic(errors.Wrapf(serverErr, "cannot read Notary Server configuration template %s", "/assets/templates/notary/server.json"))
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
				serverConfigKey: serverConfig,
			},
		},
	}, nil
}

func (m *Manager) GetConfigMapsCheckSum() string {
	sum := sha256.New().Sum(serverConfig)

	return fmt.Sprintf("%x", sum)
}
