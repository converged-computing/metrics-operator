# OSU Benchmarks Example

This will demonstrate running [OSU Benchmarks](https://mvapich.cse.ohio-state.edu/benchmarks/) with the Metrics Operator.
For running the example, parsing, and plotting output, see [the corresponding Python directory](../../python/network-osu-benchmark/).

## Usage

### 1. Prepare Cluster

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

### 2. Create the Metric Set

Then create the metrics set. This is going to run a simple sysstat tool to collect metrics
as lammps runs.

```bash
kubectl apply -f metrics.yaml
```

Wait until you see pods created by the job and then running (there should be two, as OSU benchmarks need exactly two!):

```bash
kubectl get pods
```
```diff
NAME                           READY   STATUS    RESTARTS   AGE
metricset-sample-l-0-0-lt782   1/1     Running   0          3s
metricset-sample-w-0-0-4s5p9   1/1     Running   0          3s
```

In the above, "w" is a worker pod, and "l" is the launcher that runs the mpirun commands.
If you inspect the log for the launcher you'll see a short sleep (the network isn't up immediately)
and then the benchmarks running. The output is structured in a predictable format by the Metrics
Operator so it can easily be parsed by the metricsoperator Python module.
By default, we don't include a list in metrics.yaml so we run them all!

```bash
kubectl logs metricset-sample-l-0-0-lt782 -f
```

<details>

<summary>Output of OSU Benchmarks Launcher</summary>

```console
Sleeping for 10 seconds waiting for network...
METADATA START {"pods":2,"completions":2,"metricName":"network-osu-benchmark","metricDescription":"point to point MPI benchmarks","metricType":"standalone","metricOptions":{"completions":0,"rate":10},"metricListOptions":{"commands":["/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_get_acc_latency","/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_acc_latency","/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_fop_latency","/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_get_latency","/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_put_latency","/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/collective/osu_allreduce","/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt/osu_latency","/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt/osu_bibw","/opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt/osu_bw"]}}
METADATA END
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_get_acc_latency
# OSU MPI_Get_accumulate latency Test v5.8
# Window creation: MPI_Win_create
# Synchronization: MPI_Win_lock/unlock
# Size          Latency (us)
1                       2.78
2                       1.65
4                       1.56
8                       1.55
16                      1.53
32                      1.61
64                      1.55
128                     1.61
256                     1.98
512                     1.87
1024                    2.23
2048                    2.97
4096                    4.06
8192                    7.58
16384                  15.03
32768                  29.06
65536                  57.17
131072                113.65
262144                217.74
524288                442.54
1048576               873.24
2097152              1838.53
4194304              4118.13
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_acc_latency
# OSU MPI_Accumulate latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                       0.48
2                       0.45
4                       0.50
8                       0.49
16                      0.46
32                      0.36
64                      0.31
128                     0.34
256                     0.43
512                     0.52
1024                    0.79
2048                    1.34
4096                    2.37
8192                    4.84
16384                   9.60
32768                  18.92
65536                  34.48
131072                 69.30
262144                135.87
524288                264.04
1048576               548.64
2097152              1416.32
4194304              2195.60
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_fop_latency
# OSU MPI_Fetch_and_op latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
8                       0.62
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_get_latency
# OSU MPI_Get latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                       0.20
2                       0.16
4                       0.16
8                       0.14
16                      0.14
32                      0.16
64                      0.13
128                     0.13
256                     0.12
512                     0.13
1024                    0.11
2048                    0.12
4096                    0.21
8192                    0.23
16384                   0.41
32768                   0.86
65536                   1.82
131072                  3.49
262144                  9.07
524288                 17.50
1048576                35.53
2097152                70.06
4194304               269.72
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_put_latency
# OSU MPI_Put Latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                       0.44
2                       0.42
4                       0.31
8                       0.37
16                      0.28
32                      0.24
64                      0.22
128                     0.22
256                     0.17
512                     0.18
1024                    0.15
2048                    0.14
4096                    0.17
8192                    0.23
16384                   0.59
32768                   1.11
65536                   2.78
131072                  3.92
262144                  8.96
524288                 20.62
1048576                36.66
2097152                72.93
4194304               275.19
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/collective/osu_allreduce

# OSU MPI Allreduce Latency Test v5.8
# Size       Avg Latency(us)
4                       1.34
8                       1.17
16                      0.89
32                      0.96
64                      0.94
128                     0.94
256                     1.03
512                     1.41
1024                    1.59
2048                    1.80
4096                    4.42
8192                    5.30
16384                   6.44
32768                  10.36
65536                  16.93
131072                 33.91
262144                 59.87
524288                110.83
1048576               230.90
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt/osu_latency
# OSU MPI Latency Test v5.8
# Size          Latency (us)
0                       0.80
1                       0.46
2                       0.32
4                       0.29
8                       0.24
16                      0.22
32                      0.26
64                      0.26
128                     0.29
256                     0.33
512                     0.41
1024                    0.51
2048                    0.61
4096                    1.31
8192                    1.85
16384                   2.59
32768                   2.93
65536                   5.62
131072                  9.28
262144                 16.55
524288                 30.79
1048576                58.15
2097152               184.02
4194304               584.22
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt/osu_bibw
# OSU MPI Bi-Directional Bandwidth Test v5.8
# Size      Bandwidth (MB/s)
1                       2.41
2                       3.94
4                      11.50
8                      32.48
16                     57.95
32                    140.93
64                    280.55
128                   562.17
256                   799.07
512                  2059.54
1024                 3287.79
2048                 7635.69
4096                 5608.39
8192                 9489.32
16384                9812.43
32768               17272.36
65536               26445.32
131072              31649.61
262144              32948.58
524288              32199.06
1048576             31963.81
2097152             17770.85
4194304              9360.80
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt/osu_bw
# OSU MPI Bandwidth Test v5.8
# Size      Bandwidth (MB/s)
1                       1.83
2                       3.45
4                       5.59
8                      13.45
16                     33.92
32                     75.76
64                    165.74
128                   326.71
256                   683.73
512                  1303.41
1024                 3053.38
2048                 5100.77
4096                 3959.97
8192                 6322.19
16384                8227.68
32768               12834.55
65536               16996.58
131072              17358.56
262144              17617.89
524288              16931.85
1048576             18476.89
2097152             18718.87
4194304             11256.10
METRICS OPERATOR COLLECTION END
```

