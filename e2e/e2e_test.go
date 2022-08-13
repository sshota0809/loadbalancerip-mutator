package e2e_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe("e2e tests", func() {
	BeforeEach(func() {
		time.Sleep(10 * time.Second)
	})

	AfterEach(func() {
		services, err := clientSet.CoreV1().Services("default").List(context.Background(), metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())
		for _, service := range services.Items {
			err = clientSet.CoreV1().Services("default").Delete(context.Background(), service.GetName(), metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("available IP exists", func() {
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
			},
		}
		svc, err := clientSet.CoreV1().Services("default").Create(context.Background(), svc, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(svc.Spec.LoadBalancerIP).Should(Equal("10.10.10.10"))
	})

	It("available IP doesn't exists", func() {
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
			},
		}
		svc, err := clientSet.CoreV1().Services("default").Create(context.Background(), svc, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		svc = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "svc02",
				Annotations: map[string]string{},
			},
			Spec: corev1.ServiceSpec{
				Type: "LoadBalancer",
				Selector: map[string]string{
					"label": "test02",
				},
				Ports: []corev1.ServicePort{
					corev1.ServicePort{
						Protocol: "TCP",
						Port:     443,
					},
				},
			},
		}
		_, err = clientSet.CoreV1().Services("default").Create(context.Background(), svc, metav1.CreateOptions{})
		Expect(err).To(HaveOccurred())
	})
})
