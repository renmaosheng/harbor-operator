package notaryserverresources

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
	"github.com/markbates/pkger"
	"github.com/pkg/errors"
)

const (
	configName = "app.conf"
)

var (
	once   sync.Once
	config []byte
)

func InitConfigMaps() {
	file, err := pkger.Open("/assets/templates/core/app.conf")
	if err != nil {
		panic(errors.Wrapf(err, "cannot open Core configuration template %s", "/assets/templates/core/app.conf"))
	}
	defer file.Close()

	config, err = ioutil.ReadAll(file)
	if err != nil {
		panic(errors.Wrapf(err, "cannot read Core configuration template %s", "/assets/templates/core/app.conf"))
	}
}

func (m *Manager) GetConfigMaps(ctx context.Context) ([]*corev1.ConfigMap, error) { // nolint:funlen
	once.Do(InitConfigMaps)

	operatorName := application.GetName(ctx)

	return []*corev1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Core.Name,
				Namespace: m.Core.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.CoreName,
					"operator": operatorName,
				},
			},

			BinaryData: map[string][]byte{
				configName: config,
			},
		},
	}, nil
}

func (m *Manager) GetConfigMapsCheckSum() string {
	sum := sha256.New().Sum(config)

	return fmt.Sprintf("%x", sum)
}
