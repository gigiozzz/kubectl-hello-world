package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog/v2"
)

// GreetingsOptions holds the configuration for the greetings command
type GreetingsOptions struct {
	*CommonOptions

	// Plugin-specific flags
	name    string
	verbose bool
}

// NewGreetingsOptions creates a new GreetingsOptions with default values
func NewGreetingsOptions(streams genericclioptions.IOStreams) *GreetingsOptions {
	return &GreetingsOptions{
		CommonOptions: NewCommonOptions(streams),
	}
}

// Complete sets up the plugin options based on the command arguments
func (o *GreetingsOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	// Complete common options first
	if err = o.CommonOptions.Complete(); err != nil {
		return err
	}

	// Handle positional arguments
	if len(args) > 0 {
		o.name = args[0]
	} else if o.name == "" {
		o.name = "World"
	}

	return nil
}

// Validate checks that the greetings options are valid
func (o *GreetingsOptions) Validate() error {
	// Validate common options
	if err := o.CommonOptions.Validate(); err != nil {
		return err
	}

	// Log context if verbose
	o.CommonOptions.LogContext(o.verbose)

	return nil
}

// Run executes the greetings command logic
func (o *GreetingsOptions) Run() error {
	if o.verbose {
		klog.InfoS("Starting greetings command", "name", o.name, "namespace", o.Namespace)
	}

	// Print greeting
	fmt.Fprintf(o.Out, "Hello, %s! 👋\n", o.name)

	return nil
}

// NewCmdGreetings creates the greetings subcommand
func NewCmdGreetings(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewGreetingsOptions(streams)

	cmd := &cobra.Command{
		Use:   "greetings [name]",
		Short: "Greet someone and show cluster information",
		Example: `  # Say hello to World
  kubectl hello-world greetings
  
  # Say hello to a specific name
  kubectl hello-world greetings Alice
  
  # Use flag instead of positional argument
  kubectl hello-world greetings --name Bob
  
  # Use verbose output
  kubectl hello-world greetings --verbose Alice
  
  # Use specific context
  kubectl hello-world greetings --context=my-context Alice`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(cmd, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run()
		},
	}

	// Add common kubectl config flags
	o.CommonOptions.AddConfigFlags(cmd)

	// Add my flags
	cmd.Flags().StringVar(&o.name, "name", "", "Name to greet (can also be provided as argument)")
	//cmd.MarkFlagRequired("name")
	cmd.Flags().BoolVar(&o.verbose, "verbose", false, "Enable verbose output")

	return cmd
}
