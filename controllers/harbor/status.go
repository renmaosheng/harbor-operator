package harbor

import (
	"context"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/factories/logger"
)

func (r *Reconciler) GetCondition(ctx context.Context, harbor *goharborv1alpha2.Harbor, conditionType goharborv1alpha2.HarborConditionType) goharborv1alpha2.HarborCondition {
	for _, condition := range harbor.Status.Conditions {
		if condition.Type == conditionType {
			return condition
		}
	}

	return goharborv1alpha2.HarborCondition{
		Type:   conditionType,
		Status: corev1.ConditionUnknown,
	}
}

func (r *Reconciler) GetConditionStatus(ctx context.Context, harbor *goharborv1alpha2.Harbor, conditionType goharborv1alpha2.HarborConditionType) corev1.ConditionStatus {
	return r.GetCondition(ctx, harbor, conditionType).Status
}

func (r *Reconciler) UpdateCondition(ctx context.Context, harbor *goharborv1alpha2.Harbor, conditionType goharborv1alpha2.HarborConditionType, status corev1.ConditionStatus, reasons ...string) error {
	var reason, message string

	switch len(reasons) {
	case 0: // nolint:mnd
	case 1: // nolint:mnd
		reason = reasons[0]
	case 2: // nolint:mnd
		reason = reasons[0]
		message = reasons[1]
	default:
		return errors.Errorf("expecting reason and message, got %d parameters", len(reasons))
	}

	now := metav1.Now()

	for i, condition := range harbor.Status.Conditions {
		if condition.Type == conditionType {
			now.DeepCopyInto(&condition.LastUpdateTime)

			if condition.LastTransitionTime.IsZero() || condition.Status != status {
				now.DeepCopyInto(&condition.LastTransitionTime)
			}

			condition.Status = status
			condition.Reason = reason
			condition.Message = message

			harbor.Status.Conditions[i] = condition

			return nil
		}
	}

	condition := goharborv1alpha2.HarborCondition{
		Type:    conditionType,
		Status:  status,
		Reason:  reason,
		Message: message,
	}
	now.DeepCopyInto(&condition.LastUpdateTime)
	now.DeepCopyInto(&condition.LastTransitionTime)

	harbor.Status.Conditions = append(harbor.Status.Conditions, condition)

	return nil
}

// UpdateStatus applies current in-memory statuses to the remote resource
// https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#status-subresource
func (r *Reconciler) UpdateStatus(ctx context.Context, result *ctrl.Result, harbor *goharborv1alpha2.Harbor) error {
	err := r.Status().Update(ctx, harbor)
	if err != nil {
		result.Requeue = true

		seconds, needWait := apierrors.SuggestsClientDelay(err)
		if needWait {
			result.RequeueAfter = time.Second * time.Duration(seconds)
		}

		if apierrors.IsConflict(err) {
			// the object has been modified; please apply your changes to the latest version and try again
			logger.Get(ctx).Error(err, "cannot update status field")
			return nil
		}

		return errors.Wrap(err, "cannot update status field")
	}

	return nil
}
