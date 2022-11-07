package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

const (
	KUBECONFIG_ENV_NAME = "kubeconfig"
)

var (
	defaultRefreshInterval = time.Second * 5
)

type Filter struct {
	// filter by namespace
	Namespace string
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

func newConfig(v *viper.Viper, args []string) (*Config, error) {
	flagSet := pflag.NewFlagSet("murre", pflag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Println("Usage: murre [options]")
		flagSet.PrintDefaults()
	}

	flagSet.Duration("interval", defaultRefreshInterval, "seconds to wait between updates (default '5s')")
	flagSet.String("namespace", "", "filter by namespace")
	flagSet.String("pod", "", "filter by pod")
	flagSet.String("container", "", "filter by container")
	flagSet.Bool("sortby-cpu", false, "sort by cpu")
	flagSet.Bool("sortby-cpu-util", false, "sort by cpu utilization")
	flagSet.Bool("sortby-mem", false, "sort by memory")
	flagSet.Bool("sortby-mem-util", false, "sort by memory utilization")
	flagSet.Bool("sortby-pod-name", false, "sort by pod name")
	flagSet.Bool("help", false, "show help")
	if home := homedir.HomeDir(); home != "" {
		flagSet.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flagSet.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	v.BindEnv(KUBECONFIG_ENV_NAME)

	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	err = v.BindPFlags(flagSet)
	if err != nil {
		return nil, err
	}

	if v.GetBool("help") {
		flagSet.Usage()
		return nil, nil
	}

	return &Config{
		RefreshInterval: v.GetDuration("interval"),
		Filters: Filter{
			Namespace: v.GetString("namespace"),
			Pod:       v.GetString("pod"),
			Container: v.GetString("container"),
		},
		SortBy: SortBy{
			Cpu:            v.GetBool("sortby-cpu"),
			CpuUtilization: v.GetBool("sortby-cpu-util"),
			Mem:            v.GetBool("sortby-mem"),
			MemUtilization: v.GetBool("sortby-mem-util"),
			PodName:        v.GetBool("sortby-pod-name"),
		},
		Kubeconfig: v.GetString("kubeconfig"),
	}, nil
}
