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

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
)

const (
	configKey = "config.yaml"
)

var (
	once   sync.Once
	config []byte
)

func InitConfigMaps() {
	file, err := pkger.Open("/assets/templates/clair/config.yaml")
	if err != nil {
		panic(errors.Wrapf(err, "cannot open Clair configuration template %s", "/assets/templates/clair/config.yaml"))
	}
	defer file.Close()

	config, err = ioutil.ReadAll(file)
	if err != nil {
		panic(errors.Wrapf(err, "cannot read Clair configuration template %s", "/assets/templates/clair/config.yaml"))
	}
}

func (m *Manager) GetConfigMaps(ctx context.Context) ([]*corev1.ConfigMap, error) {
	once.Do(InitConfigMaps)

	operatorName := application.GetName(ctx)

	return []*corev1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Clair.Name,
				Namespace: m.Clair.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.ClairName,
					"operator": operatorName,
				},
			},
			BinaryData: map[string][]byte{
				configKey: config,
			},
			// https://github.com/goharbor/harbor-scanner-clair#configuration
			// https://github.com/goharbor/harbor/blob/master/make/photon/prepare/templates/clair/clair_env.jinja
			Data: map[string]string{
				"SCANNER_CLAIR_URL":                   fmt.Sprintf("http://%s", m.Clair.Name),
				"SCANNER_LOG_LEVEL":                   "debug",
				"SCANNER_STORE_REDIS_POOL_MAX_ACTIVE": "5",
				"SCANNER_STORE_REDIS_POOL_MAX_IDLE":   "5",
				"SCANNER_STORE_REDIS_SCAN_JOB_TTL":    "1h",
				"SCANNER_API_SERVER_ADDR":             fmt.Sprintf(":%d", adapterPort),
			},
		},
	}, nil
}

func (m *Manager) GetConfigMapsCheckSum() string {
	value := fmt.Sprintf("%d\n%+v\n%x", adapterPort, m.Clair.Spec.VulnerabilitySources, config)
	sum := sha256.New().Sum([]byte(value))

	return fmt.Sprintf("%x", sum)
}
