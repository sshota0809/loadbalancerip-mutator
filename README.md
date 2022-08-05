# loadbalancerip-mutator
Mutation webhook to attach loadBalancerIP param to Service resource.

## Motivation

There are some service controllers that doesn't automatically attach `loadBalancerIP` param to Service resouces with `type: LoadBalancer` param. e.g. Anthos on VMware integrated with Seesaw load balancer. This mutation webhook allows to automatically attach it from IP address pool specified by option.

## Usage

```
This application is MutationWebhook to attach loadBalancerIP to Service resource from a IP pool if not presented

Usage:
  loadbalancerip-mutator [flags]

Flags:
  -h, --help                   help for loadbalancerip-mutator
  -v, --level string           [OPTIONAL] Log level. Valid value is debug, info, warn and error (default "info")
  -p, --pool string            [REQUIRED] specify ip pool that will be attached through this MutationWebhook. Valid value is comma separated CIDR list e.g. "10.10.100.10/32,10.10.10.128/25,10.10.100.0/24"
  -c, --tls-cert-file string   [REQUIRED] path of TLS cert file
  -k, --tls-key-file string    [REQUIRED] path of TLS key file
```

## Getting Started

Mutation webhook needs to be attached TLS certificate. I can recommend to use [cert-manager](https://github.com/cert-manager/cert-manager) to manage it. Once you prepare TLS certificate you can deploy loadbalancerip-mutator through helm. Here is a helmfile example.

```
repositories:
  - name: loadbalancerip-mutator
    url: https://sshota0809.github.io/loadbalancerip-mutator

releases:
  - name: loadbalancerip-mutator
    namespace: loadbalancerip-mutator
    chart: loadbalancerip-mutator/loadbalancerip-mutator
    version: 0.1.0
```

## LICENSE

[MIT](LICENSE)