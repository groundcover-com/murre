package k8s

import (
	"bytes"

	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

const (
	CONTAINER_CPU_METRICS = "container_cpu_user_seconds_total"
	CONTAINER_MEM_METRICS = "container_memory_usage_bytes"
)

const (
	METRIC_POD_LABEL        = "pod"
	METRIC_CONTAINER_LABEL  = "container"
	METRIC_NAME_LABEL       = "name"
	METRICS_NAMESPACE_LABEL = "namespace"
	METRICS_ID_LABEL        = "id"
	METRICS_IMAGE_LABEL     = "image"
	SORT_BY_CPU             = 0
	SORT_BY_MEM             = 1
	SORT_BY_POD             = 2
	SORT_BY_NAMESPACE       = 3
)

type Cpu struct {
	Name                 string
	Image                string
	PodName              string
	Namespace            string
	CpuUsageSecondsTotal float64
}

type Memory struct {
	Name             string
	Image            string
	PodName          string
	Namespace        string
	MemoryUsageBytes float64
}

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(b []byte, emptyContainer bool) ([]*Cpu, []*Memory, error) {
	reader := bytes.NewReader(b)

	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		panic(err.Error())
	}
	cpuMetrics := make([]*Cpu, 0)
	memoryMetrics := make([]*Memory, 0)

	for k, v := range mf {
		if k == CONTAINER_CPU_METRICS {
			metrics := v.GetMetric()
			if len(metrics) == 0 {
				panic(0)
			}
			cpuMetrics = append(cpuMetrics, p.parseCpuMetrics(metrics, emptyContainer)...)
		}
		if k == CONTAINER_MEM_METRICS {
			metrics := v.GetMetric()
			if len(metrics) == 0 {
				panic(0)
			}
			memoryMetrics = append(memoryMetrics, p.parseMemoryMetrics(metrics)...)
		}
	}
	return cpuMetrics, memoryMetrics, nil
}

func (p *Parser) parseCpuMetrics(metrics []*io_prometheus_client.Metric, emptyContainer bool) []*Cpu {
	cpuMetrics := make([]*Cpu, 0, len(metrics))
	for _, metric := range metrics {
		labels := metric.GetLabel()
		if len(labels) == 0 {
			panic(0)
		}
		cpuMetric := &Cpu{}
		for _, label := range labels {
			switch label.GetName() {
			case METRIC_POD_LABEL:
				cpuMetric.PodName = label.GetValue()
			case METRIC_CONTAINER_LABEL:
				cpuMetric.Name = label.GetValue()
			case METRIC_NAME_LABEL:
				continue
			case METRICS_NAMESPACE_LABEL:
				cpuMetric.Namespace = label.GetValue()
			case METRICS_IMAGE_LABEL:
				cpuMetric.Image = label.GetValue()
			case METRICS_ID_LABEL:
				continue
			default:
				panic(label.GetName())
			}
		}

		cpuMetric.CpuUsageSecondsTotal = metric.GetCounter().GetValue()
		if (!emptyContainer && cpuMetric.Name == "")  || cpuMetric.PodName == "" || cpuMetric.Namespace == "" {
			//done - I know why this happens
			//see https://github.com/google/cadvisor/issues/1873
			continue
		}
		cpuMetrics = append(cpuMetrics, cpuMetric)
	}
	return cpuMetrics
}

func (p *Parser) parseMemoryMetrics(metrics []*io_prometheus_client.Metric) []*Memory {
	memoryMetrics := make([]*Memory, 0, len(metrics))
	for _, metric := range metrics {
		labels := metric.GetLabel()
		if len(labels) == 0 {
			panic(0)
		}
		memoryMetric := &Memory{}
		for _, label := range labels {
			switch label.GetName() {
			case METRIC_POD_LABEL:
				memoryMetric.PodName = label.GetValue()
			case METRIC_CONTAINER_LABEL:
				memoryMetric.Name = label.GetValue()
			case METRIC_NAME_LABEL:
				continue
			case METRICS_NAMESPACE_LABEL:
				memoryMetric.Namespace = label.GetValue()
			case METRICS_IMAGE_LABEL:
				memoryMetric.Image = label.GetValue()
			case METRICS_ID_LABEL:
				continue
			default:
				panic(label.GetName())
			}
		}
		memoryMetric.MemoryUsageBytes = metric.GetGauge().GetValue()
		memoryMetrics = append(memoryMetrics, memoryMetric)
	}
	return memoryMetrics
}
