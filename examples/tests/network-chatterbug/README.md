# Chatterbug Networking Example

This will demonstrate running a [Chatterbug](https://github.com/hpcgroup/chatterbug) metric.
This metric is experimental and not working in all contexts.

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

Then create the metrics set. This is going to run metrics to assess networking with netmark.

```bash
kubectl apply -f metrics.yaml
```

Wait until you see pods created by the job and then running (there should be one).
Note that although this is a launcher/worker setup that can handle 2+ nodes (a launcher and workers)
for this example we only need one node.

```bash
kubectl get pods
```
```console
NAME                         READY   STATUS    RESTARTS   AGE
metricset-sample-l-0-5r9vv   1/1     Running   0          4s
```

**NOTE** this is not currently working, I haven't gotten it to work beyond one node, and one node isn't
really an option given that [the hostnames are not predictable](https://github.com/kubernetes-sigs/jobset/issues/290).

```bash
kubectl logs metricset-sample-l-0-5r9vv -f
```
```console
```

Note that the printing of tasks/pods in the beginning is for your FYI only - it's up to you to decide 
We can do this with JobSet logic that the entire set is done when the launcher is done.

```bash
$ kubectl get pods
```
```console
NAME                           READY   STATUS        RESTARTS   AGE
metricset-sample-l-0-5r9vv   0/1     Completed     0          49s
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
