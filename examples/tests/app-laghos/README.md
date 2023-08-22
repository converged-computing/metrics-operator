# Laghos Example

This is an example of a metric app, Laghos. 
We have not yet added a Python example as we want a use case first, but can and will when it is warranted.

## Usage

Create a cluster

```bash
kind create cluster
```

and install JobSet to it.

```bash
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

Then create the metrics set. This is going to run a single run of LAMMPS over MPI!
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
and then the example running, and the log is printed to the console. Note this is example2 
provided in the container.

```bash
kubectl logs metricset-sample-l-0-0-lt782 -f
```
```console
METADATA START {"pods":2,"completions":2,"metricName":"app-laghos","metricDescription":"High-order Lagrangian Hydrodynamics Miniapp","metricType":"standalone","metricOptions":{"command":"mpirun -np 4 --hostfile ./hostlist.txt ./laghos","prefix":"/bin/bash","workdir":"/workflow/laghos"}}
METADATA END
Sleeping for 10 seconds waiting for network...
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT

       __                __                 
      / /   ____  ____  / /_  ____  _____   
     / /   / __ `/ __ `/ __ \/ __ \/ ___/ 
    / /___/ /_/ / /_/ / / / / /_/ (__  )    
   /_____/\__,_/\__, /_/ /_/\____/____/  
               /____/                       

Options used:
   --dimension 3
   --mesh default
   --refine-serial 2
   --refine-parallel 0
   --cartesian-partitioning ''
   --problem 1
   --order-kinematic 2
   --order-thermo 1
   --order-intrule -1
   --ode-solver 4
   --t-final 0.6
   --cfl 0.5
   --cg-tol 1e-08
   --ftz-tol 0
   --cg-max-steps 300
   --max-steps -1
   --partial-assembly
   --no-impose-viscosity
   --no-visualization
   --visualization-steps 5
   --no-visit
   --no-print
   --outputfilename results/Laghos
   --partition 0
   --device cpu
   --no-checks
   --no-mem
   --no-fom
   --no-gpu-aware-mpi
   --dev 0
