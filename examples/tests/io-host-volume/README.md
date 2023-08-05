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

By default, non human readable output is presented in blocks of json:

<details>

<summary>Default output in JSON</summary>

```console
IOSTAT TIMEPOINT 2
{"sysstat": {
        "hosts": [
                {
                        "nodename": "metricset-sample-m-0-vht9l.ms.default.svc.cluster.local",
                        "sysname": "Linux",
                        "release": "5.15.0-78-generic",
                        "machine": "x86_64",
                        "number-of-cpus": 8,
                        "date": "08/05/23",
                        "statistics": [
                                {
                                        "disk": [
                                                {"disk_device": "loop0", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.81, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop1", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.29, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 15.26, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop10", "r/s": 0.07, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.29, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 29.14, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.02},
                                                {"disk_device": "loop11", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.15, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 17.01, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop12", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.47, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 30.58, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop13", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 17.12, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop14", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 17.69, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop15", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.30, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 18.67, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop16", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.25, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 12.75, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop17", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.29, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop18", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.24, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 14.75, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop19", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.24, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.98, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop2", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.00, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 1.21, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop20", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.32, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 14.83, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop21", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.17, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.83, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop22", "r/s": 0.30, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.01, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 22.91, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.02},
                                                {"disk_device": "loop23", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.20, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 19.47, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop24", "r/s": 0.01, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.33, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 36.41, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop25", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 3.31, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop26", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.07, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.84, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop27", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.15, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 1.91, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop28", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.11, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.58, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop29", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 18.67, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop3", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.11, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 12.48, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop30", "r/s": 0.01, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.03, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 5.45, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop31", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.15, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.81, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop32", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.32, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 14.91, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop33", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.12, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 8.08, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop34", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.56, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 29.21, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop35", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.05, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 2.70, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop36", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.10, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 10.85, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop37", "r/s": 0.69, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.04, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.27, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 55.50, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.08},
                                                {"disk_device": "loop4", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.23, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 15.54, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop5", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.19, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 20.85, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop6", "r/s": 0.06, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.19, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 37.32, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.02},
                                                {"disk_device": "loop7", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.14, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 8.12, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop8", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.24, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 6.02, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "loop9", "r/s": 0.00, "w/s": 0.00, "d/s": 0.00, "f/s": 0.00, "rMB/s": 0.00, "wMB/s": 0.00, "dMB/s": 0.00, "rrqm/s": 0.00, "wrqm/s": 0.00, "drqm/s": 0.00, "rrqm": 0.00, "wrqm": 0.00, "drqm": 0.00, "r_await": 0.16, "w_await": 0.00, "d_await": 0.00, "f_await": 0.00, "rareq-sz": 7.76, "wareq-sz": 0.00, "dareq-sz": 0.00, "aqu-sz": 0.00, "util": 0.00},
                                                {"disk_device": "nvme0n1", "r/s": 34.53, "w/s": 38.88, "d/s": 0.00, "f/s": 4.12, "rMB/s": 0.73, "wMB/s": 0.88, "dMB/s": 0.00, "rrqm/s": 17.17, "wrqm/s": 31.12, "drqm/s": 0.00, "rrqm": 33.21, "wrqm": 44.45, "drqm": 0.00, "r_await": 0.33, "w_await": 5.25, "d_await": 0.00, "f_await": 0.38, "rareq-sz": 21.51, "wareq-sz": 23.11, "dareq-sz": 0.00, "aqu-sz": 0.22, "util": 3.54}
                                        ]
                                }
                        ]
                }
        ]
}}
```

</details>


But if you add an option for human readable, e.g.,

```yaml
metrics:
  - name: io-sysstat
    rate: 10
    completions: 2
    options:
      human: true
```

You'll see a more tabular format:

<details>

<summary>Output with options->human set to "true"</summary>

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