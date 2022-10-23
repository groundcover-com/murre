package main

import (
	"fmt"
	"sort"
	"time"

	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	FETCH_CONTAINERS_SPEC_RATIO = 5
)

type DataFetcher interface {
	GetMetrics() ([]*NodeMetrics, error)
	GetContainers() ([]*ContainerResources, error)
}

type UI interface {
	Update(stats []*Stats)
}

type ContainerStats struct {
	Namespace     string
	PodName       string
	ContainerName string
	CpuUsage      float64
	MemoryBytes   float64
	LastUpdateTs  time.Time
}

type Murre struct {
	fetcher      DataFetcher
	ui           UI
	config       *Config
	containers   map[string]*Container
	fetchCounter int
	stopCh       chan struct{}
}

func NewMurre(ui UI, config *Config) (*Murre, error) {
	// use the current context in kubeconfig
	kubecfg, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(kubecfg)
	if err != nil {
		return nil, err
	}

	fetcher := NewFetcher(clientset)
	if fetcher == nil {
		return nil, err
	}

	return &Murre{
		fetcher:      fetcher,
		ui:           ui,
		config:       config,
		containers:   make(map[string]*Container),
		stopCh:       make(chan struct{}),
		fetchCounter: 0,
	}, nil

}

func (m *Murre) Run() error {
	// first tick
	err := m.tick()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(m.config.RefreshInterval)

	for {
		select {
		case <-ticker.C:
			err := m.tick()
			if err != nil {
				return err
			}
		case <-m.stopCh:
			return nil
		}
	}
}

func (m *Murre) Stop() {
	close(m.stopCh)
}

func (m *Murre) tick() error {
	defer func() {
		m.fetchCounter++
	}()
	err := m.updateContainers()
	if err != nil {
		return err
	}

	err = m.updateMetrics()
	if err != nil {
		return err
	}
	stats := m.getStats()
	stats = m.filter(stats)
	m.sort(stats)
	m.ui.Update(stats)
	return nil
}

func (m *Murre) updateContainers() error {
	if m.fetchCounter%FETCH_CONTAINERS_SPEC_RATIO != 0 {
		return nil
	}

	containers, err := m.fetcher.GetContainers()
	if err != nil {
		return err
	}

	for _, c := range containers {
		container := m.getOrCreateContainer(c.Name, c.Image, c.PodName, c.Namespace)
		container.UpdateResources(c)
	}

	return nil
}

func (m *Murre) filter(stats []*Stats) []*Stats {
	filterdStats := make([]*Stats, 0)
	for _, s := range stats {
		isNamespaceMatch := m.config.Filters.Namespace == "" || m.config.Filters.Namespace == s.Namespace
		isPodMatch := m.config.Filters.Pod == "" || m.config.Filters.Pod == s.PodName
		isContainerMatch := m.config.Filters.Container == "" || m.config.Filters.Container == s.ContainerName
		if isNamespaceMatch && isPodMatch && isContainerMatch {
			filterdStats = append(filterdStats, s)
		}
	}
	return filterdStats
}

func (m *Murre) sort(stats []*Stats) {
	if m.config.SortBy.Mem {
		sort.Slice(stats, func(i, j int) bool {
			return stats[i].MemoryBytes > stats[j].MemoryBytes
		})
		return
	}

	if m.config.SortBy.Cpu {
		sort.Slice(stats, func(i, j int) bool {
			return stats[i].CpuUsageMilli > stats[j].CpuUsageMilli
		})
		return
	}

	if m.config.SortBy.CpuUtilization {
		sort.Slice(stats, func(i, j int) bool {
			return stats[i].CpuUsagePercent > stats[j].CpuUsagePercent
		})
		return
	}

	if m.config.SortBy.MemUtilization {
		sort.Slice(stats, func(i, j int) bool {
			return stats[i].MemoryUsagePercent > stats[j].MemoryUsagePercent
		})
		return
	}

	//default is to sort by cpu
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].CpuUsageMilli > stats[j].CpuUsageMilli
	})
}

func (m *Murre) getStats() []*Stats {
	containersStats := make([]*Stats, 0)
	for _, c := range m.containers {
		stats := c.GetStats()
		if stats == nil {
			continue
		}

		if time.Since(stats.LastUpdateTs) > 2*time.Minute {
			delete(m.containers, c.Id)
			continue
		}

		containersStats = append(containersStats, stats)
	}
	return containersStats
}

func (m *Murre) updateMetrics() error {
	metrics, err := m.fetcher.GetMetrics()
	if err != nil {
		return err
	}
	for _, node := range metrics {
		m.updateCpu(node.Cpu, node.Timestamp)
		m.updateMemory(node.Memory, node.Timestamp)
	}
	return nil
}

func (m *Murre) updateCpu(cpu []*Cpu, fetchTime time.Time) {
	for _, c := range cpu {
		container := m.getOrCreateContainerFromCpu(c)
		container.UpdateCpu(c, fetchTime)
	}
}

func (m *Murre) updateMemory(memory []*Memory, fetchTime time.Time) {
	for _, mem := range memory {
		container := m.getOrCreateContainerFromMemory(mem)
		container.UpdateMemory(mem, fetchTime)
	}
}

func (m *Murre) getOrCreateContainerFromCpu(cpu *Cpu) *Container {
	return m.getOrCreateContainer(cpu.Name, cpu.Image, cpu.PodName, cpu.Namespace)
}

func (m *Murre) getOrCreateContainerFromMemory(mem *Memory) *Container {
	return m.getOrCreateContainer(mem.Name, mem.Image, mem.PodName, mem.Namespace)
}

func (m *Murre) getOrCreateContainer(name, image, podName, namespace string) *Container {
	id := fmt.Sprintf("%s/%s/%s", namespace, podName, name)
	if _, ok := m.containers[id]; !ok {
		m.containers[id] = &Container{
			Id:        id,
			Name:      name,
			Image:     image,
			PodName:   podName,
			Namespace: namespace,
		}
	}

	return m.containers[id]
}
