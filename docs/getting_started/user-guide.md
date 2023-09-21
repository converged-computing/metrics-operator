# User Guide

Welcome to the Metrics Operator user guide! If you come here, we are assuming you have a cluster
with the Metrics Operator installed and are interested to submit your own [custom resource](custom-resource-definition.md) to create a MetricSet, or that someone has already done it for you. If you are a developer wanting to work on new functionality or features, see our [Developer Guides](../development/index.md) instead.

## Usage

### Overview

Our "MetricSet" is mirroring the design of a [JobSet](https://github.com/kubernetes-sigs/jobset/), which can simply be defined as follows:

> A Metric Set is a collection of metrics to measure IO, performance, or networking that can be customized with addons.

When you create a MetricSet using this operator, we assume that you are primarily interested in measuring an application performance, collecting storage metrics, or
using a custom metric provided by the operator that has less stringent requirements.

Each metric provided by the operator (ranging from network to applications to io) has a prebuilt container, and knows how to launch one or more replicated jobs
to measure or assess the performance of something. A MetricSet itself is just a single shell for some metric, which can be further customized with addons.
A MetricAddon "addon" is flexible to be any kind of "extra" that is needed to supplement a metric run - e.g., applications, volumes/storage, or
even extra containers that add logic. High level, this includes:

 - Add extra containers (and config maps for their entrypoints)
 - Add custom logic to entrypoints for specific jobs and/or containers
 - Add additional volumes that range the gamut from empty to persistent disk.

And specific examples might include:

 - Every kind of volume is provided as a volume addon, this way you can run a storage metric against some kind of mounted storage.
 - A container (application) addon makes it easy to add your custom container to run alongside a metric that shares (and monitors) the process namespace
 - A monitoring tool provided via a modular install for a container can be provided as an addon, and it works by creating container, and sharing assets via an empty volume shared with some metric container(s) of interest. The sharing and setup of the volume happens via customizing the main metric entrypoint(s) and also adding a custom config map volume (for the addon container entrypoint).

Within this space, we can easily customize the patterns of metrics by way of shared interfaces. Common patterns for shared interfaces currently include a `LauncherWorker`, `SingleApplication`, and `StorageGeneric` design.

### Install

#### Quick Install

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

#### Helm Install

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

### Getting Started

Let's first review how this works.

1. We provide metrics here to assess performance, storage, networking, or other custom cases (e.g run an HPC application).
2. You can choose to supplement a metric with addons (e.g., add a volume to an IO metric)
3. The metric output is printed in pod logs with a standard packaging (e.g., sections and headers) to distinguish output sections.
4. We provide a Python module [metricsoperator](https://pypi.org/project/metricsoperator/) that can help you run an experiment, applying the metrics.yaml and then retrieving and parsing logs.

For the last step, this is important because every metric tool is a special snowflake, outputting some custom format that is hard to parse and then plot. By providing a parser paired with each metric, we hope to provide an easy means to do this so you can go from data collection to results more quickly. Now let's review a suggested set of steps for you as a new user! You can:

1. First choose one or more [metrics](metrics.md), [request a metric be added](https://github.com/converged-computing/metrics-operator/issues), or start with a pre-created [examples](https://github.com/converged-computing/metrics-operator/tree/main/examples). Often if you want to measure an application or storage or "other" (e.g., networking) we already have a metrics.yaml and associated parser suited for your needs.
2. Run the metric directly from the metrics.yaml, or use the Python module [metricsoperator](https://pypi.org/project/metricsoperator/) to run and collect output.
3. Plot the results, and you're done!

For step 2, you can always store output logs and then parse them later if desired.
For a quick start, you can likely explore our [examples](https://github.com/converged-computing/metrics-operator/tree/main/examples) directory,
which has both examples we use in [testing](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests) along with full
[Python examples](https://github.com/converged-computing/metrics-operator/tree/main/examples/python) that will use the associated metrics.yaml files, but also submit and parse the output.
Once you are comfortable with the basics, you can browse our [available metrics](metrics.md) and either design your
own, or [request a metric be added](https://github.com/converged-computing/metrics-operator/issues). Our goals are to
make these easy to deploy with minimal complexity for you, so we are happy to help. We also encourage you to share examples
and experiments that you put together here for others to use.

## Metrics

For all metric types, the following applies:

1. You can create more than one pod (scale the metric) as you see fit.
2. There is always a headless service provided for metrics within the JobSet to make use of.

For another overview of these designs, please see the [developer docs](../development/designs/index.md).

## Containers Available

All containers are provided under [ghcr.io/converged-computing/metrics-operator](https://github.com/converged-computing/metrics-operator/pkgs/container/metrics-operator). The latest tag is the current main branch, a "bleeding edge" version, and we will provide releases when the operator is more stable.
