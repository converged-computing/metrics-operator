# Addon MPITrace

MPI Trace can be added to an app and then used to generate files to parse with the library here.
This addon is a bit different in that it doesn't expect output in logs, but rather provided on the local
system (saved via kubectl or via the [ORAS Operator](https://github.com/converged-computing/oras-operator)).
For this example, you should first run the example in the [addon-mpitrace](../../addons/mpitrace-lammps) directory
and generate the output files:

```bash
 ls ../../addons/mpitrace-lammps/mpi_profile.114.* -l
-rw-rw-r-- 1 vanessa vanessa 5636 Oct 22 18:52 ../../addons/mpitrace-lammps/mpi_profile.114.0
-rw-rw-r-- 1 vanessa vanessa 4684 Oct 22 18:52 ../../addons/mpitrace-lammps/mpi_profile.114.1
-rw-rw-r-- 1 vanessa vanessa 4623 Oct 22 18:52 ../../addons/mpitrace-lammps/mpi_profile.114.2
```

Then run the script, and target the files:

```bash
$ python parse-metric.py ../../addons/mpitrace-lammps/ --prefix mpi_profile
```



## Example Plots

Here is a small example - 5 run with 2 pods!

![img/lammps_lammps.png](img/lammps_lammps.png)
