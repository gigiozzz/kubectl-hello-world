package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/klog/v2"
)

// scheme is needed for the TypeSetter
var scheme = runtime.NewScheme()

func init() {
	// Register core v1 types with the scheme
	corev1.AddToScheme(scheme)
}

// ListPodsOptions holds the configuration for the list-pods command
type ListPodsOptions struct {
	*CommonOptions
	PrintFlags *genericclioptions.PrintFlags
}

func NewListPodsOptions(streams genericclioptions.IOStreams) *ListPodsOptions {
	return &ListPodsOptions{
		CommonOptions: NewCommonOptions(streams),
		PrintFlags:    genericclioptions.NewPrintFlags("").WithTypeSetter(scheme),
	}
}

// Complete sets up the list-pods options based on the command arguments
func (o *ListPodsOptions) Complete(cmd *cobra.Command, args []string) error {
	// Complete common options
	return o.CommonOptions.Complete()
}

// Validate checks that the list-pods options are valid
func (o *ListPodsOptions) Validate() error {
	// Validate common options
	return o.CommonOptions.Validate()
}

// Run executes the list-pods command logic
func (o *ListPodsOptions) Run() error {
	ctx := context.Background()

	// Get pods from the current namespace
	pods, err := o.Clientset.CoreV1().Pods(o.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}

	if len(pods.Items) == 0 {
		fmt.Fprintf(o.Out, "No pods found in namespace '%s'\n", o.Namespace)
		return nil
	}

	// Get the output format
	outputFormat := ""
	if o.PrintFlags.OutputFormat != nil {
		outputFormat = *o.PrintFlags.OutputFormat
	}

	// Print pods using switch statement
	return o.printPods(pods, outputFormat)
}

// printPods handles different output formats for pod listing
func (o *ListPodsOptions) printPods(pods *corev1.PodList, format string) error {
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
		// NamePrinter needs individual objects, not lists
		for _, pod := range pods.Items {
			if err := namePrinter.PrintObj(&pod, o.Out); err != nil {
				return fmt.Errorf("failed to print pod name: %w", err)
			}
		}
		return nil

	case "wide":
		return o.printPodsWide(pods)

	default:
		// Default table format
		tablePrinter := printers.NewTablePrinter(printers.PrintOptions{})
		return tablePrinter.PrintObj(pods, o.Out)
	}
}

// printPodsWide prints pods in wide format (similar to kubectl get pods -o wide)
func (o *ListPodsOptions) printPodsWide(pods *corev1.PodList) error {
	// Print header with improved column spacing
	fmt.Fprintf(o.Out, "%-52s %-8s %-8s %-8s %-12s %-15s %-20s %-16s %-16s\n",
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

		// Print the row with improved column spacing
		fmt.Fprintf(o.Out, "%-52s %-8s %-8s %-8d %-12s %-15s %-20s %-16s %-16s\n",
			truncateString(pod.Name, 52),
			readyStr,
			string(pod.Status.Phase),
			restarts,
			age,
			podIP,
			truncateString(nodeName, 20),
			truncateString(nominatedNode, 16),
			truncateString(readinessGates, 16),
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

// NewCmdListPods creates the list-pods subcommand
func NewCmdListPods(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewListPodsOptions(streams)

	cmd := &cobra.Command{
		Use:   "list-pods",
		Short: "List pods with various output formats",
		Example: `  # Default table output
  kubectl hello-world list-pods
  
  # Wide output (like kubectl get pods -o wide)
  kubectl hello-world list-pods -o wide
  
  # JSON output
  kubectl hello-world list-pods -o json
  
  # YAML output  
  kubectl hello-world list-pods -o yaml
  
  # Name only
  kubectl hello-world list-pods -o name
  
  # Different namespace
  kubectl hello-world list-pods --namespace kube-system`,
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

	// Add print flags (-o, --output, etc.)
	o.PrintFlags.AddFlags(cmd)

	return cmd
}
