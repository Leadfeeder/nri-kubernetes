package metric

import (
	"errors"
	"fmt"
	"time"

	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/definition"
	ksmMetric "github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/ksm/metric"
	kubeletMetric "github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/kubelet/metric"
	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/prometheus"
	sdkMetric "github.com/newrelic/infra-integrations-sdk/metric"
)

// KSMSpecs are the metric specifications we want to collect from KSM.
var KSMSpecs = definition.SpecGroups{
	"replicaset": {
		IDGenerator:   prometheus.FromLabelValueEntityIDGenerator("kube_replicaset_created", "replicaset"),
		TypeGenerator: prometheus.FromLabelValueEntityTypeGenerator("kube_replicaset_created"),
		Specs: []definition.Spec{
			{Name: "createdAt", ValueFunc: prometheus.FromValue("kube_replicaset_created"), Type: sdkMetric.GAUGE},
			{Name: "podsDesired", ValueFunc: prometheus.FromValue("kube_replicaset_spec_replicas"), Type: sdkMetric.GAUGE},
			{Name: "podsReady", ValueFunc: prometheus.FromValue("kube_replicaset_status_ready_replicas"), Type: sdkMetric.GAUGE},
			{Name: "podsTotal", ValueFunc: prometheus.FromValue("kube_replicaset_status_replicas"), Type: sdkMetric.GAUGE},
			{Name: "podsFullyLabeled", ValueFunc: prometheus.FromValue("kube_replicaset_status_fully_labeled_replicas"), Type: sdkMetric.GAUGE},
			{Name: "observedGeneration", ValueFunc: prometheus.FromValue("kube_replicaset_status_observed_generation"), Type: sdkMetric.GAUGE},
			{Name: "replicasetName", ValueFunc: prometheus.FromLabelValue("kube_replicaset_created", "replicaset"), Type: sdkMetric.ATTRIBUTE},
			{Name: "namespace", ValueFunc: prometheus.FromLabelValue("kube_replicaset_created", "namespace"), Type: sdkMetric.ATTRIBUTE},
			{Name: "deploymentName", ValueFunc: ksmMetric.GetDeploymentNameForReplicaSet(), Type: sdkMetric.ATTRIBUTE},
		},
	},
	"namespace": {
		TypeGenerator: prometheus.FromLabelValueEntityTypeGenerator("kube_namespace_created"),
		Specs: []definition.Spec{
			{Name: "createdAt", ValueFunc: prometheus.FromValue("kube_namespace_created"), Type: sdkMetric.GAUGE},
			{Name: "namespace", ValueFunc: prometheus.FromLabelValue("kube_namespace_created", "namespace"), Type: sdkMetric.ATTRIBUTE},
			{Name: "status", ValueFunc: prometheus.FromLabelValue("kube_namespace_status_phase", "phase"), Type: sdkMetric.ATTRIBUTE},
			{Name: "label.*", ValueFunc: prometheus.InheritAllLabelsFrom("namespace", "kube_namespace_labels"), Type: sdkMetric.ATTRIBUTE},
		},
	},
	"deployment": {
		IDGenerator:   prometheus.FromLabelValueEntityIDGenerator("kube_deployment_created", "deployment"),
		TypeGenerator: prometheus.FromLabelValueEntityTypeGenerator("kube_deployment_created"),
		Specs: []definition.Spec{
			{Name: "podsDesired", ValueFunc: prometheus.FromValue("kube_deployment_spec_replicas"), Type: sdkMetric.GAUGE},
			{Name: "createdAt", ValueFunc: prometheus.FromValue("kube_deployment_created"), Type: sdkMetric.GAUGE},
			{Name: "podsTotal", ValueFunc: prometheus.FromValue("kube_deployment_status_replicas"), Type: sdkMetric.GAUGE},
			{Name: "podsAvailable", ValueFunc: prometheus.FromValue("kube_deployment_status_replicas_available"), Type: sdkMetric.GAUGE},
			{Name: "podsUnavailable", ValueFunc: prometheus.FromValue("kube_deployment_status_replicas_unavailable"), Type: sdkMetric.GAUGE},
			{Name: "podsUpdated", ValueFunc: prometheus.FromValue("kube_deployment_status_replicas_updated"), Type: sdkMetric.GAUGE},
			{Name: "podsMaxUnavailable", ValueFunc: prometheus.FromValue("kube_deployment_spec_strategy_rollingupdate_max_unavailable"), Type: sdkMetric.GAUGE},
			{Name: "namespace", ValueFunc: prometheus.FromLabelValue("kube_deployment_labels", "namespace"), Type: sdkMetric.ATTRIBUTE},
			{Name: "deploymentName", ValueFunc: prometheus.FromLabelValue("kube_deployment_labels", "deployment"), Type: sdkMetric.ATTRIBUTE},
			// Important: The order of these lines is important: we could have the same label in different entities, and we would like to keep the value closer to deployment
			{Name: "label.*", ValueFunc: prometheus.InheritAllLabelsFrom("namespace", "kube_namespace_labels"), Type: sdkMetric.ATTRIBUTE},
			{Name: "label.*", ValueFunc: prometheus.InheritAllLabelsFrom("deployment", "kube_deployment_labels"), Type: sdkMetric.ATTRIBUTE},
		},
	},
	// We get Pod metrics from kube-state-metrics for those pods that are in
	// "Pending" status and are not scheduled. We can't get the data from Kubelet because
	// they aren't running in any node and the information about them is only
	// present in the API.
	"pod": {
		IDGenerator:   prometheus.FromLabelsValueEntityIDGeneratorForPendingPods(),
		TypeGenerator: prometheus.FromLabelValueEntityTypeGenerator("kube_pod_status_phase"),
		Specs: []definition.Spec{
			{Name: "createdAt", ValueFunc: prometheus.FromValue("kube_pod_created"), Type: sdkMetric.GAUGE},
			{Name: "startTime", ValueFunc: prometheus.FromValue("kube_pod_start_time"), Type: sdkMetric.GAUGE},
			{Name: "createdKind", ValueFunc: prometheus.FromLabelValue("kube_pod_info", "created_by_kind"), Type: sdkMetric.ATTRIBUTE},
			{Name: "createdBy", ValueFunc: prometheus.FromLabelValue("kube_pod_info", "created_by_name"), Type: sdkMetric.ATTRIBUTE},
			{Name: "nodeIP", ValueFunc: prometheus.FromLabelValue("kube_pod_info", "host_ip"), Type: sdkMetric.ATTRIBUTE},
			{Name: "namespace", ValueFunc: prometheus.FromLabelValue("kube_pod_info", "namespace"), Type: sdkMetric.ATTRIBUTE},
			{Name: "nodeName", ValueFunc: prometheus.FromLabelValue("kube_pod_info", "node"), Type: sdkMetric.ATTRIBUTE},
			{Name: "podName", ValueFunc: prometheus.FromLabelValue("kube_pod_info", "pod"), Type: sdkMetric.ATTRIBUTE},
			{Name: "isReady", ValueFunc: definition.Transform(prometheus.FromLabelValue("kube_pod_status_ready", "condition"), toNumericBoolean), Type: sdkMetric.GAUGE},
			{Name: "status", ValueFunc: prometheus.FromLabelValue("kube_pod_status_phase", "phase"), Type: sdkMetric.ATTRIBUTE},
			{Name: "isScheduled", ValueFunc: definition.Transform(prometheus.FromLabelValue("kube_pod_status_scheduled", "condition"), toNumericBoolean), Type: sdkMetric.GAUGE},
			{Name: "deploymentName", ValueFunc: ksmMetric.GetDeploymentNameForPod(), Type: sdkMetric.ATTRIBUTE},
			{Name: "label.*", ValueFunc: prometheus.InheritAllLabelsFrom("pod", "kube_pod_labels"), Type: sdkMetric.ATTRIBUTE},
		},
	},
}

