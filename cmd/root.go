package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCmdHello creates the cobra command
func NewCmdHello(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hello-world",
		Short: "A hello world kubectl plugin",
		Long: `kubectl-hello-world is a demonstration plugin that shows:
- How to build kubectl plugins with Cobra
- Integration with kubectl cli-runtime library
- Multiple output formats
- Subcommand architecture
- Best practices for kubectl plugin development

This plugin demonstrates proper project structure with:
- Shared utilities for common kubectl operations
- Modular command organization
- Consistent error handling and logging
- Professional CLI patterns`,
		Example: `  # Greet someone
  kubectl hello-world greetings --name Alice
  
  # List pods in current namespace
  kubectl hello-world list-pods
  
  # List pods with wide output
  kubectl hello-world list-pods -o wide
  
  # List pods in specific namespace
  kubectl hello-world list-pods --namespace kube-system
  
  # Use different context
  kubectl hello-world greetings --context=production Bob`,
	}

	cmd.AddCommand(NewCmdGreetings(streams))
	cmd.AddCommand(NewCmdListPods(streams))
	cmd.AddCommand(NewCmdVersion(streams))

	return cmd
}
