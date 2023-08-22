# Nekbone Example

This is an example of a metric app, Nekbone, which is part of the [coral 2 benchmarks](https://asc.llnl.gov/coral-2-benchmarks). 
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
and then the example running, and the log is printed to the console. Note this is example2 
provided in the container.

```bash
kubectl logs metricset-sample-l-0-0-lt782 -f
```
```console
METADATA START {"pods":2,"completions":2,"metricName":"app-nekbone","metricDescription":"A mini-app derived from the Nek5000 CFD code which is a high order, incompressible Navier-Stokes CFD solver based on the spectral element method. The conjugate gradiant solve is compute intense, contains small messages and frequent allreduces.","metricType":"standalone","metricOptions":{"command":"mpiexec --hostfile ./hostlist.txt -np 2 ./nekbone","prefix":"/bin/bash","workdir":"/root/nekbone-3.0/test/example2"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
 Number of processors:           2
 REAL    wdsize      :           8
 INTEGER wdsize      :           4
 ifmgrid    : F     ifbrick    : T

 Processor Distribution:  npx,npy,npz=           2           1           1
 Element Distribution: nelx,nely,nelz=          10           5           2
 Local Element Distribution: mx,my,mz=           5           5           2
gs_setup: 874 unique labels shared
   pairwise times (avg, min, max): 4.41642e-05 4.41619e-05 4.41666e-05
   crystal router                : 4.37791e-05 4.37076e-05 4.38505e-05
   all reduce                    : 8.83144e-05 8.82852e-05 8.83436e-05
   handle bytes (avg, min, max): 277128 277128 277128
   buffer bytes (avg, min, max): 27968 27968 27968
   used all_to_all method: crystal router

cg:   0  1.5055E+02
cg: 101  1.4517E-08  4.0957E-01  6.4959E-01  8.0727E-16
cg:   0  1.6454E+02
cg: 101  1.8142E-08  4.1163E-01  6.5230E-01  1.2515E-15

nelt =      50, np =         2, nx1 =      10, elements =       100
Tot MFlops =   6.2108E+03, MFlops      =   3.1054E+03
Setup Flop =   6.9500E+08, Solver Flop =   7.5750E+07
Solve Time =   0.2482E+00
Avg MFlops =   6.2108E+03
 Exitting....
METRICS OPERATOR COLLECTION END
```

The above shows the structured output that is done in a way for our Python parsing script to easily
find sections of data. Also note that the worker will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
METADATA START {"pods":2,"completions":2,"metricName":"app-nekbone","metricDescription":"A mini-app derived from the Nek5000 CFD code which is a high order, incompressible Navier-Stokes CFD solver based on the spectral element method. The conjugate gradiant solve is compute intense, contains small messages and frequent allreduces.","metricType":"standalone","metricOptions":{"command":"mpiexec --hostfile ./hostlist.txt -np 2 ./nekbone","prefix":"/bin/bash","workdir":"/root/nekbone-3.0/test/example2"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
```

Nekbone currently runs a set of scoped examples and metrics that it provides, and you can see [the docs](https://converged-computing.github.io/metrics-operator/getting_started/metrics.html#app-nekbone) for details on working directories and commands.
When you are done, the pods should be completed.

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