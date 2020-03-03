package harborresources

import (
	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
)

type Manager struct {
	Harbor *goharborv1alpha2.Harbor
}
