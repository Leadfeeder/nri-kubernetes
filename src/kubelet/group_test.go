package kubelet

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/kubelet/metric"
	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/kubelet/metric/testdata"
	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/src/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type testClient struct {
	handler http.HandlerFunc
}

func (c *testClient) Do(method, path string) (*http.Response, error) {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	c.handler(w, req)

	return w.Result(), nil
}

func (c *testClient) NodeIP() string {
	// nothing to do
	return ""
}

func rawGroupsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case metric.KubeletPodsPath:
		f, err := os.Open("metric/testdata/kubelet_pods_payload.json") // TODO move fetch and testdata to just kubelet package.
		if err != nil {
			panic(err)
		}

		defer f.Close() // nolint: errcheck

		io.Copy(w, f)
	case metric.StatsSummaryPath:
		f, err := os.Open("metric/testdata/kubelet_stats_summary_payload.json") // TODO move fetch and testdata to just kubelet package.
		if err != nil {
			panic(err)
		}

		defer f.Close() // nolint: errcheck

		io.Copy(w, f)
	case metric.KubeletCAdvisorMetricsPath:
		f, err := os.Open("metric/testdata/kubelet_metrics_cadvisor_payload_plain.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close() // nolint: errcheck

		prometheus.TextToProtoHandleFunc(f)(w, r)
	}

}

func TestGroup(t *testing.T) {
	c := testClient{
		handler: rawGroupsHandlerFunc,
	}

	queries := []prometheus.Query{
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

	grouper := NewGrouper(&c, logrus.StandardLogger(), metric.PodsFetchFunc(&c), metric.CadvisorFetchFunc(&c, queries))
	r, errGroup := grouper.Group(nil)
	assert.Nil(t, errGroup)

	assert.Equal(t, testdata.ExpectedGroupData, r)
}
