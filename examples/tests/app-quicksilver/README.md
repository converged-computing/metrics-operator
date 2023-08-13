# Quicksilver Example

This is an example of a metric app, Quicksilver, which is part of the [coral 2 benchmarks](https://asc.llnl.gov/coral-2-benchmarks). 
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
and then the example running, and the log is printed to the console.

```bash
kubectl logs metricset-sample-l-0-0-lt782 -f
```
```console
METADATA START {"pods":2,"completions":2,"metricName":"app-quicksilver","metricDescription":"A proxy app for the Monte Carlo Transport Code","metricType":"standalone","metricOptions":{"command":"qs /opt/quicksilver/Examples/CORAL2_Benchmark/Problem1/Coral2_P1.inp","completions":0,"mpirun":"mpirun --hostfile ./hostlist.txt","rate":10,"workdir":"/opt/quicksilver/Examples"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
MPI Initialized         : MPI_THREAD_FUNNELED
Copyright (c) 2016
Lawrence Livermore National Security, LLC
All Rights Reserved
Quicksilver Version     : 2023-Jun-2-22:49:20
Quicksilver Git Hash    : 6f14107d15824a757fa13abcb42399ff60c020e6
MPI Version             : 4.0
Number of MPI ranks     : 2
Number of OpenMP Threads: 8
Number of OpenMP CPUs   : 8

Simulation:
   dt: 1e-08
   fMax: 0.1
   inputFile: 
   energySpectrum: 
   boundaryCondition: reflect
   loadBalance: 0
   cycleTimers: 0
   debugThreads: 0
   lx: 100
   ly: 100
   lz: 100
   nParticles: 1000000
   batchSize: 0
   nBatches: 10
   nSteps: 10
   nx: 10
   ny: 10
   nz: 10
   seed: 1029384756
   xDom: 0
   yDom: 0
   zDom: 0
   eMax: 20
   eMin: 1e-09
   nGroups: 230
   lowWeightCutoff: 0.001
   bTally: 1
   fTally: 1
   cTally: 1
   coralBenchmark: 0
   crossSectionsOut:

Geometry:
   material: sourceMaterial
   shape: brick
   xMax: 100
   xMin: 0
   yMax: 100
   yMin: 0
   zMax: 100
   zMin: 0

Material:
   name: sourceMaterial
   mass: 1000
   nIsotopes: 10
   nReactions: 9
   sourceRate: 1e+10
   totalCrossSection: 1
   absorptionCrossSection: flat
   fissionCrossSection: flat
   scatteringCrossSection: flat
   absorptionCrossSectionRatio: 1
   fissionCrossSectionRatio: 0.1
   scatteringCrossSectionRatio: 1

CrossSection:
   name: flat
   A: 0
   B: 0
   C: 0
   D: 0
   E: 1
   nuBar: 2.4
Building partition 0
done building
Building MC_Domain 0
Finished initMesh
cycle           start       source           rr        split       absorb      scatter      fission      produce      collisn       escape       census      num_seg   scalar_flux      cycleInit  cycleTracking  cycleFinalize
       0            0       100000            0       900000      1078182      1076792       107133       257364      2262107            0        72049      2670386  2.264064e+08   4.083900e-02   1.450656e+00   3.485000e-03
       1        72049       100000            0       828008      1107255      1106235       110306       264657      2323796            0        47153      2719702  2.438830e+08   3.487200e-02   1.549426e+00   2.600000e-05
       2        47153       100000            0       852712      1086097      1088696       108334       259738      2283127            0        65172      2687840  2.435394e+08   3.052400e-02   1.394607e+00   4.700000e-05
       3        65172       100000        68015       834785      1017555      1018659       101778       244593      2137992            0        57202      2517378  2.450517e+08   5.353000e-02   1.464542e+00   2.600000e-05
       4        57202       100000        62214       842934      1020418      1019522       101687       244038      2141627            0        59855      2522163  2.434017e+08   3.349800e-02   1.500948e+00   4.100000e-05
       5        59855       100000        56726       840345      1029994      1029682       103183       247672      2162859            0        57969      2545713  2.451216e+08   3.957600e-02   1.338280e+00   3.200000e-05
       6        57969       100000        52439       841925      1032190      1032180       102801       246877      2167171            0        59341      2551468  2.446226e+08   5.347400e-02   1.278508e+00   1.457800e-02
       7        59341       100000        59663       840635      1023444      1022593       102792       246649      2148829            0        60726      2531066  2.441845e+08   5.303800e-02   1.555290e+00   4.500000e-05
       8        60726       100000        68187       839357      1013501      1014287       101238       243112      2129026            0        60269      2508491  2.440307e+08   4.993400e-02   2.005018e+00   4.700000e-05
       9        60269       100000        71159       839953      1012439      1011892       101368       243262      2125699            0        58518      2500968  2.444142e+08   5.591000e-02   1.903332e+00   3.900000e-05

Timer                       Cumulative   Cumulative   Cumulative   Cumulative   Cumulative   Cumulative
Name                            number    microSecs    microSecs    microSecs    microSecs   Efficiency
                              of calls          min          avg          max       stddev       Rating
main                                 1    1.592e+07    1.592e+07    1.592e+07    2.852e+03        99.98
cycleInit                           10    4.452e+05    6.596e+05    8.741e+05    2.144e+05        75.47
cycleTracking                       10    1.501e+07    1.522e+07    1.544e+07    2.167e+05        98.60
cycleTracking_Kernel              1470    6.538e+06    1.012e+07    1.369e+07    3.578e+06        73.87
cycleTracking_MPI                 1617    1.309e+06    5.103e+06    8.896e+06    3.794e+06        57.36
cycleTracking_Test_Done            157    1.132e+06    4.980e+06    8.829e+06    3.849e+06        56.41
cycleFinalize                       20    1.011e+04    1.464e+04    1.917e+04    4.527e+03        76.38
Figure Of Merit              1.668e+06 [Num Segments / Cycle Tracking Time]
[WARNING] yaksa: 1 leaked handle pool objects
[WARNING] yaksa: 1 leaked handle pool objects
METRICS OPERATOR COLLECTION END
```

The above shows the structured output that is done in a way for our Python parsing script to easily
find sections of data. Also note that the worker will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
METADATA START {"pods":2,"completions":2,"metricName":"app-quicksilver","metricDescription":"A proxy app for the Monte Carlo Transport Code","metricType":"standalone","metricOptions":{"command":"qs /opt/quicksilver/Examples/CORAL2_Benchmark/Problem1/Coral2_P1.inp","completions":0,"mpirun":"mpirun --hostfile ./hostlist.txt","rate":10,"workdir":"/opt/quicksilver/Examples"}}
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