package mutation

import (
	"context"
	"fmt"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	"github.com/sshota0809/loadbalancerip-mutator/pkg/ip"
	"github.com/sshota0809/loadbalancerip-mutator/pkg/k8s"
	"github.com/sshota0809/loadbalancerip-mutator/pkg/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"net/http"
)

const (
	webhookID = "loadbalancerip-mutator"
)

type LoadBalancerIpMutator struct {
	ipPool    *ip.IPPool
	handler   http.Handler
	ClientSet *kubernetes.Clientset
}

func NewLoadBalancerIpHandler(pool string) (*LoadBalancerIpMutator, error) {
	// initialize ipPool
	ipPool, err := ip.NewIpPool(pool)
	if err != nil {
		return nil, err
	}

	// initialize ClientSet
	clientSet, err := k8s.NewClientSet()
	if err != nil {
		return nil, fmt.Errorf("Error creating clientSet: %w", err)
	}

	return &LoadBalancerIpMutator{
		ipPool:    ipPool,
		ClientSet: clientSet,
	}, nil
}

func (h *LoadBalancerIpMutator) GenerateHandler() (http.Handler, error) {
	// initialize mutator
	mt := kwhmutating.MutatorFunc(h.mutate)
	l := kwhlogrus.NewLogrus(logger.Log)

	// create webhook
	mcfg := kwhmutating.WebhookConfig{
		ID:      webhookID,
		Mutator: mt,
		Logger:  l,
	}
	wh, err := kwhmutating.NewWebhook(mcfg)
	if err != nil {
		return nil, fmt.Errorf("Error creating webhook: %w", err)
	}

	// get HTTP handler from webhook
	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{Webhook: wh, Logger: l})
	if err != nil {
		return nil, fmt.Errorf("Error creating webhook handler: %w", err)
	}

	return whHandler, nil
}

func (h *LoadBalancerIpMutator) getAvailableIP() (ip.IpAddr, error) {
	var none ip.IpAddr = ""

	// get services from all namespaces
	services, err := h.ClientSet.CoreV1().Services(metav1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return none, err
	}

	// deepCopy ipPool.IPs
	ips := make(map[ip.IpAddr]ip.IsUsed)
	for ipAddr, _ := range h.ipPool.IPs {
		ips[ipAddr] = false
	}

	// Check if IP addr is attached to existing Service resource as loadBalancerIP
	for _, service := range services.Items {
		if service.Spec.LoadBalancerIP != "" {
			ips[ip.IpAddr(service.Spec.LoadBalancerIP)] = true
		}
	}

	// extract first IP addr that is not used yet
	var availableIP ip.IpAddr
	for ipAddr, isUsed := range ips {
		if isUsed {
			continue
		}
		availableIP = ipAddr
		break
	}

	if availableIP == "" {
		return none, &NoAvailableIPError{}
	}

	return availableIP, nil
}

func (h *LoadBalancerIpMutator) mutate(_ context.Context, _ *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
	service, ok := obj.(*corev1.Service)
	if !ok {
		return &kwhmutating.MutatorResult{}, nil
	}

	// check if loadBalancerIP parameter is presented and LoadBalancer type
	if service.Spec.Type != "LoadBalancer" {
		logger.Log.Debug(fmt.Sprintf("%s is not LoadBalancer Type but %s.", service.GetName(), service.Spec.Type))
		return &kwhmutating.MutatorResult{MutatedObject: service}, nil
	}
	if service.Spec.LoadBalancerIP != "" {
		logger.Log.Debug(fmt.Sprintf("%s already has loadBalancerIP.", service.GetName()))
		return &kwhmutating.MutatorResult{MutatedObject: service}, nil
	}

	// Check available IP addr if loadBalancerIP parameter is not presented
	availableIP, err := h.getAvailableIP()
	if err != nil {
		return &kwhmutating.MutatorResult{MutatedObject: service}, err
	}

	logger.Log.Info(fmt.Sprintf("Attaching %s as available IP addr to %s.", availableIP, service.GetName()))

	service.Spec.LoadBalancerIP = availableIP.ToString()
	return &kwhmutating.MutatorResult{MutatedObject: service}, nil
}
