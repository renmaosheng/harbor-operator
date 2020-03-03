package registryctlresources

import (
	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
)

type Manager struct {
	RegistryController *goharborv1alpha2.RegistryController
}
