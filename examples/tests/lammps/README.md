# LAMMPS Example

This is (our first!) example to run a simple LAMMPS application in one container, and with one metric
to dump out a bunch of stats. Note that we don't yet have an intelligent way to store these, we are
just seeing them in logs.

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

Then create the metrics set. This is going to run a simple sysstat tool to collect metrics
as lammps runs.

```bash
kubectl apply -f metrics.yaml
```

Wait until you see pods created by the job and then running (there should be one with two containers, one for the app lammps and the other for the stats):

```diff
NAME                           READY   STATUS              RESTARTS   AGE
- metricset-sample-m-0-0-mkwrh   0/2     ContainerCreating   0          2m20s
+ metricset-sample-m-0-0-mkwrh   2/2     Running             0          3m10s
```

You should always be able to get logs for the application container in a pod, which is named "app" ! Here we
see it running LAMMPS.

```bash
$ kubectl logs metricset-sample-m-0-0-kd28s -c app
```
```console
LAMMPS (29 Sep 2021 - Update 2)
OMP_NUM_THREADS environment is not set. Defaulting to 1 thread. (src/comm.cpp:98)
  using 1 OpenMP thread(s) per MPI task
Reading data file ...
  triclinic box = (0.0000000 0.0000000 0.0000000) to (22.326000 11.141200 13.778966) with tilt (0.0000000 -5.0260300 0.0000000)
  1 by 1 by 1 MPI processor grid
  reading atoms ...
  304 atoms
  reading velocities ...
  304 velocities
  read_data CPU = 0.001 seconds
Replicating atoms ...
  triclinic box = (0.0000000 0.0000000 0.0000000) to (22.326000 11.141200 13.778966) with tilt (0.0000000 -5.0260300 0.0000000)
  1 by 1 by 1 MPI processor grid
  bounding box image = (0 -1 -1) to (0 1 1)
  bounding box extra memory = 0.03 MB
  average # of replicas added to proc = 1.00 out of 1 (100.00%)
  304 atoms
  replicate CPU = 0.000 seconds
Neighbor list info ...
  update every 20 steps, delay 0 steps, check no
  max neighbors/atom: 2000, page size: 100000
  master list distance cutoff = 11
  ghost atom cutoff = 11
  binsize = 5.5, bins = 5 3 3
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
Per MPI rank memory allocation (min/avg/max) = 78.15 | 78.15 | 78.15 Mbytes
Step Temp PotEng Press E_vdwl E_coul Volume 
       0          300   -113.27833    427.09094   -111.57687   -1.7014647    3427.3584 
      10    298.13784   -113.27279    1855.1535   -111.57169   -1.7011017    3427.3584 
      20    294.02745   -113.25991    3911.5126     -111.559   -1.7009101    3427.3584 
      30    293.61692   -113.25867    7296.5076   -111.55793   -1.7007375    3427.3584 
      40    301.40293   -113.28175    9622.4058   -111.58127   -1.7004797    3427.3584 
      50    310.92489   -113.31003    10101.225   -111.60982   -1.7002117    3427.3584 
      60    311.37774   -113.31149    9274.1322   -111.61144   -1.7000446    3427.3584 
      70    302.58347   -113.28582     6350.705   -111.58587   -1.6999549    3427.3584 
      80    295.34242   -113.26406    6795.0642   -111.56427   -1.6997975    3427.3584 
      90    299.15724   -113.27518    9198.0327   -111.57566   -1.6995238    3427.3584 
     100    307.63997   -113.30058    9424.4991   -111.60129   -1.6992878    3427.3584 
Loop time of 4.11984 on 1 procs for 100 steps with 304 atoms

Performance: 0.210 ns/day, 114.440 hours/ns, 24.273 timesteps/s
99.9% CPU use with 1 MPI tasks x 1 OpenMP threads

MPI task timing breakdown:
Section |  min time  |  avg time  |  max time  |%varavg| %total
---------------------------------------------------------------
Pair    | 3.3866     | 3.3866     | 3.3866     |   0.0 | 82.20
Neigh   | 0.10642    | 0.10642    | 0.10642    |   0.0 |  2.58
Comm    | 0.0027361  | 0.0027361  | 0.0027361  |   0.0 |  0.07
Output  | 0.00023328 | 0.00023328 | 0.00023328 |   0.0 |  0.01
Modify  | 0.62338    | 0.62338    | 0.62338    |   0.0 | 15.13
Other   |            | 0.0004632  |            |       |  0.01

Nlocal:        304.000 ave         304 max         304 min
Histogram: 1 0 0 0 0 0 0 0 0 0
Nghost:        4443.00 ave        4443 max        4443 min
Histogram: 1 0 0 0 0 0 0 0 0 0
Neighs:        123880.0 ave      123880 max      123880 min
Histogram: 1 0 0 0 0 0 0 0 0 0

Total # of neighbors = 123880
Ave neighs/atom = 407.50000
Neighbor list builds = 5
Dangerous builds not checked
Total wall time: 0:00:04
```

