package harbor

import (
	"context"

	"github.com/pkg/errors"

	"github.com/goharbor/harbor-operator/controllers/harbor"
	"github.com/goharbor/harbor-operator/pkg/controllers/config"
)

const (
	ConfigPrefix = "harbor-controller"
)

func New(ctx context.Context, name, version string) (*harbor.Reconciler, error) {
	config, err := config.GetConfig(ConfigPrefix)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get configuration")
	}

	return harbor.New(ctx, name, version, config)
}
