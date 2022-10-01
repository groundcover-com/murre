package main

import "time"

type Container struct {
	Id                         string
	Name                       string
	Image                      string
	PodName                    string
	Namespace                  string
	cpuUsage                   float64
	lastCpuUsageSecondsTotal   float64
	lastCpuUsageSecondsTotalTs time.Time
	memoryUsageBytes           float64
}

type Stats struct {
	Namespace     string
	PodName       string
	ContainerName string
	CpuUsage      float64
	MemoryBytes   float64
	LastUpdateTs  time.Time
}

func (c *Container) GetStats() *Stats {
	if c.cpuUsage == 0 {
		return nil
	}

	return &Stats{
		Namespace:     c.Namespace,
		PodName:       c.PodName,
		ContainerName: c.Name,
		CpuUsage:      c.cpuUsage,
		MemoryBytes:   c.memoryUsageBytes,
		LastUpdateTs:  c.lastCpuUsageSecondsTotalTs,
	}
}

func (c *Container) UpdateCpu(cpu *Cpu, fetchTime time.Time) {
	if c.lastCpuUsageSecondsTotal == cpu.CpuUsageSecondsTotal {
		return
	}

	if !c.lastCpuUsageSecondsTotalTs.IsZero() {
		increaseInCpu := cpu.CpuUsageSecondsTotal - c.lastCpuUsageSecondsTotal
		timeDiff := fetchTime.Sub(c.lastCpuUsageSecondsTotalTs)
		c.cpuUsage = (increaseInCpu) / timeDiff.Seconds()
	}

	c.lastCpuUsageSecondsTotal = cpu.CpuUsageSecondsTotal
	c.lastCpuUsageSecondsTotalTs = fetchTime
}

func (c *Container) UpdateMemory(memory *Memory, fetchTime time.Time) {
	c.memoryUsageBytes = memory.MemoryUsageBytes
}
