package rapt

import (
	"context"
	"fmt"

	"codeberg.org/lig/rapt/internal/k8s"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func InitCmd(namespace string) error {
	k8sClient, err := k8s.InitClient(namespace)
	if err != nil {
		return err
	}

	crd, err := k8s.LoadToolCRD()
	if err != nil {
		return fmt.Errorf("failed to unmarshal tool.yaml: %w", err)
	}

	_, err = k8sClient.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), crd, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			fmt.Println("CRD already exists. Skipping creation.")
			return nil
		}
		return fmt.Errorf("failed to create CRD: %w", err)
	}

	fmt.Println("CRD created successfully.")
	return nil
}
