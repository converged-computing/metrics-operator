# OSU Benchmarks Example

This will demonstrate running [OSU Benchmarks](https://mvapich.cse.ohio-state.edu/benchmarks/) with the Metrics Operator.
For running the example, parsing, and plotting output, see [the corresponding Python directory](../../python/network-osu-benchmark/).

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
root
Sleeping for 10 seconds waiting for network...
METADATA START {"pods":2,"completions":2,"metricName":"network-osu-benchmark","metricDescription":"point to point MPI benchmarks","metricType":"standalone","metricOptions":{"completions":0,"rate":10},"metricListOptions":{"commands":["osu_fop_latency","osu_get_acc_latency","osu_get_latency","osu_put_latency","osu_acc_latency"]}}
METADATA END
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_fop_latency
# OSU MPI_Fetch_and_op latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
8                       0.45
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_get_acc_latency
# OSU MPI_Get_accumulate latency Test v5.8
# Window creation: MPI_Win_create
# Synchronization: MPI_Win_lock/unlock
# Size          Latency (us)
1                       3.16
2                       2.12
4                       1.48
8                       1.43
16                      1.42
32                      1.46
64                      1.48
128                     1.51
256                     1.57
512                     1.73
1024                    2.01
2048                    2.65
4096                    3.86
8192                    7.86
16384                  14.87
32768                  27.92
65536                  67.16
131072                115.28
262144                218.97
524288                440.29
1048576               899.82
2097152              2056.84
4194304              4217.38
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_get_latency
# OSU MPI_Get latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                       0.20
2                       0.19
4                       0.17
8                       0.15
16                      0.15
32                      0.14
64                      0.13
128                     0.13
256                     0.13
512                     0.13
1024                    0.17
2048                    0.13
4096                    0.15
8192                    0.22
16384                   0.52
32768                   0.97
65536                   1.91
131072                  3.98
262144                  8.56
524288                 16.80
1048576                32.75
2097152                79.51
4194304               578.12
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_put_latency
# OSU MPI_Put Latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                       0.30
2                       0.33
4                       0.31
8                       0.24
16                      0.22
32                      0.18
64                      0.18
128                     0.17
256                     0.17
512                     0.16
1024                    0.14
2048                    0.16
4096                    0.20
8192                    0.25
16384                   0.54
32768                   1.10
65536                   1.95
131072                  3.80
262144                  7.97
524288                 48.36
1048576                38.53
2097152                71.64
4194304               307.45
METRICS OPERATOR TIMEPOINT
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_acc_latency
# OSU MPI_Accumulate latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                       0.25
2                       0.24
4                       0.25
8                       0.25
16                      0.26
32                      0.28
64                      0.31
128                     0.38
256                     0.51
512                     0.82
1024                    1.29
2048                    2.39
4096                    4.21
8192                    8.22
16384                  15.93
32768                  30.47
65536                  61.28
131072                122.85
262144                252.14
524288                509.32
1048576              1035.88
2097152              2249.48
4194304              5677.54
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

And then you can cleanup!

```bash
kubectl delete -f metrics.yaml
```

Note that you can tweak the metrics.yaml to ask for specific metrics for OSU.
If you don't define any, you'll get the default list we provide. See
the metrics.yaml for details.
