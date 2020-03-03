package jobserviceresources

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

const (
	logsDirectory = "/var/log/jobs"
)

var (
	once         sync.Once
	config       []byte
	hookMaxRetry = 5
)

func InitConfigMaps() {
	file, err := pkger.Open("/assets/templates/jobservice/config.yaml")
	if err != nil {
		panic(errors.Wrapf(err, "cannot open JobService configuration template %s", "/assets/templates/jobservice/config.yaml"))
	}
	defer file.Close()

	config, err = ioutil.ReadAll(file)
	if err != nil {
		panic(errors.Wrapf(err, "cannot read JobService configuration template %s", "/assets/templates/jobservice/config.yaml"))
	}
}

func (m *Manager) GetConfigMaps(ctx context.Context) ([]*corev1.ConfigMap, error) {
	once.Do(InitConfigMaps)

	operatorName := application.GetName(ctx)
	harborName := m.JobService.Name

	return []*corev1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.JobService.Name,
				Namespace: m.JobService.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.JobServiceName,
					"harbor":   harborName,
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

	return fmt.Sprintf("%x", sum)
}