// KSMQueries are the queries we will do to KSM in order to fetch all the raw metrics.
var KSMQueries = []prometheus.Query{
	{
		MetricName: "kube_replicaset_spec_replicas",
	},
	{
		MetricName: "kube_replicaset_status_ready_replicas",
	},
	{
		MetricName: "kube_replicaset_status_replicas",
	},
	{
		MetricName: "kube_replicaset_status_fully_labeled_replicas",
	},
	{
		MetricName: "kube_replicaset_status_observed_generation",
	},
	{
		MetricName: "kube_replicaset_created",
	},
	{
		MetricName: "kube_namespace_labels",
		Value: prometheus.QueryValue{
			Value: prometheus.GaugeValue(1),
		},
	},
	{
		MetricName: "kube_namespace_created",
	},
	{
		MetricName: "kube_namespace_status_phase",
		Value: prometheus.QueryValue{
			Value: prometheus.GaugeValue(1),
		},
	},
	{
		MetricName: "kube_namespace_created",
	},
	{
		MetricName: "kube_deployment_labels",
		Value: prometheus.QueryValue{
			Value: prometheus.GaugeValue(1),
		},
	},
	{
		MetricName: "kube_deployment_created",
	},
	{
		MetricName: "kube_deployment_spec_replicas",
	},
	{
		MetricName: "kube_deployment_status_replicas",
	},
	{
		MetricName: "kube_deployment_status_replicas_available",
	},
	{
		MetricName: "kube_deployment_status_replicas_unavailable",
	},
	{
		MetricName: "kube_deployment_status_replicas_updated",
	},
	{
		MetricName: "kube_deployment_spec_strategy_rollingupdate_max_unavailable",
	},
	{
		MetricName: "kube_pod_status_phase",
		Labels: prometheus.QueryLabels{
			Labels: prometheus.Labels{"phase": "Pending"},
		},
		Value: prometheus.QueryValue{
			Value: prometheus.GaugeValue(1),
		},
	},
	{
		MetricName: "kube_pod_info",
	},
	{
		MetricName: "kube_pod_created",
	},
	{
		MetricName: "kube_pod_labels",
	},
	{
		MetricName: "kube_pod_status_scheduled",
		Value: prometheus.QueryValue{
			Value: prometheus.GaugeValue(1),
		},
	},
	{
		MetricName: "kube_pod_status_ready",
		Value: prometheus.QueryValue{
			Value: prometheus.GaugeValue(1),
		},
	},
	{
		MetricName: "kube_pod_start_time",
	},
}

