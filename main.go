package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// PluginOptions holds the configuration for our plugin
type PluginOptions struct {
	configFlags *genericclioptions.ConfigFlags
	clientset   kubernetes.Interface
	namespace   string

	// Plugin-specific flags
	name    string
	verbose bool

	genericclioptions.IOStreams
}

// NewCmdHello creates the cobra command
func NewCmdHello(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{}
	return cmd
}

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

	cmd := NewCmdHello(streams)

	// Execute the command
	if err := cmd.Execute(); err != nil {
		// Cobra handles error printing, but we can add additional context
		klog.ErrorS(err, "Plugin execution failed")
		os.Exit(1)
	}
	fmt.Fprintf(streams.Out, "Hello World 👋\n")
}
