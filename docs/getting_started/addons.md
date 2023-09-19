# Addons

## Existing Volumes

An existing volume addon can be provided to a metric. As an example, it would make sense to run an IO benchmarks with
different kinds of volume addons. The addons for volumes currently include:

 - a persistent volume claim (PVC) and persistent volume (PV) that you've created
 - a secret that you've created
 - a config map that you've created
 - a host volume (typically for testing)
 - an empty volume

and for all of the above, you want to create it and provide metadata for the addon to the operator, which will ensure the volume is available for your metric. We will provide examples here to do that.

#### persistent volume claim addon

As an example, here is how to provide the name of an existing claim (you created separately) to a metric container:
TODO add support to specify a specific metric container or replicated job container, if applicable.

```yaml
spec:
  metrics:
    - name: app-lammps
      addons:
        # This name is a unique identifier for this addon
        - name: volume-pvc
          options:
            name: data
            claimName: data
            path: /workflow
```

The above would add a claim named "data" to the metric container(s).

#### config map addon example

Here is an example of providing a config map to an application container In layman's terms, we are deploying vanilla nginx, but adding a configuration file
to `/etc/nginx/conf.d`

```yaml
spec:
  metrics:
    - name: app-lammps
      addons:
        # This name is a unique identifier for this addon
        - name: volume-cm
          options:
            name: nginx-conf
            configMapName: nginx-conf
            path: /etc/nginx/conf.d
          mapOptions:
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

#### secret addon example

Here is an example of providing an existing secret (in the metrics-operator namespace)
to the metric container(s):

```yaml
spec:
  metrics:
    - name: app-lammps
      addons:
        # This name is a unique identifier for this addon
        - name: volume-secret
          options:
            name: certs
            path: /etc/certs
            secretName: certs
```

The above shows an existing secret named "certs" that we will mount into `/etc/certs`.

#### hostpath volume addon example

Here is how to use a host path:

```yaml
spec:
  metrics:
    - name: app-lammps
      addons:
        # This name is a unique identifier for this addon
        - name: volume-hostpath
          options:
            name: data
            hostPath: /path/on/host
            path: /path/in/container
```


TODO convert to addon logic

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

#### workingDir

To add a working directory for your application:

```yaml
spec:
  application:
    image: ghcr.io/rse-ops/vanilla-lammps:tag-latest
    command: mpirun lmp -v x 1 -v y 1 -v z 1 -in in.reaxc.hns -nocite
    workingDir: /opt/lammps/examples/reaxff/HNS
```

#### volumes

An application is allowed to have one or more existing volumes. An existing volume can be any of the types described in [existing volumes](#existing-volumes)

#### resources

You can define resources for an application or a metric container. Known keys include "memory" and "cpu" (should be provided in some string format that can be parsed) and all others are considered some kind of quantity request.

```yaml
application:
  resources:
    memory: 500M
    cpu: 4
```

Metrics can also take resource requests.

```yaml
metrics:
  - name: io-fio
    resources:
      memory: 500M
      cpu: 4
```

If you wanted to, for example, request a GPU, that might look like:

```yaml
resources:
  limits:
    gpu-vendor.example/example-gpu: 1
```

Or for a particular type of networking fabric:

```yaml
resources:
  limits:
    vpc.amazonaws.com/efa: 1
```

Both limits and resources are flexible to accept a string or an integer value, and you'll get an error if you
provide something else. If you need something else, [let us know](https://github.com/converged-computing/metrics-operator/issues).
If you are requesting GPU, [this documentation](https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/) is helpful.

### storage

When you want to measure some storage performance, you'll want to add a "storage" section to your MetricSet. This will typically just be a reference to some existing storage (see [existing volumes](#existing-volumes)) that we want to measure, and can also be done for some number of completions and metrics for storage.

#### commands

If you need to add some special logic to create or cleanup for a storage volume, you are free to define them for storage in each of pre and post sections, which will happen before and after the metric runs, respectively.

```yaml
storage:
  volume:
    claimName: data 
    path: /data
  commands:
    pre: |
      apt-get update && apt-get install -y mymounter-tool
      mymounter-tool mount /data
    post: mymounter-tool unmount /data
    # Wrap the storage metric in this prefix
    prefix: myprefix
```

All of the above are strings. The pipe allows for multiple lines, if appropriate.
Note that while a "volume" is typical, you might have a storage setup that is done via a set of custom commands, in which case
you don't need to define the volume too.
