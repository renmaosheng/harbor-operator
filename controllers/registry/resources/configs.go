package registryresources

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
	defaultRegistryConfigName = "config.yml"
	registryConfigName        = "config.yaml"
	registryCtlConfigName     = "ctl-config.yaml"
)

var (
	once              sync.Once
	registryConfig    []byte
	registryCtlConfig []byte
)

func InitConfigMaps() {
	{
		file, err := pkger.Open("/assets/templates/registry/config.yaml")
		if err != nil {
			panic(errors.Wrapf(err, "cannot open Registry configuration template %s", "/assets/templates/registry/config.yaml"))
		}
		defer file.Close()

		registryConfig, err = ioutil.ReadAll(file)
		if err != nil {
			panic(errors.Wrapf(err, "cannot read Registry configuration template %s", "/assets/templates/registry/config.yaml"))
		}
	}
	{
		file, err := pkger.Open("/assets/templates/registry/ctl-config.yaml")
		if err != nil {
			panic(errors.Wrapf(err, "cannot open Registry configuration template %s", "/assets/templates/registry/ctl-config.yaml"))
		}
		defer file.Close()

		registryCtlConfig, err = ioutil.ReadAll(file)
		if err != nil {
			panic(errors.Wrapf(err, "cannot read Registry configuration template %s", "/assets/templates/registry/ctl-config.yaml"))
		}
	}
}

func (m *Manager) GetConfigMaps(ctx context.Context) ([]*corev1.ConfigMap, error) {
	once.Do(InitConfigMaps)

	operatorName := application.GetName(ctx)
	harborName := m.Registry.Name

	return []*corev1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Registry.Name,
				Namespace: m.Registry.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.RegistryName,
					"harbor":   harborName,
					"operator": operatorName,
				},
			},

			BinaryData: map[string][]byte{
				registryConfigName:    registryConfig,
				registryCtlConfigName: registryCtlConfig,
			},
		},
	}, nil
}

func (m *Manager) GetConfigMapsCheckSum() string {
	value := fmt.Sprintf("%x\n%x", registryCtlConfig, registryConfig)
	sum := sha256.New().Sum([]byte(value))

	return fmt.Sprintf("%x", sum)
}
