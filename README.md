## Usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

helm repo add sshota0809 https://sshota0809.github.io/loadbalancerip-mutator

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
sshota0809` to see the charts.

To install the  chart:

    helm install my-loadbalancerip-mutator sshota0809/loadbalancerip-mutator

To uninstall the chart:

    helm delete my-loadbalancerip-mutator
