# In Progress Metrics

These are metrics that are consistered under development (and likely need more eyes) to get fully working.

## Network

### network-chatterbug

 - *[network-chatterbug](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/network-chatterbug)*

Chatterbug provides a [suite of communication proxy applications](https://github.com/hpcgroup/chatterbug) for HPC.
We use a launcher/worker design.

|Name | Description  | Type | Default |
|-----|--------------|------|---------|
| mpirun | The options to give to mpirun (includes tasks) | string | `-N 8` |
| command | The chatterbug command (subdirectory) to run, see options below | string | stencil3d |
| args | Arguments for the command | string | `1 2 2 10 10 10 4 1` |
| sole-tenancy | Require sole tenancy | string ("true" or "false") | "true" |

By default, we require sole-tenancy, but you can disable this. Note that the best place to look for "documentation"
on the commands seems to be [the source code]((https://github.com/hpcgroup/chatterbug)). The following command options
are available for `command`:

- pairs
- ping-ping
- spread
- stencil3d
- stencil4d
- subcom2d-coll
- subcom2d-a2a
- unstr-mesh

We have tested mostly stencil3d. Note that the mpirun command is parsed as follows:

```bash
$ mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 4 /root/chatterbug/${command}/${executable} ${args}
```

Thus for the defaults, you'd get this command (on one pod):

```bash
$ mpirun --hostfile ./hostfile.txt --allow-run-as-root -N 4 /root/chatterbug/stencil3d/stencil3d.x 1 2 2 10 10 10 4 1
```

See the example linked in the header for a metrics.yaml example.

## Standalone

### app-hpl

 - *[app-hpl](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/app-hpl)*

The [Linpack](https://ulhpc-tutorials.readthedocs.io/en/production/parallel/mpi/HPL/) benchmark is used for the [Top500](https://www.top500.org/project/linpack/),
and generally is solving a dense system of linear equations. Arguments to customize include the following:

| Name | Description  | Type | Default |
|-----|-------------|------|---------|
| mpiargs | Arguments to give to mpi | string | empty string |
| tasks | Number of tasks per node | int32 | detected used nproc | 
| ratio | target memory occupation | string (but as a float, e.g., "0.3") | "0.3" |
| memory | memory in GiB | int32 | detected from proc |
| blocksize | blocksize is the NBs "number blocks" value | int32 | |
| pfact | | int32 | |
| nbmin | | int32 | |
| ndiv | | int32 | |
| row_or_colmajor_pmapping | PMAP process mapping (0=Row-,1=Column-major) | int32 | 0 |
|	rfact | (0=left, 1=Crout, 2=Right) | int32 | 0 |
| bcast | (0=1rg,1=1rM,2=2rg,3=2rM,4=Lng,5=LnM) | int32 | 0 |
| depth | number of lookahead depth | int32 | 0 |
| swap | (0=bin-exch,1=long,2=mix) | int32 | 0 |
| swappingThreshold | | int32 | 64 |
| l1transposed | (0=transposed,1=no-transposed) | int32 | 0 |
| utransposed | (0=transposed,1=no-transposed) | int32 | 0 |
| memAlignment | memory alignment in double (> 0) (4,8,16) | int32 | |

For the meaning of each of these, see [this documentation](https://ulhpc-tutorials.readthedocs.io/en/production/parallel/mpi/HPL/#hpl-main-parameters)
and how they are used in [hpl.go](https://github.com/converged-computing/metrics-operator/tree/main/pkg/metrics/app/hpl.go)
I made an effort to define them above, but you should consult the documentation above, because I don't fully
understand these yet.

We provide a simple build here, as typically vendors spend a lot of time custom-compiling the code
for their architectures (and we are compiling for general use). We will use a script `compute_N` from the OLHPC Tutorials to generate input data for a particular
problem size, and you can vary the input to this script via the `computeArgs` parameters. We use a default, and you can inspect the
script help below:

<details>

<summary>compute_N --help</summary>

```console
# compute_N -h
Compute N for HPL runs.

SYNOPSIS
  compute_N [-v] [--mem <SIZE_IN_GB>] [-N <NODES>] [-r <RATIO>] [-NB <NB>]
  compute_N [-v] [--mem <SIZE_IN_GB>] [-N <NODES>] [-p <PERCENTAGE_MEM>] [-NB <NB>]

  The following formulae is used (when using '-r <ratio>'):
     N = <ratio>*SQRT( Total Memory Size in bytes / sizeof(double) )
       = <ratio>*SQRT( <nnodes> * <ram_size> / 8)

  Alternatively you may wish to specify a memory usage ratio (with -p <percentage_mem>),
  in which case the following formulae is used:
      N = SQRT( <percentage_mem>/100 * Total Memory Size in bytes / sizeof(doubl)

OPTIONS
  -m --mem --ramsize <SIZE>
     Specify the total memory size per node, in GiB.
     Default RAM size consider (yet in KiB): 16051112 KiB
  -N --nodes <N>
     Number of compute nodes
  -NB <NB>
     NB parameters to use. Default: 192 (384 for skylake)
  -p --memshare <PERCENTAGE_MEM>
     Percentage of the total memory size to use.
     Derived from the below global ratio (i.e. 0% since RATIO=0.8)
  -r --ratio <RATIO>
     Global ratio to apply. Default: 0.8

EXAMPLE
  For 2 broadwell nodes on iris cluster, using 30% of the total memory per node:
     compute_N -N 2 -p 30 -m 128 -NB 192
  For 4 skylake nodes on iris cluster, using 85% of the total memory per node:
     compute_N -N 4 -p 85 -m 128 -NB 384

AUTHORS
  Sebastien Varrette <Sebastien.Varrette@uni.lu> and UL HPC Team

COPYRIGHT
  This is free software; see the source for copying conditions.  There is
  NO warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
```

</details>

The following examples are [provided](https://ulhpc-tutorials.readthedocs.io/en/production/parallel/mpi/HPL/) to generate the HPL.dat for the analysis:

```bash
/opt/tutorials/benchmarks/HPL/scripts/compute_N -h
# 1 Broadwell node, alpha = 0.3
/opt/tutorials/benchmarks/HPL/scripts/compute_N -m 128 -NB 192 -r 0.3 -N 1
# 2 Skylake (regular) nodes, alpha = 0.3
/opt/tutorials/benchmarks/HPL/scripts/compute_N -m 128 -NB 384 -r 0.3 -N 2
# 4 bigmem (skylake) nodes, beta = 0.85
/opt/tutorials/benchmarks/HPL/scripts/compute_N -m 3072 -NB 384 -p 85 -N 4
```

Here is a tiny setup I created for a testing case:

```bash
/opt/tutorials/benchmarks/HPL/scripts/compute_N -m 128 -NB 192 -r 0.3 -N 2
```

Next, you might care about the input data, a file called `hpl.dat`. By default we use 
a template that is populated by the above variables, and here is another example that I found
in the repository:

<details>

<summary>Default hpl.dat</summary>

```console
HPLinpack benchmark input file
Innovative Computing Laboratory, University of Tennessee
HPL.out      output file name (if any)
6            device out (6=stdout,7=stderr,file)
1            # of problems sizes (N)
24650         Ns
1            # of NBs
192           NBs
0            PMAP process mapping (0=Row-,1=Column-major)
2            # of process grids (P x Q)
2 4             Ps
14 7            Qs
16.0         threshold
1            # of panel fact
2            PFACTs (0=left, 1=Crout, 2=Right)
1            # of recursive stopping criterium
4            NBMINs (>= 1)
1            # of panels in recursion
2            NDIVs
1            # of recursive panel fact.
1            RFACTs (0=left, 1=Crout, 2=Right)
1            # of broadcast
1            BCASTs (0=1rg,1=1rM,2=2rg,3=2rM,4=Lng,5=LnM)
1            # of lookahead depth
1            DEPTHs (>=0)
2            SWAP (0=bin-exch,1=long,2=mix)
64           swapping threshold
0            L1 in (0=transposed,1=no-transposed) form
0            U  in (0=transposed,1=no-transposed) form
1            Equilibration (0=no,1=yes)
8            memory alignment in double (> 0)
##### This line (no. 32) is ignored (it serves as a separator). ######
0                               Number of additional problem sizes for PTRANS
1200 10000 30000                values of N
0                               number of additional blocking sizes for PTRANS
40 9 8 13 13 20 16 32 64        values of NB
```

</details>

If there is something above not properly exposed please [let us know](https://github.com/converged-computing/metrics-operator/issues).
