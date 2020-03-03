package harborresources

import (
	"context"
	"net/url"
	"strings"

	"github.com/goharbor/harbor-operator/pkg/factories/application"
	"github.com/pkg/errors"
	netv1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/api/v1alpha2"
	coreresources "github.com/goharbor/harbor-operator/controllers/core/resources"
	notaryserverresources "github.com/goharbor/harbor-operator/controllers/notaryserver/resources"
	portalresources "github.com/goharbor/harbor-operator/controllers/portal/resources"
	registryresources "github.com/goharbor/harbor-operator/controllers/registry/resources"
)

func (m *Manager) GetIngresses(ctx context.Context) ([]*netv1.Ingress, error) {
	ingresses, err := m.GetCoreIngresses(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "core")
	}

	if m.Harbor.Spec.Components.NotaryServer != nil {
		notaryIngresses, err := m.GetNotaryServerIngresses(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "notary-server")
		}

		ingresses = append(ingresses, notaryIngresses...)
	}

	return ingresses, nil
}

func (m *Manager) GetCoreIngresses(ctx context.Context) ([]*netv1.Ingress, error) {
	operatorName := application.GetName(ctx)

	u, err := url.Parse(m.Harbor.Spec.PublicURL)
	if err != nil {
		panic(errors.Wrap(err, "invalid url"))
	}

	host := strings.SplitN(u.Host, ":", 1)

	var tls []netv1.IngressTLS
	if u.Scheme == "https" {
		tls = []netv1.IngressTLS{
			{
				SecretName: m.Harbor.Spec.TLSSecretName,
			},
		}
	}

	rules := []netv1.HTTPIngressPath{
		{
			Path: "/api",
			Backend: netv1.IngressBackend{
				ServiceName: m.Harbor.NormalizeComponentName(goharborv1alpha2.CoreName),
				ServicePort: intstr.FromInt(coreresources.PublicPort),
			},
		}, {
			Path: "/c",
			Backend: netv1.IngressBackend{
				ServiceName: m.Harbor.NormalizeComponentName(goharborv1alpha2.CoreName),
				ServicePort: intstr.FromInt(coreresources.PublicPort),
			},
		}, {
			Path: "/chartrepo",
			Backend: netv1.IngressBackend{
				ServiceName: m.Harbor.NormalizeComponentName(goharborv1alpha2.CoreName),
				ServicePort: intstr.FromInt(coreresources.PublicPort),
			},
		}, {
			Path: "/service",
			Backend: netv1.IngressBackend{
				ServiceName: m.Harbor.NormalizeComponentName(goharborv1alpha2.CoreName),
				ServicePort: intstr.FromInt(coreresources.PublicPort),
			},
		},
	}

	if m.Harbor.Spec.Components.Portal != nil {
		rules = append(rules, netv1.HTTPIngressPath{
			Path: "/",
			Backend: netv1.IngressBackend{
				ServiceName: m.Harbor.NormalizeComponentName(goharborv1alpha2.PortalName),
				ServicePort: intstr.FromInt(portalresources.PublicPort),
			},
		})
	}

	if m.Harbor.Spec.Components.Registry != nil {
		rules = append(rules, netv1.HTTPIngressPath{
			Path: "/v2",
			Backend: netv1.IngressBackend{
				ServiceName: m.Harbor.NormalizeComponentName(goharborv1alpha2.RegistryName),
				ServicePort: intstr.FromInt(registryresources.PublicPort),
			},
		})
	}

	return []*netv1.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Harbor.Name,
				Namespace: m.Harbor.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.HarborName,
					"operator": operatorName,
				},
			},
			Spec: netv1.IngressSpec{
				TLS: tls,
				Rules: []netv1.IngressRule{
					{
						Host: host[0],
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: rules,
							},
						},
					},
				},
			},
		},
	}, nil
}

func (m *Manager) GetNotaryServerIngresses(ctx context.Context) ([]*netv1.Ingress, error) {
	operatorName := application.GetName(ctx)

	u, err := url.Parse(m.Harbor.Spec.Components.NotaryServer.PublicURL)
	if err != nil {
		panic(errors.Wrap(err, "invalid url"))
	}

	notaryHost := strings.SplitN(u.Host, ":", 1)

	var tls []netv1.IngressTLS
	if u.Scheme == "https" {
		tls = []netv1.IngressTLS{
			{
				SecretName: m.Harbor.Spec.TLSSecretName,
			},
		}
	}

	return []*netv1.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      m.Harbor.Name,
				Namespace: m.Harbor.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha2.HarborName,
					"operator": operatorName,
				},
			},
			Spec: netv1.IngressSpec{
				TLS: tls,
				Rules: []netv1.IngressRule{
					{
						Host: notaryHost[0],
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path: "/",
										Backend: netv1.IngressBackend{
											ServiceName: m.Harbor.NormalizeComponentName(goharborv1alpha2.NotaryServerName),
											ServicePort: intstr.FromInt(notaryserverresources.PublicPort),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}
