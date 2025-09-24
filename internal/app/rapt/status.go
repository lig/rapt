package rapt

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"codeberg.org/lig/rapt/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	yamlv2 "sigs.k8s.io/yaml"
)

// StatusInfo represents the status information for display
type StatusInfo struct {
	ClusterConnected bool              `json:"cluster_connected"`
	CRDInstalled     bool              `json:"crd_installed"`
	CRDName          string            `json:"crd_name,omitempty"`
	CRDCreated       time.Time         `json:"crd_created,omitempty"`
	ToolsCount       int               `json:"tools_count"`
	ToolsByNamespace map[string]int    `json:"tools_by_namespace,omitempty"`
	CurrentNamespace string            `json:"current_namespace"`
	AllNamespaces    bool              `json:"all_namespaces"`
	Errors           []string          `json:"errors,omitempty"`
}

// ShowStatus displays the current status of Rapt installation and cluster health
func ShowStatus(namespace, outputFormat string, allNamespaces bool) error {
	status := StatusInfo{
		CurrentNamespace: namespace,
		AllNamespaces:    allNamespaces,
		ToolsByNamespace: make(map[string]int),
	}

	// Check cluster connectivity and CRD status
	err := checkClusterStatus(&status)
	if err != nil {
		status.Errors = append(status.Errors, fmt.Sprintf("Cluster connectivity: %v", err))
	}

	// Check tools count
	err = checkToolsStatus(&status, namespace, allNamespaces)
	if err != nil {
		status.Errors = append(status.Errors, fmt.Sprintf("Tools status: %v", err))
	}

	// Output based on format
	switch outputFormat {
	case "json":
		return outputStatusJSON(status)
	case "yaml":
		return outputStatusYAML(status)
	case "table":
		fallthrough
	default:
		return outputStatusTable(status)
	}
}

// checkClusterStatus checks cluster connectivity and CRD installation
func checkClusterStatus(status *StatusInfo) error {
	// Initialize clients
	k8sClient, err := k8s.InitClient(status.CurrentNamespace)
	if err != nil {
		status.ClusterConnected = false
		return err
	}
	status.ClusterConnected = true

	// Check if CRD exists
	crd, err := k8s.LoadToolCRD()
	if err != nil {
		status.CRDInstalled = false
		return fmt.Errorf("failed to load CRD definition: %w", err)
	}

	crdName := crd.GetName()
	existingCRD, err := k8sClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), crdName, metav1.GetOptions{})
	if err != nil {
		status.CRDInstalled = false
		return fmt.Errorf("CRD not found: %w", err)
	}

	status.CRDInstalled = true
	status.CRDName = crdName
	status.CRDCreated = existingCRD.GetCreationTimestamp().Time

	return nil
}

// checkToolsStatus checks the status of tools in the cluster
func checkToolsStatus(status *StatusInfo, namespace string, allNamespaces bool) error {
	dynClient, err := k8s.InitDynamicClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to initialize dynamic client: %w", err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "rapt.dev",
		Version:  "v1alpha1",
		Resource: "tools",
	}

	var tools interface{}
	if allNamespaces {
		tools, err = dynClient.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	} else {
		tools, err = dynClient.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	}

	if err != nil {
		// If CRD is not installed, this is expected
		if !status.CRDInstalled {
			status.ToolsCount = 0
			return nil
		}
		return fmt.Errorf("failed to list tools: %w", err)
	}

	// Count tools
	if toolList, ok := tools.(*unstructured.UnstructuredList); ok {
		status.ToolsCount = len(toolList.Items)
		
		// Count by namespace if listing all namespaces
		if allNamespaces {
			for _, tool := range toolList.Items {
				ns := tool.GetNamespace()
				status.ToolsByNamespace[ns]++
			}
		}
	}

	return nil
}

// outputStatusTable outputs status information in table format
func outputStatusTable(status StatusInfo) error {
	fmt.Println("Rapt Status")
	fmt.Println("===========")
	
	// Cluster connectivity
	clusterStatus := "❌ Not Connected"
	if status.ClusterConnected {
		clusterStatus = "✅ Connected"
	}
	fmt.Printf("Cluster:     %s\n", clusterStatus)

	// CRD status
	crdStatus := "❌ Not Installed"
	if status.CRDInstalled {
		crdStatus = "✅ Installed"
		if !status.CRDCreated.IsZero() {
			crdStatus += fmt.Sprintf(" (created: %s)", status.CRDCreated.Format("2006-01-02 15:04"))
		}
	}
	fmt.Printf("CRD:         %s\n", crdStatus)
	if status.CRDName != "" {
		fmt.Printf("CRD Name:    %s\n", status.CRDName)
	}

	// Namespace info
	fmt.Printf("Namespace:   %s\n", status.CurrentNamespace)
	if status.AllNamespaces {
		fmt.Println("Scope:       All Namespaces")
	} else {
		fmt.Println("Scope:       Current Namespace")
	}

	// Tools count
	fmt.Printf("Tools:       %d\n", status.ToolsCount)
	
	// Tools by namespace if applicable
	if status.AllNamespaces && len(status.ToolsByNamespace) > 0 {
		fmt.Println("\nTools by Namespace:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAMESPACE\tCOUNT")
		for ns, count := range status.ToolsByNamespace {
			fmt.Fprintf(w, "%s\t%d\n", ns, count)
		}
		w.Flush()
	}

	// Errors
	if len(status.Errors) > 0 {
		fmt.Println("\nIssues:")
		for _, err := range status.Errors {
			fmt.Printf("  ❌ %s\n", err)
		}
	}

	// Overall status
	fmt.Println("\nOverall Status:")
	if status.ClusterConnected && status.CRDInstalled {
		fmt.Println("  ✅ Rapt is properly installed and ready to use")
	} else {
		fmt.Println("  ❌ Rapt is not properly installed")
		if !status.ClusterConnected {
			fmt.Println("    - Check cluster connectivity")
		}
		if !status.CRDInstalled {
			fmt.Println("    - Run 'rapt init' to install the CRD")
		}
	}

	return nil
}

// outputStatusJSON outputs status information in JSON format
func outputStatusJSON(status StatusInfo) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(status)
}

// outputStatusYAML outputs status information in YAML format
func outputStatusYAML(status StatusInfo) error {
	yamlBytes, err := yamlv2.Marshal(status)
	if err != nil {
		return err
	}
	fmt.Print(string(yamlBytes))
	return nil
}