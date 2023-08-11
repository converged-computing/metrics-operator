# LAMMPS Example

This is an example of a metric app, lammps, which is part of the [coral 2 benchmarks](https://asc.llnl.gov/coral-2-benchmarks) and technically
isn't a metric, but we use it often to assess time and thus the MPI latency. A Python example (parsing the output data)
is provided in [python/app-lammps](../../python/app-lammps).

## Usage

Create a cluster and install JobSet to it.

```bash
kind create cluster
VERSION=v0.2.0
kubectl apply --server-side -f https://github.com/kubernetes-sigs/jobset/releases/download/$VERSION/manifests.yaml
```

Install the operator (from the development manifest here):

```bash
$ kubectl apply -f ../../dist/metrics-operator-dev.yaml
```

How to see metrics operator logs:

```bash
$ kubectl logs -n metrics-system metrics-controller-manager-859c66464c-7rpbw
```

Then create the metrics set. This is going to run a single run of LAMMPS over MPI!
as lammps runs.

```bash
kubectl apply -f metrics.yaml
```

Wait until you see pods created by the job and then running (there should be two - a launcher and worker for LAMMPS):

```bash
kubectl get pods
```
```diff
NAME                           READY   STATUS    RESTARTS   AGE
metricset-sample-l-0-0-lt782   1/1     Running   0          3s
metricset-sample-w-0-0-4s5p9   1/1     Running   0          3s
```

In the above, "l" is a launcher pod, and "w" is a worker node.
If you inspect the log for the launcher you'll see a short sleep (the network isn't up immediately)
and then LAMMPS running, and the log is printed to the console.

