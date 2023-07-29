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