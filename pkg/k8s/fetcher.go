package k8s

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

const (
	CADVISOR_PATH_TEMPLATE = "/api/v1/nodes/%s/proxy/metrics/cadvisor"
)

type Resources struct {
	Cpu    float64
	Memory float64
}

type ContainerResources struct {
	PodName   string
	Name      string
	Namespace string
	NodeName  string
	Image     string
	Request   Resources
	Limit     Resources
}
type NodeMetrics struct {
	NodeName  string
	Cpu       []*Cpu
	Memory    []*Memory
	Timestamp time.Time
}

type Fetcher struct {
	clientset     *kubernetes.Clientset
	metricsParser *Parser
	nodes         []string
}

func NewFetcher(clientset *kubernetes.Clientset) *Fetcher {
	return &Fetcher{
		clientset:     clientset,
		metricsParser: NewParser(),
	}
}

func (f *Fetcher) GetMetrics() ([]*NodeMetrics, error) {
	nodes, err := f.getNodes()
	if err != nil {
		return nil, err
	}

	metrics := make([]*NodeMetrics, 0)
	for _, node := range nodes {
		nodeMetric, err := f.fetchMetricsFromNode(node)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, nodeMetric)
	}

	return metrics, nil
}

func (f *Fetcher) GetContainers() ([]*ContainerResources, error) {
	pods, err := f.clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	containers := make([]*ContainerResources, 0)
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			requestCpu := container.Resources.Requests.Cpu()
			requestMemory := container.Resources.Requests.Memory()
			limitCpu := container.Resources.Limits.Cpu()
			limitMemory := container.Resources.Limits.Memory()
			containerResource := &ContainerResources{
				PodName:   pod.Name,
				Name:      container.Name,
				Namespace: pod.Namespace,
				NodeName:  pod.Spec.NodeName,
				Image:     container.Image,
				Request: Resources{
					Cpu:    0,
					Memory: 0,
				},
				Limit: Resources{
					Cpu:    0,
					Memory: 0,
				},
			}
			if requestCpu != nil {
				containerResource.Request.Cpu = float64(requestCpu.MilliValue())
			}

			if requestMemory != nil {
				containerResource.Request.Memory = float64(requestMemory.Value())
			}
			if limitCpu != nil {
				containerResource.Limit.Cpu = float64(limitCpu.MilliValue())
			}
			if limitMemory != nil {
				containerResource.Limit.Memory = float64(limitMemory.Value())
			}

			containers = append(containers, containerResource)
		}
	}
	return containers, nil
}

func (f *Fetcher) getNodes() ([]string, error) {
	if len(f.nodes) > 0 {
		return f.nodes, nil
	}

	nodes, err := f.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	f.nodes = make([]string, len(nodes.Items))
	for i, node := range nodes.Items {
		f.nodes[i] = node.Name
	}
	return f.nodes, nil
}

func (f *Fetcher) fetchMetricsFromNode(node string) (*NodeMetrics, error) {
	fetchTime := time.Now()
	path := fmt.Sprintf(CADVISOR_PATH_TEMPLATE, node)
	b, err := f.clientset.RESTClient().Get().AbsPath(path).Do(context.Background()).Raw()
	if err != nil {
		return nil, err
	}

	cpu, memory, err := f.metricsParser.Parse(b)
	if err != nil {
		return nil, err
	}

	return &NodeMetrics{
		NodeName:  node,
		Cpu:       cpu,
		Memory:    memory,
		Timestamp: fetchTime,
	}, nil
}
