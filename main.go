package main

import (
	"os"

	"github.com/spf13/viper"
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

func main() {
	v := viper.New()
	config, err := newConfig(v, os.Args[1:])
	if err != nil {
		panic(err)
	}

	table := CreateNewTable()
	murre, err := NewMurre(table, config)
	if err != nil {
		panic(err)
	}

	go murre.Run()

	table.Draw()
	murre.Stop()
}
