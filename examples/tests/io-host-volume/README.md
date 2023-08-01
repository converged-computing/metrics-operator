# Host Volume IO Example

This is our first example of running a storage metric on a simple host volume.
We will first create the volume, and then reference it for our metric of interest.

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

We would typically apply a manifest for storage, but since we are using a host volume, we don't need to do that (the host is the kind container that we will bind to). We currently have a bug that previous config maps aren't cleaned up, so do that first (if you've run other examples):

```bash
$ kubectl delete cm metricset-sample
```

Next, let's create the metrics set. This is going to run a simple sysstat io tool to collect metrics about our volumes.

```bash
kubectl apply -f metrics.yaml
```

Note that the metrics.yaml is asking to run a sysstat io metric for 10 completions, with a delay of 10 seconds between each check (10). Wait until you see a pod created by the job and then running (there should be one since we've only asked for one metric):

```bash
kubectl get pods 
NAME                           READY   STATUS              RESTARTS   AGE
```

If you peek at logs, you'll see the storage metric running once every 10 seconds, and for a total of 10 times.

```bash
host-volume-io$ kubectl logs metricset-sample-m-0-9pq6w 
```
```console
IOSTAT TIMEPOINT 10
Linux 5.15.0-76-generic (metricset-sample-m-0-c9gqj)    08/01/23        _x86_64_        (8 CPU)

avg-cpu:  %user   %nice %system %iowait  %steal   %idle
          38.41    0.09   13.63    0.09    0.00   47.79

Device             tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
loop0             0.00         0.09         0.00         0.00     206733          0          0
loop1             0.01         0.11         0.00         0.00     269890          0          0
loop10            0.01         0.13         0.00         0.00     313378          0          0
loop11            0.00         0.04         0.00         0.00      91891          0          0
loop12            0.01         0.14         0.00         0.00     345082          0          0
loop13            0.03         0.77         0.00         0.00    1841667          0          0
loop14            0.00         0.11         0.00         0.00     261155          0          0
loop15            0.01         0.18         0.00         0.00     424935          0          0
loop16            0.02         0.40         0.00         0.00     969406          0          0
loop17            0.01         0.08         0.00         0.00     190895          0          0
loop18            0.01         0.08         0.00         0.00     191004          0          0
loop19            0.01         0.10         0.00         0.00     243290          0          0
loop2             0.00         0.00         0.00         0.00         17          0          0
loop20            0.01         0.10         0.00         0.00     243371          0          0
loop21            0.01         0.16         0.00         0.00     378857          0          0
loop22            0.03         0.74         0.00         0.00    1781497          0          0
loop23            0.01         0.21         0.00         0.00     497402          0          0
loop24            0.02         0.58         0.00         0.00    1400747          0          0
loop25            0.00         0.00         0.00         0.00       2328          0          0
loop26            0.00         0.00         0.00         0.00       2458          0          0
loop27            0.00         0.00         0.00         0.00        268          0          0
loop28            0.00         0.00         0.00         0.00        187          0          0
loop29            0.02         0.06         0.00         0.00     154997          0          0
loop3             0.00         0.07         0.00         0.00     171119          0          0
loop30            0.03         0.11         0.00         0.00     253868          0          0
loop31            0.01         0.12         0.00         0.00     299770          0          0
loop32            0.01         0.14         0.00         0.00     330542          0          0
loop33            0.00         0.02         0.00         0.00      59062          0          0
loop34            0.00         0.06         0.00         0.00     154983          0          0
loop35            0.00         0.00         0.00         0.00        504          0          0
loop36            0.00         0.00         0.00         0.00       2090          0          0
loop37            0.12         3.07         0.00         0.00    7380549          0          0
loop38            0.29        16.43         0.00         0.00   39541096          0          0
loop4             0.00         0.07         0.00         0.00     164613          0          0
loop5             0.00         0.05         0.00         0.00     130233          0          0
loop6             0.04         1.36         0.00         0.00    3270977          0          0
loop7             0.00         0.03         0.00         0.00      63369          0          0
loop8             0.00         0.03         0.00         0.00      63564          0          0
loop9             0.02         0.96         0.00         0.00    2318794          0          0
nvme0n1          53.05       512.13       640.93       109.00 1232624135 1542636186  262339888
```

The jobset, associated jobs, and pods will be completed:

```bash
$ kubectl get jobset
```
```console
NAME               RESTARTS   COMPLETED   AGE
metricset-sample              True        2m28s
```
```bash
$ kubectl get jobs
```
```console
NAME                   COMPLETIONS   DURATION   AGE
metricset-sample-m-0   1/1           13s        2m51s
```
```bash
$ kubectl get pods
```
```console
NAME                           READY   STATUS      RESTARTS   AGE
metricset-sample-m-0-0-rq4q9   0/1     Completed   0          3m19s
```

When you are done, cleanup!

```bash
kubectl delete -f metrics.yaml
kubectl delete cm metricset-sample
```