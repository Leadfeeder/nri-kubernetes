package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Kubernetes provides an interface to common Kubernetes API operations
type Kubernetes interface {
	// FindNode returns a Node reference containing the pod named as the argument, if any
	FindNode(name string) (*v1.Node, error)
	// FindPodsByLabel returns a PodList reference containing the pods matching the provided name/value label pair
	FindPodsByLabel(name, value string) (*v1.PodList, error)
	// FindServiceByLabel returns a ServiceList containing the services matching the provided name/value label pair
	// name/value pairs
	FindServiceByLabel(name, value string) (*v1.ServiceList, error)
	// Config returns a config of API client
	Config() *rest.Config
	// SecureHTTPClient returns http.Client configured with timeout and CA Cert
	SecureHTTPClient(time.Duration) (*http.Client, error)
}

type goClientImpl struct {
	client *kubernetes.Clientset
	config *rest.Config
}

func (ka *goClientImpl) Config() *rest.Config {
	return ka.config
}

func (ka *goClientImpl) FindNode(name string) (*v1.Node, error) {
	return ka.client.CoreV1().Nodes().Get(name, metav1.GetOptions{})
}

func (ka *goClientImpl) FindPodsByLabel(name, value string) (*v1.PodList, error) {
	return ka.client.CoreV1().Pods("").List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", name, value),
	})
}

func (ka *goClientImpl) FindServiceByLabel(name, value string) (*v1.ServiceList, error) {
	return ka.client.CoreV1().Services("").List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", name, value),
	})
}

func (ka *goClientImpl) SecureHTTPClient(t time.Duration) (*http.Client, error) {
	c, ok := ka.client.RESTClient().(*rest.RESTClient)
	if !ok {
		return nil, errors.New("failed to set up a client for connecting to Kubelet through API proxy")
	}
	return c.Client, nil
}

// BasicHTTPClient returns http.Client configured with timeout
func BasicHTTPClient(t time.Duration) *http.Client {
	return &http.Client{
		Timeout: t,
	}
}

// InsecureHTTPClient returns http.Client configured with timeout
// and InsecureSkipVerify flag enabled
func InsecureHTTPClient(t time.Duration) *http.Client {
	client := BasicHTTPClient(t)
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return client
}

// NewKubernetes instantiates a Kubernetes API client
func NewKubernetes() (Kubernetes, error) {
	ka := new(goClientImpl)
	var err error

	ka.config, err = rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	ka.client, err = kubernetes.NewForConfig(ka.config)
	if err != nil {
		return nil, err
	}

	return ka, nil
}
