# DLIO Example

This is an example of using the IO tool[DLIO](https://dlio-profiler.readthedocs.io/en/latest/build.html#build-dlio-profiler-with-pip-recommended) that can 
be added on the fly with pip.

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

Then create the metrics set. This is going to run a single run of LAMMPS over MPI.
as lammps runs.

```bash
kubectl apply -f metrics.yaml
```

Wait until you see pods created by the job and then running (there should be two - a launcher and worker for LAMMPS):

```bash
kubectl get pods
```
```diff
NAME                           READY   STATUS    RESTARTS   AGE
metricset-sample-l-0-0-lt782   1/1     Running   0          3s
metricset-sample-w-0-0-4s5p9   1/1     Running   0          3s
```

In the above, "l" is a launcher pod, and "w" is a worker node.
If you inspect the log for the launcher you'll see a short sleep (the network isn't up immediately)
and then LAMMPS running, and the log is printed to the console.

```bash
kubectl logs metricset-sample-l-0-0-lt782 -f
```
```console
Installing collected packages: pybind11, dlio-profiler-py
Successfully installed dlio-profiler-py-0.0.1 pybind11-2.11.1
WARNING: Running pip as the 'root' user can result in broken permissions and conflicting behaviour with the system package manager. It is recommended to use a virtual environment instead: https://pip.pypa.io/warnings/venv
preload path is /usr/lib/python3/dist-packages/dlio_profiler/lib/libdlio_profiler_preload.so
[DLIO_PROFILER INFO]: Extracted process_name ior
[DLIO_PROFILER INFO]: created log file /ior/logs/trace-ior-2315.pfw with fd 3
[DLIO_PROFILER INFO]: Writing trace to /ior/logs/trace-ior-2315.pfw
[DLIO_PROFILER INFO]: Preloading DLIO Profiler with log_file /ior/logs/trace-ior-2315.pfw data_dir testfile and process 2315
IOR-4.0.0rc1: MPI Coordinated Test of Parallel I/O
Began               : Fri Nov  3 01:42:25 2023
Command line        : ior -k -w -r -o testfile
Machine             : Linux metricset-sample-m-0-0.ms.default.svc.cluster.local
TestID              : 0
StartTime           : Fri Nov  3 01:42:25 2023
Path                : testfile
FS                  : 1.8 TiB   Used FS: 13.4%   Inodes: 116.4 Mi   Used Inodes: 4.9%

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
[DLIO_PROFILER WARN]: Profiler Intercepted POSIX tracing file testfile for func access
[DLIO_PROFILER WARN]: Profiler Intercepted POSIX tracing file testfile for func open64
write     1893.59    8248       0.000121    1024.00    256.00     0.000039   0.000485   0.000004   0.000528   0   
read      10305      42582      0.000023    1024.00    256.00     0.000002   0.000094   0.000001   0.000097   0   

Summary of all tests:
Operation   Max(MiB)   Min(MiB)  Mean(MiB)     StdDev   Max(OPs)   Min(OPs)  Mean(OPs)     StdDev    Mean(s) Stonewall(s) Stonewall(MiB) Test# #Tasks tPN reps fPP reord reordoff reordrand seed segcnt   blksiz    xsize aggs(MiB)   API RefNum
write        1893.59    1893.59    1893.59       0.00    7574.36    7574.36    7574.36       0.00    0.00053         NA            NA     0      1   1    1   0     0        1         0    0      1  1048576   262144       1.0 POSIX      0
read        10305.42   10305.42   10305.42       0.00   41221.66   41221.66   41221.66       0.00    0.00010         NA            NA     0      1   1    1   0     0        1         0    0      1  1048576   262144       1.0 POSIX      0
Finished            : Fri Nov  3 01:42:25 2023
[DLIO_PROFILER INFO]: Calling finalize on pid 2315
[DLIO_PROFILER INFO]: Release Prefix Tree
[DLIO_PROFILER INFO]: Release I/O bindings
[DLIO_PROFILER INFO]: Profiler finalizing writer /ior/logs/trace-ior-2315.pfw
[DLIO_PROFILER INFO]: Profiler writing the final symbol
[DLIO_PROFILER INFO]: Extracted process_name sh
[DLIO_PROFILER INFO]: created log file /ior/logs/trace-sh-2318.pfw with fd 3
[DLIO_PROFILER INFO]: Writing trace to /ior/logs/trace-sh-2318.pfw
[DLIO_PROFILER INFO]: Preloading DLIO Profiler with log_file /ior/logs/trace-sh-2318.pfw data_dir testfile and process 2318
[DLIO_PROFILER INFO]: Applying Gzip compression on file /ior/logs/trace-ior-2315.pfw
[DLIO_PROFILER INFO]: Extracted process_name sh
[DLIO_PROFILER INFO]: created log file /ior/logs/trace-sh-2320.pfw with fd 3
[DLIO_PROFILER INFO]: Writing trace to /ior/logs/trace-sh-2320.pfw
[DLIO_PROFILER INFO]: Preloading DLIO Profiler with log_file /ior/logs/trace-sh-2320.pfw data_dir testfile and process 2320
[DLIO_PROFILER INFO]: Extracted process_name gzip
[DLIO_PROFILER INFO]: created log file /ior/logs/trace-gzip-2321.pfw with fd 4
[DLIO_PROFILER INFO]: Writing trace to /ior/logs/trace-gzip-2321.pfw
[DLIO_PROFILER INFO]: Preloading DLIO Profiler with log_file /ior/logs/trace-gzip-2321.pfw data_dir testfile and process 2321
[DLIO_PROFILER INFO]: Calling finalize on pid 2321
[DLIO_PROFILER INFO]: Release Prefix Tree
[DLIO_PROFILER INFO]: Release I/O bindings
[DLIO_PROFILER INFO]: Profiler finalizing writer /ior/logs/trace-gzip-2321.pfw
[DLIO_PROFILER INFO]: No trace data written. Deleting file /ior/logs/trace-gzip-2321.pfw
[DLIO_PROFILER INFO]: Released Logger
[DLIO_PROFILER INFO]: Successfully compressed file /ior/logs/trace-ior-2315.pfw.gz
[DLIO_PROFILER INFO]: Released Logger
METRICS OPERATOR COLLECTION END
trace-ior-2315.pfw.gz  trace-sh-2318.pfw  trace-sh-2319.pfw  trace-sh-2320.pfw
[
 {"id":"0","name":"access","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745976947","dur":"2","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","fname":"testfile"}}
{"id":"1","name":"open64","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745976978","dur":"29","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","ret":15,"flags":66,"mode":436,"fname":"testfile"}}
{"id":"2","name":"lseek64","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745977017","dur":"1","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","ret":"0","whence":0,"offset":"0","fd":15,"fname":"testfile"}}
{"id":"3","name":"write","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745977021","dur":"205","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","ret":"262144","count":"262144","fd":15,"fname":"testfile"}}
{"id":"4","name":"lseek64","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745977233","dur":"0","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","ret":"262144","whence":0,"offset":"262144","fd":15,"fname":"testfile"}}
{"id":"5","name":"write","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745977235","dur":"91","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","ret":"262144","count":"262144","fd":15,"fname":"testfile"}}
{"id":"6","name":"lseek64","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745977328","dur":"1","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","ret":"524288","whence":0,"offset":"524288","fd":15,"fname":"testfile"}}
{"id":"7","name":"write","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745977330","dur":"81","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","ret":"262144","count":"262144","fd":15,"fname":"testfile"}}
{"id":"8","name":"lseek64","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745977413","dur":"1","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","ret":"786432","whence":0,"offset":"786432","fd":15,"fname":"testfile"}}
{"id":"9","name":"write","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745977415","dur":"83","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","ret":"262144","count":"262144","fd":15,"fname":"testfile"}}
{"id":"10","name":"close","cat":"POSIX","pid":"2315","tid":"4630","ts":"1698975745977500","dur":"2","ph":"X","args":{"hostname":"metricset-sample-m-0-0.ms.default.svc.cluster.local","ret":0,"fd":15,"fname":"testfile"}}
```

There is purposefully a sleep infinity at the end to give you a chance to copy over data.

```bash
mkdir logs
kubectl  cp metricset-sample-m-0-0-xk28x:/ior/logs ./logs/
```

You can open the tiny file in [https://ui.perfetto.dev/](https://ui.perfetto.dev/).

![img/ior.png](img/ior.png)

Other applications of interest might be related to AI/ML - we will try more soon!
Cleanup when you are done.

```bash
kubectl delete -f metrics.yaml
```
