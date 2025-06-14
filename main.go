package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// PluginOptions holds the configuration for our plugin
type PluginOptions struct {
	configFlags *genericclioptions.ConfigFlags
	printFlags  *genericclioptions.PrintFlags
	clientset   kubernetes.Interface
	namespace   string
	context     string
	clusterName string

	// Plugin-specific flags
	name        string
	verbose     bool
	skipPodlist bool

	genericclioptions.IOStreams
}

var scheme = runtime.NewScheme()

func init() {
	// Register core v1 types with the scheme
	corev1.AddToScheme(scheme)
}

// NewPluginOptions creates a new PluginOptions with default values
func NewPluginOptions(streams genericclioptions.IOStreams) *PluginOptions {
	return &PluginOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		printFlags:  genericclioptions.NewPrintFlags("").WithTypeSetter(scheme),
		IOStreams:   streams,
	}
}

// Complete sets up the plugin options based on the command arguments
func (o *PluginOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	// Build kubernetes client
	config, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("failed to create REST config: %w", err)
	}

	o.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Get namespace
	o.namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return fmt.Errorf("failed to get namespace: %w", err)
	}

	// Get namespace, context, and kubeconfig using cli-runtime
	rawConfig, err := o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	// Get current context
	o.context = rawConfig.CurrentContext
	if o.context == "" {
		return fmt.Errorf("no current context set in kubeconfig")
	}

	// Get cluster name
	if o.clusterName == "" {
		currentContext := rawConfig.Contexts[o.context]
		if currentContext != nil {
			o.clusterName = currentContext.Cluster
		}
	}

	return nil
}

// Validate checks that the plugin options are valid
func (o *PluginOptions) Validate() error {
	if o.clientset == nil {
		return fmt.Errorf("kubernetes client not initialized")
	}
	// Try to connect to the cluster
	ctx := context.Background()
	namespace, err := o.clientset.CoreV1().Namespaces().Get(ctx, o.namespace, metav1.GetOptions{})
	if err != nil {
		if o.verbose {
			klog.ErrorS(err, "Failed to get namespace")
		}
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}
	fmt.Fprintf(o.Out, "✅ Connected to cluster to namespace '%s'\n", namespace.Name)
	fmt.Fprintf(o.Out, "✅ Context to use '%s'\n", o.context)
	fmt.Fprintf(o.Out, "✅ Cluster name to use '%s'\n", o.clusterName)

	return nil
}

// Run executes the plugin logic
func (o *PluginOptions) Run() error {
	if o.verbose {
		klog.InfoS("Starting hello plugin", "name", o.name, "namespace", o.namespace)
	}

	fmt.Fprintf(o.Out, "Hello, %s👋\n", o.name)

	if !o.skipPodlist {
		// Get some pods to demonstrate different printing options
		ctx := context.Background()
		pods, err := o.clientset.CoreV1().Pods(o.namespace).List(ctx, metav1.ListOptions{
			Limit: 5, // Limit to 5 pods for demo
		})
		if err != nil {
			return fmt.Errorf("failed to list pods: %w", err)
		}

		// Get the output format
		outputFormat := ""
		if o.printFlags.OutputFormat != nil {
			outputFormat = *o.printFlags.OutputFormat
		}

		// Print using switch statement
		return o.printPods(pods, outputFormat)
		/*
			// 1. Create a printer based on output format flags
			printer, err := o.printFlags.ToPrinter()
			if err != nil {
				return fmt.Errorf("failed to create printer: %w", err)
			}

			fmt.Fprintf(o.Out, "=== Printing with format: %s ===\n", *o.printFlags.OutputFormat)
			if err := printer.PrintObj(pods, o.Out); err != nil {
				return fmt.Errorf("failed to print pods: %w", err)
			}

			// 2. Table printer (like kubectl get)
			tablePrinter := printers.NewTablePrinter(printers.PrintOptions{})
			fmt.Fprintf(o.Out, "\n--- Table format ---\n")
			if err := tablePrinter.PrintObj(pods, o.Out); err != nil {
				return fmt.Errorf("failed to print table: %w", err)
			}

			// 6. Custom table printer with specific columns
			fmt.Fprintf(o.Out, "\n--- Custom table format ---\n")
			customTablePrinter := printers.NewTablePrinter(printers.PrintOptions{
				NoHeaders: true,
				Wide:      true,
			})
			if err := customTablePrinter.PrintObj(pods, o.Out); err != nil {
				return fmt.Errorf("failed to print custom table: %w", err)
			}

		*/

	}
	return nil
}

