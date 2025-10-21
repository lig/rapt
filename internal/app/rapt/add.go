package rapt

import (
	"context"
	"fmt"
	"strings"

	"codeberg.org/lig/rapt/internal/k8s"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Add registers a new tool definition in the Kubernetes cluster.
func Add(namespace, name, image, command string, env []string, dryRun bool) error {
	// Validate required fields
	if name == "" {
		return fmt.Errorf("tool name is required")
	}
	if image == "" {
		return fmt.Errorf("container image is required")
	}

	// Prepare env variables for the Tool spec
	envVars := make([]map[string]any, len(env))
	for i, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid environment variable format: %s (expected NAME=VALUE)", e)
		}
		envVars[i] = map[string]any{
			"name":  strings.TrimSpace(parts[0]),
			"value": strings.TrimSpace(parts[1]),
		}
	}

	// Prepare the Tool object
	jobTemplate := map[string]any{
		"image": image,
		"env":   envVars,
	}
	
	// Only add command if it's not empty
	if command != "" {
		jobTemplate["command"] = strings.Fields(command)
	}

	tool := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "rapt.dev/v1alpha1",
			"kind":       "Tool",
			"metadata": map[string]any{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]any{
				"jobTemplate": jobTemplate,
			},
		},
	}

	// If dry-run mode, print YAML and exit
	if dryRun {
		return k8s.PrintToolYAML(tool)
	}

	// Initialize dynamic client
	dynClient, err := k8s.InitDynamicClient(namespace)
	if err != nil {
		return err
	}

	gvr := schema.GroupVersionResource{
		Group:    "rapt.dev",
		Version:  "v1alpha1",
		Resource: "tools",
	}

	_, err = dynClient.Resource(gvr).Namespace(namespace).Create(context.Background(), tool, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create tool %s: %w", name, err)
	}

	fmt.Printf("Successfully created tool '%s' in namespace '%s'\n", name, namespace)
	return nil
}
