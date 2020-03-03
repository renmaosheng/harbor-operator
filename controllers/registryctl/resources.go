package registryctl

import (
	"context"
	"errors"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
)

func (r *Reconciler) AddResources(ctx context.Context, registryctl *goharborv1alpha2.RegistryController) error {
	cm, err := r.GetConfigMap(ctx, registryctl)
	if err != nil {
		return errors.Wrap(err, "cannot get configMap")
	}

	err = r.Controller.AddResourceToManage(ctx, cm)
	if err != nil {
		return errors.Wrapf(err, "cannot add resource %+v", cm)
	}

	return errors.New("not yet implemented")
}
