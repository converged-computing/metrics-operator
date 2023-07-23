# Testing Communication between Containers

I want to test the extent to which I can [share process namespace between containers in a pod](https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace/).
If this is possible it might be a cool way to wrap or monitor (get metrics for) things going on in one container from another container.

## Pod Experiment

Create a kind cluster.

```
kind create cluster
```

Create the single nginx pod. Note that we have two containers - one with nginx and one with bash.

```
kubectl apply -f pod.yaml
```

Shell into the "shell" container.

```
kubectl exec -it nginx -c shell -- sh
```

Wow we can see all the nginx processes!

```bash
ps ax
```
```console
PID   USER     TIME  COMMAND
    1 65535     0:00 /pause
    7 root      0:00 nginx: master process nginx -g daemon off;
   41 101       0:00 nginx: worker process
   42 101       0:00 nginx: worker process
   43 101       0:00 nginx: worker process
   44 101       0:00 nginx: worker process
   45 101       0:00 nginx: worker process
   46 101       0:00 nginx: worker process
   47 101       0:00 nginx: worker process
   48 101       0:00 nginx: worker process
   49 root      0:00 sh
   67 root      0:00 sh
   75 root      0:00 ps ax
```

We can take a process id, and use it to look at the nginx container filesystem (that's pretty nuts, but these are just containers, so not crazy)

```
head /proc/7/root/etc/nginx/nginx.conf
```
```console
user  nginx;
worker_processes  auto;

error_log  /var/log/nginx/error.log notice;
pid        /var/run/nginx.pid;


events {
    worker_connections  1024;
```

Clean up when you are done.

```
kubectl delete -f pod.yaml
```

## JobSet Experiment

Let's instead try running something more real, like ML training. We are going to use a similar strategy, but with a Jobset. You'll
need to install it to your cluster.

```
VERSION=v0.2.0
kubectl apply --server-side -f https://github.com/kubernetes-sigs/jobset/releases/download/$VERSION/manifests.yaml
```

Note that this container is large and it takes 6- minutes to pull!

```
kubectl apply -f jobset.yaml
kubectl exec -it nginx-workers-0-0-ck29z -c shell -- sh
```

Here are processes again:

```
/ # ps ax
PID   USER     TIME  COMMAND
    1 65535     0:00 /pause
    7 root      0:00 nginx: master process nginx -g daemon off;
   41 101       0:00 nginx: worker process
   42 101       0:00 nginx: worker process
   43 101       0:00 nginx: worker process
   44 101       0:00 nginx: worker process
   45 101       0:00 nginx: worker process
   46 101       0:00 nginx: worker process
   47 101       0:00 nginx: worker process
   48 101       0:00 nginx: worker process
   49 root      0:00 sh
   67 root      0:00 sh
   73 root      0:00 ps ax
```

A general mapping I'm coming up with:

```console
/proc/pid                  # The virtual address space of a PID, e.g., can be the process in the other container, nginx!
/proc/pid/attr             # API for security modules
/proc/pid/attr/current     # Current security attributes of the process
....
/proc/7/autogroup          # I see "nice" so I think this is CPU scheduling related?
/proc/7/auxv               # "Contents of the ELF interpreter"
/proc/7/cgroup             # control group (e.g., 0::/../cri-containerd-a3cd74dfcc8f1ac5d4d999129414222c5f7ea921223ae63deca19a6e0c0f9e66.scope)
/proc/7/clear_refs         # Write only file that (says) it can be used for assessing memory usage - need to read about this more
/proc/7/cmdline            # Nuts! This is the command line! (I see the nginx start command)
/proc/7/comm               # Command name (e.g., I see "nginx")
/proc/7/coredump_filter    # This controls (filters) what parts of a coredump are written to file (if it dumps) see https://man7.org/linux/man-pages/man5/core.5.html
/proc/7/cpuset             # "confine memory and processes to memory subsets" I see the same cgroup here
/proc/7/cwd/               # Wow, symbolic link to CWD of process - I see the root of the nginx container (again)
/proc/7/environ            # Environment when process executed (need to test this more, it doesn't include the global envars it seems)
/proc/7/exe                # Actual symbolic link to binary executed
/proc/7/fd                 # Directory that shows all the file descriptors the process has open (again, wow)
/proc/7/fdinfo

STOPPED HERE - will add more later, want to make bad life decisions :)
```


Here is the [man page](https://man7.org/linux/man-pages/man5/proc.5.html) that I found useful for learning about the below. 

### Filesystem

Here is the filesystem!

```
# ls /proc/7/root/
bin                   docker-entrypoint.d   home                  lib64                 mnt                   product_name          run                   sys                   var
boot                  docker-entrypoint.sh  lib                   libx32                opt                   product_uuid          sbin                  tmp
dev                   etc                   lib32                 media                 proc                  root                  srv                   usr
```

### Status

And in "status" there are metrics from cpu to memor to pretty much everything I could think of!

<details>

<summary>Contents of /proc/7/status</summary>

```
/ # cat /proc/7/status 
Name:   nginx
Umask:  0022
State:  S (sleeping)
Tgid:   7
Ngid:   0
Pid:    7
PPid:   0
TracerPid:      0
Uid:    0       0       0       0
Gid:    0       0       0       0
FDSize: 64
Groups:  
NStgid: 7
NSpid:  7
NSpgid: 7
NSsid:  7
VmPeak:    11376 kB
VmSize:    11376 kB
VmLck:         0 kB
VmPin:         0 kB
VmHWM:      7688 kB
VmRSS:      7688 kB
RssAnon:            1248 kB
RssFile:            6436 kB
RssShmem:              4 kB
VmData:     1592 kB
VmStk:       132 kB
VmExe:       932 kB
VmLib:      5012 kB
VmPTE:        56 kB
VmSwap:        0 kB
HugetlbPages:          0 kB
CoreDumping:    0
THP_enabled:    1
Threads:        1
SigQ:   1/62394
SigPnd: 0000000000000000
ShdPnd: 0000000000000000
SigBlk: 0000000000000000
SigIgn: 0000000040001000
SigCgt: 0000000018016a07
CapInh: 0000000000000000
CapPrm: 00000000a80425fb
CapEff: 00000000a80425fb
CapBnd: 00000000a80425fb
CapAmb: 0000000000000000
NoNewPrivs:     0
Seccomp:        0
Seccomp_filters:        0
Speculation_Store_Bypass:       thread vulnerable
SpeculationIndirectBranch:      conditional enabled
Cpus_allowed:   ff
Cpus_allowed_list:      0-7
Mems_allowed:   00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000000,00000001
Mems_allowed_list:      0
voluntary_ctxt_switches:        22
nonvoluntary_ctxt_switches:     4
```

</details>

Here are environment variables, so we could theoretically pass information between containers based on them.

### Maps

This looks to be the [virtual address space](https://www.baeldung.com/linux/proc-id-maps) (VAS) for the other container!

<details>

<summary>VAS of nginx container</summary>

```
# cat /proc/7/maps
55e5ead27000-55e5ead55000 r--p 00000000 00:309 28874585                  /usr/sbin/nginx
55e5ead55000-55e5eae3e000 r-xp 0002e000 00:309 28874585                  /usr/sbin/nginx
55e5eae3e000-55e5eae7c000 r--p 00117000 00:309 28874585                  /usr/sbin/nginx
55e5eae7c000-55e5eae7f000 r--p 00155000 00:309 28874585                  /usr/sbin/nginx
55e5eae7f000-55e5eae9d000 rw-p 00158000 00:309 28874585                  /usr/sbin/nginx
55e5eae9d000-55e5eaf5d000 rw-p 00000000 00:00 0 
55e5eb047000-55e5eb0ca000 rw-p 00000000 00:00 0                          [heap]
7fd91b12d000-7fd91b130000 rw-p 00000000 00:00 0 
7fd91b130000-7fd91b156000 r--p 00000000 00:309 28871052                  /usr/lib/x86_64-linux-gnu/libc.so.6
7fd91b156000-7fd91b2ab000 r-xp 00026000 00:309 28871052                  /usr/lib/x86_64-linux-gnu/libc.so.6
7fd91b2ab000-7fd91b2fe000 r--p 0017b000 00:309 28871052                  /usr/lib/x86_64-linux-gnu/libc.so.6
7fd91b2fe000-7fd91b302000 r--p 001ce000 00:309 28871052                  /usr/lib/x86_64-linux-gnu/libc.so.6
7fd91b302000-7fd91b304000 rw-p 001d2000 00:309 28871052                  /usr/lib/x86_64-linux-gnu/libc.so.6
7fd91b304000-7fd91b313000 rw-p 00000000 00:00 0 
7fd91b313000-7fd91b316000 r--p 00000000 00:309 28871151                  /usr/lib/x86_64-linux-gnu/libz.so.1.2.13
7fd91b316000-7fd91b329000 r-xp 00003000 00:309 28871151                  /usr/lib/x86_64-linux-gnu/libz.so.1.2.13
7fd91b329000-7fd91b330000 r--p 00016000 00:309 28871151                  /usr/lib/x86_64-linux-gnu/libz.so.1.2.13
7fd91b330000-7fd91b331000 r--p 0001c000 00:309 28871151                  /usr/lib/x86_64-linux-gnu/libz.so.1.2.13
7fd91b331000-7fd91b332000 rw-p 0001d000 00:309 28871151                  /usr/lib/x86_64-linux-gnu/libz.so.1.2.13
7fd91b332000-7fd91b3f7000 r--p 00000000 00:309 28874495                  /usr/lib/x86_64-linux-gnu/libcrypto.so.3
7fd91b3f7000-7fd91b66f000 r-xp 000c5000 00:309 28874495                  /usr/lib/x86_64-linux-gnu/libcrypto.so.3
7fd91b66f000-7fd91b74c000 r--p 0033d000 00:309 28874495                  /usr/lib/x86_64-linux-gnu/libcrypto.so.3
7fd91b74c000-7fd91b7ad000 r--p 0041a000 00:309 28874495                  /usr/lib/x86_64-linux-gnu/libcrypto.so.3
7fd91b7ad000-7fd91b7b0000 rw-p 0047b000 00:309 28874495                  /usr/lib/x86_64-linux-gnu/libcrypto.so.3
7fd91b7b0000-7fd91b7b3000 rw-p 00000000 00:00 0 
7fd91b7b3000-7fd91b7d2000 r--p 00000000 00:309 28874567                  /usr/lib/x86_64-linux-gnu/libssl.so.3
7fd91b7d2000-7fd91b82f000 r-xp 0001f000 00:309 28874567                  /usr/lib/x86_64-linux-gnu/libssl.so.3
7fd91b82f000-7fd91b84e000 r--p 0007c000 00:309 28874567                  /usr/lib/x86_64-linux-gnu/libssl.so.3
7fd91b84e000-7fd91b858000 r--p 0009b000 00:309 28874567                  /usr/lib/x86_64-linux-gnu/libssl.so.3
7fd91b858000-7fd91b85c000 rw-p 000a5000 00:309 28874567                  /usr/lib/x86_64-linux-gnu/libssl.so.3
7fd91b85c000-7fd91b85e000 r--p 00000000 00:309 28871115                  /usr/lib/x86_64-linux-gnu/libpcre2-8.so.0.11.2
7fd91b85e000-7fd91b8c9000 r-xp 00002000 00:309 28871115                  /usr/lib/x86_64-linux-gnu/libpcre2-8.so.0.11.2
7fd91b8c9000-7fd91b8f4000 r--p 0006d000 00:309 28871115                  /usr/lib/x86_64-linux-gnu/libpcre2-8.so.0.11.2
7fd91b8f4000-7fd91b8f5000 r--p 00098000 00:309 28871115                  /usr/lib/x86_64-linux-gnu/libpcre2-8.so.0.11.2
7fd91b8f5000-7fd91b8f6000 rw-p 00099000 00:309 28871115                  /usr/lib/x86_64-linux-gnu/libpcre2-8.so.0.11.2
7fd91b8f6000-7fd91b8f8000 r--p 00000000 00:309 28871061                  /usr/lib/x86_64-linux-gnu/libcrypt.so.1.1.0
7fd91b8f8000-7fd91b90e000 r-xp 00002000 00:309 28871061                  /usr/lib/x86_64-linux-gnu/libcrypt.so.1.1.0
7fd91b90e000-7fd91b928000 r--p 00018000 00:309 28871061                  /usr/lib/x86_64-linux-gnu/libcrypt.so.1.1.0
7fd91b928000-7fd91b929000 r--p 00031000 00:309 28871061                  /usr/lib/x86_64-linux-gnu/libcrypt.so.1.1.0
7fd91b929000-7fd91b92a000 rw-p 00032000 00:309 28871061                  /usr/lib/x86_64-linux-gnu/libcrypt.so.1.1.0
7fd91b92a000-7fd91b932000 rw-p 00000000 00:00 0 
7fd91b935000-7fd91b936000 rw-s 00000000 00:01 2090488                    /dev/zero (deleted)
7fd91b936000-7fd91b938000 rw-p 00000000 00:00 0 
7fd91b938000-7fd91b939000 r--p 00000000 00:309 28871034                  /usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
7fd91b939000-7fd91b95e000 r-xp 00001000 00:309 28871034                  /usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
7fd91b95e000-7fd91b968000 r--p 00026000 00:309 28871034                  /usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
7fd91b968000-7fd91b96a000 r--p 00030000 00:309 28871034                  /usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
7fd91b96a000-7fd91b96c000 rw-p 00032000 00:309 28871034                  /usr/lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
7ffc8fb41000-7ffc8fb62000 rw-p 00000000 00:00 0                          [stack]
7ffc8fbd6000-7ffc8fbda000 r--p 00000000 00:00 0                          [vvar]
7ffc8fbda000-7ffc8fbdc000 r-xp 00000000 00:00 0                          [vdso]
ffffffffff600000-ffffffffff601000 --xp 00000000 00:00 0                  [vsyscall]
```

</details>

THIS IS SO COOL! There is so much we can do with this! I need to pause to do some wonky experiments with Flux.


## Flux

Let's create a Flux jobset - two containers (each with Flux) and one starting a test broker.

```bash
kubectl apply -f flux.yaml
```

What I want to try doing is connecting to the socket of one from the other. This might not make sense,
but abstractly I'm wondering if we would have a way to interact with Flux (eventually from a container that doesn't have it)
to maybe do cool things (that I need to think more about). Shell in to flux1 (running the instance)

```bash
kubectl exec -it flux-flux-0-0-cjd96 -c flux1 bash
```

We can see the sockets in tmp - since we did flux start with test size 4, we see a bunch!

```
 fluxuser@flux-flux-0-0:~$ ls /tmp/flux-jhgpJK/
content.sqlite  local-0         local-1         local-2         local-3         start           tbon-0          tbon-1
```
Note this for later - the other container won't have these. Let's proxy to one and run a job.

```console
$ flux proxy local:///tmp/flux-jhgpJK/local-1 
$ flux resource list
     STATE NNODES   NCORES    NGPUS NODELIST
      free      4       16        0 flux-flux-0-[0,0,0,0]
 allocated      0        0        0 
      down      0        0        0 
```

Let's submit a very useless job.

```
flux submit sleep 900
```
```
$ flux jobs -a
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
   ƒ26MYGwXM fluxuser sleep       R      1      1   1.493s flux-flux-0-0
```

Exit from the instance and container and shell into flux2

```
$ kubectl exec -it flux-flux-0-0-cjd96 -c flux2 bash
```

We didn't start anything here, so there is only the local-0 in tmp that is started by the entrypoint:

```
$ ls /tmp/flux-4CqwWk/
content.sqlite  local-0
```

BUT we can see the other filesystem right? And processes?

```
 fluxuser@flux-flux-0-0:~$ ps aux
USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
65535          1  0.0  0.0    972     4 ?        Ss   21:46   0:00 /pause
fluxuser       7  0.0  0.0  22004  3380 ?        Ss   21:48   0:00 start --test-size=4 sleep infinity
fluxuser      19  0.1  0.1 1253084 17824 ?       Sl   21:48   0:00 /usr/libexec/flux/cmd/flux-broker --setattr=rundir=/tmp/flux-jhgpJK sleep infinity
fluxuser      20  0.0  0.0 630192 10992 ?        Sl   21:48   0:00 /usr/libexec/flux/cmd/flux-broker --setattr=rundir=/tmp/flux-jhgpJK
fluxuser      21  0.0  0.0 621828 11044 ?        Sl   21:48   0:00 /usr/libexec/flux/cmd/flux-broker --setattr=rundir=/tmp/flux-jhgpJK
fluxuser      22  0.0  0.0 621960 11008 ?        Sl   21:48   0:00 /usr/libexec/flux/cmd/flux-broker --setattr=rundir=/tmp/flux-jhgpJK
fluxuser      55  0.0  0.1 1187532 17296 pts/0   Ssl  21:48   0:00 /usr/libexec/flux/cmd/flux-broker /bin/bash
munge         73  0.0  0.0  71180  1940 ?        Sl   21:48   0:00 /usr/sbin/munged
fluxuser     196  0.0  0.0   8968  3828 pts/0    S+   21:48   0:00 /bin/bash
fluxuser     256  0.0  0.0   7236   580 ?        S    21:48   0:00 sleep infinity
fluxuser     288  0.0  0.0  51920 10508 ?        S    21:50   0:00 /usr/libexec/flux/flux-shell 2411808687104
fluxuser     289  0.0  0.0   7236   516 ?        S    21:50   0:00 sleep 900
fluxuser     291  0.0  0.0   8968  3884 pts/1    Ss   21:51   0:00 bash
```

I found the sockets for the other container, flux1!

```
$ ls /proc/7/root/tmp/flux-jhgpJK/
content.sqlite  jobtmp-0-ƒ26MYGwXM  local-0  local-1  local-2  local-3  start  tbon-0  tbon-1
```

Moment of truth - can we connect?

```
 fluxuser@flux-flux-0-0:~$ flux proxy local:///proc/7/root/tmp/flux-jhgpJK/local-1 bash
ƒ(s=4,d=0) fluxuser@flux-flux-0-0:~$ flux jobs -a
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
   ƒ26MYGwXM fluxuser sleep       R      1      1   3.654m flux-flux-0-0
```

OH MAH GOSH!! My head just exploded. It is astounding. 🤯️🤣️ There is so much cool ideas we can try with this!

Clean up when you are done, meaning the container and your exploded brains.

```
$ kubectl delete -f flux.yaml 
```

:)