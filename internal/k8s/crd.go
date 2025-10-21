package k8s

import (
	_ "embed"
	"fmt"

	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

//go:embed data/tool.yaml
var toolCRDYAML []byte

func LoadToolCRD() (*apiv1.CustomResourceDefinition, error) {
	var crd apiv1.CustomResourceDefinition
	if err := yaml.Unmarshal(toolCRDYAML, &crd); err != nil {
		return nil, err
	}
	return &crd, nil
}

// PrintCRDYAML prints a CustomResourceDefinition as YAML to stdout
func PrintCRDYAML(crd *apiv1.CustomResourceDefinition) error {
	yamlBytes, err := yaml.Marshal(crd)
	if err != nil {
		return fmt.Errorf("failed to marshal CRD to YAML: %w", err)
	}
	fmt.Print(string(yamlBytes))
	return nil
}

// PrintToolYAML prints a Tool custom resource as YAML to stdout
func PrintToolYAML(tool *unstructured.Unstructured) error {
	yamlBytes, err := yaml.Marshal(tool.Object)
	if err != nil {
		return fmt.Errorf("failed to marshal Tool to YAML: %w", err)
	}
	fmt.Print(string(yamlBytes))
	return nil
}
