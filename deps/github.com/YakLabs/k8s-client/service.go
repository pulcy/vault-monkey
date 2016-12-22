package client

import "github.com/YakLabs/k8s-client/intstr"

const (
	// ServiceTypeClusterIP means a service will only be accessible inside the
	// cluster, via the ClusterIP.
	ServiceTypeClusterIP ServiceType = "ClusterIP"

	// ServiceTypeNodePort means a service will be exposed on one port of
	// every node, in addition to 'ClusterIP' type.
	ServiceTypeNodePort ServiceType = "NodePort"

	// ServiceTypeLoadBalancer means a service will be exposed via an
	// external load balancer (if the cloud provider supports it), in addition
	// to 'NodePort' type.
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"

	// ServiceAffinityClientIP is the Client IP based.
	ServiceAffinityClientIP ServiceAffinity = "ClientIP"

	// ServiceAffinityNone - no session affinity.
	ServiceAffinityNone ServiceAffinity = "None"
)

type (
	// ServiceInterface has methods to work with Service resources.
	ServiceInterface interface {
		CreateService(namespace string, item *Service) (*Service, error)
		GetService(namespace, name string) (result *Service, err error)
		ListServices(namespace string, opts *ListOptions) (*ServiceList, error)
		DeleteService(namespace, name string) error
		UpdateService(namespace string, item *Service) (*Service, error)
	}

	// Service is a named abstraction of software service (for example, mysql) consisting of local port (for example 3306) that the proxy listens on, and the selector that determines which pods will answer requests sent through the proxy.
	Service struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`

		// Spec defines the behavior of a service
		Spec ServiceSpec `json:"spec,omitempty"`

		// Status represents the current status of a service.
		Status ServiceStatus `json:"status,omitempty"`
	}

	// ServiceList holds a list of services.
	ServiceList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		Items []Service `json:"items"`
	}

	// ServiceStatus represents the current status of a service
	ServiceStatus struct {
		// LoadBalancer contains the current status of the load-balancer,
		// if one is present.
		LoadBalancer LoadBalancerStatus `json:"loadBalancer,omitempty"`
	}

	// LoadBalancerStatus represents the status of a load-balancer
	LoadBalancerStatus struct {
		// Ingress is a list containing ingress points for the load-balancer;
		// traffic intended for the service should be sent to these ingress points.
		Ingress []LoadBalancerIngress `json:"ingress,omitempty"`
	}

	// LoadBalancerIngress represents the status of a load-balancer ingress point: traffic intended for the service should be sent to an ingress point.
	LoadBalancerIngress struct {
		// IP is set for load-balancer ingress points that are IP based
		// (typically GCE or OpenStack load-balancers)
		IP string `json:"ip,omitempty"`

		// Hostname is set for load-balancer ingress points that are DNS based
		// (typically AWS load-balancers)
		Hostname string `json:"hostname,omitempty"`
	}

	// ServiceType describes ingress methods for a service
	ServiceType string
	// ServiceAffinity is the type of affinity for the service
	ServiceAffinity string

	// ServiceSpec describes the attributes that a user creates on a service
	ServiceSpec struct {
		// Type determines how the service will be exposed.  Valid options: ClusterIP, NodePort, LoadBalancer
		Type ServiceType `json:"type,omitempty"`

		// Required: The list of ports that are exposed by this service.
		Ports []ServicePort `json:"ports"`

		// This service will route traffic to pods having labels matching this selector. If empty or not present,
		// the service is assumed to have endpoints set by an external process and Kubernetes will not modify
		// those endpoints.
		Selector map[string]string `json:"selector"`

		// ClusterIP is usually assigned by the master.  If specified by the user
		// we will try to respect it or else fail the request.  This field can
		// not be changed by updates.
		// Valid values are None, empty string (""), or a valid IP address
		// None can be specified for headless services when proxying is not required
		ClusterIP string `json:"clusterIP,omitempty"`

		// ExternalIPs are used by external load balancers, or can be set by
		// users to handle external traffic that arrives at a node.
		ExternalIPs []string `json:"externalIPs,omitempty"`

		// Only applies to Service Type: LoadBalancer
		// LoadBalancer will get created with the IP specified in this field.
		// This feature depends on whether the underlying cloud-provider supports specifying
		// the loadBalancerIP when a load balancer is created.
		// This field will be ignored if the cloud-provider does not support the feature.
		LoadBalancerIP string `json:"loadBalancerIP,omitempty"`

		// Optional: Supports "ClientIP" and "None".  Used to maintain session affinity.
		SessionAffinity ServiceAffinity `json:"sessionAffinity,omitempty"`

		// Optional: If specified and supported by the platform, this will restrict traffic through the cloud-provider
		// load-balancer will be restricted to the specified client IPs. This field will be ignored if the
		// cloud-provider does not support the feature."
		LoadBalancerSourceRanges []string `json:"loadBalancerSourceRanges,omitempty"`
	}

	// ServicePort is a single exposed port.
	ServicePort struct {
		// Optional if only one ServicePort is defined on this service: The
		// name of this port within the service.  This must be a DNS_LABEL.
		// All ports within a ServiceSpec must have unique names.  This maps to
		// the 'Name' field in EndpointPort objects.
		Name string `json:"name"`

		// The IP protocol for this port.  Supports "TCP" and "UDP".
		Protocol Protocol `json:"protocol"`

		// The port that will be exposed on the service.
		Port int32 `json:"port"`

		// Optional: The target port on pods selected by this service.  If this
		// is a string, it will be looked up as a named port in the target
		// Pod's container ports.  If this is not specified, the value
		// of the 'port' field is used (an identity map).
		// This field is ignored for services with clusterIP=None, and should be
		// omitted or set equal to the 'port' field.
		TargetPort intstr.IntOrString `json:"targetPort"`

		// The port on each node on which this service is exposed.
		// Default is to auto-allocate a port if the ServiceType of this Service requires one.
		NodePort int32 `json:"nodePort"`
	}
)

// NewService creates a new service struct
func NewService(namespace, name string) *Service {
	return &Service{
		TypeMeta: TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: ObjectMeta{
			Namespace:   namespace,
			Name:        name,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
		Spec: ServiceSpec{
			Ports:    make([]ServicePort, 0),
			Selector: make(map[string]string),
		},
	}
}
