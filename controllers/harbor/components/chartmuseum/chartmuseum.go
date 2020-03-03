package chartmuseum

import (
	"context"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
)

type ChartMuseum struct {
	harbor *goharborv1alpha2.Harbor
	Option Option
}

type Option interface {
	GetPriority() *int32
}

func New(ctx context.Context, harbor *goharborv1alpha2.Harbor, opt Option) (*ChartMuseum, error) {
	return &ChartMuseum{
		harbor: harbor,
		Option: opt,
	}, nil
}
