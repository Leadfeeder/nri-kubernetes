package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/client"
	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/data"
	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/definition"

	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/ksm"
	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/kubelet"

	clientKsm "github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/ksm/client"
	clientKubelet "github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/kubelet/client"
	metric2 "github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/kubelet/metric"
	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/metric"
	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/storage"
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/sdk"
	"github.com/sirupsen/logrus"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Timeout             int    `default:"5000" help:"timeout in milliseconds for calling metrics sources"`
	ClusterName         string `help:"Identifier of your cluster. You could use it later to filter data in your New Relic account"`
	DiscoveryCacheDir   string `default:"/var/cache/nr-kubernetes" help:"The location of the cached values for discovered endpoints."`
	DiscoveryCacheTTL   string `default:"1h" help:"Duration since the discovered endpoints are stored in the cache until they expire. Valid time units: 'ns', 'us', 'ms', 's', 'm', 'h'"`
	KubeStateMetricsURL string `help: kube-state-metrics URL. If it is not provided, it will be discovered.`
}

const (
	integrationName    = "com.newrelic.kubernetes"
	integrationVersion = "1.0.0-beta2.2"
	nodeNameEnvVar     = "NRK8S_NODE_NAME"
)

var args argumentList

func populate(grouper data.Grouper, specs definition.SpecGroups, i *sdk.IntegrationProtocol2, clusterName string, logger *logrus.Logger) *data.PopulateErr {
	groups, errs := grouper.Group(specs)
	if errs != nil && len(errs.Errors) > 0 {
		if !errs.Recoverable {
			return &data.PopulateErr{
				Errs:      errs.Errors,
				Populated: false,
			}
		}

		logger.Warnf("%s", errs)
	}

	return metric.NewK8sPopulator().Populate(groups, specs, i, clusterName)
}

func main() {
	integration, err := sdk.NewIntegrationProtocol2(integrationName, integrationVersion, &args)
	exitLog := fmt.Sprintf("Integration %q exited", integrationName)
	if err != nil {
		defer log.Debug(exitLog)
		log.Fatal(err) // Global logs used as args processed inside NewIntegrationProtocol2
	}

	logger := log.New(args.Verbose)
	defer func() {
		if r := recover(); r != nil {
			recErr, ok := r.(*logrus.Entry)
			if ok {
				recErr.Fatal(recErr.Message)
			} else {
				panic(r)
			}
		}
	}()

	defer logger.Debug(exitLog)
	logger.Debugf("Integration %q with version %s started", integrationName, integrationVersion)
	if args.ClusterName == "" {
		logger.Panic(errors.New("cluster_name argument is mandatory"))
	}

	nodeName := os.Getenv(nodeNameEnvVar)
	if nodeName == "" {
		logger.Panicf("%s env var should be provided by Kubernetes and is mandatory", nodeNameEnvVar)
	}

	if args.All || args.Metrics {
		ttl, err := time.ParseDuration(args.DiscoveryCacheTTL)
		if err != nil {
			logger.WithError(err).Error("while parsing the cache TTL value. Defaulting to 1h")
			ttl = time.Hour
		}

		timeout := time.Millisecond * time.Duration(args.Timeout)

		innerKubeletDiscoverer, err := clientKubelet.NewDiscoverer(nodeName, logger)
		if err != nil {
			logger.Panicf("error during Kubelet auto discovering process. %s", err)
		}
		cacheStorage := storage.NewJSONDiskStorage(args.DiscoveryCacheDir)
		kubeletDiscoverer := clientKubelet.NewDiscoveryCacher(innerKubeletDiscoverer, cacheStorage, ttl, logger)

		kubeletClient, err := kubeletDiscoverer.Discover(timeout)
		if err != nil {
			logger.Panic(err)
		}
		kubeletNodeIP := kubeletClient.NodeIP()
		logger.Debugf("Kubelet Node = %s", kubeletNodeIP)

		var innerKSMDiscoverer client.Discoverer
		if args.KubeStateMetricsURL != "" {
			innerKSMDiscoverer, err = clientKsm.NewDiscovererForNodeIP(args.KubeStateMetricsURL, logger)
		} else {
			innerKSMDiscoverer, err = clientKsm.NewDiscoverer(logger)
		}
		if err != nil {
			logger.Panic(err)
		}
		ksmDiscoverer := clientKsm.NewDiscoveryCacher(innerKSMDiscoverer, cacheStorage, ttl, logger)
		ksmClient, err := ksmDiscoverer.Discover(timeout)
		if err != nil {
			logger.Panic(err)
		}
		ksmNodeIP := ksmClient.NodeIP()
		logger.Debugf("KSM Node = %s", ksmNodeIP)

		// setting role by auto discovery
		var role string
		if kubeletNodeIP == ksmNodeIP {
			role = "leader"
		} else {
			role = "follower"
		}
		logger.Debugf("Auto-discovered role = %s", role)

		kubeletGrouper := kubelet.NewGrouper(kubeletClient, logger, metric2.PodsFetchFunc(kubeletClient), metric2.CadvisorFetchFunc(kubeletClient, metric.CadvisorQueries))

		switch role {
		case "leader":
			kubeletErr := populate(kubeletGrouper, metric.KubeletSpecs, integration, args.ClusterName, logger)
			if kubeletErr != nil {
				// We don't panic as we want to try populating ksm metrics.
				logger.Errorf("Error populating Kubelet metrics: %s", kubeletErr)
			}

			ksmGrouper := ksm.NewGrouper(ksmClient, metric.KSMQueries, logger)
			ksmErr := populate(ksmGrouper, metric.KSMSpecs, integration, args.ClusterName, logger)
			if ksmErr != nil {
				logger.Errorf("Error populating KSM metrics: %s", ksmErr)
			}

			if !ksmErr.Populated && !kubeletErr.Populated {
				// We panic since both populate processes failed.
				logger.Panic("No data was populated")
			}
		case "follower":
			populateErr := populate(kubeletGrouper, metric.KubeletSpecs, integration, args.ClusterName, logger)
			if err != nil {
				logger.Errorf("Error populating Kubelet metrics: %s", err)
			}

			if !populateErr.Populated {
				// We panic since the only populate process failed.
				logger.Panic("No data was populated")
			}
		}

		err = integration.Publish()
		if err != nil {
			logger.Panic(err)
		}
	}
}
