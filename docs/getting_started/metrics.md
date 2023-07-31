# Metrics

The following metrics are under development (or being planned). These will be easier to add once I have the basic infrastructure for the different types!

## Implemented Metrics

### Performance

 - *perf-sysstat*: the "pidstat" executable of the sysstat library.

### Storage

 - *io-sysstat*: the "iostat" executable of the sysstat library.

### Apps to be Measured

 - LAMMPS (already in tests)
 - Laghos
 - https://github.com/LLNL/Kripke

### Metrics To Be Added

 -  https://github.com/glennklockwood/bioinformatics-profile


## Examples

The following examples are provided alongside the operator. Each directory has a README with complete instructions for usage.

 - [perf-lammps](https://github.com/converged-computing/metrics-operator/tree/main/examples/tests/perf-lammps)
 
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
