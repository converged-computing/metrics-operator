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

Install the operator (from the development manifest here) or the production release:

```bash
kubectl apply -f ../../dist/metrics-operator-dev.yaml
kubectl apply -f https://raw.githubusercontent.com/converged-computing/metrics-operator/main/examples/dist/metrics-operator.yaml
```

How to see metrics operator logs:

```bash
$ kubectl logs -n metrics-system metrics-controller-manager-859c66464c-7rpbw
```

Then create the metrics set. This is going to run a single run of LAMMPS over MPI!
as lammps runs.

```bash
kubectl apply -f metrics-rocky.yaml
```

Wait until you see pods created by the job and then running.

```bash
kubectl get pods
```

And then you can shell in and look at the output, which should be named with the pattern `mpi_profile.<proc>.<rank>`.
I use kubectl copy to copy examples to the present working directory here.

When you are done, cleanup.

```bash
kubectl delete -f metrics.yaml
```