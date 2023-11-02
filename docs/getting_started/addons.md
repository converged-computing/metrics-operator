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

## Workload

### workload-flux

If you need to "throw in" Flux Framework into your container to use as a scheduler, you can do that with an addon!

> Yes, it's astounding. 🦩️

This works by way of the same trick that we use for other addons that have a complex (and/or large) install setup. We:

- Build the software into an isolated spack "copy" view
- The software is then (generally) at some `/opt/view` and `/opt/software`
- The flux container is added as a sidecar container to your pod for your replicated job
  - Additional setup / configuration is done here
- We can then create an empty volume that is shared by your metric or scaled application
- The entire tree is copied over into the empty volume
- When the copy is done, indicated by the final touch of a file, the updated container entrypoint is run
- This typically means we have taken your metric command, and wrapped it in a Flux submit.

It's really cool because it means you can run a metric / application with Flux without needing
to install it into your container to begin with. The one important detail is a matching of
general operating system. The current view uses rocky, however the image is customizable
(and we can provide other bases if/when requested). Here are the arguments you can customize
under the metric -> options.

| Name | Description | Type | Default |
|-----|-------------|------------|------|
| mount | Path to mount flux view in application container | string | /opt/share |
| tasks | Number of tasks `-n` to give to flux (not provided if not set) | string | unset |
| image | Customize the container image | string | `ghcr.io/rse-ops/spack-flux-rocky-view:tag-8` |
| fluxUser  | The flux user (currently not used, but TBA)  | string | flux |
| fluxUid  | The flux user ID (currently not used, but TBA)  | string | 1004 |
| interactive  | Run flux in interactive mode  | string | "false" |
| connectTimeout | How long zeroMQ should wait to retry | string | "5s" |
| quorum | The number of brokers to require before starting the cluster | string | (total brokers or pods) |
| debugZeroMQ | Turn on zeroMQ debugging | string | "false" |
| logLevel | Customize the flux log level | string | "6" |
| queuePolicy | Queue policy for flux to use | string | fcfs |
| workerLetter | The letter that the worker job is expected to have | string | w |
| launcherLetter | The letter that the launcher job is expected to have | string | w |
| workerIndex | The index of the replicated job for the worker | string | 0 |
| launcherIndex | The index of the replicated job for the launcher | string | 0 |
| preCommand | Pre-command logic to run in launcher/workers before flux is started (after setup in flux container) | string | unset |

Note that the number of pods for flux defaults to the number in your MetricSet, along 
with the namespace and service name.

**Important** the flux addon is currently supported for metric types that:

1. have the launcher / worker design (so the hostlist.txt is present in the PWD)
2. Have scp installed, as the shared certificate needs to be copied from the lead broker to all followers
3. Ideally have munge installed - we do try to install it (but better to already be there)

We also currently run flux as root. This is considered bad practice, but probably OK
for this early development work. We don't see a need to have shared namespace / operator
environments at this point, which is why I didn't add it.

## Performance

### perf-hpctoolkit

 - *[perf-hpctoolkit](https://github.com/converged-computing/metrics-operator/tree/main/examples/addons/hpctoolkit-lammps)*

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


### perf-mpitrace

 - *[perf-mpitrace](https://github.com/converged-computing/metrics-operator/tree/main/examples/addons/perf-mpitrace)*

This metric provides [mpitrace](https://github.com/IBM/mpitrace) to wrap an MPI application. The setup is the same as hpctoolkit, and we
currently only provide a rocky base (please let us know if you need another). It works by way of wrapping the mpirun command with `LD_PRELOAD`.
See the link above for an example that uses LAMMPS.

Here are the acceptable parameters.

| Name | Description | Type | Default |
|-----|-------------|------------|------|
| mount | Path to mount hpctoolview view in application container | string | /opt/share |
| image | Customize the container image | string | `ghcr.io/converged-computing/metric-mpitrace:rocky` |


