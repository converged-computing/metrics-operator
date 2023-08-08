# Flux Operator Example

This will be an example (experiment) to create Flux Operator assets (config maps, service, etc) that can be deployed alongside
the Metrics operator, meaning that:

 - The Flux Operator (Flux MiniCluster) is the application of interest
 - We can collect metrics alongside it.

My first idea is to save all the assets for a small LAMMPS run and then package the actual cluster into
the metrics operator setup, and then allow the application to take a custom network instead of the one
that we create.


## Usage

Create a local kind cluster:

```bash
$ kind create cluster
```

### 1. Create Flux Operator Assets

Let's first create assets for the Flux Operator. We will install the Flux Operator to start with a LAMMPS [minicluster.yaml](minicluster.yaml).

```bash
$ kubectl apply -f https://github.com/flux-framework/flux-operator/releases/download/0.1.0/flux-operator.yaml
$ kubectl apply -f minicluster.yaml
```

Now let's save the yaml assets for inspection / later use

```bash
kubectl get jobs flux-sample -o yaml > yaml/minicluster-job.yaml
kubectl get svc flux-service -o yaml > yaml/minicluster-service.yaml
kubectl get cm flux-sample-curve-mount -o yaml > yaml/minicluster-curve-cm.yaml
kubectl get cm flux-sample-entrypoint -o yaml > yaml/minicluster-entrypoint-cm.yaml
kubectl get cm flux-sample-flux-config -o yaml > yaml/minicluster-config-cm.yaml
```

We can then assemble a custom yaml that includes everything except for the Job, which we will create
via the Metrics Operator. We shouldn't need the flux operator headless service, because the metrics operator will
create one for the application. However, we do want to give it a custom name. Note that I assembled this into
[yaml/minicluster.yaml](yaml/minicluster.yaml).  Note that the tweaks needed are:

- hostname in the broker config / entrypoint for R generation needs to be tweaked for a JobSet, e.g., 
  - `flux-sample-0.flux-service.default.svc.cluster.local` changed to `flux-sample-m-0-0.flux-service.default.svc.cluster.local`
- metrics.yaml size needs to correspond with MiniCluster size
- tasks/cores (and when added) resource limits too.

Let's next delete the flux operator to free up our tiny
resources.

```bash
$ kubectl delete -f https://github.com/flux-framework/flux-operator/releases/download/0.1.0/flux-operator.yaml
```

Now let's install JobSet and the metrics operator:

```bash
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

Then create the flux operator assets. This includes the entrypoint, the flux broker config,
and the curve certificate.

```bash
kubectl apply -f yaml/minicluster.yaml
```

That should only need to create three config maps:

```bash
 kubectl get cm
NAME                      DATA   AGE
flux-sample-curve-mount   1      14s
flux-sample-entrypoint    1      14s
flux-sample-flux-config   1      14s
kube-root-ca.crt          1      59m
```

Then create the metrics set. This is going to run a simple sysstat tool to collect metrics
as lammps runs (we hope)!

```bash
kubectl apply -f metrics.yaml
```

Wait until you see pods created by the job and then running (there should be four. Each pod has flux running LAMMPS,
and sidecar containers to watch that process:

```bash
kubectl get pods
```
```console
NAME                      READY   STATUS      RESTARTS   AGE
flux-sample-m-0-0-ggc64   0/2     Completed   0          2m24s
flux-sample-m-0-1-nhqd7   0/2     Completed   0          2m24s
flux-sample-m-0-2-ph99v   0/2     Completed   0          2m24s
flux-sample-m-0-3-qw77b   0/2     Completed   0          2m24s
```

You can look at an application "app" container in any pod to see flux (and LAMMPS)

```bash
kubectl flux-sample-m-0-0-ggc64 -c app -f
```

<details>

<summary>Output of LAMMPS</summary>

```console
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
  read_data CPU = 0.002 seconds
Replicating atoms ...
  triclinic box = (0.0000000 0.0000000 0.0000000) to (44.652000 22.282400 27.557932) with tilt (0.0000000 -10.052060 0.0000000)
  2 by 1 by 1 MPI processor grid
  bounding box image = (0 -1 -1) to (0 1 1)
  bounding box extra memory = 0.03 MB
  average # of replicas added to proc = 5.00 out of 8 (62.50%)
  2432 atoms
  replicate CPU = 0.000 seconds
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
Loop time of 12.2637 on 2 procs for 100 steps with 2432 atoms

