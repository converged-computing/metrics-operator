# LAMMPS HPCToolkit Example

This is an example to show running HPCToolkit. This application metric takes a design
of providing a modular install of HPCToolkit (via a spack copy view) and then moving it
into a shared volume, shared by the metrics container.

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

Then create the metrics set. 

```bash
kubectl apply -f metrics.yaml
```

Wait until you see pods created by the job and then running (there should be one with two containers, one for the app lammps and the other for HPCToolkit, which is just used to provision the volume):

```bash
kubectl get pods
```
```diff
NAME                         READY   STATUS    RESTARTS   AGE
metricset-sample-m-0-8n5wf   2/2     Running   0          90m
```

Note that this metrics example has interactive true, so it won't exit. We do this so we can shell in and explore the data! Let's do that first - looking at the app log:

```bash
kubectl logs metricset-sample-m-0-czxrq -c app
```

<details>

<summary>Output of lammps app</summary>

```console
--2023-09-06 20:45:22--  https://github.com/converged-computing/goshare/releases/download/2023-09-06/wait-fs
Resolving github.com (github.com)... 140.82.112.4
Connecting to github.com (github.com)|140.82.112.4|:443... connected.
HTTP request sent, awaiting response... 302 Found
Location: https://objects.githubusercontent.com/github-production-release-asset-2e65be/670447414/dac45779-6f67-4c45-9a94-1d4ab9dc7331?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20230906%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20230906T204522Z&X-Amz-Expires=300&X-Amz-Signature=de24651fb52e5badb739f2985e7d04ac045c18d23d0e94d43d1f73cd0f0ffba9&X-Amz-SignedHeaders=host&actor_id=0&key_id=0&repo_id=670447414&response-content-disposition=attachment%3B%20filename%3Dwait-fs&response-content-type=application%2Foctet-stream [following]
--2023-09-06 20:45:22--  https://objects.githubusercontent.com/github-production-release-asset-2e65be/670447414/dac45779-6f67-4c45-9a94-1d4ab9dc7331?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAIWNJYAX4CSVEH53A%2F20230906%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20230906T204522Z&X-Amz-Expires=300&X-Amz-Signature=de24651fb52e5badb739f2985e7d04ac045c18d23d0e94d43d1f73cd0f0ffba9&X-Amz-SignedHeaders=host&actor_id=0&key_id=0&repo_id=670447414&response-content-disposition=attachment%3B%20filename%3Dwait-fs&response-content-type=application%2Foctet-stream
Resolving objects.githubusercontent.com (objects.githubusercontent.com)... 185.199.109.133, 185.199.108.133, 185.199.111.133, ...
Connecting to objects.githubusercontent.com (objects.githubusercontent.com)|185.199.109.133|:443... connected.
HTTP request sent, awaiting response... 200 OK
Length: 2116087 (2.0M) [application/octet-stream]
Saving to: 'wait-fs'

wait-fs             100%[===================>]   2.02M  5.45MB/s    in 0.4s    

2023-09-06 20:45:23 (5.45 MB/s) - 'wait-fs' saved [2116087/2116087]

🟧️  wait-fs: 2023/09/06 20:45:23 wait-fs.go:40: /opt/share/software
🟧️  wait-fs: 2023/09/06 20:45:23 wait-fs.go:53: Path /opt/share/software does not exist yet, sleeping 5
🟧️  wait-fs: 2023/09/06 20:45:28 wait-fs.go:53: Path /opt/share/software does not exist yet, sleeping 5
🟧️  wait-fs: 2023/09/06 20:45:33 wait-fs.go:53: Path /opt/share/software does not exist yet, sleeping 5
🟧️  wait-fs: 2023/09/06 20:45:38 wait-fs.go:49: Found existing path /opt/share/software
🟧️  wait-fs: 2023/09/06 20:45:52 wait-fs.go:40: /opt/share/view/bin/hpcrun
🟧️  wait-fs: 2023/09/06 20:45:52 wait-fs.go:49: Found existing path /opt/share/view/bin/hpcrun
METADATA START {"pods":1,"completions":1,"applicationImage":"ghcr.io/rse-ops/vanilla-lammps:tag-latest","applicationCommand":"lmp -v x 1 -v y 1 -v z 1 -in in.reaxc.hns -nocite","metricName":"perf-hpctoolkit","metricDescription":"performance tools for measurement and analysis","metricType":"application","metricOptions":{"events":"-e IO","mount":"/opt/share"}}
METADATA END
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
LAMMPS (29 Sep 2021 - Update 2)
OMP_NUM_THREADS environment is not set. Defaulting to 1 thread. (src/comm.cpp:98)
  using 1 OpenMP thread(s) per MPI task
Reading data file ...
  triclinic box = (0.0000000 0.0000000 0.0000000) to (22.326000 11.141200 13.778966) with tilt (0.0000000 -5.0260300 0.0000000)
  1 by 1 by 1 MPI processor grid
  reading atoms ...
  304 atoms
  reading velocities ...
  304 velocities
  read_data CPU = 0.003 seconds
Replicating atoms ...
  triclinic box = (0.0000000 0.0000000 0.0000000) to (22.326000 11.141200 13.778966) with tilt (0.0000000 -5.0260300 0.0000000)
  1 by 1 by 1 MPI processor grid
  bounding box image = (0 -1 -1) to (0 1 1)
  bounding box extra memory = 0.03 MB
  average # of replicas added to proc = 1.00 out of 1 (100.00%)
  304 atoms
  replicate CPU = 0.001 seconds
Neighbor list info ...
  update every 20 steps, delay 0 steps, check no
  max neighbors/atom: 2000, page size: 100000
  master list distance cutoff = 11
  ghost atom cutoff = 11
  binsize = 5.5, bins = 5 3 3
  2 neighbor lists, perpetual/occasional/extra = 2 0 0
  (1) pair reax/c, perpetual
      attributes: half, newton off, ghost
      pair build: half/bin/newtoff/ghost
      stencil: full/ghost/bin/3d
      bin: standard
  (2) fix qeq/reax, perpetual, copy from (1)
      attributes: half, newton off, ghost
      pair build: copy
      stencil: none
      bin: none
Setting up Verlet run ...
  Unit style    : real
  Current step  : 0
  Time step     : 0.1
Per MPI rank memory allocation (min/avg/max) = 78.15 | 78.15 | 78.15 Mbytes
Step Temp PotEng Press E_vdwl E_coul Volume 
       0          300   -113.27833    427.09094   -111.57687   -1.7014647    3427.3584 
      10    298.13784   -113.27279    1855.1535   -111.57169   -1.7011017    3427.3584 
      20    294.02745   -113.25991    3911.5126     -111.559   -1.7009101    3427.3584 
      30    293.61692   -113.25867    7296.5076   -111.55793   -1.7007375    3427.3584 
      40    301.40293   -113.28175    9622.4058   -111.58127   -1.7004797    3427.3584 
      50    310.92489   -113.31003    10101.225   -111.60982   -1.7002117    3427.3584 
      60    311.37774   -113.31149    9274.1322   -111.61144   -1.7000446    3427.3584 
      70    302.58347   -113.28582     6350.705   -111.58587   -1.6999549    3427.3584 
      80    295.34242   -113.26406    6795.0642   -111.56427   -1.6997975    3427.3584 
      90    299.15724   -113.27518    9198.0327   -111.57566   -1.6995238    3427.3584 
     100    307.63997   -113.30058    9424.4991   -111.60129   -1.6992878    3427.3584 
Loop time of 3.81015 on 1 procs for 100 steps with 304 atoms

Performance: 0.227 ns/day, 105.838 hours/ns, 26.246 timesteps/s
99.9% CPU use with 1 MPI tasks x 1 OpenMP threads

MPI task timing breakdown:
Section |  min time  |  avg time  |  max time  |%varavg| %total
---------------------------------------------------------------
Pair    | 3.1212     | 3.1212     | 3.1212     |   0.0 | 81.92
Neigh   | 0.094558   | 0.094558   | 0.094558   |   0.0 |  2.48
Comm    | 0.0025178  | 0.0025178  | 0.0025178  |   0.0 |  0.07
Output  | 0.00027439 | 0.00027439 | 0.00027439 |   0.0 |  0.01
Modify  | 0.59109    | 0.59109    | 0.59109    |   0.0 | 15.51
Other   |            | 0.0005321  |            |       |  0.01

Nlocal:        304.000 ave         304 max         304 min
Histogram: 1 0 0 0 0 0 0 0 0 0
Nghost:        4443.00 ave        4443 max        4443 min
Histogram: 1 0 0 0 0 0 0 0 0 0
Neighs:        123880.0 ave      123880 max      123880 min
Histogram: 1 0 0 0 0 0 0 0 0 0

Total # of neighbors = 123880
Ave neighs/atom = 407.50000
Neighbor list builds = 5
Dangerous builds not checked
Total wall time: 0:00:03
METRICS OPERATOR COLLECTION END
```

