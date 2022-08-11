package mutation

import (
	"github.com/sshota0809/loadbalancerip-mutator/pkg/ip"
	"github.com/sshota0809/loadbalancerip-mutator/pkg/logger"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"testing"
)

func initLogger() {
	logger.Init("error")
}

func TestGetAvailableIP(t *testing.T) {
	var tests = []struct {
		description string
		ipPool      *ip.IPPool
		clientSet   *fake.Clientset
		expIpAdder  ip.IpAddr
		expErr      error
	}{
		{
			description: "AvailableIP exists",
			ipPool: &ip.IPPool{
				IPs: map[ip.IpAddr]ip.IsUsed{
					"10.10.10.10": false,
					"10.10.10.11": false,
				},
			},
			clientSet: fake.NewSimpleClientset(
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "svc01",
						Namespace:   "ns01",
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
				}),
			expIpAdder: "10.10.10.11",
			expErr:     nil,
		},
		{
			description: "AvailableIP doesn't exist",
			ipPool: &ip.IPPool{
				IPs: map[ip.IpAddr]ip.IsUsed{
					"10.10.10.10": false,
					"10.10.10.11": false,
				},
			},
			clientSet: fake.NewSimpleClientset(
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "svc01",
						Namespace:   "ns01",
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
				},
				&corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "svc02",
						Namespace:   "ns02",
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
						LoadBalancerIP: "10.10.10.11",
					},
				}),
			expErr: &NoAvailableIPError{},
		},
	}

	initLogger()

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {

			mutator := &LoadBalancerIpMutator{
				ipPool:    tt.ipPool,
				ClientSet: tt.clientSet,
			}

			availableIP, err := mutator.getAvailableIP()
			if err != nil {
				assert.ErrorAs(t, err, &tt.expErr, "Error should be equal to expErr")
				return
			}
			assert.Equal(t, availableIP, tt.expIpAdder, "IPs should be equal to expIpAdder")
		})
	}
}
