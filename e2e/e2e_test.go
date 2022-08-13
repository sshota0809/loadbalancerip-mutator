package e2e_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("e2e tests", func() {
	It("test step", func() {
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "svc01",
				Annotations: map[string]string{},
			},
			Spec: corev1.ServiceSpec{
				Type: "LoadBalancer",
				Selector: map[string]string{
					"label": "test01",
				},
				Ports: []corev1.ServicePort{
					corev1.ServicePort{
						Protocol: "TCP",
						Port:     443,
					},
				},
				LoadBalancerIP: "10.10.10.10",
			},
		}
		_, err := k8sClientSet.CoreV1().Services("default").Create(context.Background(), svc, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
	})
})
