# DLIO Example

This is an example of using the IO tool[DLIO](https://dlio-profiler.readthedocs.io/en/latest/build.html#build-dlio-profiler-with-pip-recommended) that can 
be added on the fly with pip.

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

Then create the metrics set. This is going to run a single run of LAMMPS over MPI.
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

There is purposefully a sleep infinity at the end to give you a chance to copy over data.

```bash
mkdir logs
kubectl cp metricset-sample-l-0-0-mm2dx:/opt/logs/ logs/
```

Cleanup when you are done.

```bash
kubectl delete -f metrics.yaml
```