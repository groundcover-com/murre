package config

import (
	"time"
)

const (
	KUBECONFIG_ENV_NAME = "kubeconfig"
)

var (
	DefaultRefreshInterval = time.Second * 5
)

type Filter struct {
	// filter by namespace
	Namespace string
	// filter by node
	Node string
	// filter by pod
	Pod string
	// filter by container
	Container string
}

type SortBy struct {
	// sort by cpu
	Cpu bool
	// sort by cpu utilization
	CpuUtilization bool
	// sort by memory
	Mem bool
	// sort by memory utilization
	MemUtilization bool
	// sort by pod name
	PodName bool
}

type Config struct {
	RefreshInterval time.Duration
	Filters         Filter
	SortBy          SortBy
	Kubeconfig      string
}