</details>

At this point we can shell into the app and interact with hpctoolkit!

```bash
# Shell into the application container "app"
kubectl exec -it metricset-sample-m-0-npbc9 -c app bash
```

You'll want to add the view bin to the path:

```bash
export PATH=/opt/share/view/bin:$PATH
```

You will see output in the present working directory where lammps was run:

```bash
# ls
README.txt      ffield.reax.hns              in.reaxc.hns                log.8Mar18.reaxc.hns.g++.4
data.hns-equil  hpctoolkit-lmp-measurements  log.8Mar18.reaxc.hns.g++.1  log.lammps
```

And we can use hpctoolkit commands to assemble our database, etc.

```bash
hpcstruct hpctoolkit-lmp-measurements
```
```console
ADVICE: See the usage message for how to use a structure cache to accelerate analysis of CPU and GPU binaries

INFO: Using a pool of 4 threads to analyze binaries in a measurement directory
INFO: Analyzing each large binary of >= 100000000 bytes in parallel using 4 threads
INFO: Analyzing each small binary using 2 threads

 begin parallel analysis of CPU binary libfabric.so.1.14.0 (size = 1389960, threads = 2)
 begin parallel analysis of CPU binary libhwloc.so.15.5.2 (size = 372064, threads = 2)
   end parallel analysis of CPU binary libfabric.so.1.14.0 (Cache disabled by user)
 begin parallel analysis of CPU binary libmpich.so.12.2.0 (size = 38697968, threads = 2)
   end parallel analysis of CPU binary libhwloc.so.15.5.2 (Cache disabled by user)
 begin parallel analysis of CPU binary libudev.so.1.7.2 (size = 166240, threads = 2)
   end parallel analysis of CPU binary libudev.so.1.7.2 (Cache disabled by user)
 begin parallel analysis of CPU binary lmp (size = 57324664, threads = 2)
   end parallel analysis of CPU binary libmpich.so.12.2.0 (Cache disabled by user)
   end parallel analysis of CPU binary lmp (Cache disabled by user)
```

And then "the professor!" 🤓️ (that's what I call this executable)

```bash
hpcprof hpctoolkit-lmp-measurements
```

This generates a database:

```bash
root@metricset-sample-m-0-npbc9:/opt/lammps/examples/reaxff/HNS# ls hpctoolkit-lmp-database
FORMATS.md  cct.db  meta.db  metrics  profile.db  src
```

That very likely we can do things with! E.g.,

```bash
apt-get update && apt-get install libgtk-3-dev
```

And this starts a viewer (but it won't work in the container)

```bash
hpcviewer ./hpctoolkit-lmp-database
```

What I did (from outside the container) is copy the database directory for later inspection

```bash
$ kubectl cp -c app metricset-sample-m-0-npbc9:/opt/lammps/examples/reaxff/HNS/hpctoolkit-lmp-database hpctoolkit-lmp-database
```

When you are done, clean up!

```bash
$ kubectl delete -f metrics.yaml
```