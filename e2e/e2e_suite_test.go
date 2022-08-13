package e2e_test

import (
	"context"
	"fmt"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmanagermmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	certmanager "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	arv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	clientSet         kubernetes.Interface
	certManagerClient *certmanager.Clientset
)

var _ = BeforeSuite(func() {
	initK8sClients()
	createMutationWebhook()
})

func initK8sClients() {
	home := homedir.HomeDir()
	if home == "" {
		Fail("Home dir is not found")
	}
	kubeconfig := filepath.Join(home, ".kube", "config")
	c, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	Expect(err).NotTo(HaveOccurred())

	cs, err := kubernetes.NewForConfig(c)
	Expect(err).NotTo(HaveOccurred())
	clientSet = cs

	cmc, err := certmanager.NewForConfig(c)
	Expect(err).NotTo(HaveOccurred())
	certManagerClient = cmc
}

func createMutationWebhook() {
	createNamespace("webhook")
	createIssuer("selfsigned", "webhook")
	createCert("webhook-certificate", "webhook", "selfsigned")
	createWebhook("webhook", "webhook", "webhook-certificate", "10.10.10.10/32")
	createWebhookConfig("webhook", "webhook", "webhook-certificate")
}

func createNamespace(name string) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := clientSet.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())
}

func createIssuer(name, ns string) {
	issuer := &certmanagerv1.Issuer{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: certmanagerv1.IssuerSpec{
			IssuerConfig: certmanagerv1.IssuerConfig{
				SelfSigned: &certmanagerv1.SelfSignedIssuer{},
			},
		},
	}
	_, err := certManagerClient.CertmanagerV1().Issuers(ns).Create(context.Background(), issuer, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())
}

func createCert(name, ns, issuer string) {
	cert := &certmanagerv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: certmanagerv1.CertificateSpec{
			SecretName: name,
			DNSNames:   []string{"webhook", "webhook.webhook", "webhook.webhook.svc", "webhook.webhook.svc.cluster.local"},
			IssuerRef: certmanagermmetav1.ObjectReference{
				Name: issuer,
			},
		},
	}
	_, err := certManagerClient.CertmanagerV1().Certificates(ns).Create(context.Background(), cert, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())
}

func createWebhook(name, ns, cert, pool string) {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := clientSet.CoreV1().ServiceAccounts(ns).Create(context.Background(), sa, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	cr := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				APIGroups: []string{""},
				Resources: []string{"services"},
				Verbs: []string{
					"get",
					"watch",
					"list",
				},
			},
		},
	}
	_, err = clientSet.RbacV1().ClusterRoles().Create(context.Background(), cr, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      name,
				Namespace: ns,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     name,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	_, err = clientSet.RbacV1().ClusterRoleBindings().Create(context.Background(), crb, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"name": name,
			},
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Protocol:   "TCP",
					Port:       443,
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
	}
	_, err = clientSet.CoreV1().Services(ns).Create(context.Background(), svc, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": name,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: name,
					Containers: []corev1.Container{
						{
							Name:            "loadbalancerip-mutator",
							Image:           "ghcr.io/sshota0809/loadbalancerip-mutator:latest",
							ImagePullPolicy: corev1.PullAlways,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "webhook-certificate",
									MountPath: "/etc/cert",
								},
							},
							Args: []string{
								"--pool",
								pool,
								"-v",
								"debug",
								"--tls-key-file",
								"/etc/cert/tls.key",
								"--tls-cert-file",
								"/etc/cert/tls.crt",
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "webhook-certificate",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName:  "webhook-certificate",
									DefaultMode: int32Ptr(0644),
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = clientSet.AppsV1().Deployments(ns).Create(context.Background(), deployment, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())
}

func createWebhookConfig(name, ns, cert string) {
	fp := arv1.Fail
	spoce := arv1.NamespacedScope
	es := arv1.SideEffectClassNone
	mwc := &arv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"cert-manager.io/inject-ca-from": fmt.Sprintf("%s/%s", ns, cert),
			},
		},
		Webhooks: []arv1.MutatingWebhook{
			{
				Name:          fmt.Sprintf("%s.%s.svc.cluster.local", name, ns),
				FailurePolicy: &fp,
				Rules: []arv1.RuleWithOperations{
					{
						Operations: []arv1.OperationType{
							arv1.Create,
							arv1.Update,
						},
						Rule: arv1.Rule{
							APIGroups:   []string{""},
							APIVersions: []string{"v1"},
							Resources:   []string{"services"},
							Scope:       &spoce,
						},
					},
				},
				ClientConfig: arv1.WebhookClientConfig{
					Service: &arv1.ServiceReference{
						Namespace: ns,
						Name:      name,
						Path:      stringPtr("/mutate"),
					},
				},
				AdmissionReviewVersions: []string{
					"v1",
					"v1beta1",
				},
				TimeoutSeconds: int32Ptr(10),
				SideEffects:    &es,
			},
		},
	}
	_, err := clientSet.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.Background(), mwc, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())
}

func int32Ptr(i int32) *int32    { return &i }
func stringPtr(s string) *string { return &s }
