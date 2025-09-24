package k8s

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"

	clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func InitClient(namespace string) (*clientset.Clientset, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	configOverrides.Context.Namespace = namespace
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cluster := config.Host

	namespace, _, err = kubeConfig.Namespace()
	if err != nil {
		return nil, err
	}
	if namespace == "" {
		namespace = "unknown"
	}

	prompt := &survey.Confirm{
		Message: "Rapt will work on the cluster: " + cluster + "\nNamespace: " + namespace + "\nProceed?",
		Default: false,
	}
	proceed := false
	err = survey.AskOne(prompt, &proceed)
	if err != nil {
		return nil, err
	}
	if !proceed {
		return nil, errors.New("operation aborted by user")
	}

	return client, nil
}

// InitDynamicClient initializes and returns a dynamic Kubernetes client using the same kubeconfig logic as InitClient.
func InitDynamicClient(namespace string) (dynamic.Interface, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	configOverrides.Context.Namespace = namespace
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dynClient, nil
}
