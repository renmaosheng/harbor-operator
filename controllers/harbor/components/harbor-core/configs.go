package core

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"sync"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/controllers/harbor/components/clair"
	"github.com/goharbor/harbor-operator/controllers/harbor/components/notary"
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

func (c *HarborCore) GetConfigMaps(ctx context.Context) []*corev1.ConfigMap { // nolint:funlen
	once.Do(InitConfigMaps)

	operatorName := application.GetName(ctx)
	harborName := c.harbor.Name

	return []*corev1.ConfigMap{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      c.harbor.NormalizeComponentName(goharborv1alpha2.CoreName),
				Namespace: c.harbor.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.CoreName,
					"harbor":   harborName,
					"operator": operatorName,
				},
			},

			BinaryData: map[string][]byte{
				configName: config,
			},

			// https://github.com/goharbor/harbor/blob/master/make/photon/prepare/templates/core/env.jinja
			Data: map[string]string{
				"CONFIG_PATH": path.Join(coreConfigPath, configFileName),

				"AUTH_MODE":                      "db_auth",
				"CFG_EXPIRATION":                 "5",
				"CHART_CACHE_DRIVER":             "memory",
				"EXT_ENDPOINT":                   c.harbor.Spec.PublicURL,
				"LOG_LEVEL":                      "debug",
				"MAX_JOB_WORKERS":                fmt.Sprintf("%d", c.harbor.Spec.Components.JobService.WorkerCount),
				"READ_ONLY":                      fmt.Sprintf("%+v", c.harbor.Spec.ReadOnly),
				"REGISTRY_STORAGE_PROVIDER_NAME": "memory",
				"RELOAD_KEY":                     "true",
				"SYNC_QUOTA":                     "true",
				"SYNC_REGISTRY":                  "false",

				"_REDIS_URL":                    "", // For session purpose
				"ADMIRAL_URL":                   "NA",
				"CHART_REPOSITORY_URL":          fmt.Sprintf("http://%s", c.harbor.NormalizeComponentName(goharborv1alpha2.ChartMuseumName)),
				"CLAIR_HEALTH_CHECK_SERVER_URL": fmt.Sprintf("http://%s:6061", c.harbor.NormalizeComponentName(goharborv1alpha2.ClairName)),
				"CLAIR_URL":                     fmt.Sprintf("http://%s", c.harbor.NormalizeComponentName(goharborv1alpha2.ClairName)),
				"CLAIR_ADAPTER_URL":             fmt.Sprintf("http://%s:%d", c.harbor.NormalizeComponentName(goharborv1alpha2.ClairName), clair.AdapterPublicPort),
				"CORE_LOCAL_URL":                fmt.Sprintf("http://%s", c.harbor.NormalizeComponentName(goharborv1alpha2.CoreName)),
				"CORE_URL":                      fmt.Sprintf("http://%s", c.harbor.NormalizeComponentName(goharborv1alpha2.CoreName)),
				"JOBSERVICE_URL":                fmt.Sprintf("http://%s", c.harbor.NormalizeComponentName(goharborv1alpha2.JobServiceName)),
				"NOTARY_URL":                    fmt.Sprintf("http://%s", c.harbor.NormalizeComponentName(notary.NotaryServerName)),
				"PORTAL_URL":                    fmt.Sprintf("http://%s", c.harbor.NormalizeComponentName(goharborv1alpha2.PortalName)),
				"REGISTRY_URL":                  fmt.Sprintf("http://%s", c.harbor.NormalizeComponentName(goharborv1alpha2.RegistryName)),
				"REGISTRYCTL_URL":               fmt.Sprintf("http://%s:8080", c.harbor.NormalizeComponentName(goharborv1alpha2.RegistryName)),
				"TOKEN_SERVICE_URL":             fmt.Sprintf("http://%s/service/token", c.harbor.NormalizeComponentName(goharborv1alpha2.CoreName)),

				"DATABASE_TYPE":             "postgresql",
				"POSTGRESQL_MAX_IDLE_CONNS": fmt.Sprintf("%d", maxIdleConns),
				"POSTGRESQL_MAX_OPEN_CONNS": fmt.Sprintf("%d", maxOpenConns),

				"WITH_CHARTMUSEUM": strconv.FormatBool(c.harbor.Spec.Components.ChartMuseum != nil),
				"WITH_CLAIR":       strconv.FormatBool(c.harbor.Spec.Components.Clair != nil),
				"WITH_NOTARY":      strconv.FormatBool(c.harbor.Spec.Components.Notary != nil),
			},
		},
	}
}

func (c *HarborCore) GetConfigMapsCheckSum() string {
	value := fmt.Sprintf("%s\n%+v\n%x", c.harbor.Spec.PublicURL, c.harbor.Spec.Components.Clair != nil, config)
	sum := sha256.New().Sum([]byte(value))

	// todo get generation of the secret
	return fmt.Sprintf("%x", sum)
}
