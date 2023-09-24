# Addons

An addon is a generic way to customize a metric. An addon can do everything from:

- generating new application or sidecar containers
- adding volumes, including writing new config maps
- customizing entrypoints, to the granularity of a jobset or a container

And as an example, if you wanted to use an IO benchmark, you would test that against different storage
solutions by way of using a volume added. The different groups available are discussed below, and if you
have a request for an addon please [let us know](https://github.com/converged-computing/metrics-operator/issues). 

<iframe src="../_static/data/addons.html" style="width:100%; height:500px;" frameBorder="0"></iframe>

## Command Addons

The Commands group of addons are some of my favorites, because they allow you to customize entrypoints for existing metrics! 

### Commands

> Use addon with name "commands"

The basic "commands" addon allows you to customize:

 - **preBlock**: A custom block of commands to run before the primary entrypoint command.
 - **prefix**: a wrapping prefix to the entrypoint
 - **suffix**: a wrapping suffix to the entrypoint
 - **postBlock**: a block of commands to run after the primary entrypoint command.

For example, you might want to time something by adding "time" as the prefix. You may want to
install something special to the container (or otherwise customize files or content) before running
the entrypoint. You might also want to run some kind of cleanup or save in the postBlock. The
reason "commands" is so cool is because it's flexible to so many ideas! Here is an example:

 - *[metrics-time.yaml](https://github.com/converged-computing/metrics-operator/tree/main/examples/addons/commands/metrics-time.yaml)*

### Perf

> Use addon with name "perf-commands"

Per commands has the same arguments as [commands](#commands) above, but will additionally add CAP_PTRACE and CAP_SYSADMIN
to your container, which are typically needed for performance benchmarking tools. As an example here, you might
install a performance tool in the preBlock, run it using the "prefix" and then use "suffix" optionally to pipe to
a file, and postBlock to upload somewhere.


## Existing Volumes

An existing volume addon can be provided to a metric. As an example, it would make sense to run an IO benchmarks with
different kinds of volume addons. The addons for volumes currently include:

 - a persistent volume claim (PVC) and persistent volume (PV) that you've created
 - a secret that you've created
 - a config map that you've created
 - a host volume (typically for testing)
 - an empty volume

and for all of the above, you want to create it and provide metadata for the addon to the operator, which will ensure the volume is available for your metric. We will provide examples here to do that.

### persistent volume claim addon

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

### config map addon example

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

### secret addon example

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

### hostpath volume addon example

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

**Note that we have support for a custom application container, but haven't written any good examples yet!**

## Performance

### perf-hpctoolkit

 - *[perf-hpctoolkit](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/perf-lammps-hpctoolkit)*

This metric provides [HPCToolkit](https://gitlab.com/hpctoolkit/hpctoolkit) for your application to use. This is the first metric of its type
to use a shared volume approach. Specifically, we:

- add a new ability for an application metric to define an empty volume, and have the metrics container copy stuff to it
- also add an ability for this kind of application metric to customize the application entrypoint (e.g., copy volume contents to destinations)
- build a spack copy view into the [hpctoolkit metrics container](https://github.com/converged-computing/metrics-containers/blob/main/hpctoolkit-containerize/Dockerfile)
- move the `/opt/software` and `/opt/views/view` roots into the application container, this is a modular install of HPCToolkit.
- copy over `/opt/share/software` (provided via the shared empty volume) to `/opt/software`` where spack expects it. We also add `/opt/share/view/bin` to the path (where hpcrun is)

After those steps are done, HPCToolkit is essentially installed, on the fly, in the application container. Since the `hpcrun` command is using `LD_AUDIT` we need
all libraries to be in the same system (the shared process namespace would not work). We can then run it, and generate a database. Also note that by default,
we run the post-analysis steps (shown below) and also provide them in each container as `post-run.sh`, which the addon will run for you, unless you
set `postAnalysis` to "false." Finally, if you need to run it manually, here is an example
given `hpctoolkit-lmp-measurements` in the present working directory of the container.


```bash
hpcstruct hpctoolkit-lmp-measurements

# Run "the professor!" 🤓️
hpcprof hpctoolkit-lmp-measurements
```

The above generates a database, `hpctoolkit-lmp-database` that you can copy to your machine for further interaction with hpcviewer
(or some future tool that doesn't use Java)!

```bash
kubectl cp -c app metricset-sample-m-0-npbc9:/opt/lammps/examples/reaxff/HNS/hpctoolkit-lmp-database hpctoolkit-lmp-database
hpcviewer ./hpctoolkit-lmp-database
```

Here are the acceptable parameters.

| Name | Description | Type | Default |
|-----|-------------|------------|------|
| mount | Path to mount hpctoolview view in application container | string | /opt/share |
| events | Events for hpctoolkit | string |  `-e IO` |
| image | Customize the container image | string | `ghcr.io/converged-computing/metric-hpctoolkit-view:ubuntu` |
| output | The output directory for hpcrun (database will generate to *-database) | string | hpctoolkit-result |

Note that for image we also provide a rocky build base, `ghcr.io/converged-computing/metric-hpctoolkit-view:rocky`. 
You can also see events available with `hpcrun -L`, and use the container for this metric.
There is a brief listing on [this page](https://hpc.llnl.gov/software/development-environment-software/hpc-toolkit).
We recommend that you do not pair hpctoolkit with another metric, primarily because it is customizing the application
entrypoint. If you add a process-namespace based metric, you likely need to account for the hpcrun command being the
wrapper to the actual executable.