// NewCmdHello creates the cobra command
func NewCmdHello(streams genericclioptions.IOStreams) *cobra.Command {
	options := NewPluginOptions(streams)
	cmd := &cobra.Command{
		Use:     "hello-world",
		Short:   "A hello world kubectl plugin",
		Version: "0.1.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(cmd, args); err != nil {
				return err
			}
			if err := options.Validate(); err != nil {
				return err
			}
			return options.Run()
		},
	}

	// Add kubectl flags to my cmd (--kubeconfig, --context, --namespace, etc.)
	options.configFlags.AddFlags(cmd.Flags())
	// Add print flags (-o, --output, --no-headers, etc.)
	options.printFlags.AddFlags(cmd)

	// Add my flags
	cmd.Flags().StringVar(
		&options.name,
		"name",
		"",
		"Name to greet")
	cmd.MarkFlagRequired("name")

	cmd.Flags().StringVar(
		&options.clusterName,
		"cluster-name",
		"",
		"Cluster name to use (default from current context)")
	cmd.Flags().BoolVar(&options.verbose, "verbose", false, "Enable verbose output")
	cmd.Flags().BoolVar(&options.skipPodlist, "skip-pod-list", false, "Disable pod list output")

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
		// klog.ErrorS(err, "Plugin execution failed")
		os.Exit(1)
	}

}

func (o *PluginOptions) printPods(pods *corev1.PodList, format string) error {
	pods.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{Kind: "List"})
	klog.Info("GetObjectKind:", pods.GetObjectKind(), "GroupVersionKind:", pods.GetObjectKind().GroupVersionKind())
	switch format {
	case "json":
		jsonPrinter := &printers.JSONPrinter{}
		return jsonPrinter.PrintObj(pods, o.Out)

	case "yaml":
		yamlPrinter := &printers.YAMLPrinter{}
		return yamlPrinter.PrintObj(pods, o.Out)

	case "name":
		namePrinter := printers.NamePrinter{
			Operation: "listed",
		}
		return namePrinter.PrintObj(pods, o.Out)

	case "wide":
		return o.printPodsWide(pods)

	default:
		// Default table format - let cli-runtime handle wide vs normal
		// The PrintFlags will automatically handle -o wide
		tablePrinter := printers.NewTablePrinter(printers.PrintOptions{
			Wide: format == "wide",
		})
		return tablePrinter.PrintObj(pods, o.Out)
	}
}

func (o *PluginOptions) printPodsWide(pods *corev1.PodList) error {
	// Print header similar to kubectl get pods -o wide
	fmt.Fprintf(o.Out, "%-30s %-8s %-8s %-8s %-12s %-15s %-15s %-12s %s\n",
		"NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP", "NODE", "NOMINATED NODE", "READINESS GATES")

	for _, pod := range pods.Items {
		// Calculate ready containers
		ready := 0
		total := len(pod.Spec.Containers)
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Ready {
				ready++
			}
		}
		readyStr := fmt.Sprintf("%d/%d", ready, total)

		// Get restart count
		restarts := int32(0)
		for _, containerStatus := range pod.Status.ContainerStatuses {
			restarts += containerStatus.RestartCount
		}

		// Calculate age
		age := ""
		if !pod.CreationTimestamp.IsZero() {
			age = translateTimestampSince(pod.CreationTimestamp)
		}

		// Get pod IP
		podIP := pod.Status.PodIP
		if podIP == "" {
			podIP = "<none>"
		}

		// Get node name
		nodeName := pod.Spec.NodeName
		if nodeName == "" {
			nodeName = "<none>"
		}

		// Get nominated node (for scheduling)
		nominatedNode := "<none>"
		if pod.Status.NominatedNodeName != "" {
			nominatedNode = pod.Status.NominatedNodeName
		}

		// Get readiness gates info
		readinessGates := "<none>"
		if len(pod.Spec.ReadinessGates) > 0 {
			readinessGates = fmt.Sprintf("%d", len(pod.Spec.ReadinessGates))
		}

		// Print the row
		fmt.Fprintf(o.Out, "%-30s %-8s %-8s %-8d %-12s %-15s %-15s %-12s %s\n",
			truncateString(pod.Name, 30),
			readyStr,
			string(pod.Status.Phase),
			restarts,
			age,
			podIP,
			truncateString(nodeName, 15),
			truncateString(nominatedNode, 12),
			readinessGates,
		)
	}

	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// Helper function to calculate age (simplified version)
func translateTimestampSince(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	duration := metav1.Now().Sub(timestamp.Time)

	if duration.Hours() >= 24 {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	} else if duration.Hours() >= 1 {
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh", hours)
	} else if duration.Minutes() >= 1 {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%dm", minutes)
	} else {
		seconds := int(duration.Seconds())
		return fmt.Sprintf("%ds", seconds)
	}
}
