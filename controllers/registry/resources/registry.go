package registryresources

import (
	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
)

type Manager struct {
	Registry *goharborv1alpha2.Registry
}
