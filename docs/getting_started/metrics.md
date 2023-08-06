# Metrics

The following metrics are under development (or being planned).

 - [Examples](https://converged-computing.github.io/metrics-operator/getting_started/metrics.html#examples)
 - [Storage Metrics](https://converged-computing.github.io/metrics-operator/getting_started/metrics.html#storage)
 - [Application Metrics](https://converged-computing.github.io/metrics-operator/getting_started/metrics.html#application)
 - [Standalone Metrics](https://converged-computing.github.io/metrics-operator/getting_started/metrics.html#standalone)

<iframe src="../_static/data/table.html" style="width:100%; height:400px;" frameBorder="0"></iframe>

All metrics can be customized with the following variables

|Name | Description | Type | Default |
|-----|-------------|------------|------|
| completions | Number of times to run metric | int32 | unset (runs for lifetime of application or indefinitely) |
| rate | Seconds to pause between measurements | int32 | 10 |

## Implemented Metrics

### Performance

These metrics are intended to assess application performance.

#### perf-sysstat

 - [Application Metric Set](user-guide.md#application-metric-set)

This metric provides the "pidstat" executable of the sysstat library.

### Storage

These metrics are intended to assess storage volumes.

#### io-sysstat

 - [Storage Metric Set](user-guide.md#application-metric-set)

This is the "iostat" executable of the sysstat library.

|Name | Description | Type | Default |
|-----|-------------|------------|------|
| human | Show tabular, human-readable output inside of json | string "true" or "false" | "false" |

### Standalone

#### network-netmark

 - [Standalone Metric Set](user-guide.md#application-metric-set)

This is currently a private container/software, but we have support for it when it's ready to be made public (networking)
Variables to customize include:

|Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| tasks | Total number of tasks across pods | options->tasks | string | nproc * pods |
| warmups | Number of warmups | options->warmups | int32 | 10 |
| trials | Number of trials | options->trials | int32 | 20 |
| sendReceiveCycles | Number of send-receive cycles | options-sendReceiveCycles | int32 | 20 |
| messageSize | Message size in bytes | options->messageSize | int32 | 0 |
| storeEachTrial | Flag to indicate storing each trial data | options->storeEachTrial | string (true/false) | "true" |

#### network-osu-benchmark

Point to point benchmarks for MPI (networking). If listOptions->commands not set, will use all one-point commands.
Variables to customize include:

|Name | Description | Option Key | Type | Default |
|-----|-------------|------------|------|---------|
| commands | Custom list of osu-benchmark one-sided commands to run | listOptions->commands | array | unset uses default set |


## Examples

The following examples are provided alongside the operator. Each directory has a README with complete instructions for usage.

 - [perf-hello-world](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/perf-hello-world)
 - [perf-lammps](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/perf-lammps)
 - [io-host-volume](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/io-host-volume)
 - [network-netmark](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/network-netmark) (code still private)
 - [network-osu-benchmarks](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/network-osubenchmarks) (code still private)

### Apps to be Measured

 - Laghos
 - https://github.com/LLNL/Kripke

### Metrics To Be Added

 - https://github.com/glennklockwood/bioinformatics-profile
 - HPCToolkit
 - https://dl.acm.org/doi/pdf/10.1145/3611007
 - https://hpc.fau.de/research/tools/likwid/
 - https://www.vi-hps.org/tools/tools.html
 - https://open.xdmod.org/10.0/index.html


## Containers

The following tools are folded into the metrics above. Often, one tool can be built into one container and used across multiple metrics.

### Sysstat

 - [ghcr.io/converged-computing/metric-sysstat](https://github.com/converged-computing/metrics-operator/pkgs/container/metric-sysstat)

text.startsWith("Hello");

Sysstat is stored as a general metrics analyzer, as it provides several different metric types; It generally provides utils to monitor system performance and usage, including:

- *iostat* reports CPU statistics and input/output statistics for block devices and partitions.
- *mpstat* reports individual or combined processor related statistics.
- *pidstat* reports statistics for Linux tasks (processes) : I/O, CPU, memory, etc.
- *tapestat* reports statistics for tape drives connected to the system.
- *cifsiostat* reports CIFS statistics.

## LLNL Storage / Filesystems

 - NFS
 - Vast
 - Lustre
