package rapt

import (
	"context"
	"fmt"
	"strings"

	"codeberg.org/lig/rapt/internal/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"github.com/AlecAivazis/survey/v2"
)

// DeleteTools deletes one or more tool definitions from the cluster
func DeleteTools(namespace string, toolNames []string, force bool) error {
	// Initialize dynamic client
	dynClient, err := k8s.InitDynamicClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to initialize dynamic client: %w", err)
	}

	// Confirm deletion unless forced
	if !force {
		confirmMessage := fmt.Sprintf("Are you sure you want to delete tool(s): %s?", strings.Join(toolNames, ", "))
		prompt := &survey.Confirm{
			Message: confirmMessage,
			Default: false,
		}
		confirmed := false
		err = survey.AskOne(prompt, &confirmed)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	// Delete each tool
	gvr := schema.GroupVersionResource{
		Group:    "rapt.dev",
		Version:  "v1alpha1",
		Resource: "tools",
	}

	var deletedTools []string
	var failedTools []string

	for _, toolName := range toolNames {
		err := dynClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), toolName, metav1.DeleteOptions{})
		if err != nil {
			failedTools = append(failedTools, toolName)
			fmt.Printf("Failed to delete tool '%s': %v\n", toolName, err)
		} else {
			deletedTools = append(deletedTools, toolName)
			fmt.Printf("Successfully deleted tool '%s'\n", toolName)
		}
	}

	// Summary
	if len(deletedTools) > 0 {
		fmt.Printf("\nDeleted %d tool(s): %s\n", len(deletedTools), strings.Join(deletedTools, ", "))
	}
	if len(failedTools) > 0 {
		return fmt.Errorf("failed to delete %d tool(s): %s", len(failedTools), strings.Join(failedTools, ", "))
	}

	return nil
}

// DeleteAllTools deletes all tool definitions from the namespace
func DeleteAllTools(namespace string, force bool) error {
	// Initialize dynamic client
	dynClient, err := k8s.InitDynamicClient(namespace)
	if err != nil {
		return fmt.Errorf("failed to initialize dynamic client: %w", err)
	}

	// Get all tools first
	gvr := schema.GroupVersionResource{
		Group:    "rapt.dev",
		Version:  "v1alpha1",
		Resource: "tools",
	}

	tools, err := dynClient.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list tools: %w", err)
	}

	if len(tools.Items) == 0 {
		fmt.Printf("No tools found in namespace '%s'.\n", namespace)
		return nil
	}

	// Extract tool names
	toolNames := make([]string, len(tools.Items))
	for i, tool := range tools.Items {
		toolNames[i] = tool.GetName()
	}

	// Confirm deletion unless forced
	if !force {
		confirmMessage := fmt.Sprintf("Are you sure you want to delete ALL %d tool(s) in namespace '%s'?\nTools: %s", 
			len(toolNames), namespace, strings.Join(toolNames, ", "))
		prompt := &survey.Confirm{
			Message: confirmMessage,
			Default: false,
		}
		confirmed := false
		err = survey.AskOne(prompt, &confirmed)
		if err != nil {
			return fmt.Errorf("failed to get confirmation: %w", err)
		}
		if !confirmed {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	// Delete all tools
	var deletedTools []string
	var failedTools []string

	for _, toolName := range toolNames {
		err := dynClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), toolName, metav1.DeleteOptions{})
		if err != nil {
			failedTools = append(failedTools, toolName)
			fmt.Printf("Failed to delete tool '%s': %v\n", toolName, err)
		} else {
			deletedTools = append(deletedTools, toolName)
			fmt.Printf("Successfully deleted tool '%s'\n", toolName)
		}
	}

	// Summary
	fmt.Printf("\nDeleted %d out of %d tool(s) in namespace '%s'\n", len(deletedTools), len(toolNames), namespace)
	if len(failedTools) > 0 {
		return fmt.Errorf("failed to delete %d tool(s): %s", len(failedTools), strings.Join(failedTools, ", "))
	}

	return nil
}