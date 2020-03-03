package registryctl

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	certv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	"github.com/goharbor/harbor-operator/pkg/controllers/common"
	"github.com/goharbor/harbor-operator/pkg/controllers/config"
	"github.com/goharbor/harbor-operator/pkg/event-filter/class"
	"github.com/goharbor/harbor-operator/pkg/factories/logger"
)

const (
	DefaultRequeueWait = 2 * time.Second
)

// Reconciler reconciles a Harbor object
type Reconciler struct {
	common.Controller

	Log    logr.Logger
	Scheme *runtime.Scheme

	Config config.Config
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Scheme = mgr.GetScheme()

	err := r.Init()
	if err != nil {
		return errors.Wrap(err, "failed to init controller")
	}

	return ctrl.NewControllerManagedBy(mgr).
		WithEventFilter(&class.Filter{
			ClassName: r.Config.ClassName,
		}).
		For(&goharborv1alpha2.RegistryController{}).
		Owns(&appsv1.Deployment{}).
		Owns(&certv1.Certificate{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&netv1.Ingress{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: r.Config.ConcurrentReconciles,
		}).
		Complete(r)
}

func (r *Reconciler) Init() error {
	var g errgroup.Group

	g.Go(func() error {
		err := r.InitConfigMaps()
		return errors.Wrap(err, "configMaps")
	})

	return g.Wait()
}

// +kubebuilder:rbac:groups=containerregistryctl.ovhcloud.com,resources=registries,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=containerregistryctl.ovhcloud.com,resources=registries/status,verbs=get;update;patch

func New(ctx context.Context, name, version string, config *config.Config) (*Reconciler, error) {
	return &Reconciler{
		Controller: common.Controller{
			Name:    name,
			Version: version,
		},
		Log:    logger.Get(ctx).WithName("controller").WithName(name),
		Config: *config,
	}, nil
}
