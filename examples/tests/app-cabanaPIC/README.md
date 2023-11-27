# Laghos Example

This is an example of a metric app, Laghos. 
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

Then create the metrics set. This is going to run cabanaPIC on a single node.

```bash
kubectl apply -f metrics.yaml
```

Wait until you see the pod created by the job and then running.

```bash
kubectl get pods
```
```diff
NAME                           READY   STATUS    RESTARTS   AGE
metricset-sample-l-0-0-lt782   1/1     Running   0          3s
```

And the output is the simulation. There are output files generated but we aren't retrieving them for this demo.

```bash
kubectl logs metricset-sample-l-0-0-lt782 -f
```
```console
...
5988 117.057419 6.957814e-05 2.644137e-03
5989 117.076973 7.096451e-05 2.644661e-03
5990 117.096519 7.223449e-05 2.645145e-03
5991 117.116066 7.343685e-05 2.645631e-03
5992 117.135612 7.469198e-05 2.646064e-03
5993 117.155167 7.599744e-05 2.646421e-03
5994 117.174713 7.732580e-05 2.646714e-03
5995 117.194260 7.854241e-05 2.646956e-03
5996 117.213814 7.978067e-05 2.647144e-03
5997 117.233360 8.094098e-05 2.647246e-03
5998 117.252907 8.202127e-05 2.647287e-03
5999 117.272453 8.307652e-05 2.647311e-03
6000 117.292007 8.416179e-05 2.647350e-03
METRICS OPERATOR COLLECTION END
```

When you are done, cleanup.

```bash
kubectl delete -f metrics.yaml
```