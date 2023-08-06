# Networking Example

The container for this example (and code) is private, so you won't be able to run it!
I'm using it for development purposes.

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

Make sure to load your private image into the node:

```bash
kind load docker-image vanessa/netmark:latest
```

Then create the metrics set. This is going to run a simple sysstat tool to collect metrics
as lammps runs.

```bash
kubectl apply -f metrics.yaml
```

Wait until you see pods created by the job and then running (there should be two):

```bash
kubectl get pods
```
```diff
NAME                           READY   STATUS    RESTARTS   AGE
metricset-sample-n-0-0-lt782   1/1     Running   0          3s
metricset-sample-w-0-0-4s5p9   1/1     Running   0          3s
```

In the above, "w" is a worker pod, and "n" is the netmark launcher.
If you inspect the log for the launcher you'll see a short sleep (the network isn't up immediately)
and then netmark running.

```bash
kubectl logs metricset-sample-n-0-0-lt782 -f
```
```console
root
#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
whoami
# Show ourselves!
cat ${0}

# If we have zero tasks, default to workers * nproc
np=2
pods=2
if [[ $np -eq 0 ]]; then
    np=$(nproc)
    np=$(( $pods*$np ))
fi

# Write the hosts file
cat <<EOF > ./hostlist.txt
metricset-sample-m-0-0.ms.default.svc.cluster.local
metricset-sample-m-0-1.ms.default.svc.cluster.local

EOF

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10

if [ $JOB_COMPLETION_INDEX = 0 ]; then
   mpirun -f ./hostlist.txt -np $np /usr/local/bin/netmark.x -w 10 -t 20 -c 20 -b 0 -s
else
   sleep infinity
fi
Sleeping for 10 seconds waiting for network...
=========== SETUP ===========
warmups                  10
trials                   20
send_recv_cycles         20
bytes                     0
store_per_trial           1
=============================
size 2 rank 1 on host metricset-sample-m-0-1.ms.default.svc.cluster.local ip 10.244.0.43
size 2 rank 0 on host metricset-sample-m-0-0.ms.default.svc.cluster.local ip 10.244.0.42
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 25.250 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 25.632 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 25.849 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 25.824 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 25.651 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 25.801 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 25.572 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 25.048 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 24.780 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 24.627 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 24.506 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 24.293 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 24.758 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 24.687 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 24.592 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 24.488 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 24.391 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 24.093 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 23.828 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 23.643 micro-seconds
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.500 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.252 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.728 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.488 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.410 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.358 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.264 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.222 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.438 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.485 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.474 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.477 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.422 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.341 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.321 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.313 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.262 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.190 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.111 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 17.079 micro-seconds
```
The worker will only be alive long enough for the main job to
finish, and once it does, the worker goes away! Here is what you'll see in its brief life:

```console
#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
whoami
# Show ourselves!
cat ${0}

# If we have zero tasks, default to workers * nproc
np=2
pods=2
if [[ $np -eq 0 ]]; then
        np=$(nproc)
        np=$(( $pods*$np ))
fi

# Write the hosts file
cat <<EOF > ./hostlist.txt
metricset-sample-n-0-0.ms.default.svc.cluster.local
metricset-sample-w-0-0.ms.default.svc.cluster.local

EOF

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10

sleep infinity
Sleeping for 10 seconds waiting for network...
```

We can do this with JobSet logic that the entire set is done when the launcher is done.

```bash
$ kubectl get pods
```
```console
NAME                           READY   STATUS        RESTARTS   AGE
metricset-sample-n-0-0-bqqf4   0/1     Completed     0          49s
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
metricset-sample-n-0   1/1           18s        84s
```

And then you can cleanup!

```bash
kubectl delete -f metrics.yaml
```

If you want to see how it scales, try increasing the number of pods for Netmark. You
might want to change the ranks too. Here is an example with 4 pods total (so one launcher
and 3 workers):

```bash
$ kubectl get pods
NAME                           READY   STATUS        RESTARTS   AGE
metricset-sample-n-0-0-2zprc   0/1     Completed     0          23s
metricset-sample-w-0-0-76wx9   1/1     Terminating   0          23s
metricset-sample-w-0-1-5j2kh   1/1     Terminating   0          23s
metricset-sample-w-0-2-lxg7w   1/1     Terminating   0          23s
```

Again, the worker jobs clean up nicely! We will next need to figure out a good strategy
for saving the outside, aside from parsing the main pod logs.
