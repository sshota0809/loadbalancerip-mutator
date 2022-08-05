package k8s

import (
	"fmt"
	"github.com/sshota0809/loadbalancerip-mutator/pkg/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func NewClientSet() (*kubernetes.Clientset, error) {
	config, err := getK8sConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getK8sConfig() (*rest.Config, error) {
	// get InClusterConfig
	config, err := rest.InClusterConfig()

	if err == nil {
		logger.Log.Debug(fmt.Sprintf("Using Kubernetes configuration: %s", config))
		return config, nil
	}

	var kubeconfig string
	home := homedir.HomeDir()
	if home == "" {
		return nil, &NoHomeDirError{}
	}
	kubeconfig = filepath.Join(home, ".kube", "config")
	logger.Log.Debug(fmt.Sprintf("Kubernetes configuration: %s", kubeconfig))

	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}
