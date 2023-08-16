# Host Volume IO Example

This is our first example of running a storage metric on a simple host volume.
We will first create the volume, and then reference it for our metric of interest.
For running the example and parsing output, see [the corresponding Python directory](../../python/io-host-volume/).

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

We would typically apply a manifest for storage, but since we are using a host volume, we don't need to do that (the host is the kind container that we will bind to). 
Next, let's create the metrics set. This is going to run a simple sysstat io tool to collect metrics about our volumes.

```bash
kubectl apply -f metrics.yaml
```

Note that the metrics.yaml is asking to run a sysstat io metric for 10 completions, with a delay of 10 seconds between each check (10). Wait until you see a pod created by the job and then running (there should be one since we've only asked for one metric):

```bash
kubectl get pods
```
```diff
NAME                         READY   STATUS              RESTARTS   AGE
metricset-sample-m-0-tjxsj   0/1     ContainerCreating   0          25s
metricset-sample-m-0-tjxsj   1/1     Running             0          70s
```

If you peek at logs, you'll see the storage metric running once every 10 seconds, and for a total of 10 times.

```bash
$ kubectl logs metricset-sample-m-0-9pq6w
```

By default, non human readable output is presented in blocks of json, with appropriate sections
added by the Metrics operator for easy parsing by the metricsoperator Python module.

<details>

<summary>Default output in JSON</summary>

```console
root
METADATA START {"pods":2,"completions":2,"metricName":"network-netmark","metricDescription":"point to point networking tool","metricType":"standalone","metricOptions":{"messageSize":0,"rate":10,"sendReceiveCycles":20,"storeEachTrial":"true","tasks":2,"trials":20,"warmups":10}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
(env) (base) vanessa@vanessa-ThinkPad-T490s:~/Desktop/Code/metrics-operator$ kubectl logs metricset-sample-m-0-tr58z -f
Error from server (BadRequest): container "io-sysstat" in pod "metricset-sample-m-0-tr58z" is waiting to start: ContainerCreating
(env) (base) vanessa@vanessa-ThinkPad-T490s:~/Desktop/Code/metrics-operator$ kubectl logs metricset-sample-m-0-tr58z -f
Error from server (BadRequest): container "io-sysstat" in pod "metricset-sample-m-0-tr58z" is waiting to start: ContainerCreating
(env) (base) vanessa@vanessa-ThinkPad-T490s:~/Desktop/Code/metrics-operator$ kubectl logs metricset-sample-m-0-tr58z -f
Error from server (BadRequest): container "io-sysstat" in pod "metricset-sample-m-0-tr58z" is waiting to start: ContainerCreating
(env) (base) vanessa@vanessa-ThinkPad-T490s:~/Desktop/Code/metrics-operator$ kubectl logs metricset-sample-m-0-tr58z -f
Error from server (BadRequest): container "io-sysstat" in pod "metricset-sample-m-0-tr58z" is waiting to start: ContainerCreating
(env) (base) vanessa@vanessa-ThinkPad-T490s:~/Desktop/Code/metrics-operator$ kubectl logs metricset-sample-m-0-tr58z -f
METADATA START {"pods":1,"completions":1,"storageVolumePath":"/workflow","storageVolumeHostPath":"/tmp/workflow","metricName":"io-sysstat","metricDescription":"statistics for Linux tasks (processes) : I/O, CPU, memory, etc.","metricType":"storage","metricOptions":{"completions":2,"human":"false"}}
METADATA END
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
{"sysstat": {
        "hosts": [
                {
                        "nodename": "metricset-sample-m-0-tr58z.ms.default.svc.cluster.local",
                        "sysname": "Linux",
                        "release": "5.15.0-78-generic",
                        "machine": "x86_64",
                        "number-of-cpus": 8,
                        "date": "08/07/23",
                        "statistics": [
                                {
                                        "disk": [
                                                {"disk_device": "loop0", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.81, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop1", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.30, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 14.25, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop10", "r/s": 0.08, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.27, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 28.98, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.02},
                                                {"disk_device": "loop11", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 16.35, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop12", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.44, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 36.08, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop13", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 17.12, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop14", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 17.69, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop15", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.30, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 18.67, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop16", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.26, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 12.18, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop17", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.29, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop18", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.26, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 13.82, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop19", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.24, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.98, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop2", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.00, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 1.21, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop20", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.33, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 13.65, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop21", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.83, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop22", "r/s": 0.38, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.01, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 24.02, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.03},
                                                {"disk_device": "loop23", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.20, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop24", "r/s": 0.01, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.33, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 42.26, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop25", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 3.31, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop26", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.06, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.76, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop27", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.18, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 1.79, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop28", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.11, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.58, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop29", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 18.67, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop3", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.11, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 12.24, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop30", "r/s": 0.01, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.03, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 5.88, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop31", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.15, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.81, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop32", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.32, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 13.89, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop33", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.12, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 8.08, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop34", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.52, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 31.38, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop35", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.05, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.70, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop36", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 11.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop37", "r/s": 0.86, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.05, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.29, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 55.61, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.10},
                                                {"disk_device": "loop4", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.23, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 15.54, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop5", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.19, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.85, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop6", "r/s": 0.07, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.21, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 37.43, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.02},
                                                {"disk_device": "loop7", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 8.12, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop8", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.26, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 5.57, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop9", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 7.76, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "nvme0n1", "r/s": 43.43, "w/s": 39.03, "d/s": 0.72, "f/s": 4.33, "rMB/s": 0.86, "wMB/s": 0.83, "dMB/s": 0.21, "rrqm/s": 23.20, "wrqm/s": 30.94, "drqm/s": 0.00, "rrqm": 34.82, "wrqm": 44.22, "drqm": 0.00, "r_await": 0.34, "w_await": 5.05, "d_await": 0.31, "f_await": 0.38, "rareq-sz": 20.31, "wareq-sz": 21.80, "dareq-sz": 298.96, "aqu-sz": 0.21, "util": 3.93}
                                        ]
                                }
                        ]
                }
        ]
}}
METRICS OPERATOR TIMEPOINT
{"sysstat": {
        "hosts": [
                {
                        "nodename": "metricset-sample-m-0-tr58z.ms.default.svc.cluster.local",
                        "sysname": "Linux",
                        "release": "5.15.0-78-generic",
                        "machine": "x86_64",
                        "number-of-cpus": 8,
                        "date": "08/07/23",
                        "statistics": [
                                {
                                        "disk": [
                                                {"disk_device": "loop0", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.81, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop1", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.30, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 14.25, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop10", "r/s": 0.08, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.27, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 28.98, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.02},
                                                {"disk_device": "loop11", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 16.35, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop12", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.44, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 36.08, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop13", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 17.12, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop14", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 17.69, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop15", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.30, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 18.67, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop16", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.26, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 12.18, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop17", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.29, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop18", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.26, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 13.82, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop19", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.24, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.98, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop2", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.00, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 1.21, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop20", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.33, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 13.65, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop21", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.83, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop22", "r/s": 0.38, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.01, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 24.02, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.03},
                                                {"disk_device": "loop23", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.20, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop24", "r/s": 0.01, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.33, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 42.26, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop25", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 3.31, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop26", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.06, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.76, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop27", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.18, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 1.79, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop28", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.11, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.58, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop29", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 18.67, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop3", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.11, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 12.24, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop30", "r/s": 0.01, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.03, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 5.88, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop31", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.15, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.81, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop32", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.32, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 13.89, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop33", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.12, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 8.08, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop34", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.52, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 31.38, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop35", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.05, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.70, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop36", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 11.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop37", "r/s": 0.86, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.05, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.29, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 55.61, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.10},
                                                {"disk_device": "loop4", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.23, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 15.54, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop5", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.19, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.85, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop6", "r/s": 0.07, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.21, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 37.43, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.02},
                                                {"disk_device": "loop7", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 8.12, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop8", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.26, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 5.57, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop9", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 7.76, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "nvme0n1", "r/s": 43.43, "w/s": 39.03, "d/s": 0.72, "f/s": 4.33, "rMB/s": 0.86, "wMB/s": 0.83, "dMB/s": 0.21, "rrqm/s": 23.20, "wrqm/s": 30.94, "drqm/s": 0.00, "rrqm": 34.82, "wrqm": 44.22, "drqm": 0.00, "r_await": 0.34, "w_await": 5.05, "d_await": 0.31, "f_await": 0.38, "rareq-sz": 20.31, "wareq-sz": 21.81, "dareq-sz": 298.96, "aqu-sz": 0.21, "util": 3.93}
                                        ]
                                }
                        ]
                }
        ]
}}
METRICS OPERATOR TIMEPOINT
{"sysstat": {
        "hosts": [
                {
                        "nodename": "metricset-sample-m-0-tr58z.ms.default.svc.cluster.local",
                        "sysname": "Linux",
                        "release": "5.15.0-78-generic",
                        "machine": "x86_64",
                        "number-of-cpus": 8,
                        "date": "08/07/23",
                        "statistics": [
                                {
                                        "disk": [
                                                {"disk_device": "loop0", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.81, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop1", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.30, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 14.25, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop10", "r/s": 0.08, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.27, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 28.98, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.02},
                                                {"disk_device": "loop11", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 16.35, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop12", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.44, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 36.08, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop13", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 17.12, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop14", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 17.69, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop15", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.30, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 18.67, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop16", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.26, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 12.18, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop17", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.29, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop18", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.26, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 13.82, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop19", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.24, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.98, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop2", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.00, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 1.21, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop20", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.33, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 13.65, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop21", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.83, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop22", "r/s": 0.38, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.01, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 24.02, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.03},
                                                {"disk_device": "loop23", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.20, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop24", "r/s": 0.01, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.33, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 42.26, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop25", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 3.31, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop26", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.06, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.76, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop27", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.18, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 1.79, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop28", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.11, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.58, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop29", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 18.67, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop3", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.11, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 12.24, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop30", "r/s": 0.01, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.03, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 5.88, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop31", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.15, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.81, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop32", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.32, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 13.89, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop33", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.12, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 8.08, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop34", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.52, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 31.38, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop35", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.05, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.70, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop36", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 11.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop37", "r/s": 0.86, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.05, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.29, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 55.61, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.10},
                                                {"disk_device": "loop4", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.23, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 15.54, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop5", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.19, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.85, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop6", "r/s": 0.07, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.21, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 37.43, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.02},
                                                {"disk_device": "loop7", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 8.12, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop8", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.26, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 5.57, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop9", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 7.76, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "nvme0n1", "r/s": 43.43, "w/s": 39.03, "d/s": 0.72, "f/s": 4.33, "rMB/s": 0.86, "wMB/s": 0.83, "dMB/s": 0.21, "rrqm/s": 23.20, "wrqm/s": 30.94, "drqm/s": 0.00, "rrqm": 34.82, "wrqm": 44.22, "drqm": 0.00, "r_await": 0.34, "w_await": 5.05, "d_await": 0.31, "f_await": 0.38, "rareq-sz": 20.31, "wareq-sz": 21.81, "dareq-sz": 298.96, "aqu-sz": 0.21, "util": 3.93}
                                        ]
                                }
                        ]
                }
        ]
}}
METRICS OPERATOR COLLECTION END
```

</details>


But if you add an option for human readable, e.g.,

```yaml
metrics:
  - name: io-sysstat
    rate: 10
    completions: 2
    options:
      human: "true"
```

You'll see a more tabular format:

<details>

<summary>Output with options->human set to "true"</summary>

```console
METADATA START {"pods":1,"completions":1,"storageVolumePath":"/workflow","storageVolumeHostPath":"/tmp/workflow","metricName":"io-sysstat","metricDescription":"statistics for Linux tasks (processes) : I/O, CPU, memory, etc.","metricType":"storage","metricOptions":{"completions":2,"human":"true"}}
METADATA END
METRICS OPERATOR COLLECTION START
...
METRICS OPERATOR TIMEPOINT
Linux 5.15.0-78-generic (metricset-sample-m-0-xc8v6.ms.default.svc.cluster.local)       08/07/23        _x86_64_        (8 CPU)

Device            r/s     rMB/s   rrqm/s  %rrqm r_await rareq-sz     w/s     wMB/s   wrqm/s  %wrqm w_await wareq-sz     d/s     dMB/s   drqm/s  %drqm d_await dareq-sz     f/s f_await  aqu-sz  %util
loop0            0.00      0.00     0.00   0.00    0.17    19.81    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop1            0.00      0.00     0.00   0.00    0.30    14.25    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop10           0.08      0.00     0.00   0.00    0.27    28.98    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.02
loop11           0.00      0.00     0.00   0.00    0.16    16.35    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop12           0.00      0.00     0.00   0.00    0.44    36.08    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop13           0.00      0.00     0.00   0.00    0.17    17.12    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop14           0.00      0.00     0.00   0.00    0.17    17.69    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop15           0.00      0.00     0.00   0.00    0.30    18.67    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop16           0.00      0.00     0.00   0.00    0.26    12.18    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop17           0.00      0.00     0.00   0.00    0.29    19.47    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop18           0.00      0.00     0.00   0.00    0.26    13.82    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop19           0.00      0.00     0.00   0.00    0.24    19.98    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop2            0.00      0.00     0.00   0.00    0.00     1.21    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop20           0.00      0.00     0.00   0.00    0.33    13.65    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop21           0.00      0.00     0.00   0.00    0.17    20.83    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop22           0.38      0.01     0.00   0.00    0.17    24.02    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.03
loop23           0.00      0.00     0.00   0.00    0.20    19.47    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop24           0.01      0.00     0.00   0.00    0.33    42.26    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop25           0.00      0.00     0.00   0.00    0.14     3.31    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop26           0.00      0.00     0.00   0.00    0.06     2.76    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop27           0.00      0.00     0.00   0.00    0.18     1.79    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop28           0.00      0.00     0.00   0.00    0.11     2.58    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop29           0.00      0.00     0.00   0.00    0.16    18.67    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop3            0.00      0.00     0.00   0.00    0.11    12.24    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop30           0.01      0.00     0.00   0.00    0.03     5.88    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop31           0.00      0.00     0.00   0.00    0.15    20.81    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop32           0.00      0.00     0.00   0.00    0.32    13.89    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop33           0.00      0.00     0.00   0.00    0.12     8.08    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop34           0.00      0.00     0.00   0.00    0.52    31.38    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop35           0.00      0.00     0.00   0.00    0.05     2.70    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop36           0.00      0.00     0.00   0.00    0.14    11.47    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop37           0.86      0.05     0.00   0.00    0.29    55.61    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.10
loop4            0.00      0.00     0.00   0.00    0.23    15.54    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop5            0.00      0.00     0.00   0.00    0.19    20.85    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop6            0.07      0.00     0.00   0.00    0.21    37.43    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.02
loop7            0.00      0.00     0.00   0.00    0.14     8.12    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop8            0.00      0.00     0.00   0.00    0.26     5.57    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
loop9            0.00      0.00     0.00   0.00    0.16     7.76    0.00      0.00     0.00   0.00    0.00     0.00    0.00      0.00     0.00   0.00    0.00     0.00    0.00    0.00    0.00   0.00
nvme0n1         43.43      0.86    23.19  34.82    0.34    20.31   39.03      0.83    30.95  44.22    5.05    21.81    0.72      0.21     0.00   0.00    0.31   298.96    4.33    0.38    0.21   3.93
METRICS OPERATOR COLLECTION END
```

</details>

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
```
