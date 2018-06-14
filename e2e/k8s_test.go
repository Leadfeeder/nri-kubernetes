package e2e

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/newrelic/infra-integrations-beta/integrations/kubernetes/e2e/jsonschema"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	nrLabelKey   = "name"
	nrLabelValue = "newrelic-infra"
	namespace    = "default"
	nrContainer  = "newrelic-infra"
)

type integrationData struct {
	role    string
	podName string
	stdOut  []byte
	stdErr  []byte
	err     error
}

var dataChannel = make(chan integrationData)

func getNRPods(clientset *kubernetes.Clientset, config *rest.Config) ([]string, error) {
	sv, err := clientset.ServerVersion()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Executing our integration in %q cluster. K8s version: %s\n", config.Host, sv.String())
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", nrLabelKey, nrLabelValue),
	})
	if err != nil {
		return nil, fmt.Errorf("error retrieving pod list by label %s = %s: %v", nrLabelKey, nrLabelValue, err)
	}
	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("pods not found by label: %s=%s", nrLabelKey, nrLabelValue)
	}
	podsName := make([]string, 0)
	for i := 0; i < len(pods.Items); i++ {
		podsName = append(podsName, pods.Items[i].Name)
	}
	return podsName, nil
}

func execIntegration(clientset *kubernetes.Clientset, config *rest.Config, podName string, wg *sync.WaitGroup) {
	defer wg.Done()
	d := integrationData{
		podName: podName,
	}

	execReq := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", nrContainer).
		Param("command", "/var/db/newrelic-infra/newrelic-integrations/bin/nr-kubernetes").
		// Param("command", "-pretty").
		Param("command", "-verbose").
		Param("stdin", "false").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
	)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", execReq.URL())
	if err != nil {
		d.err = fmt.Errorf("failed to init executor for pod %s: %v", podName, err)
		dataChannel <- d
		return

	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &execOut,
		Stderr: &execErr,
	})

	if err != nil {
		d.err = fmt.Errorf("could not execute command inside pod %s: %v. Integration error output:\n\n%v", podName, err, execErr.String())
		dataChannel <- d
		return
	}

	re, err := regexp.Compile("Auto-discovered role = (\\w*)")
	if err != nil {
		d.err = fmt.Errorf("cannot compile regex and determine role for pod %s, err: %v", podName, err)
		dataChannel <- d
		return
	}

	matches := re.FindStringSubmatch(execErr.String())
	role := matches[1]
	if role == "" {
		d.err = fmt.Errorf("cannot find a role for pod %s", podName)
		dataChannel <- d
		return
	}

	d.role = role
	d.stdOut = execOut.Bytes()
	d.stdErr = execErr.Bytes()

	dataChannel <- d

}

func TestBasic(t *testing.T) {
	if os.Getenv("RUN_TESTS") == "" {
		t.Skip("Flag RUN_TESTS is not specified, skipping tests")
	}
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	podsName, err := getNRPods(clientset, config)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	output := make(map[string]integrationData)

	var wg sync.WaitGroup
	wg.Add(len(podsName))
	go func() {
		wg.Wait()
		close(dataChannel)
	}()

	for _, pName := range podsName {
		fmt.Printf("Executing integration inside pod: %s\n", pName)
		go execIntegration(clientset, config, pName, &wg)
	}

	for d := range dataChannel {
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		output[d.podName] = d
	}

	leaderMap := jsonschema.EventTypeToSchemaFilepath{
		"K8sReplicasetSample": "schema/replicaset.json",
		"K8sNamespaceSample":  "schema/namespace.json",
		"K8sDeploymentSample": "schema/deployment.json",
		"K8sPodSample":        "schema/pod.json",
		"K8sContainerSample":  "schema/container.json",
		"K8sNodeSample":       "schema/node.json",
	}

	followerMap := jsonschema.EventTypeToSchemaFilepath{
		"K8sPodSample":       leaderMap["K8sPodSample"],
		"K8sContainerSample": leaderMap["K8sContainerSample"],
		"K8sNodeSample":      leaderMap["K8sNodeSample"],
	}

	var errs []error
	var lcount int
	var fcount int

	for podName, o := range output {
		var m jsonschema.EventTypeToSchemaFilepath
		switch o.role {
		case "leader":
			lcount++
			m = leaderMap
		case "follower":
			fcount++
			m = followerMap
		}

		err := jsonschema.Match(o.stdOut, m)
		if err != nil {
			errs = append(errs, fmt.Errorf("\n------ %s pod %s ------\n%s", o.role, podName, err))
		}
	}
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	assert.NoError(t, err)

	assert.Equal(t, (lcount + fcount), len(nodes.Items))
	assert.Equal(t, 1, lcount)
	assert.Empty(t, errs)
}