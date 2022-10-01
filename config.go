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

type filter struct {
	// filter by namespace
	namespace string
	// filter by pod
	pod string
	// filter by container
	container string
}

type sortBy struct {
	// sort by cpu
	cpu bool
	// sort by memory
	mem bool
}

type config struct {
	refreshInterval time.Duration
	filters         filter
	sortBy          sortBy
	kubeconfig      string
}

func newConfig(v *viper.Viper, args []string) (*config, error) {
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

	return &config{
		refreshInterval: v.GetDuration("interval"),
		filters: filter{
			namespace: v.GetString("namespace"),
			pod:       v.GetString("pod"),
			container: v.GetString("container"),
		},
		sortBy: sortBy{
			cpu: v.GetBool("cpu"),
			mem: v.GetBool("mem"),
		},
		kubeconfig: v.GetString("kubeconfig"),
	}, nil
}
