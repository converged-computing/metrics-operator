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
metricset-sample-m-0-0-mwjns   1/1     Running   0          4s
metricset-sample-m-0-1-tlknf   1/1     Running   0          4s
```

If you inspect the log, you'll see a short sleep (the network isn't up immediately)
and then netmark running.

```bash
kubectl logs metricset-sample-m-0-0-mwjns -f
```
```console
root
#!/bin/bash
# Start ssh daemon
/usr/sbin/sshd -D &
whoami
# Show ourselves!
cat ${0}

# Write the hosts file
cat <<EOF > ./hostlist.txt
metricset-sample-m-0-0.ms.default.svc.cluster.local
metricset-sample-m-0-1.ms.default.svc.cluster.local

EOF

# Allow network to ready
echo "Sleeping for 10 seconds waiting for network..."
sleep 10

if [ $JOB_COMPLETION_INDEX = 0 ]; then
   mpirun -f ./hostlist.txt -np 2 /usr/local/bin/netmark.x -w 10 -t 20 -c 20 -b 0 -s     
else
   sleep infinity
fi
Sleeping for 10 seconds waiting for network...
size 2 rank 1 on host metricset-sample-m-0-1.ms.default.svc.cluster.local ip 10.244.0.32
=========== SETUP ===========
warmups                  10
trials                   20
send_recv_cycles         20
bytes                     0
store_per_trial           1
=============================
size 2 rank 0 on host metricset-sample-m-0-0.ms.default.svc.cluster.local ip 10.244.0.31
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
RTT between rank 0 and rank 1 is 30.801 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 31.778 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 32.152 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 32.356 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 31.308 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 31.311 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 31.020 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 30.661 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 30.670 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 30.625 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 31.342 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 31.160 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 30.879 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 30.523 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 30.734 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 30.849 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 30.857 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 30.907 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 30.931 micro-seconds
Rank 0 sends to rank 1
RTT between rank 0 and rank 1 is 31.246 micro-seconds
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
RTT between rank 1 and rank 0 is 29.244 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 27.696 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 29.390 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 30.313 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 30.798 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.261 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.378 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.359 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.598 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.370 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.478 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.423 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.268 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.459 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.754 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.828 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 33.008 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 33.039 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.989 micro-seconds
Rank 1 sends to rank 0
RTT between rank 1 and rank 0 is 32.692 micro-seconds
```

We are currently still adding support for custom completion, so the jobset/pods
won't be completed.

When you are done, cleanup!

```bash
kubectl delete -f metrics.yaml
kubectl delete cm metricset-sample
```