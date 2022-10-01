package main

import (
	"sort"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type DataFetcher interface {
	GetMetrics() ([]*NodeMetrics, error)
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
	// ...
	fetcher    DataFetcher
	ui         UI
	config     *Config
	containers map[string]*Container
	stopCh     chan struct{}
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
		fetcher:    fetcher,
		ui:         ui,
		config:     config,
		containers: make(map[string]*Container),
		stopCh:     make(chan struct{}),
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
	err := m.updateMetrics()
	if err != nil {
		return err
	}
	stats := m.getStats()
	stats = m.filter(stats)
	m.sort(stats)
	m.ui.Update(stats)
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
	}

	if m.config.SortBy.Cpu {
		sort.Slice(stats, func(i, j int) bool {
			return stats[i].CpuUsage > stats[j].CpuUsage
		})
	}

	//default is to sort by cpu
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].CpuUsage > stats[j].CpuUsage
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
	return m.getOrCreateContainer(cpu.Id, cpu.Name, cpu.Image, cpu.PodName, cpu.Namespace)
}

func (m *Murre) getOrCreateContainerFromMemory(mem *Memory) *Container {
	return m.getOrCreateContainer(mem.Id, mem.Name, mem.Image, mem.PodName, mem.Namespace)
}

func (m *Murre) getOrCreateContainer(id, name, image, podName, namespace string) *Container {
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
