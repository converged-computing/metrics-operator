# OSU Benchmarks Example

This will demonstrate running [OSU Benchmarks](https://mvapich.cse.ohio-state.edu/benchmarks/) with the Metrics Operator.

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
and then the benchmarks running. By default, we don't include a list in metrics.yaml so we run them all!

```bash
kubectl logs metricset-sample-l-0-0-lt782 -f
```

<details>

<summary>Output of OSU Benchmarks Launcher</summary>

```console
root
#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
whoami
# Show ourselves!
cat ${0}

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10

# Write the hosts file
launcher=$(getent hosts metricset-sample-l-0-0.ms.default.svc.cluster.local | awk '{ print $1 }')
worker=$(getent hosts metricset-sample-w-0-0.ms.default.svc.cluster.local | awk '{ print $1 }')
echo "${launcher}" >> ./hostfile.txt
echo "${worker}" >> ./hostfile.txt

sleep 5
echo "mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_acc_latency"
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_acc_latency
echo "mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_fop_latency"
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_fop_latency
echo "mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_get_acc_latency"
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_get_acc_latency
echo "mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_get_latency"
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_get_latency
echo "mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_put_latency"
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_put_latency
Sleeping for 10 seconds waiting for network...
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_acc_latency
# OSU MPI_Accumulate latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                       0.58
2                       0.41
4                       0.37
8                       0.30
16                      0.28
32                      0.25
64                      0.26
128                     0.32
256                     0.39
512                     0.54
1024                    0.96
2048                    1.61
4096                    2.88
8192                    5.57
16384                  11.43
32768                  21.93
65536                  41.58
131072                 81.81
262144                157.80
524288                278.97
1048576               548.83
2097152              1311.45
4194304              2484.51
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_fop_latency
# OSU MPI_Fetch_and_op latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
8                       0.40
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_get_acc_latency
# OSU MPI_Get_accumulate latency Test v5.8
# Window creation: MPI_Win_create
# Synchronization: MPI_Win_lock/unlock
# Size          Latency (us)
1                       2.09
2                       1.53
4                       1.40
8                       1.39
16                      1.46
32                      1.45
64                      1.60
128                     1.58
256                     1.64
512                     1.75
1024                    2.03
2048                    2.69
4096                    4.17
8192                    7.30
16384                  14.18
32768                  27.59
65536                  54.19
131072                113.30
262144                248.55
524288                457.77
1048576              1000.72
2097152              2149.51
4194304              4332.33
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_get_latency
# OSU MPI_Get latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                       0.22
2                       0.23
4                       0.24
8                       0.23
16                      0.22
32                      0.24
64                      0.18
128                     0.18
256                     0.16
512                     0.18
1024                    0.16
2048                    0.16
4096                    0.17
8192                    0.24
16384                   0.45
32768                   0.92
65536                   1.91
131072                  3.39
262144                  7.82
524288                 15.79
1048576                32.33
2097152                67.27
4194304               331.57
mpirun --hostfile ./hostfile.txt --allow-run-as-root -np 2 ./osu_put_latency
# OSU MPI_Put Latency Test v5.8
# Window creation: MPI_Win_allocate
# Synchronization: MPI_Win_flush
# Size          Latency (us)
1                       0.27
2                       0.29
4                       0.31
8                       0.27
16                      0.24
32                      0.21
64                      0.20
128                     0.16
256                     0.16
512                     0.19
1024                    0.15
2048                    0.17
4096                    0.23
8192                    0.29
16384                   0.63
32768                   1.19
65536                   2.37
131072                  4.30
262144                  8.75
524288                 16.65
1048576                33.12
2097152                70.90
4194304               293.05
```

</details>

The worker comes up and sleeps, and will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
root
#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
whoami
# Show ourselves!
cat ${0}

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10

# Write the hosts file
launcher=$(getent hosts metricset-sample-l-0-0.ms.default.svc.cluster.local | awk '{ print $1 }')
worker=$(getent hosts metricset-sample-w-0-0.ms.default.svc.cluster.local | awk '{ print $1 }')
echo "${launcher}" >> ./hostfile.txt
echo "${worker}" >> ./hostfile.txt

sleep infinity
Sleeping for 10 seconds waiting for network...
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