```bash
kubectl logs metricset-sample-l-0-0-lt782 -f
```
```console
METADATA START {"pods":2,"completions":2,"metricName":"app-lammps","metricDescription":"LAMMPS molecular dynamic simulation","metricType":"standalone","metricOptions":{"command":"mpirun --hostfile ./hostlist.txt -np 2 --map-by socket lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite","completions":0,"rate":10,"workdir":"/opt/lammps/examples/reaxff/HNS"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
LAMMPS (29 Sep 2021 - Update 2)
OMP_NUM_THREADS environment is not set. Defaulting to 1 thread. (src/comm.cpp:98)
  using 1 OpenMP thread(s) per MPI task
Reading data file ...
  triclinic box = (0.0000000 0.0000000 0.0000000) to (22.326000 11.141200 13.778966) with tilt (0.0000000 -5.0260300 0.0000000)
  2 by 1 by 1 MPI processor grid
  reading atoms ...
  304 atoms
  reading velocities ...
  304 velocities
  read_data CPU = 0.003 seconds
Replicating atoms ...
  triclinic box = (0.0000000 0.0000000 0.0000000) to (44.652000 22.282400 27.557932) with tilt (0.0000000 -10.052060 0.0000000)
  2 by 1 by 1 MPI processor grid
  bounding box image = (0 -1 -1) to (0 1 1)
  bounding box extra memory = 0.03 MB
  average # of replicas added to proc = 5.00 out of 8 (62.50%)
  2432 atoms
  replicate CPU = 0.001 seconds
Neighbor list info ...
  update every 20 steps, delay 0 steps, check no
  max neighbors/atom: 2000, page size: 100000
  master list distance cutoff = 11
  ghost atom cutoff = 11
  binsize = 5.5, bins = 10 5 6
  2 neighbor lists, perpetual/occasional/extra = 2 0 0
  (1) pair reax/c, perpetual
      attributes: half, newton off, ghost
      pair build: half/bin/newtoff/ghost
      stencil: full/ghost/bin/3d
      bin: standard
  (2) fix qeq/reax, perpetual, copy from (1)
      attributes: half, newton off, ghost
      pair build: copy
      stencil: none
      bin: none
Setting up Verlet run ...
  Unit style    : real
  Current step  : 0
  Time step     : 0.1
Per MPI rank memory allocation (min/avg/max) = 143.9 | 143.9 | 143.9 Mbytes
Step Temp PotEng Press E_vdwl E_coul Volume 
       0          300   -113.27833    437.52118   -111.57687   -1.7014647    27418.867 
      10    299.38517   -113.27631    1439.2824   -111.57492   -1.7013813    27418.867 
      20    300.27107   -113.27884     3764.342   -111.57762   -1.7012247    27418.867 
      30    302.21063   -113.28428    7007.6629   -111.58335   -1.7009363    27418.867 
      40    303.52265   -113.28799    9844.8245   -111.58747   -1.7005186    27418.867 
      50    301.87059   -113.28324    9663.0973   -111.58318   -1.7000523    27418.867 
      60    296.67807   -113.26777    7273.8119   -111.56815   -1.6996137    27418.867 
      70    292.19999   -113.25435    5533.5522   -111.55514   -1.6992158    27418.867 
      80    293.58677   -113.25831    5993.4438   -111.55946   -1.6988533    27418.867 
      90    300.62635   -113.27925    7202.8369   -111.58069   -1.6985592    27418.867 
     100    305.38276   -113.29357    10085.805   -111.59518   -1.6983874    27418.867 
Loop time of 11.8719 on 2 procs for 100 steps with 2432 atoms

Performance: 0.073 ns/day, 329.776 hours/ns, 8.423 timesteps/s
96.4% CPU use with 2 MPI tasks x 1 OpenMP threads

MPI task timing breakdown:
Section |  min time  |  avg time  |  max time  |%varavg| %total
---------------------------------------------------------------
Pair    | 8.1642     | 8.465      | 8.7659     |  10.3 | 71.30
Neigh   | 0.18776    | 0.18784    | 0.18792    |   0.0 |  1.58
Comm    | 0.022673   | 0.32358    | 0.62449    |  52.9 |  2.73
Output  | 0.0011288  | 0.0011759  | 0.001223   |   0.1 |  0.01
Modify  | 2.8935     | 2.8935     | 2.8935     |   0.0 | 24.37
Other   |            | 0.0007992  |            |       |  0.01

Nlocal:        1216.00 ave        1216 max        1216 min
Histogram: 2 0 0 0 0 0 0 0 0 0
Nghost:        7591.50 ave        7597 max        7586 min
Histogram: 1 0 0 0 0 0 0 0 0 1
Neighs:        432912.0 ave      432942 max      432882 min
Histogram: 1 0 0 0 0 0 0 0 0 1

Total # of neighbors = 865824
Ave neighs/atom = 356.01316
Neighbor list builds = 5
Dangerous builds not checked
Total wall time: 0:00:12
METRICS OPERATOR COLLECTION END
```

The above also shows the structured output that is done in a way for our Python parsing script to easily
find sections of data. Also note that the worker will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
METADATA START {"pods":2,"completions":2,"metricName":"app-lammps","metricDescription":"LAMMPS molecular dynamic simulation","metricType":"standalone","metricOptions":{"command":"mpirun --hostfile ./hostlist.txt -np 2 --map-by socket lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite","completions":0,"rate":10,"workdir":"/opt/lammps/examples/reaxff/HNS"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
```

We never actually parse the output of the worker, so it isn't important.
We can do this with JobSet logic that the entire set is done when the launcher is done.

```bash
$ kubectl get pods
```
```console
NAME                           READY   STATUS      RESTARTS   AGE
metricset-sample-l-0-0-vfz4w   0/1     Completed   0          68s
```

When you are done, the job and jobset will be completed.

```bash
$ kubectl get jobset
```
```console
NAME               RESTARTS   COMPLETED   AGE
metricset-sample              True        82s
```
```bash
$ kubectl get jobs
```
```console
NAME                   COMPLETIONS   DURATION   AGE
metricset-sample-n-0   1/1           18s        84s
```

And then you can cleanup!

```bash
kubectl delete -f metrics.yaml
```