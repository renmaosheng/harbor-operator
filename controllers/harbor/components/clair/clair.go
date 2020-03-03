package clair

import (
	"context"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
)

type Clair struct {
	harbor *goharborv1alpha2.Harbor
	Option Option
}

type Option interface {
	GetPriority() *int32
}

func New(ctx context.Context, harbor *goharborv1alpha2.Harbor, opt Option) (*Clair, error) {
	return &Clair{
		harbor: harbor,
		Option: opt,
	}, nil
}
