# Fio IO Example

This will test running [Fio](https://fio.readthedocs.io/en/latest/fio_doc.html) on a host volume.
We will first create the volume, and then reference it for our metric of interest.
For running the example and parsing output, see [the corresponding Python directory](../../python/io-fio/).

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

Next, let's create the metrics set. This is going to run a simple sysstat io tool to collect metrics about our volumes.

```bash
kubectl apply -f metrics.yaml
```

You should be able to see pods:

```bash
kubectl get pods
```
```diff
NAME                         READY   STATUS              RESTARTS   AGE
- metricset-sample-m-0-tjxsj   0/1     ContainerCreating   0          25s
+ metricset-sample-m-0-tjxsj   1/1     Running             0          70s
```

And then look at the logs...

```bash
$ kubectl logs metricset-sample-m-0-9pq6w
```

And see the fio result!

<details>

<summary>Default output in JSON</summary>

```console
$ kubectl logs metricset-sample-m-0-4x56g 
METADATA START {"pods":1,"completions":1,"storageVolumePath":"/workflow","storageVolumeHostPath":"/tmp/workflow","metricName":"io-fio","metricDescription":"Flexible IO Tester (FIO)","metricType":"storage","metricOptions":{"blocksize":"4k","completions":0,"directory":"/tmp","iodepth":64,"size":"4G","testname":"test"}}
METADATA END
FIO COMMAND START
fio --randrepeat=1 --ioengine=libaio --direct=1 --gtod_reduce=1 --name=test --bs=4k --iodepth=64 --readwrite=randrw --rwmixread=75 --size=4G --filename=/tmp/test-b273108fb88ca182ac07dad8b6fe4e61 --output-format=json
FIO COMMAND END
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
{
  "fio version" : "fio-3.28",
  "timestamp" : 1691554888,
  "timestamp_ms" : 1691554888193,
  "time" : "Wed Aug  9 04:21:28 2023",
  "global options" : {
    "randrepeat" : "1",
    "ioengine" : "libaio",
    "direct" : "1",
    "gtod_reduce" : "1"
  },
  "jobs" : [
    {
      "jobname" : "test",
      "groupid" : 0,
      "error" : 0,
      "eta" : 0,
      "elapsed" : 5,
      "job options" : {
        "name" : "test",
        "bs" : "4k",
        "iodepth" : "64",
        "rw" : "randrw",
        "rwmixread" : "75",
        "size" : "4G",
        "filename" : "/tmp/test-b273108fb88ca182ac07dad8b6fe4e61"
      },
      "read" : {
        "io_bytes" : 3219128320,
        "io_kbytes" : 3143680,
        "bw_bytes" : 696630235,
        "bw" : 680302,
        "iops" : 170075.741182,
        "runtime" : 4621,
        "total_ios" : 785920,
        "short_ios" : 0,
        "drop_ios" : 0,
        "slat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "clat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "lat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "bw_min" : 442936,
        "bw_max" : 784496,
        "bw_agg" : 99.808032,
        "bw_mean" : 678997.555556,
        "bw_dev" : 104924.785107,
        "bw_samples" : 9,
        "iops_min" : 110734,
        "iops_max" : 196124,
        "iops_mean" : 169749.333333,
        "iops_stddev" : 26231.194807,
        "iops_samples" : 9
      },
      "write" : {
        "io_bytes" : 1075838976,
        "io_kbytes" : 1050624,
        "bw_bytes" : 232815186,
        "bw" : 227358,
        "iops" : 56839.645098,
        "runtime" : 4621,
        "total_ios" : 262656,
        "short_ios" : 0,
        "drop_ios" : 0,
        "slat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "clat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "lat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "bw_min" : 150008,
        "bw_max" : 264112,
        "bw_agg" : 99.837886,
        "bw_mean" : 226990.111111,
        "bw_dev" : 34718.495735,
        "bw_samples" : 9,
        "iops_min" : 37502,
        "iops_max" : 66028,
        "iops_mean" : 56747.444444,
        "iops_stddev" : 8679.626506,
        "iops_samples" : 9
      },
      "trim" : {
        "io_bytes" : 0,
        "io_kbytes" : 0,
        "bw_bytes" : 0,
        "bw" : 0,
        "iops" : 0.000000,
        "runtime" : 0,
        "total_ios" : 0,
        "short_ios" : 0,
        "drop_ios" : 0,
        "slat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "clat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "lat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "bw_min" : 0,
        "bw_max" : 0,
        "bw_agg" : 0.000000,
        "bw_mean" : 0.000000,
        "bw_dev" : 0.000000,
        "bw_samples" : 0,
        "iops_min" : 0,
        "iops_max" : 0,
        "iops_mean" : 0.000000,
        "iops_stddev" : 0.000000,
        "iops_samples" : 0
      },
      "sync" : {
        "total_ios" : 0,
        "lat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        }
      },
      "job_runtime" : 4620,
      "usr_cpu" : 16.709957,
      "sys_cpu" : 64.523810,
      "ctx" : 236491,
      "majf" : 0,
      "minf" : 9,
      "iodepth_level" : {
        "1" : 0.100000,
        "2" : 0.100000,
        "4" : 0.100000,
        "8" : 0.100000,
        "16" : 0.100000,
        "32" : 0.100000,
        ">=64" : 99.993992
      },
      "iodepth_submit" : {
        "0" : 0.000000,
        "4" : 100.000000,
        "8" : 0.000000,
        "16" : 0.000000,
        "32" : 0.000000,
        "64" : 0.000000,
        ">=64" : 0.000000
      },
      "iodepth_complete" : {
        "0" : 0.000000,
        "4" : 99.999905,
        "8" : 0.000000,
        "16" : 0.000000,
        "32" : 0.000000,
        "64" : 0.100000,
        ">=64" : 0.000000
      },
      "latency_ns" : {
        "2" : 0.000000,
        "4" : 0.000000,
        "10" : 0.000000,
        "20" : 0.000000,
        "50" : 0.000000,
        "100" : 0.000000,
        "250" : 0.000000,
        "500" : 0.000000,
        "750" : 0.000000,
        "1000" : 0.000000
      },
      "latency_us" : {
        "2" : 0.000000,
        "4" : 0.000000,
        "10" : 0.000000,
        "20" : 0.000000,
        "50" : 0.000000,
        "100" : 0.000000,
        "250" : 0.000000,
        "500" : 0.000000,
        "750" : 0.000000,
        "1000" : 0.000000
      },
      "latency_ms" : {
        "2" : 0.000000,
        "4" : 0.000000,
        "10" : 0.000000,
        "20" : 0.000000,
        "50" : 0.000000,
        "100" : 0.000000,
        "250" : 0.000000,
        "500" : 0.000000,
        "750" : 0.000000,
        "1000" : 0.000000,
        "2000" : 0.000000,
        ">=2000" : 0.000000
      },
      "latency_depth" : 64,
      "latency_target" : 0,
      "latency_percentile" : 100.000000,
      "latency_window" : 0
    }
  ]
}
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
