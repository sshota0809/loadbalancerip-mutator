package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sshota0809/loadbalancerip-mutator/pkg/logger"
	"github.com/sshota0809/loadbalancerip-mutator/pkg/mutation"
	"github.com/sshota0809/loadbalancerip-mutator/webhook"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "loadbalancerip-mutator",
		Short: "MutationWebhook to attach loadBalancerIP to Service resource",
		Long:  `This application is MutationWebhook to attach loadBalancerIP to Service resource from a IP pool if not presented`,
		Run: func(cmd *cobra.Command, args []string) {
			// initialize logger
			level, _ := cmd.Flags().GetString("level")
			logger.Init(level)

			logger.Log.Info("Starting server...")

			port, _ := cmd.Flags().GetInt("port")

			tlsCertFile, _ := cmd.Flags().GetString("tls-cert-file")
			if tlsCertFile == "" {
				logger.Log.Error("The tls-cert-file option is not specified but required.")
				os.Exit(1)
			}

			tlsKeyFile, _ := cmd.Flags().GetString("tls-key-file")
			if tlsKeyFile == "" {
				logger.Log.Error("The tls-key-file option is not specified but required.")
				os.Exit(1)
			}

			pool, _ := cmd.Flags().GetString("pool")
			if pool == "" {
				logger.Log.Error("The pool option is not specified but required.")
				os.Exit(1)
			}

			h, err := mutation.NewLoadBalancerIpHandler(pool)
			if err != nil {
				logger.Log.Error(err.Error())
				os.Exit(1)
			}

			ws, err := webhook.NewWebhookServer(port, tlsCertFile, tlsKeyFile, h)
			if err != nil {
				logger.Log.Error(err.Error())
				os.Exit(1)
			}

			if err = ws.Run(); err != nil {
				logger.Log.Error(err.Error())
				os.Exit(1)
			}
		},
	}
)

func Run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("level", "v", "info", "[OPTIONAL] Log level. Valid value is debug, info, warn and error")
	rootCmd.Flags().StringP("pool", "p", "", "[REQUIRED] specify ip pool that will be attached through this MutationWebhook. Valid value is comma separated CIDR list e.g. \"10.10.100.10/32,10.10.10.128/25,10.10.100.0/24\"")
	rootCmd.Flags().StringP("tls-cert-file", "c", "", "[REQUIRED] path of TLS cert file")
	rootCmd.Flags().StringP("tls-key-file", "k", "", "[REQUIRED] path of TLS key file")
	rootCmd.Flags().IntP("port", "P", 8080, "[OPTIONAL] port number to listen")
}
