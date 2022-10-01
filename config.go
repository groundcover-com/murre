package main

import (
	"path/filepath"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

var (
	defaultRefreshInterval = time.Second * 10
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
	// sort by memory
	Mem bool
}

type Config struct {
	RefreshInterval time.Duration
	Filters         Filter
	SortBy          SortBy
	Kubeconfig      string
}

func newConfig(v *viper.Viper, args []string) (*Config, error) {
	flagSet := pflag.NewFlagSet("murre", pflag.ExitOnError)

	flagSet.Duration("interval", defaultRefreshInterval, "seconds to wait between updates (default '10s')")
	flagSet.String("namespace", "", "filter by namespace")
	flagSet.String("pod", "", "filter by pod")
	flagSet.String("container", "", "filter by container")
	flagSet.Bool("sortby-cpu", false, "sort by cpu")
	flagSet.Bool("sortby-mem", false, "sort by memory")
	if home := homedir.HomeDir(); home != "" {
		flagSet.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		flagSet.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	err = v.BindPFlags(flagSet)
	if err != nil {
		return nil, err
	}

	return &Config{
		RefreshInterval: v.GetDuration("interval"),
		Filters: Filter{
			Namespace: v.GetString("namespace"),
			Pod:       v.GetString("pod"),
			Container: v.GetString("container"),
		},
		SortBy: SortBy{
			Cpu: v.GetBool("cpu"),
			Mem: v.GetBool("mem"),
		},
		Kubeconfig: v.GetString("kubeconfig"),
	}, nil
}
