package notary

import (
	"context"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
)

const (
	NotaryServerName = "notary-server"
	NotarySignerName = "notary-signer"
)

type Notary struct {
	harbor *goharborv1alpha2.Harbor
	Option Option
}

type Option interface {
	GetPriority() *int32
}

func New(ctx context.Context, harbor *goharborv1alpha2.Harbor, opt Option) (*Notary, error) {
	return &Notary{
		harbor: harbor,
		Option: opt,
	}, nil
}
