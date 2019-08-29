package testdata

import (
	"github.com/newrelic/nri-kubernetes/src/definition"
)

// ExpectedGroupData is the expectation for main group_test tests.
var ExpectedGroupData = definition.RawGroups{
	"pod": {
		"kube-system_newrelic-infra-rz225": {
			"createdKind": "DaemonSet",
			"createdBy":   "newrelic-infra",
			"nodeIP":      "192.168.99.100",
			"namespace":   "kube-system",
			"podName":     "newrelic-infra-rz225",
			"nodeName":    "minikube",
			"startTime":   parseTime("2018-02-14T16:26:33Z"),
			"status":      "Running",
			"isReady":     "True",
			"isScheduled": "True",
			"createdAt":   parseTime("2018-02-14T16:26:33Z"),
			"labels": map[string]string{
				"controller-revision-hash": "3887482659",
				"name":                     "newrelic-infra",
				"pod-template-generation":  "1",
			},
			"errors":  uint64(0),
			"rxBytes": uint64(106175985),
			"txBytes": uint64(35714359),
		},
		"kube-system_kube-state-metrics-57f4659995-6n2qq": {
			"createdKind":    "ReplicaSet",
			"createdBy":      "kube-state-metrics-57f4659995",
			"nodeIP":         "192.168.99.100",
			"namespace":      "kube-system",
			"podName":        "kube-state-metrics-57f4659995-6n2qq",
			"nodeName":       "minikube",
			"status":         "Running",
			"isReady":        "True",
			"isScheduled":    "True",
			"createdAt":      parseTime("2018-02-14T16:27:38Z"),
			"deploymentName": "kube-state-metrics",
			"labels": map[string]string{
				"k8s-app":           "kube-state-metrics",
				"pod-template-hash": "1390215551",
			},
			"errors":  uint64(0),
			"rxBytes": uint64(32575098),
			"txBytes": uint64(27840584),
		},
		"default_sh-7c95664875-4btqh": {
			"createdKind":    "ReplicaSet",
			"createdBy":      "sh-7c95664875",
			"nodeIP":         "192.168.99.100",
			"namespace":      "default",
			"podName":        "sh-7c95664875-4btqh",
			"nodeName":       "minikube",
			"status":         "Failed",
			"reason":         "Evicted",
			"message":        "The node was low on resource: memory.",
			"createdAt":      parseTime("2019-03-13T07:59:00Z"),
			"startTime":      parseTime("2019-03-13T07:59:00Z"),
			"deploymentName": "sh",
			"labels": map[string]string{
				"pod-template-hash": "3751220431",
				"run":               "sh",
			},
		},
	},
	"container": {
		"kube-system_newrelic-infra-rz225_newrelic-infra": {
			"containerName":        "newrelic-infra",
			"containerID":          "69d7203a8f2d2d027ffa51d61002eac63357f22a17403363ef79e66d1c3146b2",
			"containerImage":       "newrelic/ohaik:1.0.0-beta3",
			"containerImageID":     "sha256:1a95d0df2997f93741fbe2a15d2c31a394e752fd942ec29bf16a44163342f6a1",
			"namespace":            "kube-system",
			"podName":              "newrelic-infra-rz225",
			"nodeName":             "minikube",
			"nodeIP":               "192.168.99.100",
			"restartCount":         int32(6),
			"isReady":              true,
			"status":               "Running",
			"startedAt":            parseTime("2018-02-27T15:21:16Z"),
			"cpuRequestedCores":    int64(100),
			"memoryRequestedBytes": int64(104857600),
			"memoryLimitBytes":     int64(104857600),
			"usageBytes":           uint64(18083840),
			"workingSetBytes":      uint64(17113088),
			"usageNanoCores":       uint64(17428240),
			"fsAvailableBytes":     uint64(14924988416),
			"fsUsedBytes":          uint64(126976),
			"fsCapacityBytes":      uint64(17293533184),
			"fsInodesFree":         uint64(9713372),
			"fsInodes":             uint64(9732096),
			"fsInodesUsed":         uint64(36),
			"labels": map[string]string{
				"controller-revision-hash": "3887482659",
				"name":                     "newrelic-infra",
				"pod-template-generation":  "1",
			},
		},
		"kube-system_kube-state-metrics-57f4659995-6n2qq_kube-state-metrics": {
			"containerName":    "kube-state-metrics",
			"containerID":      "c452821fcf6c5f594d4f98a1426e7a2c51febb65d5d50d92903f9dfb367bfba7",
			"containerImage":   "quay.io/coreos/kube-state-metrics:v1.1.0",
			"containerImageID": "quay.io/coreos/kube-state-metrics@sha256:52a2c47355c873709bb4e37e990d417e9188c2a778a0c38ed4c09776ddc54efb",
			"namespace":        "kube-system",
			"podName":          "kube-state-metrics-57f4659995-6n2qq",
			"nodeName":         "minikube",
			"nodeIP":           "192.168.99.100",
			//"restartCount": int32(7), // No restartCount since there is no restartCount in status field in the pod fetched from kubelet /pods.
			//"isReady":              false, // No isReady since there is no isReady in status field in the pod fetched from kubelet /pods.
			//"status":         "Running", // No Status since there is no ContainerStatuses field in the pod fetched from kubelet /pods.
			//"startedAt":            parseTime("2018-02-27T15:21:37Z"), // No startedAt since there is no startedAt in status field in the pod fetched from kubelet /pods.
			"deploymentName":       "kube-state-metrics",
			"cpuRequestedCores":    int64(101),
			"cpuLimitCores":        int64(101),
			"memoryRequestedBytes": int64(106954752),
			"memoryLimitBytes":     int64(106954752),
			"usageBytes":           uint64(15568896),
			"workingSetBytes":      uint64(15110144),
			"usageNanoCores":       uint64(941138),
			"fsAvailableBytes":     uint64(14924988416),
			"fsUsedBytes":          uint64(28672),
			"fsCapacityBytes":      uint64(17293533184),
			"fsInodesFree":         uint64(9713372),
			"fsInodes":             uint64(9732096),
			"fsInodesUsed":         uint64(7),
			"labels": map[string]string{
				"k8s-app":           "kube-state-metrics",
				"pod-template-hash": "1390215551",
			},
		},
		"kube-system_kube-state-metrics-57f4659995-6n2qq_addon-resizer": {
			"containerName":    "addon-resizer",
			"containerID":      "3328c17bfd22f1a82fcdf8707c2f8f040c462e548c24780079bba95d276d93e1",
			"containerImage":   "gcr.io/google_containers/addon-resizer:1.0",
			"containerImageID": "gcr.io/google_containers/addon-resizer@sha256:e77acf80697a70386c04ae3ab494a7b13917cb30de2326dcf1a10a5118eddabe",
			"namespace":        "kube-system",
			"podName":          "kube-state-metrics-57f4659995-6n2qq",
			"nodeName":         "minikube",
			"nodeIP":           "192.168.99.100",
			//"restartCount": int32(7), // No restartCount since there is no restartCount in status field in the pod fetched from kubelet /pods.
			//"isReady":              false, // No isReady since there is no isReady in status field in the pod fetched from kubelet /pods.
			//"status":         "Running", // No Status since there is no ContainerStatuses field in the pod fetched from kubelet /pods.
			//"startedAt":            parseTime("2018-02-27T15:21:37Z"), // No startedAt since there is no startedAt in status field in the pod fetched from kubelet /pods.
			"deploymentName":       "kube-state-metrics",
			"cpuRequestedCores":    int64(100),
			"cpuLimitCores":        int64(100),
			"memoryRequestedBytes": int64(31457280),
			"memoryLimitBytes":     int64(31457280),
			"usageBytes":           uint64(6373376),
			"workingSetBytes":      uint64(6270976),
			"usageNanoCores":       uint64(131742),
			"fsAvailableBytes":     uint64(14924988416),
			"fsUsedBytes":          uint64(24576),
			"fsCapacityBytes":      uint64(17293533184),
			"fsInodesFree":         uint64(9713372),
			"fsInodes":             uint64(9732096),
			"fsInodesUsed":         uint64(6),
			"labels": map[string]string{
				"k8s-app":           "kube-state-metrics",
				"pod-template-hash": "1390215551",
			},
		},
		"default_sh-7c95664875-4btqh_sh": {
			"containerName":  "sh",
			"containerImage": "python",
			"namespace":      "default",
			"podName":        "sh-7c95664875-4btqh",
			"nodeName":       "minikube",
			"nodeIP":         "192.168.99.100",
			"deploymentName": "sh",
			"labels": map[string]string{
				"pod-template-hash": "3751220431",
				"run":               "sh",
			},
		},
	},
	"node": {
		"minikube": {
			"nodeName":              "minikube",
			"errors":                uint64(0),
			"fsAvailableBytes":      uint64(14924988416),
			"fsCapacityBytes":       uint64(17293533184),
			"fsInodes":              uint64(9732096),
			"fsInodesFree":          uint64(9713372),
			"fsInodesUsed":          uint64(18724),
			"fsUsedBytes":           uint64(1355673600),
			"memoryAvailableBytes":  uint64(791736320),
			"memoryMajorPageFaults": uint64(0),
			"memoryPageFaults":      uint64(113947),
			"memoryRssBytes":        uint64(660684800),
			"memoryUsageBytes":      uint64(1843650560),
			"memoryWorkingSetBytes": uint64(1305468928),
			"runtimeAvailableBytes": uint64(14924988416),
			"runtimeCapacityBytes":  uint64(17293533184),
			"runtimeInodes":         uint64(9732096),
			"runtimeInodesFree":     uint64(9713372),
			"runtimeInodesUsed":     uint64(18724),
			"runtimeUsedBytes":      uint64(969241979),
			"rxBytes":               uint64(1507694406),
			"txBytes":               uint64(120789968),
			"usageCoreNanoSeconds":  uint64(22332102208229),
			"usageNanoCores":        uint64(228759290),
			"labels": map[string]string{
				"kubernetes.io/arch":             "amd64",
				"kubernetes.io/hostname":         "minikube",
				"kubernetes.io/os":               "linux",
				"node-role.kubernetes.io/master": "",
			},
		},
	},
	"volume": {
		"kube-system_heapster-5mz5f_default-token-mmlq2": {
			"fsAvailableBytes": uint64(1048588288),
			"fsCapacityBytes":  uint64(1048600576),
			"fsUsedBytes":      uint64(12288),
			"fsInodesFree":     uint64(255997),
			"fsInodes":         uint64(256006),
			"fsInodesUsed":     uint64(9),
			"volumeName":       "default-token-mmlq2",
			"namespace":        "kube-system",
			"podName":          "heapster-5mz5f",
		},
		"kube-system_influxdb-grafana-rsmwp_default-token-mmlq2": {
			"fsAvailableBytes": uint64(1048588288),
			"fsCapacityBytes":  uint64(1048600576),
			"fsUsedBytes":      uint64(12288),
			"fsInodesFree":     uint64(255997),
			"fsInodes":         uint64(256006),
			"fsInodesUsed":     uint64(9),
			"podName":          "influxdb-grafana-rsmwp",
			"volumeName":       "default-token-mmlq2",
			"namespace":        "kube-system",
		},
		"kube-system_kube-dns-54cccfbdf8-dznm7_default-token-mmlq2": {
			"fsAvailableBytes": uint64(1048588288),
			"fsCapacityBytes":  uint64(1048600576),
			"fsUsedBytes":      uint64(12288),
			"fsInodesFree":     uint64(255997),
			"fsInodes":         uint64(256006),
			"fsInodesUsed":     uint64(9),
			"podName":          "kube-dns-54cccfbdf8-dznm7",
			"volumeName":       "default-token-mmlq2",
			"namespace":        "kube-system",
		},
		"kube-system_newrelic-infra-rz225_default-token-mmlq2": {
			"fsAvailableBytes": uint64(1048588288),
			"fsCapacityBytes":  uint64(1048600576),
			"fsUsedBytes":      uint64(12288),
			"fsInodesFree":     uint64(255997),
			"fsInodes":         uint64(256006),
			"fsInodesUsed":     uint64(9),
			"podName":          "newrelic-infra-rz225",
			"volumeName":       "default-token-mmlq2",
			"namespace":        "kube-system",
		},
		"kube-system_storage-provisioner_default-token-mmlq2": {
			"fsAvailableBytes": uint64(1048588288),
			"fsCapacityBytes":  uint64(1048600576),
			"fsUsedBytes":      uint64(12288),
			"fsInodesFree":     uint64(255997),
			"fsInodes":         uint64(256006),
			"fsInodesUsed":     uint64(9),
			"podName":          "storage-provisioner",
			"volumeName":       "default-token-mmlq2",
			"namespace":        "kube-system",
		},
		"kube-system_kubernetes-dashboard-77d8b98585-mtjld_default-token-mmlq2": {
			"fsAvailableBytes": uint64(1048588288),
			"fsCapacityBytes":  uint64(1048600576),
			"fsUsedBytes":      uint64(12288),
			"fsInodesFree":     uint64(255997),
			"fsInodes":         uint64(256006),
			"fsInodesUsed":     uint64(9),
			"podName":          "kubernetes-dashboard-77d8b98585-mtjld",
			"volumeName":       "default-token-mmlq2",
			"namespace":        "kube-system",
		},
		"kube-system_kube-state-metrics-57f4659995-6n2qq_kube-state-metrics-token-rl9b8": {
			"fsAvailableBytes": uint64(1048588288),
			"fsCapacityBytes":  uint64(1048600576),
			"fsUsedBytes":      uint64(12288),
			"fsInodesFree":     uint64(255997),
			"fsInodes":         uint64(256006),
			"fsInodesUsed":     uint64(9),
			"volumeName":       "kube-state-metrics-token-rl9b8",
			"namespace":        "kube-system",
			"podName":          "kube-state-metrics-57f4659995-6n2qq",
		},
		"kube-system_influxdb-grafana-rsmwp_grafana-storage": {
			"fsAvailableBytes": uint64(14925000704),
			"fsCapacityBytes":  uint64(17293533184),
			"fsUsedBytes":      uint64(413696),
			"fsInodesFree":     uint64(9713372),
			"fsInodes":         uint64(9732096),
			"fsInodesUsed":     uint64(53),
			"volumeName":       "grafana-storage",
			"namespace":        "kube-system",
			"podName":          "influxdb-grafana-rsmwp",
		},
		"kube-system_influxdb-grafana-rsmwp_influxdb-storage": {
			"fsAvailableBytes": uint64(14925000704),
			"fsCapacityBytes":  uint64(17293533184),
			"fsUsedBytes":      uint64(5853184),
			"fsInodesFree":     uint64(9713372),
			"fsInodes":         uint64(9732096),
			"fsInodesUsed":     uint64(42),
			"volumeName":       "influxdb-storage",
			"namespace":        "kube-system",
			"podName":          "influxdb-grafana-rsmwp",
		},
	},
}