</details>

The worker comes up and sleeps, and will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
oot
Sleeping for 10 seconds waiting for network...
METADATA START {"pods":2,"completions":2,"metricName":"network-osu-benchmark","metricDescription":"point to point MPI benchmarks","metricType":"standalone","metricOptions":{"completions":0,"rate":10},"metricListOptions":{"commands":["osu_fop_latency","osu_get_acc_latency","osu_get_latency","osu_put_latency","osu_acc_latency"]}}
METADATA END
```

We can do this with JobSet logic that the entire set is done when the launcher is done.

```bash
$ kubectl get pods
```
```console
NAME                           READY   STATUS        RESTARTS   AGE
metricset-sample-l-0-0-bqqf4   0/1     Completed     0          49s
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
metricset-sample-l-0   1/1           18s        84s
```

### 3. Save Configs

If you want to run a single one-off Metrics set without the operator, you can save configs. We did this to provide a single YAML to run the entire thing.

```bash
kubectl get jobset metricset-sample -o yaml > ./manifests/jobset.yaml
kubectl get cm metricset-sample -o yaml > ./manifests/configmap.yaml
kubectl get svc ms -o yaml > ./manifests/service.yaml
```

**important** you will want to edit the resources (limits and requests) for your instance size!
I also tweaked the initial files to remove the status / controller reference from those. And then you could create the entire directory as follows:

```bash
$ kubectl apply -f ./manifests
```

And delete when you are done:

```bash
$ kubectl delete -f ./manifests
```

### 4. Cleanup

And then you can cleanup!

```bash
kubectl delete -f metrics.yaml
```

Note that you can tweak the metrics.yaml to ask for specific metrics for OSU.
If you don't define any, you'll get the default list we provide. See
the metrics.yaml for details.
