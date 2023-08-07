# Networking Example

The container for this example (and code) is private, so you won't be able to run it!
I'm using it for development purposes.
For running the example, parsing, and plotting output, see [the corresponding Python directory](../../python/network-netmark/).

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

Make sure to load your private image into the node:

```bash
kind load docker-image vanessa/netmark:latest
```

Then create the metrics set. This is going to run a simple sysstat tool to collect metrics
as lammps runs.

```bash
kubectl apply -f metrics.yaml
```

Wait until you see pods created by the job and then running (there should be two):

```bash
kubectl get pods
```
```diff
NAME                           READY   STATUS    RESTARTS   AGE
metricset-sample-n-0-0-lt782   1/1     Running   0          3s
metricset-sample-w-0-0-4s5p9   1/1     Running   0          3s
```

In the above, "w" is a worker pod, and "n" is the netmark launcher.
If you inspect the log for the launcher you'll see a short sleep (the network isn't up immediately)
and then netmark running, and the matrix of RTT.csv times is printed at the end.

```bash
kubectl logs metricset-sample-n-0-0-lt782 -f
```
```console
root
METADATA START {"pods":2,"completions":2,"metricName":"network-netmark","metricDescription":"point to point networking tool","metricType":"standalone","metricOptions":{"completions":0,"messageSize":0,"rate":10,"sendReceiveCycles":20,"storeEachTrial":"true","tasks":2,"trials":20,"warmups":10}}
METADATA END
Sleeping for 10 seconds waiting for network...
(env) (base) vanessa@vanessa-ThinkPad-T490s:~/Desktop/Code/metrics-operator$ kubectl logs metricset-sample-n-0-0-82jz4 -f
root
METADATA START {"pods":2,"completions":2,"metricName":"network-netmark","metricDescription":"point to point networking tool","metricType":"standalone","metricOptions":{"completions":0,"messageSize":0,"rate":10,"sendReceiveCycles":20,"storeEachTrial":"true","tasks":2,"trials":20,"warmups":10}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
=========== SETUP ===========
warmups                  10
trials                   20
send_recv_cycles         20
bytes                     0
store_per_trial           1
=============================
size 2 rank 1 on host metricset-sample-w-0-0.ms.default.svc.cluster.local ip 10.244.0.68
size 2 rank 0 on host metricset-sample-n-0-0.ms.default.svc.cluster.local ip 10.244.0.67
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 34.527 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 29.943 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 28.696 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 27.141 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 26.679 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 25.841 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 28.158 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 29.306 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 28.491 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 27.729 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 27.176 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 26.680 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 26.275 micro-seconds
Rank 0 sends to rank 1
Rank 1 sends to rank 0
RTT between rank 0 and rank 1 is 26.246 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 26.975 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 27.065 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 26.780 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 26.764 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 26.815 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 26.782 micro-seconds
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 25.498 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 24.708 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 24.268 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 24.765 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 26.183 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 27.531 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 28.310 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 28.821 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 29.143 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 29.486 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 29.838 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 29.084 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 28.213 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 27.658 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 27.058 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 27.269 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 28.099 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 27.888 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 27.467 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 27.071 micro-seconds
BW-1.csv   BW-16.csv  BW-4.csv  RTT-1.csv   RTT-16.csv  RTT-4.csv  hostlist.txt
BW-10.csv  BW-17.csv  BW-5.csv  RTT-10.csv  RTT-17.csv  RTT-5.csv  hosts.csv
BW-11.csv  BW-18.csv  BW-6.csv  RTT-11.csv  RTT-18.csv  RTT-6.csv  ips.csv
BW-12.csv  BW-19.csv  BW-7.csv  RTT-12.csv  RTT-19.csv  RTT-7.csv  scripts
BW-13.csv  BW-2.csv   BW-8.csv  RTT-13.csv  RTT-2.csv   RTT-8.csv
BW-14.csv  BW-20.csv  BW-9.csv  RTT-14.csv  RTT-20.csv  RTT-9.csv
BW-15.csv  BW-3.csv   BW.csv    RTT-15.csv  RTT-3.csv   RTT.csv
NETMARK RTT.CSV START
0.000,26.782
27.071,0.000
NETMARK RTT.CSV END
METRICS OPERATOR COLLECTION END
```
The above also shows the structured output that is done in a way for our Python parsing script to easily
find sections of data. Also note that the worker will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
root
METADATA START {"pods":2,"completions":2,"metricName":"network-netmark","metricDescription":"point to point networking tool","metricType":"standalone","metricOptions":{"completions":0,"messageSize":0,"rate":10,"sendReceiveCycles":20,"storeEachTrial":"true","tasks":2,"trials":20,"warmups":10}}
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
NAME                           READY   STATUS        RESTARTS   AGE
metricset-sample-n-0-0-bqqf4   0/1     Completed     0          49s
metricset-sample-w-0-0-97h2g   1/1     Terminating   0          49s
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

If you want to see how it scales, try increasing the number of pods for Netmark. You
might want to change the ranks too. Here is an example with 4 pods total (so one launcher
and 3 workers):

```bash
$ kubectl get pods
NAME                           READY   STATUS        RESTARTS   AGE
metricset-sample-n-0-0-2zprc   0/1     Completed     0          23s
metricset-sample-w-0-0-76wx9   1/1     Terminating   0          23s
metricset-sample-w-0-1-5j2kh   1/1     Terminating   0          23s
metricset-sample-w-0-2-lxg7w   1/1     Terminating   0          23s
```

Again, the worker jobs clean up nicely, and the output is available for us
to parse via the launcher pod.
