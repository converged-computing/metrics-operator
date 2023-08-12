# OSU Benchmarks Example

This will demonstrate running [OSU Benchmarks](https://mvapich.cse.ohio-state.edu/benchmarks/) with the Metrics Operator.
Note that I'm still tweaking the mpirun arguments for each command (I'm learning too)!
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
Number of tasks (nproc on one node) is 8
Number of tasks total (across 2 nodes) is 16
Sleeping for 10 seconds waiting for network...
METADATA START {"pods":2,"completions":2,"metricName":"network-osu-benchmark","metricDescription":"point to point MPI benchmarks","metricType":"standalone","metricOptions":{"completions":0,"rate":10,"tasks":0},"metricListOptions":{"commands":["osu_get_acc_latency","osu_acc_latency","osu_fop_latency","osu_get_latency","osu_put_latency","osu_allreduce","osu_latency","osu_bibw","osu_bw"]}}
METADATA END
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 -map-by ppr:1:node /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_get_acc_latency
# OSU MPI_Get_accumulate latency Test v5.8
# Window creation: MPI_Win_create
# Synchronization: MPI_Win_lock/unlock
# Size          Latency (us)
1                      52.67
2                      50.36
4                      50.14
8                      48.87
16                     49.71
32                     48.98
64                     49.96
128                    48.78
256                    50.28
512                    44.27
1024                   46.65
2048                   51.04
4096                   54.27
8192                   51.89
16384                  63.45
32768                  92.55
65536                 174.83
131072                239.30
262144                397.48
524288                989.71
1048576              1656.63
2097152              2842.55
4194304              5630.91
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 -map-by ppr:1:node /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_acc_latency
# OSU MPI_Accumulate latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                      25.24
2                      23.72
4                      21.61
8                      19.79
16                     20.76
32                     22.51
64                     21.24
128                    20.56
256                    19.45
512                    20.27
1024                   21.44
2048                   28.04
4096                   28.20
8192                   37.41
16384                  51.21
32768                  74.02
65536                 124.97
131072                187.47
262144                300.23
524288                584.68
1048576              1324.63
2097152              3098.60
4194304              5490.46
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 -map-by ppr:1:node /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_fop_latency
# OSU MPI_Fetch_and_op latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
8                      33.40
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 -map-by ppr:1:node /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_get_latency
# OSU MPI_Get latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                      34.70
2                      30.95
4                      32.65
8                      30.62
16                     29.24
32                     30.06
64                     31.02
128                    29.27
256                    31.05
512                    26.09
1024                   24.83
2048                   26.03
4096                   27.75
8192                   28.40
16384                  33.26
32768                  42.23
65536                  78.78
131072                 87.65
262144                100.80
524288                147.90
1048576               243.02
2097152               501.91
4194304               948.50
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 -map-by ppr:1:node /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/one-sided/osu_put_latency
# OSU MPI_Put Latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                      21.26
2                      20.65
4                      23.96
8                      25.82
16                     26.27
32                     25.65
64                     24.99
128                    25.38
256                    25.25
512                    24.81
1024                   25.73
2048                   32.73
4096                   31.24
8192                   48.13
16384                  53.23
32768                  59.94
65536                 101.21
131072                114.61
262144                137.86
524288                196.58
1048576               331.81
2097152               630.59
4194304              1378.39
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 -map-by ppr:1:node -rank-by core /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/collective/osu_allreduce

# OSU MPI Allreduce Latency Test v5.8
# Size       Avg Latency(us)
4                      17.67
8                      15.49
16                     13.70
32                     13.10
64                     12.71
128                    12.76
256                    13.45
512                    13.72
1024                   14.43
2048                   21.85
4096                   21.35
8192                   23.97
16384                  44.40
32768                  51.14
65536                  77.59
131072                154.57
262144                224.81
524288                423.96
1048576               845.90
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 -map-by ppr:1:node /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt/osu_latency
# OSU MPI Latency Test v5.8
# Size          Latency (us)
0                       9.29
1                       9.04
2                       9.38
4                       9.35
8                       9.31
16                      8.85
32                      8.91
64                      9.90
128                     9.41
256                     9.67
512                     8.69
1024                   10.16
2048                   14.38
4096                   15.40
8192                   15.47
16384                  12.44
32768                  27.44
65536                  54.52
131072                 70.90
262144                 86.34
524288                145.27
1048576               310.09
2097152               686.40
4194304              1359.87
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 -map-by ppr:1:node /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt/osu_bibw
# OSU MPI Bi-Directional Bandwidth Test v5.8
# Size      Bandwidth (MB/s)
1                       0.23
2                       0.57
4                       1.10
8                       2.23
16                      4.45
32                      9.04
64                     17.76
128                    35.90
256                    70.83
512                   141.47
1024                  283.41
2048                  521.45
4096                  938.93
8192                 1583.56
16384                2445.28
32768                2558.25
65536                2116.49
131072               2778.19
262144               2784.05
524288               3770.98
1048576              3027.00
2097152              3520.72
4194304              2959.60
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 2 -np 2 -map-by ppr:1:node /opt/osu-benchmark/build.openmpi/libexec/osu-micro-benchmarks/mpi/pt2pt/osu_bw

# OSU MPI Bandwidth Test v5.8
# Size      Bandwidth (MB/s)
1                       0.13
2                       0.30
4                       0.60
8                       1.32
16                      2.51
32                      5.32
64                     10.98
128                    21.22
256                    39.65
512                    79.35
1024                  168.48
2048                  289.70
4096                  432.01
8192                  838.33
16384                1767.14
32768                1895.18
65536                1872.09
131072               2845.75
262144               3915.30
524288               4842.98
1048576              5199.52
2097152              4566.23
4194304              4117.48
METRICS OPERATOR COLLECTION END
```

</details>

The worker comes up and sleeps, and will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
Number of tasks (nproc on one node) is 8
Number of tasks total (across 2 nodes) is 16
Sleeping for 10 seconds waiting for network...
METADATA START {"pods":2,"completions":2,"metricName":"network-osu-benchmark","metricDescription":"point to point MPI benchmarks","metricType":"standalone","metricOptions":{"completions":0,"rate":10,"tasks":0},"metricListOptions":{"commands":["osu_get_acc_latency","osu_acc_latency","osu_fop_latency","osu_get_latency","osu_put_latency","osu_allreduce","osu_latency","osu_bibw","osu_bw"]}}
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
