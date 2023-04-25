package k8s

import (
	"time"
)

type Container struct {
	Id                         string
	Name                       string
	Image                      string
	NodeName                   string
	PodName                    string
	Namespace                  string
	cpuUsage                   float64
	lastCpuUsageSecondsTotal   float64
	lastCpuUsageSecondsTotalTs time.Time
	memoryUsageBytes           float64
	cpuRequest                 float64
	cpuLimits                  float64
	memoryRequestBytes         float64
	memoryLimitBytes           float64
}

type Stats struct {
	Namespace          string
	NodeName           string
	PodName            string
	ContainerName      string
	CpuUsageMilli      float64
	MemoryBytes        float64
	LastUpdateTs       time.Time
	MemoryLimitBytes   float64
	CpuLimit           float64
	MemoryUsagePercent float64
	CpuUsagePercent    float64
}

func (c *Container) GetStats() *Stats {
	if c.cpuUsage == 0 && c.memoryUsageBytes == 0 {
		return nil
	}

	cpuUsageInMillis := c.cpuUsage * 1000
	var cpuUsagePercent float64
	if c.cpuLimits > 0 {
		cpuUsagePercent = cpuUsageInMillis / c.cpuLimits * 100
	}

	if cpuUsagePercent > 100 {
		cpuUsagePercent = 100
	}

	var memoryUsagePercent float64
	if c.memoryLimitBytes > 0 {
		memoryUsagePercent = c.memoryUsageBytes / c.memoryLimitBytes * 100
	}

	if memoryUsagePercent > 100 {
		memoryUsagePercent = 100
	}

	return &Stats{
		Namespace:          c.Namespace,
		NodeName:           c.NodeName,
		PodName:            c.PodName,
		ContainerName:      c.Name,
		CpuUsageMilli:      cpuUsageInMillis,
		MemoryBytes:        c.memoryUsageBytes,
		LastUpdateTs:       c.lastCpuUsageSecondsTotalTs,
		CpuLimit:           c.cpuLimits,
		MemoryLimitBytes:   c.memoryLimitBytes,
		CpuUsagePercent:    cpuUsagePercent,
		MemoryUsagePercent: memoryUsagePercent,
	}
}

func (c *Container) UpdateCpu(cpu *Cpu, fetchTime time.Time) {
	if c.lastCpuUsageSecondsTotal == cpu.CpuUsageSecondsTotal {
		return
	}

	timeDiff := fetchTime.Sub(c.lastCpuUsageSecondsTotalTs)
	if !c.lastCpuUsageSecondsTotalTs.IsZero() && timeDiff > 0 {
		increaseInCpu := cpu.CpuUsageSecondsTotal - c.lastCpuUsageSecondsTotal
		c.cpuUsage = (increaseInCpu) / timeDiff.Seconds()
	}

	c.lastCpuUsageSecondsTotal = cpu.CpuUsageSecondsTotal
	c.lastCpuUsageSecondsTotalTs = fetchTime
}

func (c *Container) UpdateMemory(memory *Memory, fetchTime time.Time) {
	c.memoryUsageBytes = memory.MemoryUsageBytes
}

func (c *Container) UpdateResources(resources *ContainerResources) {
	c.cpuRequest = resources.Request.Cpu
	c.cpuLimits = resources.Limit.Cpu
	c.memoryRequestBytes = resources.Request.Memory
	c.memoryLimitBytes = resources.Limit.Memory
}
