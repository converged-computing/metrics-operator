# Hello World Perf Example

This is a simple "hello world" example for a performance metric, sysstat.

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

Wait until you see pods created by the job and then running (there should be one with two containers, one for the app lammps and the other for the stats):

```bash
kubectl get pods
```
```diff
NAME                           READY   STATUS              RESTARTS   AGE
- metricset-sample-m-0-0-mkwrh   0/2     ContainerCreating   0          2m20s
+ metricset-sample-m-0-0-mkwrh   2/2     Running             0          3m10s
```

You should always be able to get logs for the application container in a pod, which is named "app". Since this is a sleep, we won't see anything interesting. However, the sidecar metrics containers will output their metrics for the lifetime of the application, and (currently) also to their logs:

```bash
kubectl logs metricset-sample-m-0-czxrq -c perf-sysstat
```
```console
05:39:24        0        20         -    0.00    0.00    0.00    0.00    0.00     4  mpirun
05:39:24        0         -        20    0.00    0.00    0.00    0.00    0.00     4  |__mpirun
KERNEL TABLES 2176
        34  pidstat -p 20 -v -h
    echo TASK SWITCHING 2176
/metrics_operator/entrypoint-0.sh: line 28: 35: command not found
CPU STATISTICS TIMEPOINT 2177
    pidstat -p 20 -u -h
    echo KERNEL STATISTICS TIMEPOINT 2177
Linux 5.15.0-76-generic (metricset-sample-m-0-0)        07/29/23        _x86_64_    (8 CPU)

# Time        UID       PID   kB_rd/s   kB_wr/s kB_ccwr/s iodelay  Command
05:39:24        0        20      0.00      0.00      0.00       0  mpirun
POLICY TIMEPOINT 2177
    pidstat -p 20 -R -h
    echo PAGEFAULTS and MEMORY 2177
Linux 5.15.0-76-generic (metricset-sample-m-0-0)        07/29/23        _x86_64_    (8 CPU)

# Time        UID       PID  minflt/s  majflt/s     VSZ     RSS   %MEM  Command
05:39:24        0        20      0.00      0.00    6548    3160   0.02  mpirun
STACK UTILIZATION 2177
        pidstat -p 20 -s -h
    echo THREADS 2177
Linux 5.15.0-76-generic (metricset-sample-m-0-0)        07/29/23        _x86_64_    (8 CPU)

# Time        UID      TGID       TID    %usr %system  %guest   %wait    %CPU   CPU  Command
05:39:24        0        20         -    0.00    0.00    0.00    0.00    0.00     4  mpirun
05:39:24        0         -        20    0.00    0.00    0.00    0.00    0.00     4  |__mpirun
```

Those are just a random set of stats I am running using this tool for the shared PID - I need to think
of a better way to capture and save these! Also, right now we don't have a completion policy, but instead have
the metrics collector exit when the PID goes away. In the future we could use a success policy. When you are done, the JobsSet, associated jobs, and pods will be completed:

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
metricset-sample-m-0-0-rq4q9   0/2     Completed   0          3m19s
```

When you are done, cleanup!

```bash
kubectl delete -f metrics.yaml
```