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
	"k8s.io/client-go/dynamic"
	yamlv2 "sigs.k8s.io/yaml"
)

// ToolInfo represents information about a tool for display
type ToolInfo struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Image       string            `json:"image"`
	Command     []string          `json:"command,omitempty"`
	Arguments   []ToolArgument    `json:"arguments,omitempty"`
	Environment []ToolEnvironment `json:"environment,omitempty"`
	Help        string            `json:"help,omitempty"`
	Created     time.Time         `json:"created"`
}

// ToolArgument represents a tool argument
type ToolArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
}

// ToolEnvironment represents a tool environment variable
type ToolEnvironment struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ListTools lists all available tools in the cluster
func ListTools(namespace, outputFormat string, allNamespaces bool) error {
	// Initialize dynamic client
	dynClient, err := k8s.InitDynamicClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to initialize dynamic client: %w", err)
	}

	// If no namespace was provided, treat it as all namespaces
	showAllNamespaces := allNamespaces || namespace == ""

	// Get tools
	tools, err := getTools(dynClient, namespace, showAllNamespaces)
	if err != nil {
		return fmt.Errorf("failed to get tools: %w", err)
	}

	if len(tools) == 0 {
		if showAllNamespaces {
			fmt.Println("No tools found in any namespace.")
		} else {
			fmt.Printf("No tools found in namespace '%s'.\n", namespace)
		}
		return nil
	}

	// Convert to ToolInfo for display
	toolInfos := make([]ToolInfo, len(tools))
	for i, tool := range tools {
		toolInfo, err := convertToToolInfo(tool)
		if err != nil {
			return fmt.Errorf("failed to convert tool %s: %w", tool.GetName(), err)
		}
		toolInfos[i] = toolInfo
	}

	// Output based on format
	switch outputFormat {
	case "json":
		return outputJSON(toolInfos)
	case "yaml":
		return outputYAML(toolInfos)
	case "table":
		fallthrough
	default:
		return outputTable(toolInfos, showAllNamespaces)
	}
}

// getTools retrieves tools from Kubernetes
func getTools(dynClient dynamic.Interface, namespace string, allNamespaces bool) ([]*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    "rapt.dev",
		Version:  "v1alpha1",
		Resource: "tools",
	}

	var listOptions metav1.ListOptions
	if allNamespaces {
		listOptions = metav1.ListOptions{}
	} else {
		listOptions = metav1.ListOptions{}
	}

	var tools *unstructured.UnstructuredList
	var err error

	if allNamespaces {
		tools, err = dynClient.Resource(gvr).List(context.TODO(), listOptions)
	} else {
		tools, err = dynClient.Resource(gvr).Namespace(namespace).List(context.TODO(), listOptions)
	}

	if err != nil {
		return nil, err
	}

	// Convert []unstructured.Unstructured to []*unstructured.Unstructured
	result := make([]*unstructured.Unstructured, len(tools.Items))
	for i := range tools.Items {
		result[i] = &tools.Items[i]
	}
	return result, nil
}

// convertToToolInfo converts an unstructured tool to ToolInfo
func convertToToolInfo(tool *unstructured.Unstructured) (ToolInfo, error) {
	toolInfo := ToolInfo{
		Name:      tool.GetName(),
		Namespace: tool.GetNamespace(),
		Created:   tool.GetCreationTimestamp().Time,
	}

	// Extract spec
	spec, found, err := unstructured.NestedMap(tool.Object, "spec")
	if err != nil || !found {
		return toolInfo, fmt.Errorf("invalid tool spec")
	}

	// Extract help text
	if help, found, err := unstructured.NestedString(spec, "help"); err == nil && found {
		toolInfo.Help = help
	}

	// Extract job template
	jobTemplate, found, err := unstructured.NestedMap(spec, "jobTemplate")
	if err != nil || !found {
		return toolInfo, fmt.Errorf("tool spec missing jobTemplate")
	}

	// Extract image
	if image, found, err := unstructured.NestedString(jobTemplate, "image"); err == nil && found {
		toolInfo.Image = image
	}

	// Extract command
	if command, found, err := unstructured.NestedStringSlice(jobTemplate, "command"); err == nil && found {
		toolInfo.Command = command
	}

	// Extract arguments
	if args, found, err := unstructured.NestedSlice(spec, "arguments"); err == nil && found {
		toolInfo.Arguments = make([]ToolArgument, len(args))
		for i, argItem := range args {
			if argMap, ok := argItem.(map[string]interface{}); ok {
				name, _ := argMap["name"].(string)
				description, _ := argMap["description"].(string)
				required, _ := argMap["required"].(bool)
				defaultValue, _ := argMap["default"].(string)

				toolInfo.Arguments[i] = ToolArgument{
					Name:        name,
					Description: description,
					Required:    required,
					Default:     defaultValue,
				}
			}
		}
	}

	// Extract environment variables
	if env, found, err := unstructured.NestedSlice(jobTemplate, "env"); err == nil && found {
		toolInfo.Environment = make([]ToolEnvironment, len(env))
		for i, envItem := range env {
			if envMap, ok := envItem.(map[string]interface{}); ok {
				name, _ := envMap["name"].(string)
				value, _ := envMap["value"].(string)

				toolInfo.Environment[i] = ToolEnvironment{
					Name:  name,
					Value: value,
				}
			}
		}
	}

	return toolInfo, nil
}

// outputTable outputs tools in table format
func outputTable(tools []ToolInfo, allNamespaces bool) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	if allNamespaces {
		fmt.Fprintln(w, "NAME\tNAMESPACE\tIMAGE\tCOMMAND\tARGUMENTS\tCREATED")
	} else {
		fmt.Fprintln(w, "NAME\tIMAGE\tCOMMAND\tARGUMENTS\tCREATED")
	}

	// Print tools
	for _, tool := range tools {
		command := ""
		if len(tool.Command) > 0 {
			command = tool.Command[0]
			if len(tool.Command) > 1 {
				command += "..."
			}
		}

		args := ""
		if len(tool.Arguments) > 0 {
			args = fmt.Sprintf("%d args", len(tool.Arguments))
		}

		created := tool.Created.Format("2006-01-02 15:04")

		if allNamespaces {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				tool.Name, tool.Namespace, tool.Image, command, args, created)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				tool.Name, tool.Image, command, args, created)
		}
	}

	return nil
}

// outputJSON outputs tools in JSON format
func outputJSON(tools []ToolInfo) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tools)
}

// outputYAML outputs tools in YAML format
func outputYAML(tools []ToolInfo) error {
	yamlBytes, err := yamlv2.Marshal(tools)
	if err != nil {
		return err
	}
	fmt.Print(string(yamlBytes))
	return nil
}