Device configuration: cpu
Memory configuration: host-std
Number of zones in the serial mesh: 512
Non-Cartesian partitioning through METIS will be used.
Mesh::GeneratePartitioning(...): edgecut = 161
Zones min/max: 124 131
Number of kinematic (position, velocity) dofs: 14739
Number of specific internal energy dofs: 4096
step     5,     t = 0.0033,     dt = 0.000659,  |e| = 8.5702199098e+02
step    10,     t = 0.0066,     dt = 0.000686,  |e| = 7.0127018377e+02
step    15,     t = 0.0100,     dt = 0.000686,  |e| = 6.0421681833e+02
step    20,     t = 0.0135,     dt = 0.000686,  |e| = 5.4324701048e+02
step    25,     t = 0.0169,     dt = 0.000699,  |e| = 5.0040143051e+02
step    30,     t = 0.0205,     dt = 0.000742,  |e| = 4.6536535053e+02
Repeating step 33
Repeating step 35
step    35,     t = 0.0238,     dt = 0.000536,  |e| = 4.3916772365e+02
Repeating step 37
step    40,     t = 0.0262,     dt = 0.000456,  |e| = 4.2299947867e+02
step    45,     t = 0.0285,     dt = 0.000465,  |e| = 4.0914093543e+02
Repeating step 46
step    50,     t = 0.0305,     dt = 0.000395,  |e| = 3.9861827823e+02
step    55,     t = 0.0324,     dt = 0.000395,  |e| = 3.8907465383e+02
step    60,     t = 0.0344,     dt = 0.000395,  |e| = 3.8031603191e+02
step    65,     t = 0.0364,     dt = 0.000395,  |e| = 3.7221609042e+02
step    70,     t = 0.0384,     dt = 0.000395,  |e| = 3.6468819250e+02
Repeating step 71
step    75,     t = 0.0400,     dt = 0.000336,  |e| = 3.5869408943e+02
Repeating step 77
step    80,     t = 0.0415,     dt = 0.000286,  |e| = 3.5370225876e+02
step    85,     t = 0.0429,     dt = 0.000286,  |e| = 3.4910519380e+02
step    90,     t = 0.0444,     dt = 0.000297,  |e| = 3.4469307368e+02
step    95,     t = 0.0459,     dt = 0.000297,  |e| = 3.4032312257e+02
step   100,     t = 0.0473,     dt = 0.000297,  |e| = 3.3614681766e+02
step   105,     t = 0.0488,     dt = 0.000297,  |e| = 3.3215020871e+02
step   110,     t = 0.0503,     dt = 0.000297,  |e| = 3.2832683437e+02
Repeating step 113
step   115,     t = 0.0517,     dt = 0.000258,  |e| = 3.2495647602e+02
step   120,     t = 0.0530,     dt = 0.000258,  |e| = 3.2190798523e+02
step   125,     t = 0.0543,     dt = 0.000258,  |e| = 3.1896866967e+02
step   130,     t = 0.0555,     dt = 0.000263,  |e| = 3.1613115216e+02
step   135,     t = 0.0569,     dt = 0.000268,  |e| = 3.1331326201e+02
step   140,     t = 0.0582,     dt = 0.000279,  |e| = 3.1050474287e+02
step   145,     t = 0.0596,     dt = 0.000290,  |e| = 3.0770752528e+02
step   150,     t = 0.0611,     dt = 0.000296,  |e| = 3.0491220626e+02
step   155,     t = 0.0626,     dt = 0.000308,  |e| = 3.0213030501e+02
step   160,     t = 0.0642,     dt = 0.000327,  |e| = 2.9932940618e+02
step   165,     t = 0.0658,     dt = 0.000340,  |e| = 2.9648924185e+02
step   170,     t = 0.0676,     dt = 0.000347,  |e| = 2.9364581676e+02
step   175,     t = 0.0693,     dt = 0.000361,  |e| = 2.9084419820e+02
step   180,     t = 0.0712,     dt = 0.000375,  |e| = 2.8803967477e+02
step   185,     t = 0.0731,     dt = 0.000390,  |e| = 2.8523188687e+02
step   190,     t = 0.0751,     dt = 0.000406,  |e| = 2.8244131403e+02
step   195,     t = 0.0771,     dt = 0.000423,  |e| = 2.7964346624e+02
step   200,     t = 0.0793,     dt = 0.000440,  |e| = 2.7684052278e+02
step   205,     t = 0.0815,     dt = 0.000457,  |e| = 2.7403386214e+02
step   210,     t = 0.0838,     dt = 0.000476,  |e| = 2.7122428986e+02
step   215,     t = 0.0862,     dt = 0.000495,  |e| = 2.6841539261e+02
step   220,     t = 0.0887,     dt = 0.000505,  |e| = 2.6562272409e+02
step   225,     t = 0.0913,     dt = 0.000525,  |e| = 2.6285241373e+02
step   230,     t = 0.0939,     dt = 0.000536,  |e| = 2.6010186116e+02
step   235,     t = 0.0966,     dt = 0.000536,  |e| = 2.5743131013e+02
step   240,     t = 0.0993,     dt = 0.000569,  |e| = 2.5482506897e+02
step   245,     t = 0.1022,     dt = 0.000592,  |e| = 2.5212460051e+02
step   250,     t = 0.1053,     dt = 0.000616,  |e| = 2.4942592956e+02
step   255,     t = 0.1083,     dt = 0.000628,  |e| = 2.4676381746e+02
step   260,     t = 0.1115,     dt = 0.000640,  |e| = 2.4413966416e+02
step   265,     t = 0.1148,     dt = 0.000653,  |e| = 2.4155594336e+02
step   270,     t = 0.1181,     dt = 0.000680,  |e| = 2.3901843136e+02
step   275,     t = 0.1215,     dt = 0.000693,  |e| = 2.3651770360e+02
step   280,     t = 0.1250,     dt = 0.000707,  |e| = 2.3406528897e+02
step   285,     t = 0.1286,     dt = 0.000736,  |e| = 2.3165955702e+02
step   290,     t = 0.1322,     dt = 0.000750,  |e| = 2.2927866610e+02
step   295,     t = 0.1360,     dt = 0.000765,  |e| = 2.2693821184e+02
step   300,     t = 0.1399,     dt = 0.000796,  |e| = 2.2462847001e+02
step   305,     t = 0.1439,     dt = 0.000812,  |e| = 2.2234050930e+02
step   310,     t = 0.1480,     dt = 0.000845,  |e| = 2.2009171137e+02
step   315,     t = 0.1523,     dt = 0.000862,  |e| = 2.1786339521e+02
step   320,     t = 0.1566,     dt = 0.000897,  |e| = 2.1566337571e+02
step   325,     t = 0.1612,     dt = 0.000915,  |e| = 2.1347307972e+02
step   330,     t = 0.1658,     dt = 0.000952,  |e| = 2.1130063215e+02
step   335,     t = 0.1706,     dt = 0.000990,  |e| = 2.0914573282e+02
step   340,     t = 0.1756,     dt = 0.001010,  |e| = 2.0699799564e+02
step   345,     t = 0.1808,     dt = 0.001051,  |e| = 2.0486591370e+02
step   350,     t = 0.1861,     dt = 0.001093,  |e| = 2.0273375843e+02
step   355,     t = 0.1917,     dt = 0.001137,  |e| = 2.0061306102e+02
step   360,     t = 0.1974,     dt = 0.001160,  |e| = 1.9850591097e+02
step   365,     t = 0.2032,     dt = 0.001160,  |e| = 1.9646760495e+02
step   370,     t = 0.2091,     dt = 0.001207,  |e| = 1.9447522871e+02
step   375,     t = 0.2153,     dt = 0.001256,  |e| = 1.9246649863e+02
step   380,     t = 0.2218,     dt = 0.001332,  |e| = 1.9045188701e+02
step   385,     t = 0.2286,     dt = 0.001386,  |e| = 1.8840836032e+02
step   390,     t = 0.2355,     dt = 0.001386,  |e| = 1.8641756709e+02
step   395,     t = 0.2425,     dt = 0.001386,  |e| = 1.8450794495e+02
step   400,     t = 0.2494,     dt = 0.001386,  |e| = 1.8267410844e+02
step   405,     t = 0.2564,     dt = 0.001414,  |e| = 1.8090335785e+02
step   410,     t = 0.2635,     dt = 0.001471,  |e| = 1.7914396509e+02
step   415,     t = 0.2709,     dt = 0.001501,  |e| = 1.7740195941e+02
step   420,     t = 0.2786,     dt = 0.001592,  |e| = 1.7565682697e+02
step   425,     t = 0.2867,     dt = 0.001657,  |e| = 1.7388932766e+02
step   430,     t = 0.2951,     dt = 0.001690,  |e| = 1.7212865373e+02
step   435,     t = 0.3037,     dt = 0.001758,  |e| = 1.7039586175e+02
step   440,     t = 0.3126,     dt = 0.001793,  |e| = 1.6867112004e+02
step   445,     t = 0.3217,     dt = 0.001866,  |e| = 1.6696743007e+02
step   450,     t = 0.3312,     dt = 0.001941,  |e| = 1.6526474954e+02
step   455,     t = 0.3412,     dt = 0.002020,  |e| = 1.6355612235e+02
step   460,     t = 0.3513,     dt = 0.002060,  |e| = 1.6188102248e+02
step   465,     t = 0.3616,     dt = 0.002060,  |e| = 1.6024446396e+02
step   470,     t = 0.3719,     dt = 0.002060,  |e| = 1.5866879987e+02
step   475,     t = 0.3822,     dt = 0.002060,  |e| = 1.5715037032e+02
step   480,     t = 0.3925,     dt = 0.002060,  |e| = 1.5568577769e+02
step   485,     t = 0.4029,     dt = 0.002101,  |e| = 1.5426063598e+02
step   490,     t = 0.4136,     dt = 0.002186,  |e| = 1.5284012328e+02
step   495,     t = 0.4247,     dt = 0.002274,  |e| = 1.5141915216e+02
step   500,     t = 0.4363,     dt = 0.002366,  |e| = 1.4998743109e+02
step   505,     t = 0.4484,     dt = 0.002462,  |e| = 1.4854639097e+02
step   510,     t = 0.4610,     dt = 0.002561,  |e| = 1.4710309775e+02
step   515,     t = 0.4739,     dt = 0.002612,  |e| = 1.4568053311e+02
step   520,     t = 0.4871,     dt = 0.002665,  |e| = 1.4428281557e+02
step   525,     t = 0.5006,     dt = 0.002772,  |e| = 1.4289813255e+02
step   530,     t = 0.5147,     dt = 0.002828,  |e| = 1.4151567770e+02
step   535,     t = 0.5288,     dt = 0.002828,  |e| = 1.4017275099e+02
step   540,     t = 0.5429,     dt = 0.002828,  |e| = 1.3887794823e+02
step   545,     t = 0.5571,     dt = 0.002828,  |e| = 1.3762843809e+02
step   550,     t = 0.5712,     dt = 0.002828,  |e| = 1.3642137253e+02
step   555,     t = 0.5856,     dt = 0.002942,  |e| = 1.3523582574e+02
step   560,     t = 0.6000,     dt = 0.002449,  |e| = 1.3408616722e+02

CG (H1) total time: 293.3765060420
CG (H1) rate (megadofs x cg_iterations / second): 2.6940266849

CG (L2) total time: 12.4852572010
CG (L2) rate (megadofs x cg_iterations / second): 2.9752389880

Forces total time: 21.9565496320
Forces rate (megadofs x timesteps / second): 1.9455597858

UpdateQuadData total time: 128.8685935800
UpdateQuadData rate (megaquads x timesteps / second): 0.5787288115

Major kernels total time (seconds): 441.9408362730
Major kernels total rate (megadofs x time steps / second): 2.0538085859

Energy  diff: 6.90e-06
METRICS OPERATOR COLLECTION END
```

The above shows the structured output that is done in a way for our Python parsing script to easily
find sections of data. Also note that the worker will only be alive long enough for the main job to
finish, and once it does, the worker goes away! When you are done, the pods should be completed.

```bash
$ kubectl get pods
```
```console
NAME                           READY   STATUS      RESTARTS   AGE
metricset-sample-l-0-0-vfz4w   0/1     Completed   0          68s
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