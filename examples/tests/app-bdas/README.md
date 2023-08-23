# BDAS (Big Data Analysis Suite) Example

This is an example of a set of machine learning mini apps that are part of the [coral 2 benchmarks](https://asc.llnl.gov/coral-2-benchmarks). 
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
METADATA START {"pods":2,"completions":2,"metricName":"app-bdas","metricDescription":"The big data analytic suite contains the K-Means observation label, PCA, and SVM benchmarks.","metricType":"standalone","metricOptions":{"command":"mpirun --allow-run-as-root -np 4 --hostfile ./hostlist.txt Rscript /opt/bdas/benchmarks/r/princomp.r 250 50","prefix":"/bin/bash","workdir":"/opt/bdas/benchmarks/r"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
METADATA START {"pods":2,"completions":2,"metricName":"app-bdas","metricDescription":"The big data analytic suite contains the K-Means observation label, PCA, and SVM benchmarks.","metricType":"standalone","metricOptions":{"command":"mpirun --allow-run-as-root -np 4 --hostfile ./hostlist.txt Rscript /opt/bdas/benchmarks/r/princomp.r 250 50","prefix":"/bin/bash","workdir":"/opt/bdas/benchmarks/r"}}
METADATA END
Hostlist
10.244.0.16
10.244.0.17
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
[1] 1.472100 0.624046
COMM.RANK = 0
    min    mean     max 
0.00900 0.01425 0.02000 
METRICS OPERATOR COLLECTION END
```

The above shows the structured output that is done in a way for our Python parsing script to easily
find sections of data. Also note that the worker will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
METADATA START {"pods":2,"completions":2,"metricName":"app-bdas","metricDescription":"The big data analytic suite contains the K-Means observation label, PCA, and SVM benchmarks.","metricType":"standalone","metricOptions":{"command":"mpirun --allow-run-as-root -np 4 --hostfile ./hostlist.txt Rscript /opt/bdas/benchmarks/r/princomp.r 250 50","prefix":"/bin/bash","workdir":"/opt/bdas/benchmarks/r"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
```

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