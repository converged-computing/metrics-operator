# Metrics

The following metrics are under development (or being planned). These will be easier to add once I have the basic infrastructure for the different types!

## Implemented Metrics

### Performance

 - *perf-sysstat*: the "pidstat" executable of the sysstat library.

### Storage

 - *io-sysstat*: the "iostat" executable of the sysstat library.

### Standalone

 - *network-netmark*: this is currently a private container/software, but we have support for it when it's ready to be made public.

### Apps to be Measured

 - LAMMPS (already in tests)
 - Laghos
 - https://github.com/LLNL/Kripke

### Metrics To Be Added

 - https://github.com/glennklockwood/bioinformatics-profile
 - HPCToolkit
 - https://dl.acm.org/doi/pdf/10.1145/3611007
 - https://hpc.fau.de/research/tools/likwid/
 - https://www.vi-hps.org/tools/tools.html
 - https://open.xdmod.org/10.0/index.html


## Examples

The following examples are provided alongside the operator. Each directory has a README with complete instructions for usage.

 - [perf-lammps](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/perf-lammps)
 - [io-host-volume](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/io-host-volume)
 - [network-netmark](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/network-netmark) (code still private)
 
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
