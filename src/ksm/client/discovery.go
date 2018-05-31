package client

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"net"

	"strings"

	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/client"
	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/prometheus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ksmAppLabelName  = "k8s-app"
	ksmAppLabelValue = "kube-state-metrics"
	ksmPortName      = "http-metrics"
	k8sTCP           = "TCP"
	ksmQualifiedName = "kube-state-metrics.kube-system.svc.cluster.local"
	ksmDNSService    = "http-metrics"
	ksmDNSProto      = "tcp"
)

// discoverer implements Discoverer interface by using official Kubernetes' Go client
type discoverer struct {
	lookupSRV         func(service, proto, name string) (cname string, addrs []*net.SRV, err error)
	apiClient         client.Kubernetes
	logger            *logrus.Logger
	overridenEndpoint string
}

// ksm implements Client interface
type ksm struct {
	httpClient *http.Client
	endpoint   url.URL
	nodeIP     string
	logger     *logrus.Logger
}

func (sd *discoverer) Discover(timeout time.Duration) (client.HTTPClient, error) {

	var endpoint url.URL
	if sd.overridenEndpoint != "" {
		ep, err := url.Parse(sd.overridenEndpoint)
		if err != nil {
			return nil, fmt.Errorf("wrong user-provided KSM endpoint: %s", err)
		}
		endpoint = *ep
	} else {
		var err error
		endpoint, err = sd.dnsDiscover()
		if err != nil {
			// if DNS discovery fails, we dig into Kubernetes API to get the service data
			endpoint, err = sd.apiDiscover()
			if err != nil {
				return nil, fmt.Errorf("failed to discover kube-state-metrics endpoint, got error: %s", err)
			}
		}
	}

	// KSM and Prometheus only work with HTTP
	endpoint.Scheme = "http"
	nodeIP, err := sd.nodeIP()
	if err != nil {
		return nil, fmt.Errorf("failed to discover nodeIP with kube-state-metrics, got error: %s", err)
	}

	return &ksm{
		nodeIP:   nodeIP,
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: sd.logger,
	}, nil
}

func (c *ksm) NodeIP() string {
	return c.nodeIP
}

func (c *ksm) Do(method, path string) (*http.Response, error) {
	e := c.endpoint
	e.Path = filepath.Join(c.endpoint.Path, path)

	r, err := prometheus.NewRequest(method, e.String())
	if err != nil {
		return nil, fmt.Errorf("Error creating %s request to: %s. Got error: %s ", method, e.String(), err)
	}

	c.logger.Debugf("Calling kube-state-metrics endpoint: %s", r.URL.String())

	return c.httpClient.Do(r)
}

// dnsDiscover uses DNS to discover KSM
func (sd *discoverer) dnsDiscover() (url.URL, error) {
	var endpoint url.URL
	_, addrs, err := sd.lookupSRV(ksmDNSService, ksmDNSProto, ksmQualifiedName)
	if err == nil {
		for _, addr := range addrs {
			endpoint.Host = fmt.Sprintf("%v:%v", ksmQualifiedName, addr.Port)
			return endpoint, nil
		}
	}
	return endpoint, fmt.Errorf("can't get DNS port for %s", ksmQualifiedName)
}

// apiDiscover uses Kubernetes API to discover KSM
func (sd *discoverer) apiDiscover() (url.URL, error) {
	var endpoint url.URL

	services, err := sd.apiClient.FindServiceByLabel(ksmAppLabelName, ksmAppLabelValue)
	if err != nil {
		return endpoint, err
	}

	if len(services.Items) == 0 {
		return endpoint, fmt.Errorf("no service found by label %s=%s", ksmAppLabelName, ksmAppLabelValue)
	}

	for _, service := range services.Items {
		if service.Spec.ClusterIP != "" && len(service.Spec.Ports) > 0 {
			// Look for a port called "http-metrics"
			for _, port := range service.Spec.Ports {
				if port.Name == ksmPortName {
					endpoint.Host = fmt.Sprintf("%v:%v", service.Spec.ClusterIP, port.Port)
					return endpoint, nil
				}
			}
			// If not found, return the first TCP port
			for _, port := range service.Spec.Ports {
				if port.Protocol == k8sTCP {
					endpoint.Host = fmt.Sprintf("%v:%v", service.Spec.ClusterIP, port.Port)
					return endpoint, nil
				}
			}
		}
	}

	return endpoint, fmt.Errorf("could not guess the Kube State Metrics host/port")
}

// NodeIP discover IP of a node, where kube-state-metrics is installed
func (sd *discoverer) nodeIP() (string, error) {
	pods, err := sd.apiClient.FindPodsByLabel(ksmAppLabelName, ksmAppLabelValue)
	if err != nil {
		return "", err
	}
	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no pod found by label %s=%s", ksmAppLabelName, ksmAppLabelValue)
	}
	// In case there are multiple pods for the same service, we must be sure we always show the Node IP of the
	// same pod. So we chose, for example, the HostIp with highest precedence in alphabetical order
	var nodeIP string
	for _, pod := range pods.Items {
		if pod.Status.HostIP != "" && (nodeIP == "" || strings.Compare(pod.Status.HostIP, nodeIP) < 0) {
			nodeIP = pod.Status.HostIP
		}
	}
	if nodeIP == "" {
		return "", errors.New("no HostIP address found for KSM node")
	}
	return nodeIP, nil
}

// NewDiscoverer instantiates a new Discoverer
func NewDiscoverer(logger *logrus.Logger) (client.Discoverer, error) {
	return NewNodeIPDiscoverer("", logger)
}

func NewNodeIPDiscoverer(ksmEndpoint string, logger *logrus.Logger) (client.Discoverer, error) {
	var discoverer discoverer
	var err error

	discoverer.apiClient, err = client.NewKubernetes()
	if err != nil {
		return nil, err
	}
	discoverer.lookupSRV = net.LookupSRV
	discoverer.logger = logger
	discoverer.overridenEndpoint = ksmEndpoint
	return &discoverer, nil
}
