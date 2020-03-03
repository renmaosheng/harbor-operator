package controllers

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/goharbor/harbor-operator/pkg/controllers/harbor"
)

type Controller interface {
	SetupWithManager(mgr manager.Manager) error
	GetName() string
}

type ControllerBuilder func(context.Context, string, string) (Controller, error)

func SetupWithManager(ctx context.Context, mgr manager.Manager, version string) error {
	var g errgroup.Group

	g.Go(func() error {
		const name = "Harbor"

		controller, err := harbor.New(ctx, name, version)
		if err != nil {
			return errors.Wrap(err, "failed to create controller")
		}

		err = controller.SetupWithManager(mgr)
		return errors.Wrap(err, "failed to setup controller")
	})

	g.Go(func() error {
		const name = "Portal"

		controller, err := harbor.New(ctx, name, version)
		if err != nil {
			return errors.Wrap(err, "failed to create controller")
		}

		err = controller.SetupWithManager(mgr)
		return errors.Wrap(err, "failed to setup controller")
	})

	return g.Wait()
}
