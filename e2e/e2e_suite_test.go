package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"testing"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var (
	k8sClientSet kubernetes.Interface
)

var _ = BeforeSuite(func() {
	var kubeconfig string
	home := homedir.HomeDir()
	if home == "" {
		Fail("Home dir is not found")
	}
	kubeconfig = filepath.Join(home, ".kube", "config")
	c, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	Expect(err).NotTo(HaveOccurred())

	cs, err := kubernetes.NewForConfig(c)
	Expect(err).NotTo(HaveOccurred())

	k8sClientSet = cs
})
