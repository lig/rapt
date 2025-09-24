package rapt

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"codeberg.org/lig/rapt/internal/k8s"
	yamlv2 "sigs.k8s.io/yaml"
)

// DescribeTool shows detailed information about a specific tool
func DescribeTool(namespace, toolName, outputFormat string) error {
	// Initialize dynamic client
	dynClient, err := k8s.InitDynamicClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to initialize dynamic client: %w", err)
	}

	// Get the tool definition
	tool, err := getToolDefinition(dynClient, namespace, toolName)
	if err != nil {
		return fmt.Errorf("failed to get tool definition: %w", err)
	}

	// Convert to ToolInfo for display
	toolInfo, err := convertToToolInfo(tool)
	if err != nil {
		return fmt.Errorf("failed to convert tool: %w", err)
	}

	// Output based on format
	switch outputFormat {
	case "json":
		return outputToolJSON(toolInfo)
	case "yaml":
		return outputToolYAML(toolInfo)
	case "table":
		fallthrough
	default:
		return outputToolTable(toolInfo)
	}
}


// outputToolTable outputs tool information in table format
func outputToolTable(tool ToolInfo) error {
	fmt.Printf("Name:        %s\n", tool.Name)
	fmt.Printf("Namespace:   %s\n", tool.Namespace)
	fmt.Printf("Created:     %s\n", tool.Created.Format("2006-01-02 15:04:05"))
	fmt.Printf("Image:       %s\n", tool.Image)
	
	if tool.Help != "" {
		fmt.Printf("Help:        %s\n", tool.Help)
	}

	if len(tool.Command) > 0 {
		fmt.Printf("Command:     %s\n", strings.Join(tool.Command, " "))
	}

	if len(tool.Arguments) > 0 {
		fmt.Println("\nArguments:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tDESCRIPTION\tREQUIRED\tDEFAULT")
		for _, arg := range tool.Arguments {
			required := "No"
			if arg.Required {
				required = "Yes"
			}
			defaultValue := "-"
			if arg.Default != "" {
				defaultValue = arg.Default
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", arg.Name, arg.Description, required, defaultValue)
		}
		w.Flush()
	}

	if len(tool.Environment) > 0 {
		fmt.Println("\nEnvironment Variables:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVALUE")
		for _, env := range tool.Environment {
			fmt.Fprintf(w, "%s\t%s\n", env.Name, env.Value)
		}
		w.Flush()
	}

	return nil
}

// outputToolJSON outputs tool information in JSON format
func outputToolJSON(tool ToolInfo) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tool)
}

// outputToolYAML outputs tool information in YAML format
func outputToolYAML(tool ToolInfo) error {
	yamlBytes, err := yamlv2.Marshal(tool)
	if err != nil {
		return err
	}
	fmt.Print(string(yamlBytes))
	return nil
}