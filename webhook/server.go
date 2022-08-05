package webhook

import (
	"fmt"
	"github.com/sshota0809/loadbalancerip-mutator/pkg/logger"
	"github.com/sshota0809/loadbalancerip-mutator/pkg/mutation"
	"net/http"
)

const (
	port = ":8080"
)

type certConfig struct {
	certFile string
	keyFile  string
}

type webhookServer struct {
	certConfig *certConfig
	Mutator    mutation.Mutator
}

func NewWebhookServer(tlsCertFile, tlsKeyFile string, handler mutation.Mutator) (*webhookServer, error) {
	return &webhookServer{
		certConfig: &certConfig{
			certFile: tlsCertFile,
			keyFile:  tlsKeyFile,
		},
		Mutator: handler,
	}, nil
}

func (ws *webhookServer) Run() error {
	h, err := ws.Mutator.GenerateHandler()
	if err != nil {
		return err
	}

	logger.Log.Info(fmt.Sprintf("Listening on %s", port))
	mux := http.NewServeMux()
	mux.Handle("/mutate", h)
	mux.Handle("/health", ws)
	err = http.ListenAndServeTLS(port, ws.certConfig.certFile, ws.certConfig.keyFile, mux)
	if err != nil {
		return err
	}

	return nil
}

func (ws *webhookServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Health Check OK.")
}
