# Ior IO Example

This will show running [Ior](https://github.com/hpc/ior) on a host volume.
We will first create the volume, and then reference it for our metric of interest.

![../../../docs/getting_started/img/ior.jpeg](../../../docs/getting_started/img/ior.jpeg )

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

You should be able to see the main pod - it will be mounted to our volume of interest.

```bash
kubectl get pods
```
```console
NAME                         READY   STATUS      RESTARTS   AGE
metricset-sample-m-0-29fn8   0/1     Completed   0          38s
```

And then look at the logs...

```bash
$ kubectl logs metricset-sample-m-0-9pq6w
```

And see the ior result!

<details>

<summary>Example output for Ior</summary>

```console
METADATA START {"pods":1,"completions":1,"storageVolumePath":"/tmp/workflow","storageVolumeHostPath":"/tmp/workflow","metricName":"io-ior","metricDescription":"HPC IO Benchmark","metricType":"storage","metricOptions":{"command":"ior -w -r -o testfile","workdir":"/opt/ior"}}
METADATA END
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
IOR-4.0.0rc1: MPI Coordinated Test of Parallel I/O
Began               : Wed Aug 23 23:40:51 2023
Command line        : ior -w -r -o testfile
Machine             : Linux metricset-sample-m-0-29fn8.ms.default.svc.cluster.local
TestID              : 0
StartTime           : Wed Aug 23 23:40:51 2023
Path                : testfile
FS                  : 467.9 GiB   Used FS: 85.9%   Inodes: 29.8 Mi   Used Inodes: 21.1%

Options: 
api                 : POSIX
apiVersion          : 
test filename       : testfile
access              : single-shared-file
type                : independent
segments            : 1
ordering in a file  : sequential
ordering inter file : no tasks offsets
nodes               : 1
tasks               : 1
clients per node    : 1
repetitions         : 1
xfersize            : 262144 bytes
blocksize           : 1 MiB
aggregate filesize  : 1 MiB

Results: 

access    bw(MiB/s)  IOPS       Latency(s)  block(KiB) xfer(KiB)  open(s)    wr/rd(s)   close(s)   total(s)   iter
------    ---------  ----       ----------  ---------- ---------  --------   --------   --------   --------   ----
write     443.65     1981.01    0.000505    1024.00    256.00     0.000229   0.002019   0.000006   0.002254   0   
read      3584.88    15033      0.000067    1024.00    256.00     0.000009   0.000266   0.000004   0.000279   0   

Summary of all tests:
Operation   Max(MiB)   Min(MiB)  Mean(MiB)     StdDev   Max(OPs)   Min(OPs)  Mean(OPs)     StdDev    Mean(s) Stonewall(s) Stonewall(MiB) Test# #Tasks tPN reps fPP reord reordoff reordrand seed segcnt   blksiz    xsize aggs(MiB)   API RefNum
write         443.65     443.65     443.65       0.00    1774.62    1774.62    1774.62       0.00    0.00225         NA            NA     0      1   1    1   0     0        1         0    0      1  1048576   262144       1.0 POSIX      0
read         3584.88    3584.88    3584.88       0.00   14339.50   14339.50   14339.50       0.00    0.00028         NA            NA     0      1   1    1   0     0        1         0    0      1  1048576   262144       1.0 POSIX      0
Finished            : Wed Aug 23 23:40:51 2023
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