Performance: 0.070 ns/day, 340.659 hours/ns, 8.154 timesteps/s
99.7% CPU use with 2 MPI tasks x 1 OpenMP threads

MPI task timing breakdown:
Section |  min time  |  avg time  |  max time  |%varavg| %total
---------------------------------------------------------------
Pair    | 8.5292     | 8.87       | 9.2107     |  11.4 | 72.33
Neigh   | 0.25438    | 0.25486    | 0.25534    |   0.1 |  2.08
Comm    | 0.02493    | 0.36577    | 0.7066     |  56.4 |  2.98
Output  | 0.0002314  | 0.00025307 | 0.00027474 |   0.0 |  0.00
Modify  | 2.7709     | 2.7715     | 2.772      |   0.0 | 22.60
Other   |            | 0.001404   |            |       |  0.01

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
```

</details>

or in the metrics pod to see pidstat running for the application container entrypoint (to start flux and run LAMMPS)

```bash
kubectl flux-sample-m-0-0-ggc64 -f
```
```console
Defaulted container "perf-sysstat" out of: perf-sysstat, app
METADATA START {"pods":4,"completions":4,"applicationImage":"ghcr.io/rse-ops/lammps:flux-sched-focal","applicationCommand":"/bin/bash /flux_operator/wait-0.sh lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite","metricName":"perf-sysstat","metricDescription":"statistics for Linux tasks (processes) : I/O, CPU, memory, etc.","metricType":"application","metricOptions":{"completions":0,"rate":10}}
METADATA END
--2023-08-08 05:36:33--  https://github.com/converged-computing/goshare/releases/download/2023-07-27/wait
Resolving github.com (github.com)... 140.82.113.3
Connecting to github.com (github.com)|140.82.113.3|:443... connected.
HTTP request sent, awaiting response... 302 Found
Location: https://objects.githubusercontent.com/github-production-release-asset-2e65be/670447414/18f62ebc-64a5-483e-8e12-abab49f1d694?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20230808%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20230808T053634Z&X-Amz-Expires=300&X-Amz-Signature=b6905861de3f4dc1a75e227eac0b0c6d465b942d6ec974e77d0eeedc2f215ac2&X-Amz-SignedHeaders=host&actor_id=0&key_id=0&repo_id=670447414&response-content-disposition=attachment%3B%20filename%3Dwait&response-content-type=application%2Foctet-stream [following]
--2023-08-08 05:36:34--  https://objects.githubusercontent.com/github-production-release-asset-2e65be/670447414/18f62ebc-64a5-483e-8e12-abab49f1d694?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20230808%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20230808T053634Z&X-Amz-Expires=300&X-Amz-Signature=b6905861de3f4dc1a75e227eac0b0c6d465b942d6ec974e77d0eeedc2f215ac2&X-Amz-SignedHeaders=host&actor_id=0&key_id=0&repo_id=670447414&response-content-disposition=attachment%3B%20filename%3Dwait&response-content-type=application%2Foctet-stream
Resolving objects.githubusercontent.com (objects.githubusercontent.com)... 185.199.110.133, 185.199.108.133, 185.199.111.133, ...
Connecting to objects.githubusercontent.com (objects.githubusercontent.com)|185.199.110.133|:443... connected.
HTTP request sent, awaiting response... 200 OK
Length: 2556028 (2.4M) [application/octet-stream]
Saving to: 'wait'

wait                100%[===================>]   2.44M  9.78MB/s    in 0.2s    

2023-08-08 05:36:34 (9.78 MB/s) - 'wait' saved [2556028/2556028]

