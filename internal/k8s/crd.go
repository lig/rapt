package k8s

import (
	_ "embed"

	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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
