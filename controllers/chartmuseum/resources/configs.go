package chartmuseumresources

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

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
)

const (
	configName = "config.yaml"
)

var (
	once   sync.Once
	config []byte
)

func InitConfigMaps() {
	file, err := pkger.Open("/assets/templates/chartmuseum/config.yaml")
	if err != nil {
		panic(errors.Wrapf(err, "cannot open ChartMuseum configuration template %s", "/assets/templates/chartmuseum/config.yaml"))
	}
	defer file.Close()

	config, err = ioutil.ReadAll(file)
	if err != nil {
		panic(errors.Wrapf(err, "cannot read ChartMuseum configuration template %s", "/assets/templates/chartmuseum/config.yaml"))
	}
}

// https://github.com/goharbor/harbor/blob/master/make/photon/prepare/templates/chartserver/env.jinja

func (m *Manager) GetConfigMaps(ctx context.Context) ([]*corev1.ConfigMap, error) {
	once.Do(InitConfigMaps)

	operatorName := application.GetName(ctx)

	return []*corev1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.ChartMuseum.Name,
				Namespace: m.ChartMuseum.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.ChartMuseumName,
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
	sum := sha256.New().Sum([]byte(config))

	// TODO get generation of the secret
	return fmt.Sprintf("%x", sum)
}
