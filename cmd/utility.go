package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// CommonOptions holds shared configuration used by multiple commands
type CommonOptions struct {
	ConfigFlags *genericclioptions.ConfigFlags
	Clientset   kubernetes.Interface
	Namespace   string
	Context     string
	Kubeconfig  string
	ClusterName string

	genericclioptions.IOStreams
}

// NewCommonOptions creates a new CommonOptions with default values
func NewCommonOptions(streams genericclioptions.IOStreams) *CommonOptions {
	return &CommonOptions{
		ConfigFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

// Complete initializes the common kubernetes client and configuration
func (o *CommonOptions) Complete() error {
	var err error

	// Build kubernetes client
	config, err := o.ConfigFlags.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("failed to create REST config: %w", err)
	}

	o.Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Get namespace, context, and kubeconfig using cli-runtime
	rawConfig, err := o.ConfigFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	// Get current context
	o.Context = rawConfig.CurrentContext
	if o.Context == "" {
		return fmt.Errorf("no current context set in kubeconfig")
	}

	// Get namespace from context
	o.Namespace, _, err = o.ConfigFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return fmt.Errorf("failed to get namespace: %w", err)
	}

	// Get cluster name from current context
	currentContext := rawConfig.Contexts[o.Context]
	if currentContext != nil {
		o.ClusterName = currentContext.Cluster
	}

	// Get kubeconfig file path
	if kubeconfig := o.ConfigFlags.KubeConfig; kubeconfig != nil && *kubeconfig != "" {
		o.Kubeconfig = *kubeconfig
	} else {
		o.Kubeconfig = "default kubeconfig location"
	}

	return nil
}

// Validate performs common validation checks
func (o *CommonOptions) Validate() error {
	if o.Clientset == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}

	if o.Namespace == "" {
		return fmt.Errorf("namespace not determined")
	}

	// Test connection by retrieving the current namespace
	ctx := context.Background()
	_, err := o.Clientset.CoreV1().Namespaces().Get(ctx, o.Namespace, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to access namespace '%s': %w", o.Namespace, err)
	}

	return nil
}

// LogContext logs the current kubernetes context information
func (o *CommonOptions) LogContext(verbose bool) {
	if verbose {
		klog.InfoS("Context information",
			"namespace", o.Namespace,
			"context", o.Context,
			"cluster", o.ClusterName,
			"kubeconfig", o.Kubeconfig)
	}
}

// AddConfigFlags adds the common kubectl configuration flags to a command
func (o *CommonOptions) AddConfigFlags(cmd *cobra.Command) {
	o.ConfigFlags.AddFlags(cmd.Flags())
}
