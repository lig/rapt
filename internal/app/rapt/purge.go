package rapt

import (
	"context"
	"fmt"

	"codeberg.org/lig/rapt/internal/k8s"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Purge(namespace string) error {
	k8sClient, err := k8s.InitClient(namespace)
	if err != nil {
		return err
	}

	crd, err := k8s.LoadToolCRD()
	if err != nil {
		return fmt.Errorf("failed to unmarshal tool.yaml: %w", err)
	}

	crdName := crd.GetName()
	err = k8sClient.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), crdName, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			fmt.Println("CRD not found. Nothing to delete.")
			return nil
		}
		return fmt.Errorf("failed to delete CRD: %w", err)
	}

	fmt.Println("CRD deleted successfully.")
	return nil
}
