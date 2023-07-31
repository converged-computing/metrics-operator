# MetricSet

> The CRD "Custom Resource Definition" defines a MetricSet

A CRD is a yaml file that you can apply to your cluster (with the Metrics Operator
installed) to ask for a MetricSet to be deployed. Kubernetes has these [custom resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
to make it easy to automate tasks, and in fact this is the goal of an operator!
A Kubernetes operator is conceptually like a human operator that takes your CRD,
looks at the cluster state, and does whatever is necessary to get your cluster state
to match your request. In the case of the Metrics Operator, this means deploying one or more
metrics alongside an application or storage solution of interest, and returning the metrics to you. This document describes the spec of our custom resource definition.

## Custom Resource Definition

### Header

The yaml spec will normally have an API version, the kind `MetricSet` and then
a name and (optionally, a namespace) to identify the custom resource definition followed by the spec for it. Here is a spec that will deploy to the `default` namespace:

```yaml
apiVersion: flux-framework.org/v1alpha1
kind: MetricSet
metadata:
  labels:
    app.kubernetes.io/name: metricset
    app.kubernetes.io/instance: metricset-sample
  name: metricset-sample
spec:
  ...
```

### Spec

Under the spec, there are several variables to define. Descriptions are included below, and we
recommend that you look at [examples](https://github.com/converged-computing/metrics-operator/tree/main/examples) in the repository for more examples.

### application

The core of a MetricSet is an application. This is the container that houses some application that you want to measure performance for. This means that minimally, you are required to define the application container image and command:


```yaml
spec:
  application:
    image: ghcr.io/rse-ops/vanilla-lammps:tag-latest
    command: mpirun lmp -v x 1 -v y 1 -v z 1 -in in.reaxc.hns -nocite
```

In the above example, we target a container with LAMMPS and mpi, and we are going to run MPIrun.
The command will be used by the metrics sidecar containers to find the PID of interest to measure.

### metrics

The core of the MetricSet of course is the metrics! Since we can measure more than one thing at once, this is a list of named metrics known to the operator. As an example, here is how to run the `perf-sysstat` metric:

```yaml
spec:
  metrics:
    - name: perf-sysstat
```

To see all the metrics available, see [metrics](metrics.md). We will be adding many more as the operator is developed.