PIDSTAT COMMAND START
/bin/bash /flux_operator/wait-0.sh lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite
PIDSTAT COMMAND END
Waiting for application PID...
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
CPU STATISTICS
[{"time":53634,"uid":0,"pid":20,"percent_usr":0.0,"percent_system":0.0,"percent_guest":0.0,"percent_wait":"0.00","percent_cpu":0.0,"cpu":7,"command":"bash"}]
KERNEL STATISTICS
[{"time":53634,"uid":0,"pid":20,"kb_rd_s":0.01,"kb_wr_s":0.0,"kb_ccwr_s":0.0,"iodelay":"0","command":"bash"}]
POLICY
[{"time":53634,"uid":0,"pid":20,"prio":"0","policy":"NORMAL","command":"bash"}]
PAGEFAULTS
[{"time":53635,"uid":0,"pid":20,"minflt_s":0.0,"majflt_s":0.0,"vsz":8808,"rss":3440,"percent_mem":0.02,"command":"bash"}]
STACK UTILIZATION
[{"time":53635,"uid":0,"pid":20,"stksize":132,"stkref":12,"command":"bash"}]
THREADS
[{"time":53635,"uid":0,"tgid":"20","tid":"-","percent_usr":0.0,"percent_system":0.0,"percent_guest":0.0,"percent_wait":"0.00","percent_cpu":0.0,"cpu":7,"command":"bash"},{"time":53635,"uid":0,"tgid":"-","tid":"20","percent_usr":0.0,"percent_system":0.0,"percent_guest":0.0,"percent_wait":"0.00","percent_cpu":0.0,"cpu":7,"command":"|__bash"}]
KERNEL TABLES
[{"time":53635,"uid":0,"pid":20,"threads":"1","fd-nr":"4","command":"bash"}]
TASK SWITCHING
[{"time":53635,"uid":0,"pid":20,"cswch_s":0.0,"nvcswch_s":0.0,"command":"bash"}]
METRICS OPERATOR TIMEPOINT
CPU STATISTICS
[{"time":53645,"uid":0,"pid":20,"percent_usr":0.0,"percent_system":0.0,"percent_guest":0.0,"percent_wait":"0.00","percent_cpu":0.0,"cpu":5,"command":"bash"}]
KERNEL STATISTICS
[{"time":53645,"uid":0,"pid":20,"kb_rd_s":0.01,"kb_wr_s":0.0,"kb_ccwr_s":0.0,"iodelay":"0","command":"bash"}]
POLICY
[{"time":53645,"uid":0,"pid":20,"prio":"0","policy":"NORMAL","command":"bash"}]
PAGEFAULTS
[{"time":53645,"uid":0,"pid":20,"minflt_s":0.0,"majflt_s":0.0,"vsz":8808,"rss":552,"percent_mem":0.0,"command":"bash"}]
STACK UTILIZATION
[{"time":53646,"uid":0,"pid":20,"stksize":132,"stkref":16,"command":"bash"}]
THREADS
[{"time":53646,"uid":0,"tgid":"20","tid":"-","percent_usr":0.0,"percent_system":0.0,"percent_guest":0.0,"percent_wait":"0.00","percent_cpu":0.0,"cpu":5,"command":"bash"},{"time":53646,"uid":0,"tgid":"-","tid":"20","percent_usr":0.0,"percent_system":0.0,"percent_guest":0.0,"percent_wait":"0.00","percent_cpu":0.0,"cpu":5,"command":"|__bash"}]
KERNEL TABLES
[{"time":53646,"uid":0,"pid":20,"threads":"1","fd-nr":"4","command":"bash"}]
TASK SWITCHING
[{"time":53646,"uid":0,"pid":20,"cswch_s":0.0,"nvcswch_s":0.0,"command":"bash"}]
METRICS OPERATOR TIMEPOINT
CPU STATISTICS
[]
KERNEL STATISTICS
[]
POLICY
[]
PAGEFAULTS
[]
STACK UTILIZATION
[]
THREADS
[]
KERNEL TABLES
[]
TASK SWITCHING
[]
METRICS OPERATOR COLLECTION END
```

Note that the last measurment is empty because the process ended in the middle, and we check at the end.
Also notice the log is formatted so we should be able to add a Python parser soon.

When you are done, the job and jobset will be completed.

```bash
$ kubectl get jobset
```
```console
NAME               RESTARTS   COMPLETED   AGE
flux-sample                   True        82s
```
```bash
$ kubectl get jobs
```
```console
NAME                   COMPLETIONS   DURATION   AGE
flux-sample-m-0        1/1           18s        84s
```

And then you can cleanup!

```bash
kubectl delete -f metrics.yaml
```

We will add the Python parser for this metric soon.
