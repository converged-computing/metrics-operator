# AMG Example

This is an example of a metric app, AMG, which is part of the [coral 2 benchmarks](https://asc.llnl.gov/coral-2-benchmarks) and technically
isn't a metric, but we use it often to assess time and thus the MPI latency. A Python example (parsing the output data)
is provided in [python/app-amg](../../python/app-amg).

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
METADATA START {"pods":2,"metricName":"app-amg","metricDescription":"parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids","metricType":"standalone","metricOptions":{"command":"amg","mpirun":"mpirun --hostfile ./hostlist.txt","workdir":"/opt/AMG"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
Running with these driver parameters:
  solver ID    = 1

  Laplacian_27pt:
    (Nx, Ny, Nz) = (10, 20, 10)
    (Px, Py, Pz) = (1, 2, 1)

=============================================
Generate Matrix:
=============================================
Spatial Operator:
  wall clock time = 0.168312 seconds
  wall MFLOPS     = 0.000000
  cpu clock time  = 0.758584 seconds
  cpu MFLOPS      = 0.000000

  RHS vector has unit components
  Initial guess is 0
=============================================
IJ Vector Setup:
=============================================
RHS and Initial Guess:
  wall clock time = 0.049731 seconds
  wall MFLOPS     = 0.000000
  cpu clock time  = 0.152173 seconds
  cpu MFLOPS      = 0.000000

=============================================
Problem 1: AMG Setup Time:
=============================================
PCG Setup:
  wall clock time = 5.690908 seconds
  wall MFLOPS     = 0.000000
  cpu clock time  = 21.301134 seconds
  cpu MFLOPS      = 0.000000


FOM_Setup: nnz_AP / Setup Phase Time: 8.617078e+03

=============================================
Problem 1: AMG-PCG Solve Time:
=============================================
PCG Solve:
  wall clock time = 19.234349 seconds
  wall MFLOPS     = 0.000000
  cpu clock time  = 72.573477 seconds
  cpu MFLOPS      = 0.000000


Iterations = 14
Final Relative Residual Norm = 4.643894e-09


FOM_Solve: nnz_AP * Iterations / Solve Phase Time: 3.569375e+04



Figure of Merit (FOM_1): 2.892458e+04

METRICS OPERATOR COLLECTION END
```

The above also shows the structured output that is done in a way for our Python parsing script to easily
find sections of data. Also note that the worker will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
METADATA START {"pods":2,"metricName":"app-amg","metricDescription":"parallel algebraic multigrid solver for linear systems arising from problems on unstructured grids","metricType":"standalone","metricOptions":{"command":"amg","mpirun":"mpirun --hostfile ./hostlist.txt","workdir":"/opt/AMG"}}
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