The sidecar metrics containers will output their metrics for the lifetime of the application,
and (currently) also to their logs:

```bash
05:39:24        0        20         -    0.00    0.00    0.00    0.00    0.00     4  mpirun
05:39:24        0         -        20    0.00    0.00    0.00    0.00    0.00     4  |__mpirun
KERNEL TABLES 2176
        34  pidstat -p 20 -v -h
    echo TASK SWITCHING 2176
/metrics_operator/entrypoint-0.sh: line 28: 35: command not found
CPU STATISTICS TIMEPOINT 2177
    pidstat -p 20 -u -h
    echo KERNEL STATISTICS TIMEPOINT 2177
Linux 5.15.0-76-generic (metricset-sample-m-0-0)        07/29/23        _x86_64_    (8 CPU)

# Time        UID       PID   kB_rd/s   kB_wr/s kB_ccwr/s iodelay  Command
05:39:24        0        20      0.00      0.00      0.00       0  mpirun
POLICY TIMEPOINT 2177
    pidstat -p 20 -R -h
    echo PAGEFAULTS and MEMORY 2177
Linux 5.15.0-76-generic (metricset-sample-m-0-0)        07/29/23        _x86_64_    (8 CPU)

# Time        UID       PID  minflt/s  majflt/s     VSZ     RSS   %MEM  Command
05:39:24        0        20      0.00      0.00    6548    3160   0.02  mpirun
STACK UTILIZATION 2177
        pidstat -p 20 -s -h
    echo THREADS 2177
Linux 5.15.0-76-generic (metricset-sample-m-0-0)        07/29/23        _x86_64_    (8 CPU)

# Time        UID      TGID       TID    %usr %system  %guest   %wait    %CPU   CPU  Command
05:39:24        0        20         -    0.00    0.00    0.00    0.00    0.00     4  mpirun
05:39:24        0         -        20    0.00    0.00    0.00    0.00    0.00     4  |__mpirun
```

Those are just a random set of stats I am running using this tool for the shared PID - I need to think
of a better way to capture and save these! Also, right now we don't have a completion policy, but instead have
the metrics collector exit when the PID goes away. In the future we could use a success policy. Here they are
side by side:

![lammps.jpg](lammps.jpg)

The jobset, associated jobs, and pods will be completed:

```bash
$ kubectl get jobset
```
```console
NAME               RESTARTS   COMPLETED   AGE
metricset-sample              True        2m28s
```
```bash
$ kubectl get jobs
```
```console
NAME                   COMPLETIONS   DURATION   AGE
metricset-sample-m-0   1/1           13s        2m51s
```
```bash
$ kubectl get pods
```
```console
NAME                           READY   STATUS      RESTARTS   AGE
metricset-sample-m-0-0-rq4q9   0/2     Completed   0          3m19s
```

When you are done, cleanup!

```bash
kubectl delete -f metrics.yaml
```