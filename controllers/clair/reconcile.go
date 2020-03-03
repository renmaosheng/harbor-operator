package clair

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	notaryserverresources "github.com/goharbor/harbor-operator/controllers/clair/resources"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
	"github.com/goharbor/harbor-operator/pkg/factories/logger"
)

// +kubebuilder:rbac:groups=containerregistry.ovhcloud.com,resources=clairs,verbs=get;list;watch
// +kubebuilder:rbac:groups=containerregistry.ovhcloud.com,resources=clairs/status,verbs=get;update;patch

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.TODO()
	application.SetName(&ctx, r.GetName())
	application.SetVersion(&ctx, r.GetVersion())

	span, ctx := opentracing.StartSpanFromContext(ctx, "reconcile", opentracing.Tags{
		"Clair.Namespace": req.Namespace,
		"Clair.Name":      req.Name,
	})
	defer span.Finish()

	span.LogFields(
		log.String("Clair.Namespace", req.Namespace),
		log.String("Clair.Name", req.Name),
	)

	reqLogger := r.Log.WithValues("Request", req.NamespacedName, "Clair.Namespace", req.Namespace, "Clair.Name", req.Name)

	logger.Set(&ctx, reqLogger)

	// Fetch the Clair instance
	clair := &goharborv1alpha2.Clair{}

	err := r.Client.Get(ctx, req.NamespacedName, clair)
	if err != nil {
		if apierrs.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			reqLogger.Info("Clair does not exists")
			return reconcile.Result{}, nil
		}

		// Error reading the object
		return reconcile.Result{}, err
	}

	result := reconcile.Result{}

	if !clair.ObjectMeta.DeletionTimestamp.IsZero() {
		reqLogger.Info("Clair is being deleted")
		return result, nil
	}

	var g errgroup.Group

	g.Go(func() error {
		err = r.UpdateReadyStatus(ctx, &result, clair)
		return errors.Wrapf(err, "type=%s", goharborv1alpha2.ReadyConditionType)
	})

	g.Go(func() error {
		err = r.UpdateAppliedStatus(ctx, &result, clair)
		return errors.Wrapf(err, "type=%s", goharborv1alpha2.AppliedConditionType)
	})

	err = g.Wait()
	if err != nil {
		return result, errors.Wrap(err, "cannot set status")
	}

	return result, r.Status.UpdateStatus(ctx, &result, clair)
}

func (r *Reconciler) UpdateAppliedStatus(ctx context.Context, result *ctrl.Result, clair *goharborv1alpha2.Clair) error {
	if clair.Status.ObservedGeneration != clair.ObjectMeta.Generation {
		clair.Status.ObservedGeneration = clair.ObjectMeta.Generation

		err := r.Status.UpdateCondition(ctx, &clair.Status, goharborv1alpha2.AppliedConditionType, corev1.ConditionFalse, "new", "new generation detected")
		if err != nil {
			result.Requeue = true

			return errors.Wrapf(err, "value=%s", corev1.ConditionFalse)
		}
	}

	manager := &notaryserverresources.Manager{Clair: clair}

	switch r.Status.GetConditionStatus(ctx, &clair.Status, goharborv1alpha2.AppliedConditionType) {
	case corev1.ConditionTrue: // Already applied
		// Anyway, reconciler is triggered, so at least one child resource has been deleted
		// Try to recreate children

		err := r.Create(ctx, manager)
		if err != nil {
			result.Requeue = true

			err := r.Status.UpdateCondition(ctx, &clair.Status, goharborv1alpha2.AppliedConditionType, corev1.ConditionFalse, err.Error())
			if err != nil {
				result.Requeue = true

				return errors.Wrapf(err, "value=%s", corev1.ConditionFalse)
			}

			return nil
		}
	default: // Not yet applied
		err := r.Status.UpdateCondition(ctx, &clair.Status, goharborv1alpha2.AppliedConditionType, corev1.ConditionFalse)
		if err != nil {
			result.Requeue = true

			return errors.Wrapf(err, "value=%s", corev1.ConditionFalse)
		}

		err = r.Apply(ctx, manager)
		if err != nil {
			err := r.Status.UpdateCondition(ctx, &clair.Status, goharborv1alpha2.AppliedConditionType, corev1.ConditionFalse, err.Error())
			if err != nil {
				result.Requeue = true

				return errors.Wrapf(err, "value=%s", corev1.ConditionFalse)
			}

			return nil
		}

		err = r.Status.UpdateCondition(ctx, &clair.Status, goharborv1alpha2.AppliedConditionType, corev1.ConditionTrue)
		if err != nil {
			result.Requeue = true

			return errors.Wrapf(err, "value=%s", corev1.ConditionTrue)
		}
	}

	return nil
}

func (r *Reconciler) UpdateReadyStatus(ctx context.Context, result *ctrl.Result, clair *goharborv1alpha2.Clair) error {
	err := r.Status.UpdateCondition(ctx, &clair.Status, goharborv1alpha2.ReadyConditionType, corev1.ConditionFalse, "not-implemented", "Readiness check is not yet implemented")
	if err != nil {
		result.Requeue = true

		return errors.Wrapf(err, "value=%s", corev1.ConditionFalse)
	}

	return nil
}
