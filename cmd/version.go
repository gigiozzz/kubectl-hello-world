package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// Build information set by ldflags
var (
	Version    = "dev"
	CommitHash = "unknown"
	BuildDate  = "unknown"
)

// VersionOptions holds the configuration for the version command
type VersionOptions struct {
	Short bool

	genericclioptions.IOStreams
}

// NewVersionOptions creates a new VersionOptions with default values
func NewVersionOptions(streams genericclioptions.IOStreams) *VersionOptions {
	return &VersionOptions{
		IOStreams: streams,
	}
}

// Run executes the version command logic
func (o *VersionOptions) Run() error {
	if o.Short {
		fmt.Fprintf(o.Out, "%s\n", Version)
		return nil
	}

	fmt.Fprintf(o.Out, "kubectl-hello-world version: %s\n", Version)
	fmt.Fprintf(o.Out, "Commit: %s\n", CommitHash)
	fmt.Fprintf(o.Out, "Built: %s\n", BuildDate)
	fmt.Fprintf(o.Out, "Go version: %s\n", runtime.Version())
	fmt.Fprintf(o.Out, "OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	return nil
}

// NewCmdVersion creates the version subcommand
func NewCmdVersion(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewVersionOptions(streams)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Example: `  # Show full version information
  kubectl hello-world version
  
  # Show only version number
  kubectl hello-world version --short`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	// Add version-specific flags
	cmd.Flags().BoolVar(&o.Short, "short", false, "Print only the version number")

	return cmd
}
