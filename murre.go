package main

import (
	"time"
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
	interval   time.Duration
	containers map[string]*Container
}

func NewMurre(fetcher DataFetcher, ui UI, interval time.Duration) *Murre {
	if fetcher == nil {
		panic("fetcher is nil")
	}

	if interval <= 0 {
		panic("interval is invalid")
	}

	return &Murre{
		fetcher:    fetcher,
		ui:         ui,
		interval:   interval,
		containers: make(map[string]*Container),
	}

}

func (m *Murre) Run() error {
	ticker := time.NewTicker(m.interval)

	for range ticker.C {
		err := m.updateMetrics()
		if err != nil {
			return err
		}
		stats := m.getStats()
		//stats = m.sort(stats)
		m.ui.Update(stats)
	}

	return nil
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