// CadvisorQueries are the queries we will do to the kubelet metrics cadvisor endpoint in order to fetch all the raw metrics.
var CadvisorQueries = []prometheus.Query{
	{
		MetricName: "container_memory_usage_bytes",
		Labels: prometheus.QueryLabels{
			Operator: prometheus.QueryOpNor,
			Labels: prometheus.Labels{
				"container_name": "",
			},
		},
	},
}

// KubeletSpecs are the metric specifications we want to collect from Kubelet.
var KubeletSpecs = definition.SpecGroups{
	"pod": {
		IDGenerator:   kubeletMetric.FromRawEntityIDGroupEntityIDGenerator("namespace"),
		TypeGenerator: kubeletMetric.FromRawGroupsEntityTypeGenerator,
		Specs: []definition.Spec{
			// /stats/summary endpoint
			{Name: "net.rxBytesPerSecond", ValueFunc: definition.FromRaw("rxBytes"), Type: sdkMetric.RATE},
			{Name: "net.txBytesPerSecond", ValueFunc: definition.FromRaw("txBytes"), Type: sdkMetric.RATE},
			{Name: "net.errorsPerSecond", ValueFunc: definition.FromRaw("errors"), Type: sdkMetric.RATE},

			// /pods endpoint
			{Name: "createdAt", ValueFunc: definition.Transform(definition.FromRaw("createdAt"), toTimestamp), Type: sdkMetric.GAUGE},
			{Name: "startTime", ValueFunc: definition.Transform(definition.FromRaw("startTime"), toTimestamp), Type: sdkMetric.GAUGE},
			{Name: "createdKind", ValueFunc: definition.FromRaw("createdKind"), Type: sdkMetric.ATTRIBUTE},
			{Name: "createdBy", ValueFunc: definition.FromRaw("createdBy"), Type: sdkMetric.ATTRIBUTE},
			{Name: "nodeIP", ValueFunc: definition.FromRaw("nodeIP"), Type: sdkMetric.ATTRIBUTE},
			{Name: "namespace", ValueFunc: definition.FromRaw("namespace"), Type: sdkMetric.ATTRIBUTE},
			{Name: "nodeName", ValueFunc: definition.FromRaw("nodeName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "podName", ValueFunc: definition.FromRaw("podName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "isReady", ValueFunc: definition.Transform(definition.FromRaw("isReady"), toNumericBoolean), Type: sdkMetric.GAUGE},
			{Name: "status", ValueFunc: definition.FromRaw("status"), Type: sdkMetric.ATTRIBUTE},
			{Name: "isScheduled", ValueFunc: definition.Transform(definition.FromRaw("isScheduled"), toNumericBoolean), Type: sdkMetric.GAUGE},
			{Name: "deploymentName", ValueFunc: definition.FromRaw("deploymentName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "label.*", ValueFunc: definition.Transform(definition.FromRaw("labels"), kubeletMetric.OneMetricPerLabel), Type: sdkMetric.ATTRIBUTE},
		},
	},
	"container": {
		IDGenerator:   kubeletMetric.FromRawGroupsEntityIDGenerator("containerName"),
		TypeGenerator: kubeletMetric.FromRawGroupsEntityTypeGenerator,
		Specs: []definition.Spec{
			// /stats/summary endpoint
			{Name: "memoryUsedBytes", ValueFunc: definition.FromRaw("usageBytes"), Type: sdkMetric.GAUGE},
			{Name: "cpuUsedCores", ValueFunc: definition.Transform(definition.FromRaw("usageNanoCores"), fromNano), Type: sdkMetric.GAUGE},
			{Name: "fsAvailableBytes", ValueFunc: definition.FromRaw("fsAvailableBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsCapacityBytes", ValueFunc: definition.FromRaw("fsCapacityBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsUsedBytes", ValueFunc: definition.FromRaw("fsUsedBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsUsedPercent", ValueFunc: toComplementPercentage("fsUsedBytes", "fsAvailableBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsInodesFree", ValueFunc: definition.FromRaw("fsInodesFree"), Type: sdkMetric.GAUGE},
			{Name: "fsInodes", ValueFunc: definition.FromRaw("fsInodes"), Type: sdkMetric.GAUGE},
			{Name: "fsInodesUsed", ValueFunc: definition.FromRaw("fsInodesUsed"), Type: sdkMetric.GAUGE},

			// /metrics/cadvisor endpoint
			{Name: "containerID", ValueFunc: definition.FromRaw("containerID"), Type: sdkMetric.ATTRIBUTE},
			{Name: "containerImageID", ValueFunc: definition.FromRaw("containerImageID"), Type: sdkMetric.ATTRIBUTE},

			// /pods endpoint
			{Name: "containerName", ValueFunc: definition.FromRaw("containerName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "containerImage", ValueFunc: definition.FromRaw("containerImage"), Type: sdkMetric.ATTRIBUTE},
			{Name: "deploymentName", ValueFunc: definition.FromRaw("deploymentName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "namespace", ValueFunc: definition.FromRaw("namespace"), Type: sdkMetric.ATTRIBUTE},
			{Name: "podName", ValueFunc: definition.FromRaw("podName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "nodeName", ValueFunc: definition.FromRaw("nodeName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "nodeIP", ValueFunc: definition.FromRaw("nodeIP"), Type: sdkMetric.ATTRIBUTE},
			{Name: "restartCount", ValueFunc: definition.FromRaw("restartCount"), Type: sdkMetric.GAUGE},
			{Name: "cpuRequestedCores", ValueFunc: definition.Transform(definition.FromRaw("cpuRequestedCores"), toCores), Type: sdkMetric.GAUGE},
			{Name: "cpuLimitCores", ValueFunc: definition.Transform(definition.FromRaw("cpuLimitCores"), toCores), Type: sdkMetric.GAUGE},
			{Name: "memoryRequestedBytes", ValueFunc: definition.FromRaw("memoryRequestedBytes"), Type: sdkMetric.GAUGE},
			{Name: "memoryLimitBytes", ValueFunc: definition.FromRaw("memoryLimitBytes"), Type: sdkMetric.GAUGE},
			{Name: "status", ValueFunc: definition.FromRaw("status"), Type: sdkMetric.ATTRIBUTE},
			{Name: "isReady", ValueFunc: definition.Transform(definition.FromRaw("isReady"), toNumericBoolean), Type: sdkMetric.GAUGE},
			{Name: "reason", ValueFunc: definition.FromRaw("reason"), Type: sdkMetric.ATTRIBUTE}, // Previously called statusWaitingReason

			// Inherit from pod
			{Name: "label.*", ValueFunc: definition.Transform(definition.FromRaw("labels"), kubeletMetric.OneMetricPerLabel), Type: sdkMetric.ATTRIBUTE},
		},
	},
	"node": {
		TypeGenerator: kubeletMetric.FromRawGroupsEntityTypeGenerator,
		Specs: []definition.Spec{
			{Name: "nodeName", ValueFunc: definition.FromRaw("nodeName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "cpuUsedCores", ValueFunc: definition.Transform(definition.FromRaw("usageNanoCores"), fromNano), Type: sdkMetric.GAUGE},
			{Name: "cpuUsedCoreMilliseconds", ValueFunc: definition.Transform(definition.FromRaw("usageCoreNanoSeconds"), fromNanoToMilli), Type: sdkMetric.GAUGE},
			{Name: "memoryUsedBytes", ValueFunc: definition.FromRaw("memoryUsageBytes"), Type: sdkMetric.GAUGE},
			{Name: "memoryAvailableBytes", ValueFunc: definition.FromRaw("memoryAvailableBytes"), Type: sdkMetric.GAUGE},
			{Name: "memoryWorkingSetBytes", ValueFunc: definition.FromRaw("memoryWorkingSetBytes"), Type: sdkMetric.GAUGE},
			{Name: "memoryRssBytes", ValueFunc: definition.FromRaw("memoryRssBytes"), Type: sdkMetric.GAUGE},
			{Name: "memoryPageFaults", ValueFunc: definition.FromRaw("memoryPageFaults"), Type: sdkMetric.GAUGE},
			{Name: "memoryMajorPageFaultsPerSecond", ValueFunc: definition.FromRaw("memoryMajorPageFaults"), Type: sdkMetric.RATE},
			{Name: "net.rxBytesPerSecond", ValueFunc: definition.FromRaw("rxBytes"), Type: sdkMetric.RATE},
			{Name: "net.txBytesPerSecond", ValueFunc: definition.FromRaw("txBytes"), Type: sdkMetric.RATE},
			{Name: "net.errorsPerSecond", ValueFunc: definition.FromRaw("errors"), Type: sdkMetric.RATE},
			{Name: "fsAvailableBytes", ValueFunc: definition.FromRaw("fsAvailableBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsCapacityBytes", ValueFunc: definition.FromRaw("fsCapacityBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsUsedBytes", ValueFunc: definition.FromRaw("fsUsedBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsInodesFree", ValueFunc: definition.FromRaw("fsInodesFree"), Type: sdkMetric.GAUGE},
			{Name: "fsInodes", ValueFunc: definition.FromRaw("fsInodes"), Type: sdkMetric.GAUGE},
			{Name: "fsInodesUsed", ValueFunc: definition.FromRaw("fsInodesUsed"), Type: sdkMetric.GAUGE},
			{Name: "runtimeAvailableBytes", ValueFunc: definition.FromRaw("runtimeAvailableBytes"), Type: sdkMetric.GAUGE},
			{Name: "runtimeCapacityBytes", ValueFunc: definition.FromRaw("runtimeCapacityBytes"), Type: sdkMetric.GAUGE},
			{Name: "runtimeUsedBytes", ValueFunc: definition.FromRaw("runtimeUsedBytes"), Type: sdkMetric.GAUGE},
			{Name: "runtimeInodesFree", ValueFunc: definition.FromRaw("runtimeInodesFree"), Type: sdkMetric.GAUGE},
			{Name: "runtimeInodes", ValueFunc: definition.FromRaw("runtimeInodes"), Type: sdkMetric.GAUGE},
			{Name: "runtimeInodesUsed", ValueFunc: definition.FromRaw("runtimeInodesUsed"), Type: sdkMetric.GAUGE},
		},
	},
	"volume": {
		TypeGenerator: kubeletMetric.FromRawGroupsEntityTypeGenerator,
		Specs: []definition.Spec{
			{Name: "volumeName", ValueFunc: definition.FromRaw("volumeName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "podName", ValueFunc: definition.FromRaw("podName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "namespace", ValueFunc: definition.FromRaw("namespace"), Type: sdkMetric.ATTRIBUTE},
			{Name: "persistent", ValueFunc: isPersistentVolume(), Type: sdkMetric.ATTRIBUTE},
			{Name: "pvcName", ValueFunc: definition.FromRaw("pvcName"), Type: sdkMetric.ATTRIBUTE},
			{Name: "pvcNamespace", ValueFunc: definition.FromRaw("pvcNamespace"), Type: sdkMetric.ATTRIBUTE},
			{Name: "fsAvailableBytes", ValueFunc: definition.FromRaw("fsAvailableBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsCapacityBytes", ValueFunc: definition.FromRaw("fsCapacityBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsUsedBytes", ValueFunc: definition.FromRaw("fsUsedBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsUsedPercent", ValueFunc: toComplementPercentage("fsUsedBytes", "fsAvailableBytes"), Type: sdkMetric.GAUGE},
			{Name: "fsInodesFree", ValueFunc: definition.FromRaw("fsInodesFree"), Type: sdkMetric.GAUGE},
			{Name: "fsInodes", ValueFunc: definition.FromRaw("fsInodes"), Type: sdkMetric.GAUGE},
			{Name: "fsInodesUsed", ValueFunc: definition.FromRaw("fsInodesUsed"), Type: sdkMetric.GAUGE},
		},
	},
}

func isPersistentVolume() definition.FetchFunc {
	return func(groupLabel, entityID string, groups definition.RawGroups) (definition.FetchedValue, error) {
		name, err := definition.FromRaw("pvcName")(groupLabel, entityID, groups)
		if err == nil && name != "" {
			return "true", nil
		}
		return "false", nil
	}
}

func computePercentage(current, all uint64) (definition.FetchedValue, error) {
	if all == uint64(0) {
		return nil, errors.New("division by zero")
	}
	return ((float64(current) / float64(all)) * 100), nil
}

func toComplementPercentage(desiredMetric, complementMetric string) definition.FetchFunc {
	return func(groupLabel, entityID string, groups definition.RawGroups) (definition.FetchedValue, error) {
		complement, err := definition.FromRaw(complementMetric)(groupLabel, entityID, groups)
		if err != nil {
			return nil, err
		}
		desired, err := definition.FromRaw(desiredMetric)(groupLabel, entityID, groups)
		if err != nil {
			return nil, err
		}
		v, err := computePercentage(desired.(uint64), desired.(uint64)+complement.(uint64))
		if err != nil {
			return nil, fmt.Errorf("error computing percentage for %s & %s: %s", desiredMetric, complementMetric, err)
		}
		return v, nil
	}
}

// Used to transform from usageNanoCores to cpuUsedCores
func fromNano(value definition.FetchedValue) (definition.FetchedValue, error) {
	v, ok := value.(uint64)
	if !ok {
		return nil, errors.New("error transforming to cpu cores")
	}

	return float64(v) / 1000000000, nil
}

func fromNanoToMilli(value definition.FetchedValue) (definition.FetchedValue, error) {
	v, ok := value.(uint64)
	if !ok {
		return nil, errors.New("error transforming cpu cores to milliseconds")
	}

	return float64(v) / 1000000, nil
}

func toTimestamp(value definition.FetchedValue) (definition.FetchedValue, error) {
	v, ok := value.(time.Time)
	if !ok {
		return nil, errors.New("error transforming to timestamp")
	}

	return v.Unix(), nil
}

func toNumericBoolean(value definition.FetchedValue) (definition.FetchedValue, error) {
	switch value {
	case "true", "True", true, 1:
		return 1, nil
	case "false", "False", false, 0:
		return 0, nil
	default:
		return nil, errors.New("value can not be converted to numeric boolean")
	}
}

func toCores(value definition.FetchedValue) (definition.FetchedValue, error) {
	switch v := value.(type) {
	case int:
		return float64(v) / 1000, nil
	case int64:
		return float64(v) / 1000, nil
	default:
		return nil, errors.New("error transforming to cores")
	}
}
