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

This metric provides the "pidstat" executable of the sysstat library. The following options are available:

|Name | Description | Type | Default |
|-----|-------------|------------|------|
| color | Set to turn on color parsing | Anything set | unset |
| pids | For debugging, show consistent output of ps aux | Anything set | Unset |

By default color and pids are set to false anticipating log parsing.
And we also provide the option to see "commands" or specific commands based on a job index to the metric.
As an example, here is how we would ask to monitor two different commands for a launcher node (index 0)
and the rest (workers).

```yaml
- name: perf-sysstat
  rate: 2
  options:
    pids: "true"

  # Look for pids based on commands matched to index
  mapOptions:
    commands:
       # First set all to use the worker command, but give the lead broker a special command
       "all": /usr/libexec/flux/cmd/flux-broker --config /etc/flux/config -Scron.directory=/etc/flux/system/cron.d -Stbon.fanout
       "0": /usr/bin/python3.8 /usr/libexec/flux/cmd/flux-submit.py -n 2 --quiet --watch lmp -v x 2 -v y 2 -v z 2 -in in.reaxc.hns -nocite
```
In the map above, order matters, as the command for all indices is first set to be the flux-broker one, and then
after the index at 0 gets a custom command.

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
 - [network-osu-benchmark](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/network-osu-benchmark)

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
