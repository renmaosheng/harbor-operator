package jobserviceresources

import (
	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
)

type Manager struct {
	JobService *goharborv1alpha2.JobService
}
