package cmd

import (
	"path/filepath"

	murre "github.com/groundcover-com/murre/pkg"
	"github.com/groundcover-com/murre/pkg/config"
	"github.com/groundcover-com/murre/pkg/ui"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

var (
	murreConfig *config.Config
)

func init() {
	initMurreFlags()
}

var RootCmd = &cobra.Command{
	Use:               "murre",
	Short:             "murre is a command line tool to monitor kubernetes resources",
	SilenceErrors:     true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Long:              `murre is a command line tool to monitor kubernetes resources`,
	RunE:              run,
}

func Execute() error {
	return RootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	table := ui.CreateNewTable()
	murre, err := murre.NewMurre(table, murreConfig)
	if err != nil {
		return err
	}

	go murre.Run()

	table.Draw()
	murre.Stop()

	return nil
}

func initMurreFlags() {
	murreConfig = &config.Config{}
	RootCmd.Flags().DurationVar(
		&murreConfig.RefreshInterval,
		"interval",
		config.DefaultRefreshInterval,
		"seconds to wait between updates",
	)
	RootCmd.Flags().StringVar(
		&murreConfig.Filters.Namespace,
		"namespace",
		"",
		"filter by namespace",
	)
	RootCmd.Flags().StringVar(
		&murreConfig.Filters.Node,
		"node",
		"",
		"filter by node",
	)
	RootCmd.Flags().StringVar(
		&murreConfig.Filters.Pod,
		"pod",
		"",
		"filter by pod",
	)
	RootCmd.Flags().StringVar(
		&murreConfig.Filters.Container,
		"container",
		"",
		"filter by container",
	)
	RootCmd.Flags().BoolVar(
		&murreConfig.SortBy.Cpu,
		"sortby-cpu",
		false,
		"sort by cpu",
	)
	RootCmd.Flags().BoolVar(
		&murreConfig.SortBy.CpuUtilization,
		"sortby-cpu-utilization",
		false,
		"sort by cpu utilization",
	)
	RootCmd.Flags().BoolVar(
		&murreConfig.SortBy.Mem,
		"sortby-mem",
		false,
		"sort by memory",
	)
	RootCmd.Flags().BoolVar(
		&murreConfig.SortBy.MemUtilization,
		"sortby-mem-utilization",
		false,
		"sort by memory utilization",
	)
	RootCmd.Flags().BoolVar(
		&murreConfig.SortBy.PodName,
		"sortby-pod-name",
		false,
		"sort by pod name",
	)

	if home := homedir.HomeDir(); home != "" {
		RootCmd.Flags().StringVar(
			&murreConfig.Kubeconfig,
			"kubeconfig",
			filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file",
		)
	} else {
		RootCmd.Flags().StringVar(
			&murreConfig.Kubeconfig,
			"kubeconfig",
			"",
			"absolute path to the kubeconfig file",
		)
	}

	RootCmd.Flags()
}
