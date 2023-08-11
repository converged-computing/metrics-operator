# OSU Benchmarks Python

This is a quick example to show running and plotting a few of the benchmarks!
This assumes you have a running cluster with JobSet and the Metrics Operator installed!

```bash
$ python run-metric.py
```

The script will use the metricsoperator library to submit the associated yaml under [tests](../../tests)
and then wait for the pod to complete and parse the output in the log.

## Example Plots

![img/OSU-MPI_Accumulate-latency-Test-v5.8.png](img/OSU-MPI_Accumulate-latency-Test-v5.8.png)
![img/OSU-MPI_Get_accumulate-latency-Test-v5.8.png](img/OSU-MPI_Get_accumulate-latency-Test-v5.8.png)
![img/OSU-MPI_Get-latency-Test-v5.8.png](img/OSU-MPI_Get-latency-Test-v5.8.png)
![img/OSU-MPI_Put-Latency-Test-v5.8.png](img/OSU-MPI_Put-Latency-Test-v5.8.png)
![img/OSU-MPI-Allreduce-Latency-Test-v5.8.png](img/OSU-MPI-Allreduce-Latency-Test-v5.8.png)
![img/OSU-MPI-Bandwidth-Test-v5.8.png](img/OSU-MPI-Bandwidth-Test-v5.8.png)
![img/OSU-MPI-Bi-Directional-Bandwidth-Test-v5.8.png](img/OSU-MPI-Bi-Directional-Bandwidth-Test-v5.8.png)
![img/OSU-MPI-Latency-Test-v5.8.png](img/OSU-MPI-Latency-Test-v5.8.png)