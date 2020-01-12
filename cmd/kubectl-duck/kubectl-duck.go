package main

import (
	"flag"
	"github.com/n3wscott/kubectl-duck/pkg/commands"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

func main() {
	defer klog.Flush()
	cmds := &cobra.Command{
		Use:   "kubectl-duck",
		Short: "Ducktype support for kubectl.",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	commands.AddDuckCommands(cmds)

	if err := cmds.Execute(); err != nil {
		klog.Fatalf("error during command execution: %v", err)
	}
}

func init() {
	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// hide all glog flags except for -v
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		if f.Name != "v" {
			pflag.Lookup(f.Name).Hidden = true
		}
	})
}
