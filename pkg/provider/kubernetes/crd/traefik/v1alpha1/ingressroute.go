package v1alpha1

import (
	"errors"
	"fmt"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressRouteSpec is a specification for a IngressRouteSpec resource.
type IngressRouteSpec struct {
	Routes      []Route  `json:"routes"`
	EntryPoints []string `json:"entryPoints"`
	TLS         *TLS     `json:"tls,omitempty"`
}

// Route contains the set of routes.
type Route struct {
	Match       string          `json:"match"`
	Kind        string          `json:"kind"`
	Priority    int             `json:"priority"`
	Services    []Service       `json:"services,omitempty"`
	Middlewares []MiddlewareRef `json:"middlewares"`
}

// TLS contains the TLS certificates configuration of the routes.
// To enable Let's Encrypt, use an empty TLS struct,
// e.g. in YAML:
//
//	 tls: {} # inline format
//
//	 tls:
//	   secretName: # block format
type TLS struct {
	// SecretName is the name of the referenced Kubernetes Secret to specify the
	// certificate details.
	SecretName string `json:"secretName"`
	// PassphraseSecretName is the name of the referenced Kubernetes Secret to specify the
	// certificate passphrase.
	PassphraseSecretName string `json:"passphraseSecretName,omitempty"`
	// Options is a reference to a TLSOption, that specifies the parameters of the TLS connection.
	Options *TLSOptionRef `json:"options,omitempty"`
	// Store is a reference to a TLSStore, that specifies the parameters of the TLS store.
	Store        *TLSStoreRef   `json:"store,omitempty"`
	CertResolver string         `json:"certResolver,omitempty"`
	Domains      []types.Domain `json:"domains,omitempty"`
}

// TLSOptionRef is a ref to the TLSOption resources.
type TLSOptionRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// TLSStoreRef is a ref to the TLSStore resource.
type TLSStoreRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// LoadBalancerSpec can reference either a Kubernetes Service object (a load-balancer of servers),
// or a TraefikService object (a traefik load-balancer of services).
type LoadBalancerSpec struct {
	// Name is a reference to a Kubernetes Service object (for a load-balancer of servers),
	// or to a TraefikService object (service load-balancer, mirroring, etc).
	// The differentiation between the two is specified in the Kind field.
	Name      string          `json:"name"`
	Kind      string          `json:"kind"`
	Namespace string          `json:"namespace"`
	Sticky    *dynamic.Sticky `json:"sticky,omitempty"`

	// Port and all the fields below are related to a servers load-balancer,
	// and therefore should only be specified when Name references a Kubernetes Service.
	Port               int32                       `json:"port"`
	Scheme             string                      `json:"scheme,omitempty"`
	HealthCheck        *HealthCheck                `json:"healthCheck,omitempty"`
	Strategy           string                      `json:"strategy,omitempty"`
	PassHostHeader     *bool                       `json:"passHostHeader,omitempty"`
	ResponseForwarding *dynamic.ResponseForwarding `json:"responseForwarding,omitempty"`

	// Weight should only be specified when Name references a TraefikService object
	// (and to be precise, one that embeds a Weighted Round Robin).
	Weight *int `json:"weight,omitempty"`
}

// IsServersLB reports whether lb is a load-balancer of servers
// (as opposed to a traefik load-balancer of services).
func (lb LoadBalancerSpec) IsServersLB() (bool, error) {
	if lb.Name == "" {
		return false, errors.New("missing Name field in service")
	}
	if lb.Kind == "" || lb.Kind == "Service" {
		return true, nil
	}
	if lb.Kind != "TraefikService" {
		return false, fmt.Errorf("invalid kind value: %v", lb.Kind)
	}
	if lb.Port != 0 ||
		lb.Scheme != "" ||
		lb.HealthCheck != nil ||
		lb.Strategy != "" ||
		lb.PassHostHeader != nil ||
		lb.ResponseForwarding != nil ||
		lb.Sticky != nil {
		return false, fmt.Errorf("service of kind %v is incompatible with Kubernetes Service related fields", lb.Kind)
	}
	return false, nil
}

// Service defines an upstream to proxy traffic.
type Service struct {
	LoadBalancerSpec
}

// MiddlewareRef is a ref to the Middleware resources.
type MiddlewareRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// HealthCheck is the HealthCheck definition.
type HealthCheck struct {
	Path            string            `json:"path"`
	Host            string            `json:"host,omitempty"`
	Scheme          string            `json:"scheme"`
	IntervalSeconds int64             `json:"intervalSeconds"`
	TimeoutSeconds  int64             `json:"timeoutSeconds"`
	Headers         map[string]string `json:"headers"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressRoute is an Ingress CRD specification.
type IngressRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec IngressRouteSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IngressRouteList is a list of IngressRoutes.
type IngressRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []IngressRoute `json:"items"`
}
