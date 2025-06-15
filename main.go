package main

import (
	"os"

	"github.com/gigiozzz/kubectl-hello-world/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog/v2"
)

func main() {
	// Initialize klog (kubectl's logging library)
	klog.InitFlags(nil)
	defer klog.Flush()

	// Create the root command
	streams := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	cmd := cmd.NewCmdHello(streams)

	// Execute the command
	if err := cmd.Execute(); err != nil {
		// Cobra handles error printing, but we can add additional context
		// klog.ErrorS(err, "Plugin execution failed")
		os.Exit(1)
	}

}
