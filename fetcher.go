package main

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
