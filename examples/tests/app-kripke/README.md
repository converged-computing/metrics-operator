# Kripke Example

This is an example of a metric app, Kripke, which is part of the [coral 2 benchmarks](https://asc.llnl.gov/coral-2-benchmarks). 
We have not yet added a Python example as we want a use case first, but can and will when it is warranted.

## Usage

Create a cluster

```bash
kind create cluster
```

and install JobSet to it.

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
and then AMG running a test, and the log is printed to the console.

```bash
kubectl logs metricset-sample-l-0-0-lt782 -f
```
```console
METADATA START {"pods":2,"completions":2,"metricName":"app-kripke","metricDescription":"parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids","metricType":"standalone","metricOptions":{"command":"kripke","completions":0,"mpirun":"","rate":10,"workdir":"/opt/kripke"}}
METADATA END
/metrics_operator/kripke-launcher.sh: line 7: cd: /opt/kripke: No such file or directory
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT

   _  __       _         _
  | |/ /      (_)       | |
  | ' /  _ __  _  _ __  | | __ ___
  |  <  | '__|| || '_ \ | |/ // _ \ 
  | . \ | |   | || |_) ||   <|  __/
  |_|\_\|_|   |_|| .__/ |_|\_\\___|
                 | |
                 |_|        Version 1.2.5-dev

LLNL-CODE-775068

Copyright (c) 2014-22, Lawrence Livermore National Security, LLC

Kripke is released under the BSD 3-Clause License, please see the
LICENSE file for the full license

This work was produced under the auspices of the U.S. Department of
Energy by Lawrence Livermore National Laboratory under Contract
DE-AC52-07NA27344.

Author: Adam J. Kunen <kunen1@llnl.gov>

Compilation Options:
  Architecture:           Sequential
  Compiler:               /usr/bin/c++
  Compiler Flags:         "    "
  Linker Flags:           " "
  CHAI Enabled:           No
  CUDA Enabled:           No
  MPI Enabled:            No
  OpenMP Enabled:         No
  Caliper Enabled:        No

Input Parameters
================

  Problem Size:
    Zones:                 16 x 16 x 16  (4096 total)
    Groups:                32
    Legendre Order:        4
    Quadrature Set:        Dummy S2 with 96 points

  Physical Properties:
    Total X-Sec:           sigt=[0.100000, 0.000100, 0.100000]
    Scattering X-Sec:      sigs=[0.050000, 0.000050, 0.050000]

  Solver Options:
    Number iterations:     10

  MPI Decomposition Options:
    Total MPI tasks:       1
    Spatial decomp:        1 x 1 x 1 MPI tasks
    Block solve method:    Sweep

  Per-Task Options:
    DirSets/Directions:    8 sets, 12 directions/set
    GroupSet/Groups:       2 sets, 16 groups/set
    Zone Sets:             1 x 1 x 1
    Architecture:          Sequential
    Data Layout:           DGZ

Generating Problem
==================

  Decomposition Space:   Procs:      Subdomains (local/global):
  ---------------------  ----------  --------------------------
  (P) Energy:            1           2 / 2
  (Q) Direction:         1           8 / 8
  (R) Space:             1           1 / 1
  (Rx,Ry,Rz) R in XYZ:   1x1x1       1x1x1 / 1x1x1
  (PQR) TOTAL:           1           16 / 16

  Material Volumes=[8.789062e+03, 1.177734e+05, 2.753438e+06]

  Memory breakdown of Field variables:
  Field Variable            Num Elements    Megabytes
  --------------            ------------    ---------
  data/sigs                        15360        0.117
  dx                                  16        0.000
  dy                                  16        0.000
  dz                                  16        0.000
  ell                               2400        0.018
  ell_plus                          2400        0.018
  i_plane                         786432        6.000
  j_plane                         786432        6.000
  k_plane                         786432        6.000
  mixelem_to_fraction               4352        0.033
  phi                            3276800       25.000
  phi_out                        3276800       25.000
  psi                           12582912       96.000
  quadrature/w                        96        0.001
  quadrature/xcos                     96        0.001
  quadrature/ycos                     96        0.001
  quadrature/zcos                     96        0.001
  rhs                           12582912       96.000
  sigt_zonal                      131072        1.000
  volume                            4096        0.031
  --------                  ------------    ---------
  TOTAL                         34238832      261.222

  Generation Complete!

Steady State Solve
==================

  iter 0: particle count=3.743744e+07, change=1.000000e+00
  iter 1: particle count=5.629276e+07, change=3.349511e-01
  iter 2: particle count=6.569619e+07, change=1.431351e-01
  iter 3: particle count=7.036907e+07, change=6.640521e-02
  iter 4: particle count=7.268400e+07, change=3.184924e-02
  iter 5: particle count=7.382710e+07, change=1.548355e-02
  iter 6: particle count=7.438973e+07, change=7.563193e-03
  iter 7: particle count=7.466578e+07, change=3.697158e-03
  iter 8: particle count=7.480083e+07, change=1.805479e-03
  iter 9: particle count=7.486672e+07, change=8.801810e-04
  Solver terminated

Timers
======

  Timer                    Count       Seconds
  ----------------  ------------  ------------
  Generate                     1       0.00392
  LPlusTimes                  10       2.47684
  LTimes                      10       2.56336
  Population                  10       0.24113
  Scattering                  10       3.09063
  Solve                        1      10.14202
  Source                      10       0.00094
  SweepSolver                 10       1.50448
  SweepSubdomain             160       1.46460

TIMER_NAMES:Generate,LPlusTimes,LTimes,Population,Scattering,Solve,Source,SweepSolver,SweepSubdomain
TIMER_DATA:0.003920,2.476842,2.563358,0.241126,3.090631,10.142021,0.000936,1.504483,1.464600

Figures of Merit
================

  Throughput:         1.240671e+07 [unknowns/(second/iteration)]
  Grind time :        8.060154e-08 [(seconds/iteration)/unknowns]
  Sweep efficiency :  97.34902 [100.0 * SweepSubdomain time / SweepSolver time]
  Number of unknowns: 12582912

END
METRICS OPERATOR COLLECTION END
```

The above shows the structured output that is done in a way for our Python parsing script to easily
find sections of data. Also note that the worker will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
METADATA START {"pods":2,"completions":2,"metricName":"app-kripke","metricDescription":"parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids","metricType":"standalone","metricOptions":{"command":"kripke","completions":0,"mpirun":"mpirun --hostfile ./hostlist.txt","rate":10,"workdir":"/opt/kripke"}}
METADATA END
/metrics_operator/kripke-worker.sh: line 7: cd: /opt/kripke: No such file or directory
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