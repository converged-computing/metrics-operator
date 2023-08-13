# Pennant Example

This is an example of a metric app, Pennant, which is part of the [coral 2 benchmarks](https://asc.llnl.gov/coral-2-benchmarks). 
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
METADATA START {"pods":2,"completions":2,"metricName":"app-pennant","metricDescription":"Unstructured mesh hydrodynamics for advanced architectures ","metricType":"standalone","metricOptions":{"command":"pennant /opt/pennant/test/sedovsmall/sedovsmall.pnt","completions":0,"mpirun":"mpirun --hostfile ./hostlist.txt","rate":10,"workdir":"/opt/pennant/test"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
********************
Running PENNANT v0.9
********************

Running on 2 MPI PE(s)
Running on 8 thread(s)
--- Mesh Information ---
Points:  100
Zones:  81
Sides:  324
Edges:  189
Side chunks:  21
Point chunks:  8
Zone chunks:  6
Chunk size:  16
------------------------
Energy check:  total energy  =   2.467991e-01
(internal =   2.467991e-01, kinetic =   0.000000e+00)
End cycle      1, time = 2.50000e-03, dt = 2.50000e-03, wall = 1.64902e-01
dt limiter: Initial timestep
End cycle     10, time = 2.85593e-02, dt = 2.58849e-03, wall = 1.72612e+00
dt limiter: PE 0, Hydro dV/V limit for z = 0

Run complete
cycle =     10,         cstop =     10
time  =   2.855932e-02, tstop =   1.000000e+00

************************************
hydro cycle run time=   1.892289e+00
************************************
Energy check:  total energy  =   2.512181e-01
(internal =   1.874053e-01, kinetic =   6.381282e-02)
Writing .xy file...
METRICS OPERATOR COLLECTION END
```

The above shows the structured output that is done in a way for our Python parsing script to easily
find sections of data. Also note that the worker will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
METADATA START {"pods":2,"completions":2,"metricName":"app-pennant","metricDescription":"Unstructured mesh hydrodynamics for advanced architectures ","metricType":"standalone","metricOptions":{"command":"pennant /opt/pennant/test/sedovsmall/sedovsmall.pnt","completions":0,"mpirun":"mpirun --hostfile ./hostlist.txt","rate":10,"workdir":"/opt/pennant/test"}}
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