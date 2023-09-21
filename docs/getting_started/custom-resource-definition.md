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
apiVersion: flux-framework.org/v1alpha2
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

Under the spec, there are several variables to define. Descriptions are included below, and we recommend that you look at [examples](https://github.com/converged-computing/metrics-operator/tree/main/examples) in the repository for more. Note that the general design takes one or more metrics, and each metric can have additional addons for storage volumes, additional containers, or other addon types.
Specifically, you must choose ONE of:

### pods

The number of pods for an application or storage metric test will correspond with the parallelism of the indexed job (which comes down to pods) for the storage or application JobSet. This defaults to 1, meaning we run in a non-indexed mode. The indexed mode is determined automatically by this variable, where "1" indicates non-indexed, and >1 is indexed.

### logging

We are anticipating adding more logging options, but for not logging exposes one "interactive" option that will add a "sleep infinity" to the end of a storage, performance, or standalone metric.
This is intended for debugging purposes.

```yaml
logging:
  interactive: false
```

It is typically added to a launcher or main container, if relevant, since workers tend to sleep anyway and the JobSet completion depends on the launcher.
By default, of course, it is set to false so the metric container and JobSet will finish.

### dontSetFQDN

For more of an "expert mode" if you know you want your JobSet use fully qualified domain names (FQDN) set to false,
set this value to true.

```yaml
spec:
  dontSetFQDN: true
```

By default it is false, meaning we use fully qualified domain names.

### metrics

The core of the MetricSet of course is the metrics! Since we can measure more than one thing at once, this is a list of named metrics known to the operator. As an example, here is how to run the `perf-sysstat` metric:

```yaml
spec:
  metrics:
    - name: perf-sysstat
```

To see all the metrics available, see [metrics](metrics.md). We will be adding many more as the operator is developed.

#### options

Generally, the specific parameters for any given metric are defined via the options, including:

 - options (key value pairs, where the value is an integer/string type)
 - listOptions (key value pairs, where the value is a list of integer/string types)
 - mapOptions (key value pairs, where the value is a map (string key) of integer/string types)

Here is an example with all three:

```yaml
spec:
  metrics:
    - name: perf-dummy
      options:
        pids: true
      
      listOptions:
        pids: [1, 2, 3]

      mapOptions:
        commands:
           "1": echo hello
           "2": echo goodbye
```

Presence of absence of an option type depends on the metric. Metrics are free to use these custom
options as they see fit, and validate in the same manner.

#### addons

An addon is a flexible interface to define everything from volumes to containers to be deployed alongside the metric.
If you are curious, a metric will generate one or more replicated Jobs in a Jobset, and the addon is free to customize these.
Akin to [metric options](#options) addons support the same types:

 - options
 - listOptions
 - mapOptions

As an example, here is a metric with a few named addons - an empty volume, and adding hpctoolkit to run alongside lammps.

```yaml
metrics:
 - name: app-lammps
   addons:
     - name: volume-empty
     - name: perf-hpctoolkit
```

Each addon has its own custom options. You can look at examples and at our [addons documentation](addons.md) for more detail on how to add existing volumes
or other custom functionality.
