# User Guide

Welcome to the Metrics Operator user guide! If you come here, we are assuming you have a cluster
with the Metrics Operator installed and are interested to submit your own [custom resource](custom-resource-definition.md) to create a MetricSet, or that someone has already done it for you. If you are a developer wanting to work
on new functionality or features, see our [Developer Guides](../development/index.md) instead.

## Containers Available

All containers are provided under [ghcr.io/converged-computing/metrics-operator](https://github.com/converged-computing/metrics-operator/pkgs/container/metrics-operator). The latest tag is the current main branch, a "bleeding edge" version,
and we will provide releases when the operator is more stable.

## Install

### Quick Install

This works best for production Kubernetes clusters, and comes down to first installing JobSet (not yet part of Kubernetes)

```bash
kind create cluster
VERSION=v0.2.0
kubectl apply --server-side -f https://github.com/kubernetes-sigs/jobset/releases/download/$VERSION/manifests.yaml
```

and then downloading the latest Metrics Operator yaml config, and applying it.

```bash
kubectl apply -f https://raw.githubusercontent.com/converged-computing/metrics-operator/main/examples/dist/metrics-operator.yaml
```

Note that from the repository, this config is generated with:

```bash
$ make build-config
```

and then saved to the main branch where you retrieve it from.


### Helm Install

We optionally provide an install with helm, which you can do either from the charts in the repository:

```bash
$ git clone https://github.com/converged-computing/metrics-operator
$ cd metrics-operator
$ helm install ./chart
```

Or directly from GitHub packages (an OCI registry):

```
# helm prior to v3.8.0
$ export HELM_EXPERIMENTAL_OCI=1
$ helm pull oci://ghcr.io/converged-computing/metrics-operator-helm/chart
```
```console
Pulled: ghcr.io/converged-computing/metrics-operator-helm/chart:0.1.0
```

And install!

```bash
$ helm install chart-0.1.0.tgz 
```
```console
NAME: metrics-operator
LAST DEPLOYED: Fri Mar 24 18:36:18 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
```

## Usage

Generally, you'll be defining an application container with one or more metrics to assess performance, or a storage solution with the same, but metrics to assess IO. There are several modes of operation, depending on your choice of metrics.

### Application Metrics

An application with metrics will allow 
A storage or IO metric will simply create the volume of interest that you request, and run the tool there. Read/write is important here - e.g., if the metric needs to write to the volume, a read only volume won't work.

For storage metrics, you aren't required to 

