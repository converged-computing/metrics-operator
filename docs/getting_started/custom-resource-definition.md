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

Under the spec, there are several variables to define. Descriptions are included below, and we recommend that you look at [examples](https://github.com/converged-computing/metrics-operator/tree/main/examples) in the repository for more.
Specifically, you must choose ONE of:

 - application
 - storage

Where an application will be run for some number of pods (completions) and measured by metrics pods (separate pods) OR a storage metric will run directly, and with some
number of pods (completions) to bind to the storage and measure.

### completions

The number of completions for an application or storage metric test will correspond with the number of indexed job completions (pods) for the storage or application JobSet. This defaults to 1, meaning we run in a non-indexed mode. The indexed mode is determined automatically by this variable, where "1" indicates non-indexed, and >1 is indexed.

### application

When you want to measure application performance, you'll need to add an "application" section to your MetricSet. This is the container that houses some application that you want to measure performance for. This means that minimally, you are required to define the application container image and command:


```yaml
spec:
  application:
    image: ghcr.io/rse-ops/vanilla-lammps:tag-latest
    command: mpirun lmp -v x 1 -v y 1 -v z 1 -in in.reaxc.hns -nocite
```

In the above example, we target a container with LAMMPS and mpi, and we are going to run MPIrun.
The command will be used by the metrics sidecar containers to find the PID of interest to measure.

#### volumes

An application is allowed to have one or more existing volumes. An existing volume can be any of the types described in [existing volumes](#existing-volumes)

### storage

When you want to measure some storage performance, you'll want to add a "storage" section to your MetricSet. This will typically just be a reference to some existing storage (see [existing volumes](#existing-volumes)) that we want to measure, and can
also be done for some number of completions and metrics for storage.

### metrics

The core of the MetricSet of course is the metrics! Since we can measure more than one thing at once, this is a list of named metrics known to the operator. As an example, here is how to run the `perf-sysstat` metric:

```yaml
spec:
  metrics:
    - name: perf-sysstat
```

To see all the metrics available, see [metrics](metrics.md). We will be adding many more as the operator is developed.

#### rate

A metric will be collected at some rate (in seconds) and this defaults to 10.
To change the rate for a metric:

```yaml
spec:
  metrics:
    - name: perf-sysstat
      rate: 20
```

### completions

Completions for a metric are relevant if you are assessing storage (which doesn't have an application runtime) or a service application that will continue to run forever. When this value is set to 0, it essentially indicates no set number of completions (meaning we run forever). Any non-zero value will ensure the metric
runs for that many completions before exiting.

```yaml
spec:
  metrics:
    - name: io-sysstat
      completions: 5
```

This is usually suggested to provide for a storage metric.

## Existing Volumes

An existing volume can be provided to support an application (multiple) or one can be provided for assessing its performance (single).

 - a persistent volume claim (PVC) and persistent volume (PV) that you've created
 - a secret that you've created
 - a config map that you've created
 - a host volume (typically for testing)

and for all of the above, you want to provide it to the operator, which will ensure the volume is available for your application or storage. For an application, you'd define your volumes as such:

```yaml
spec:
  application:
    image: ghcr.io/rse-ops/vanilla-lammps:tag-latest
    command: nginx -g daemon off;
    volumes:
      data:
        path: /workflow
        claimName: data 
```

The use case above, for an application, is that it requires some kind of data or storage alongside it to function. The volumes spec above is a key value (e.g., "data" is the key) to ensure that names are unique. For storage, you'll only be defining one volume:

```yaml
spec:
  storage:
    volume:
      path: /workflow
      claimName: data 
```

And the implicit name would be "storage" (although it's probably not important for you to know that). For the remaining examples, we will provide examples for application volumes, however know that the examples are also valid for the second
storage format.

#### persistent volume claim example

As an example, here is how to provide the name of an existing claim (you created separately) to a container:

```yaml
spec:
  application:
    image: ghcr.io/rse-ops/vanilla-lammps:tag-latest
    command: nginx -g daemon off;

    # This is an existing PVC (and associated PV) we created before the MetricSet
    volumes:
      data:
        path: /workflow
        claimName: data 
```

The above would add a claim named "data" to the application container(s). 

#### config map example

Here is an example of providing a config map to an application container In layman's terms, we are deploying vanilla nginx, but adding a configuration file
to `/etc/nginx/conf.d`

```yaml
spec:
  application:
    image: nginx
    command: nginx -g daemon off;

    # This is an existing PVC (and associated PV) we created before the MetricSet
    volumes:
      nginx-conf:
        configMapName: nginx-conf 
        path: /etc/nginx/conf.d
        items:
          flux.conf: flux.conf
```


You would have created this config map first, before the MetricSet. Here is an example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-conf
  namespace: metrics-operator
data:
  flux.conf: |
    server {
        listen       80;
        server_name  localhost;
        location / {
          root   /usr/share/nginx/html;
          index  index.html index.htm;
        }        
    }
```

#### secret example

Here is an example of providing an existing secret (in the metrics-operator namespace)
to the application container(s):

```yaml
spec:
  application:
    image: nginx
    command: nginx -g daemon off;

    volumes:
      certs:
        path: /etc/certs
        secretName: certs
```

The above shows an existing secret named "certs" that we will mount into `/etc/certs`.

#### hostpath volume example

Here is how to use a host path:

```yaml
spec:
  application:
    image: ghcr.io/rse-ops/vanilla-lammps:tag-latest
    command: nginx -g daemon off;

    # This is an existing PVC (and associated PV) we created before the MetricSet
    volumes:
      data:
        hostPath: true
        path: /workflow
```