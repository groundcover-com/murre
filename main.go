package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
	_ = v.ReadInConfig()
	conf, err := newConfig(v, os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", conf.kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	fetcher := NewFetcher(clientset)
	table := CreateNewTable()
	toplite := NewMurre(fetcher, table, time.Second*10)
	go table.Draw()
	toplite.Run()
}
