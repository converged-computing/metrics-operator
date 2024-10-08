# Metrics

The following metrics are under development (or being planned).

 - [Examples](https://converged-computing.github.io/metrics-operator/getting_started/metrics.html#examples)

Each metric can be ascribed to a high level family, shown in the table below as the "Family" column. 
We likely will tweak and improve upon these categories.

<iframe src="../_static/data/table.html" style="width:100%; height:900px;" frameBorder="0"></iframe>


## Implemented Metrics

### sys-hwloc

 - *[sys-hwloc](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/sys-hwloc)*

[Hwloc](https://www.open-mpi.org/projects/hwloc/) or "portable hardware locality" can be used to look at the hardware of your system.
There is a [nice tutorial here](https://www.open-mpi.org/projects/hwloc/tutorials/20120702-POA-hwloc-tutorial.html) for the default command that is run,
"lstopo" that does exactly that - listing your hardware topology! Specifically we output a png image and machine spec for the default command, and this can be updated.
[This man page](https://manpages.ubuntu.com/manpages/impish/man1/lstopo.1.html) is recommended to see the different commands and options.

| Name | Description | Type | Default |
|-----|-------------|------------|------|
| command | Change the default command to something else. | string | lstopo architecture.png && hwloc-ls machine.xml |

The above saves a png image, and the machine data to xml. Note that if you need to copy the data post-run, you
likely want to set `interactive: true` to keep it running.


### perf-sysstat

 - *[perf-hello-world](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/perf-hello-world)*
 - *[perf-lammps](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/perf-lammps)*

This metric provides the "pidstat" executable of the sysstat library. The following options are available:


| Name | Description | Type | Default |
|-----|-------------|------------|------|
| color | Set to turn on color parsing | Anything set | unset |
| pids | For debugging, show consistent output of ps aux | Anything set | unset |
| threads | add `-t` to each pidstat command to indicate wanting thread-level output | unset |
| completions | Number of times to run metric | int32 | unset (runs for lifetime of application or indefinitely) |
| rate | Seconds to pause between measurements | int32 | 10 |

By default color and pids are set to false anticipating log parsing.
And we also provide the option to see "commands" or specific commands based on a job index to the metric.
As an example, here is how we would ask to monitor two different commands for a launcher node (index 0)
and the rest (workers).

```yaml
- name: perf-sysstat
  options:
    pids: "true"

  # Custom options
  options:
    rate: 2

# Look for pids based on commands matched to index
  mapOptions:
    commands:
       # First set all to use the worker command, but give the lead broker a special command
       "all": /usr/libexec/flux/cmd/flux-broker --config /etc/flux/config -Scron.directory=/etc/flux/system/cron.d -Stbon.fanout
       "0": /usr/bin/python3.8 /usr/libexec/flux/cmd/flux-submit.py -n 2 --quiet --watch lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite
```

In the map above, order matters, as the command for all indices is first set to be the flux-broker one, and then
after the index at 0 gets a custom command. See [pidstat](https://man7.org/linux/man-pages/man1/pidstat.1.html) for
more information on this command, and [this file](https://github.com/converged-computing/metrics-operator/blob/main/pkg/metrics/perf/sysstat.go) 
for how we use them.  If there is an option or command that is not exposed that you would like, please [open an issue](https://github.com/converged-computing/metrics-operator/issues).


### io-fio

 - *[io-host-volume](https://github.com/converged-computing/metrics-operator/tree/main/examples/storage/google/io-fusion)*

This is a nice tool that you can simply point at a path, and it measures IO stats by way of writing a file there! 
Options you can set include:

|Name | Description | Type | Default |
|-----|-------------|------------|------|
| testname | Name for the test | string | test |
| blocksize | Size of block to write. It defaults to 4k, but can be set from 256 to 8k.  | string | 4k |
| iodepth | Number of I/O units to keep in flight against the file. | int | 64 |
| size | Total size of file to write | string | 4G |
| directory | Directory (usually mounted) to test. | string | /tmp |
| pre | Custom logic / command to run before Fio | string | unset |
| post | Custom logic / command to run after Fio (e.g., cleanup) | string | unset |
| prefix | Prefix to add to running fio commands (like a wrapper) | string | unset |

For the "directory" we use this location to write a temporary file, which will be cleaned up.
This allows for testing storage mounted from multiple metric pods without worrying about a name conflict.

### io-ior

 - *[io-host-volume](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/io-ior)*

![img/ior.jpeg](img/ior.jpeg)

[Ior](https://github.com/hpc/ior) is a really nice IO tool that is now a combination of its previous self and the [mdtest](https://github.com/llnl/mdtest) 
tool. We expose a simple set of the working directory and command that you want to run, and the rest is up to you!

|Name | Description | Type | Default |
|-----|-------------|------------|------|
| command | The default ior command  |string | ior -w -r -o testfile |
| workdir | The working directory for the command | string | /opt/ior |

The [getting started](https://ior.readthedocs.io/en/latest/userDoc/tutorial.html) tutorial is great for seeing how
basic commands are done. Note that the container does have mpirun if you want to use it. We don't have support
for this across nodes, but this could be added. [Let us know](https://github.com/converged-computing/metrics-operator/issues) 
if this would be interesting to you.

### io-sysstat

 - *[io-host-volume](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/io-host-volume)*

This is the "iostat" executable of the sysstat library.

|Name | Description | Type | Default |
|-----|-------------|------------|------|
| human | Show tabular, human-readable output inside of json | string "true" or "false" | "false" |
| completions | Number of times to run metric | int32 | unset (runs for lifetime of application or indefinitely) |
| rate | Seconds to pause between measurements | int32 | 10 |
| pre | One or more commands to run before iostat | string | unset |
| post | One or more commands to run after iostat | string | unset |

This is good for mounted storage that can be seen by the operating system, but may not work for something like NFS.

### dlio

While this is a simple performance tool not coded into the Metrics Operator (it is installed on the fly to your container with pip and you minimally require hwloc)
it generates pretty cool data that can be visualized with perfetto!

 - *[dlio](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/dlio)*

You can see the full example above. It is just installing a library with pip, and then ensuring the tool `LD_PRELOAD`
is set as the prefix. I added sleep infinity to the end to copy over output data at the end.

### network-netmark

 - *[network-netmark](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/network-netmark)* (code still private)

This is currently a private container/software, but we have support for it when it's ready to be made public (networking)
Variables to customize include:

|Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| tasks | Total number of tasks across pods | options->tasks | string | nproc * pods |
| warmups | Number of warmups | options->warmups | int32 | 10 |
| trials | Number of trials | options->trials | int32 | 20 |
| sendReceiveCycles | Number of send-receive cycles | options-sendReceiveCycles | int32 | 20 |
| messageSize | Message size in bytes | options->messageSize | int32 | 0 |
| storeEachTrial | Flag to indicate storing each trial data | options->storeEachTrial | string (true/false) | "true" |
| soleTenancy | Turn off sole tenancy (one pod/node) | options->soleTenancy | string ("false" or "no") | "true" |

### network-osu-benchmark

 - *[network-osu-benchmark](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/network-osu-benchmark)*

Point to point benchmarks for MPI (networking). If listOptions->commands not set, will use all one-point commands.
Variables to customize include:

|Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| commands | Custom list of osu-benchmark one-sided commands to run | listOptions->commands | array | unset uses default set |
| soleTenancy | Turn off sole tenancy (one pod/node) | string ("false" or "no") | "true" |
| all | Run ALL benchmarks with defaults | string ("true" or "yes") | "false" |
| flags | Overwrite defaults flags (experts only!)| string | Defaults to an ideal set per metric (see [osu-benchmark.go](https://github.com/converged-computing/metrics-operator/blob/main/pkg/metrics/network/osu-benchmark.go))|
| timed | String "true" or "yes" to add time prefix to mpirun (for debugging, etc) | string | "false" |
| sleep | Number of seconds to sleep to wait for network to be ready | int32 | 60 |

By default, we run a subset of commands:

- osu_get_acc_latency
- osu_acc_latency
- osu_fop_latency
-	osu_get_latency
-	osu_put_latency
-	osu_allreduce
- osu_latency
- osu_bibw
-	osu_bw

However all of the following are available for MPI

<details>

<summary>Commands available for OSU Benchmarks</summary>

```console
.
|-- collective
|   |-- osu_allgather
|   |-- osu_allgatherv
|   |-- osu_allreduce
|   |-- osu_alltoall
|   |-- osu_alltoallv
|   |-- osu_barrier
|   |-- osu_bcast
|   |-- osu_gather
|   |-- osu_gatherv
|   |-- osu_iallgather
|   |-- osu_iallgatherv
|   |-- osu_iallreduce
|   |-- osu_ialltoall
|   |-- osu_ialltoallv
|   |-- osu_ialltoallw
|   |-- osu_ibarrier
|   |-- osu_ibcast
|   |-- osu_igather
|   |-- osu_igatherv
|   |-- osu_ireduce
|   |-- osu_iscatter
|   |-- osu_iscatterv
|   |-- osu_reduce
|   |-- osu_reduce_scatter
|   |-- osu_scatter
|   `-- osu_scatterv
|-- one-sided
|   |-- osu_acc_latency
|   |-- osu_cas_latency
|   |-- osu_fop_latency
|   |-- osu_get_acc_latency
|   |-- osu_get_bw
|   |-- osu_get_latency
|   |-- osu_put_bibw
|   |-- osu_put_bw
|   `-- osu_put_latency
|-- pt2pt
|   |-- osu_bibw
|   |-- osu_bw
|   |-- osu_latency
|   |-- osu_latency_mp
|   |-- osu_latency_mt
|   |-- osu_mbw_mr
|   `-- osu_multi_lat
`-- startup
    |-- osu_hello
    `-- osu_init
```

</details>

Note that not all of these have been tested on our setups, so
if you have any questions please [let us know](https://github.com/converged-computing/metrics-operator/issues).
Here are some useful resources for the benchmarks:

 - [HPC Council](https://hpcadvisorycouncil.atlassian.net/wiki/spaces/HPCWORKS/pages/1284538459/OSU+Benchmark+Tuning+for+2nd+Gen+AMD+EPYC+using+HDR+InfiniBand+over+HPC-X+MPI)
 - [AWS Tutorials](https://www.hpcworkshops.com/08-efa/04-complie-run-osu.html)

### app-custom

A custom application can support any application to be used as a metric app. For the following parameters, "command" and "container" are required.

| Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| command | The full mpirun command | options->command |string | unset |
| workdir | The working directory for the command | options->workdir | string | unset |
| soleTenancy | require each pod to have sole tenancy | command->soleTenancy | string | "false" |

As an example, here is running mpitrace (an addon) with a custom container.

```yaml
apiVersion: flux-framework.org/v1alpha2
kind: MetricSet
metadata:
  labels:
    app.kubernetes.io/name: metricset
    app.kubernetes.io/instance: metricset-sample
  name: metricset-sample
spec:
  # Number of pods for lammps (one launcher, the rest workers)
  pods: 4
  metrics:
   - name: app-custom
     image: ghcr.io/converged-computing/<your-container>
     options:
       command: mpirun --hostfile ./hostlist.txt -mca orte_keep_fqdn_hostnames t -np 4 --map-by socket <app> <options>
       workdir: <workdir>

     # Add on hpctoolkit, will mount a volume and wrap lammps
     addons:
       - name: perf-mpitrace
         options:
           mount: /opt/mnt
           image: ghcr.io/converged-computing/metric-mpitrace:ubuntu-jammy
           workdir: <workdir>
           # this is the target of the replicated job "l" means launcher
           target: l          
           # This is the target container, with full name "launcher"
           containerTarget: launcher
```

### app-lammps

 - *[app-lammps](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-lammps)*

Since we were using LAMMPS so often as a benchmark (and testing timing of a network) it made sense to add it here
as a standalone metric! Although we are doing MPI with communication via SSH, this can still serve as a means
to assess performance.

| Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| command | The full mpirun and lammps command | options->command |string | (see below) |
| workdir | The working directory for the command | options->workdir | string | /opt/lammps/examples/reaxff/HNS# |
| soleTenancy | require each pod to have sole tenancy | command->soleTenancy | string | "false" |

For inspection, you can see all the examples provided [in the LAMMPS GitHub repository](https://github.com/lammps/lammps/tree/develop/examples).
The default command (if you don't change it) intended as an example is:

```bash
mpirun --hostfile ./hostlist.txt -np 2 --map-by socket lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite(e
```

In the working directory `/opt/lammps/examples/reaxff/HNS#`. You should be calling `mpirun` and expecting a ./hostlist.txt in the present working directory (the "workdir" you chose above).
You should also provide the correct number of processes (np) and problem size for LAMMPS (lmp). We left this as open and flexible
anticipating that you as a user would want total control.

### app-amg

 - *[app-amg](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-amg)*

AMG means "algebraic multi-grid" and it's easy to confuse with the company [AMD](https://www.amd.com/en/solutions/supercomputing-and-hpc) "Advanced Micro Devices" ! From [the guide](https://asc.llnl.gov/sites/asc/files/2020-09/AMG_Summary_v1_7.pdf):

> AMG is a parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids. The driver provided for this benchmark builds linear systems for a 3D problem with a 27-point stencil and generates two different problems that are described in section D of the AMG.readme file in the docs directory.

Here are examples of small and medium problem sizes provided in that same guide. Each of these would be given to MPI (mpirun), but srun is provided as an example instead.

```console
# Small size problems
srun –N 32 –n 512 amg –problem 1 –n 96 96 96 –P 8 8 8
srun –N 32 –n 512 amg –problem 2 –n 40 40 40 –P 8 8 8
srun –N 64 –n 1024 amg –problem 1 –n 96 96 96 –P 16 8 8
srun –N 64 –n 1024 amg –problem 2 –n 40 40 40 –P 16 8 8

# Medium size problems
srun –N 512 –n 8192 amg –problem 1 –n 96 96 96 –P 32 16 16
srun –N 512 –n 8192 amg –problem 2 –n 40 40 40 –P 32 16 16
srun –N 1024 –n 16384 amg –problem 1 –n 96 96 96 –P 32 32 16
srun –N 1024 –n 16384 amg –problem 2 –n 40 40 40 –P 32 32 16
```

By default, akin to LAMMPS we expose the entire mpirun command along with the working directory for you to adjust.

| Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| command | The amg command (without mpirun) | options->command |string | (see below) |
| prefix | The prefix (mpirun command and arguments) | options->mpirun | string | (see below) | 
| workdir | The working directory for the command | options->workdir | string | /opt/AMG |

By default, when not set, you will just run the amg binary to get a test case run: 

```bash
# mpirun
mpirun --hostfile ./hostlist.txt

# command
amg

# Assembled into
mpirun --hostfile ./hostlist.txt ./problem.sh
```

More likely you want an actual problem size on a specific number of node and tasks, and you'll want to test this. The two problem sizes include:

 - *problem 1* (default) will use conjugate gradient preconditioned with AMG to solve a linear system with a 3D 27-point stencil of size nx*ny*nz*Px*Py*Pz.
 - *problem 2* simulates a time-dependent problem of size nx*ny*nz*Px*Py*Pz with AMG-GMRES. The linear system is also a 3D 27-point stencil. The system is sized to be 5-10% of the large problem.

**NOTE** that the Python parser was written for the test case, and likely we need to extend it to problem 2 or larger sized problems. If you
run a larger problem and the parser does not work as expected, please [send us the output](https://github.com/converged-computing/metrics-operator/issues) and we will provide an updated parser.
See [this guide](https://asc.llnl.gov/sites/asc/files/2020-09/AMG_Summary_v1_7.pdf) for more detail.

### app-cabanaPIC

 - *[app-cabanaPIC](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-cabanaPIC)*

This is a [particle in cell](https://github.com/ECP-copa/CabanaPIC/issues/3) simulation that is experimental because it does not
seem to support multiple nodes yet (but should).

| Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| command | The full command to run | options->command | string | cbnpic |
| workdir | The working directory for the command | options->workdir | string | /opt/cabanaPIC/build |

### app-quicksilver

 - *[app-quicksilver](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-quicksilver)*

Quicksilver is a proxy app for Monte Carlo simulation code. You can learn more about it on the [GitHub repository](https://github.com/LLNL/Quicksilver/).
By default, akin to other apps we expose the entire mpirun command along with the working directory for you to adjust.

| Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| command | The qs command (without mpirun) | options->command |string | (see below) |
| prefix | The prefix (mpirun command and arguments) | options->mpirun | string | (see below) | 
| workdir | The working directory for the command | options->workdir | string | /opt/AMG |

By default, when not set, you will just run the qs (quicksilver) binary on a sample problem, represented by an input text file: 

```bash
# mpirun
mpirun --hostfile ./hostlist.txt

# command
qs /opt/quicksilver/Examples/CORAL2_Benchmark/Problem1/Coral2_P1.inp

# Assembled into problem.sh as follows:
mpirun --hostfile ./hostlist.txt ./problem.sh
```

There are many problems that come in the container, and here are the fullpaths:

```console
# Example command
qs /opt/quicksilver/Examples/CORAL2_Benchmark/Problem1/Coral2_P1.inp

# All examples:
/opt/quicksilver/Examples/AllScattering/scatteringOnly.inp
/opt/quicksilver/Examples/NoCollisions/no.collisions.inp
/opt/quicksilver/Examples/NonFlatXC/NonFlatXC.inp
/opt/quicksilver/Examples/CORAL2_Benchmark/Problem2/Coral2_P2_4096.inp
/opt/quicksilver/Examples/CORAL2_Benchmark/Problem2/Coral2_P2.inp
/opt/quicksilver/Examples/CORAL2_Benchmark/Problem2/Coral2_P2_1.inp
/opt/quicksilver/Examples/CORAL2_Benchmark/Problem1/Coral2_P1.inp
/opt/quicksilver/Examples/CORAL2_Benchmark/Problem1/Coral2_P1_1.inp
/opt/quicksilver/Examples/CORAL2_Benchmark/Problem1/Coral2_P1_4096.inp
/opt/quicksilver/Examples/CTS2_Benchmark/CTS2.inp
/opt/quicksilver/Examples/CTS2_Benchmark/CTS2_36.inp
/opt/quicksilver/Examples/CTS2_Benchmark/CTS2_1.inp
/opt/quicksilver/Examples/AllAbsorb/allAbsorb.inp
/opt/quicksilver/Examples/Homogeneous/homogeneousProblem_v4_ts.inp
/opt/quicksilver/Examples/Homogeneous/homogeneousProblem_v5_ts.inp
/opt/quicksilver/Examples/Homogeneous/homogeneousProblem.inp
/opt/quicksilver/Examples/Homogeneous/homogeneousProblem_v3_wq.inp
/opt/quicksilver/Examples/Homogeneous/homogeneousProblem_v7_ts.inp
/opt/quicksilver/Examples/Homogeneous/homogeneousProblem_v4_tm.inp
/opt/quicksilver/Examples/Homogeneous/homogeneousProblem_v3.inp
/opt/quicksilver/Examples/AllEscape/allEscape.inp
/opt/quicksilver/Examples/NoFission/noFission.inp
```

You can also look more closely in the [GitHub repository](https://github.com/LLNL/Quicksilver/tree/master/Examples).

### app-pennant

 - *[app-pennant](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-pennant)*

Pennant is an unstructured mesh hydrodynamics for advanced architectures. The documentation is sparse, but you
can find the [source code on GitHub](https://github.com/llnl/pennant). 
By default, akin to other apps we expose the entire mpirun prefix and command along with the working directory for you to adjust.

| Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| command | The pennant command (without mpirun) | options->command |string | (see below) |
| prefix | The prefix (mpirun command and arguments) | options->mpirun | string | (see below) | 
| workdir | The working directory for the command | options->workdir | string | /opt/AMG |

By default, when not set, you will just run pennant on a test problem, represented by an input text file: 

```bash
# mpirun
mpirun --hostfile ./hostlist.txt

# command
pennant /opt/pennant/test/sedovsmall/sedovsmall.pnt

# Assembled into problem.sh as follows:
mpirun --hostfile ./hostlist.txt ./problem.sh
```

There are many input files that come in the container, and here are the fullpaths in `/opt/pennant/test`:

<details>

<summary>Input files available to pennant</summary>

```console
|-- leblanc
|   |-- leblanc.pnt
|   |-- leblanc.xy.std
|   `-- leblanc.xy.std4
|-- leblancbig
|   `-- leblancbig.pnt
|-- leblancx16
|   `-- leblancx16.pnt
|-- leblancx4
|   `-- leblancx4.pnt
|-- leblancx48
|   `-- leblancx48.pnt
|-- leblancx64
|   `-- leblancx64.pnt
|-- noh
|   |-- noh.pnt
|   |-- noh.xy.std
|   `-- noh.xy.std4
|-- nohpoly
|   `-- nohpoly.pnt
|-- nohsmall
|   |-- nohsmall.pnt
|   |-- nohsmall.xy.std
|   `-- nohsmall.xy.std4
|-- nohsquare
|   `-- nohsquare.pnt
|-- sample_outputs
|   |-- edison
|   |   |-- leblancbig.thr1.out
|   |   |-- leblancx16.thr1024.out
|   |   |-- leblancx4.thr16.out
|   |   |-- leblancx64.mpi2048.out
|   |   `-- nohpoly.thr1.out
|   `-- vulcan
|       |-- leblancx16.out
|       |-- leblancx48.out
|       |-- sedovflat.out
|       |-- sedovflatx16.out
|       |-- sedovflatx4.out
|       `-- sedovflatx40.out
|-- sedov
|   |-- sedov.pnt
|   |-- sedov.xy.std
|   `-- sedov.xy.std4
|-- sedovbig
|   `-- sedovbig.pnt
|-- sedovflat
|   `-- sedovflat.pnt
|-- sedovflatx120
|   `-- sedovflatx120.pnt
|-- sedovflatx16
|   `-- sedovflatx16.pnt
|-- sedovflatx4
|   `-- sedovflatx4.pnt
|-- sedovflatx40
|   `-- sedovflatx40.pnt
`-- sedovsmall
    |-- sedovsmall.pnt
    |-- sedovsmall.xy
    |-- sedovsmall.xy.std
    `-- sedovsmall.xy.std4
```

</details>

And likely you will need to adjust the mpirun parameters, etc.

### app-kripke

 - *[app-kripke](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-kripke)*

[Kripke](https://github.com/LLNL/Kripke) is (from the README):

> Kripke is a simple, scalable, 3D Sn deterministic particle transport code. Its primary purpose is to research how data layout, programming paradigms and architectures effect the implementation and performance of Sn transport. A main goal of Kripke is investigating how different data-layouts affect instruction, thread and task level parallelism, and what the implications are on overall solver performance. 

Akin to AMG, we allow you to modify each of the mpirun and kripke commands via:

| Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| command | The amg command (without mpirun) | options->command |string | (see below) |
| prefix | The prefix (mpirun command and arguments) | options->mpirun | string | (see below) | 
| workdir | The working directory for the command | options->workdir | string | /opt/AMG |

By default, when not set, you will just run the kripke binary to get a test case run, so mpirun is set to be blank.

```bash
# mpirun is blank
""
# But could be an actual mpirun command
mpirun --hostfile ./hostlist.txt

# command written to problem.sh
kripke

# Assembled into
mpirun --hostfile ./hostlist.txt ./problem.sh
```

There is a nice [guide here](https://asc.llnl.gov/sites/asc/files/2020-09/Kripke_Summary_v1.2.2-CORAL2_0.pdf) that can help you to decide
on your specific command or problem size. Also note that we expose the following executables built with it:

```console
ex1_vector-addition            ex4_atomic-histogram                ex7_nested-loop-reorder
ex1_vector-addition_solution   ex4_atomic-histogram_solution       ex7_nested-loop-reorder_solution
ex2_approx-pi                  ex5_line-of-sight                   ex8_tiled-matrix-transpose
ex2_approx-pi_solution         ex5_line-of-sight_solution          ex8_tiled-matrix-transpose_solution
ex3_colored-indexset           ex6_stencil-offset-layout           ex9_matrix-transpose-local-array
ex3_colored-indexset_solution  ex6_stencil-offset-layout_solution  ex9_matrix-transpose-local-array_solution
```
(meaning on the PATH in `/opt/Kripke/build/bin` in the container).
For apps / metrics to be added, please see [this issue](https://github.com/converged-computing/metrics-operator/issues/30).

### app-ldms

 - *[app-ldms](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-ldms)*


LDMS is "a low-overhead, low-latency framework for collecting, transferring, and storing metric data on a large distributed computer system" 
and is packaged alongside [ovis-hpc](https://github.com/ovis-hpc/ovis). While there are complex aggregator setups we could run,
for this simple metric we simply run (on each separate pod/node). The following variables are supported:

|Name | Description | Type | Default |
|-----|-------------|------|------|
| command | The command to issue to ldms_ls (or that) |string | (see below) |
| workdir | The working directory for the command |  string | /opt |
| completions | Number of times to run metric | int32 | unset (runs for lifetime of application or indefinitely) |
| rate | Seconds to pause between measurements | int32 | 10 |


The following is the default command:

```bash
ldms_ls -h localhost -x sock -p 10444 -l -v
```

### app-nekbone

 - *[app-nekbone](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-nekbone)*

Nekbone comes with a set of example that primarily depend on you choosing the correct workikng directory and command to run from.
You can do this via these primary two commands:

| Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| command | The full mpirun and nekbone command | options->command |string | (see below) |
| workdir | The working directory for the command | options->workdir | string | /root/nekbone-3.0/test |

And the following combinations are supported. Note that example1 did not build, and example2 is the default (if you don't set these variables).

| Command | Workdir |
|---------|---------|
| mpiexec --hostfile ./hostlist.txt -np 2 ./nekbone | /root/nekbone-3.0/test/example2 |
| mpiexec --hostfile ./hostlist.txt -np 2 ./nekbone | /root/nekbone-3.0/test/example3 |
| mpiexec --hostfile ./hostlist.txt -np 2 ./nekbone | /root/nekbone-3.0/test/nek_comm |
| mpiexec --hostfile ./hostlist.txt -np 2 ./nekbone | /root/nekbone-3.0/test/nek_mgrid |
| mpiexec --hostfile ./hostlist.txt -np 2 ./nekbone | /root/nekbone-3.0/test/nek_delay |

You can see the archived repository [here](https://github.com/Nek5000/Nekbone). If there are interesting metrics in this
project it would be worth bringing it back to life I think.

### app-laghos

 - *[app-laghos](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-laghos)*

From the [Laghos README](https://github.com/CEED/Laghos):

> Laghos (LAGrangian High-Order Solver) is a miniapp that solves the time-dependent Euler equations of compressible gas dynamics in a moving Lagrangian frame using unstructured high-order finite element spatial discretization and explicit high-order time-stepping.

Akin to other apps, you can customize the command and workdir. Note that the `laghos` executable is at `/workflow/laghos` and not on
the path, so the default references it as `./laghos`.

| Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| command | The full mpirun and laghos command | options->command |string | (see below) |
| workdir | The working directory for the command | options->workdir | string | /workdir/laghos |

### app-bdas

 - *[app-bdas](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-bdas)*

BDAS standards for "Big Data Analysis Suite" and you can read more about it [here](https://asc.llnl.gov/sites/asc/files/2020-09/BDAS_Summary_b4bcf27_0.pdf).
The container has machine learning analyses (provided in R) that work with MPI (openmpi),
The benchmarks are in `/opt/bdas/benchmarks/r` in the container, and we provide an example for princomp
(see default command below):

| Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| command | The full mpirun and Rscript command | options->command |string | (see below) |
| workdir | The working directory for the command | options->workdir | string | /opt/bdas/benchmarks/r |

```console
# This is the default command. You must target the --hostfile and use the allow as root flag!
mpirun --allow-run-as-root -np 4 --hostfile ./hostlist.txt Rscript /opt/bdas/benchmarks/r/princomp.r 250 50
```
Try setting the logging->interactive: true option in the spec to keep the container running and explore other benchmarks.
These are the ones I've tried:

```console
# This is the default command. You must target the --hostfile and use the allow as root flag!
mpirun --allow-run-as-root -np 4 --hostfile ./hostlist.txt Rscript /opt/bdas/benchmarks/r/kmeans.r 250 50
mpirun --allow-run-as-root -np 4 --hostfile ./hostlist.txt Rscript /opt/bdas/benchmarks/r/svm.r 250 50
```

## Containers

To see all associated app containers, look at the [converged-computing/metrics-container](https://github.com/converged-computing/metrics-containers)
repository (with `Dockerfile`s  and automation) and associated packages.
