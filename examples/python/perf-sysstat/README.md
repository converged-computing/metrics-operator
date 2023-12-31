# Performance Sysstat (pidstat) Python

This is a quick example to show running and plotting pidstat metrics!
This is done with a small LAMMPS run so we at least have a few timepoints.
This assumes you have a running cluster with JobSet and the Metrics Operator installed,

```bash
$ python run-metric.py
```

The script will use the metricsoperator library to submit the [metrics.yaml](metrics.yaml)
and then wait for the pod to complete and parse the output in the log.

## Example Plots

![img/threads_task_time-hist.png](img/threads_task_time-hist.png)
![img/pagefaults_task_rss-hist.png](img/pagefaults_task_rss-hist.png)
![img/pagefaults_task_time.png](img/pagefaults_task_time.png)
![img/cpu_statistics_task_cpu-hist.png](img/cpu_statistics_task_cpu-hist.png)
![img/cpu_statistics_child_time-hist.png](img/cpu_statistics_child_time-hist.png)
![img/pagefaults_child_time.png](img/pagefaults_child_time.png)
![img/kernel_tables_time.png](img/kernel_tables_time.png)
![img/pagefaults_task_percent_mem-hist.png](img/pagefaults_task_percent_mem-hist.png)
![img/pagefaults_child_time-hist.png](img/pagefaults_child_time-hist.png)
![img/pagefaults_task_percent_mem.png](img/pagefaults_task_percent_mem.png)
![img/threads_task_time.png](img/threads_task_time.png)
![img/stack_utilization_time-hist.png](img/stack_utilization_time-hist.png)
![img/pagefaults_task_rss.png](img/pagefaults_task_rss.png)
![img/task_switching_time.png](img/task_switching_time.png)
![img/threads_task_cpu.png](img/threads_task_cpu.png)
![img/threads_child_time-hist.png](img/threads_child_time-hist.png)
![img/pagefaults_task_time-hist.png](img/pagefaults_task_time-hist.png)
![img/kernel_tables_time-hist.png](img/kernel_tables_time-hist.png)
![img/cpu_statistics_task_time-hist.png](img/cpu_statistics_task_time-hist.png)
![img/metrics.csv](img/metrics.csv)
![img/task_switching_time-hist.png](img/task_switching_time-hist.png)
![img/cpu_statistics_task_time.png](img/cpu_statistics_task_time.png)
![img/stack_utilization_stkref-hist.png](img/stack_utilization_stkref-hist.png)
![img/cpu_statistics_task_cpu.png](img/cpu_statistics_task_cpu.png)
![img/stack_utilization_time.png](img/stack_utilization_time.png)
![img/single-task](img/single-task)
![img/policy_time-hist.png](img/policy_time-hist.png)
![img/kernel_statistics_time-hist.png](img/kernel_statistics_time-hist.png)
![img/pagefaults_task_vsz.png](img/pagefaults_task_vsz.png)
![img/pagefaults_task_vsz-hist.png](img/pagefaults_task_vsz-hist.png)
![img/stack_utilization_stksize.png](img/stack_utilization_stksize.png)
![img/threads_child_time.png](img/threads_child_time.png)
![img/threads_task_cpu-hist.png](img/threads_task_cpu-hist.png)
![img/stack_utilization_stksize-hist.png](img/stack_utilization_stksize-hist.png)
![img/policy_time.png](img/policy_time.png)
![img/stack_utilization_stkref.png](img/stack_utilization_stkref.png)
![img/kernel_statistics_time.png](img/kernel_statistics_time.png)
![img/cpu_statistics_child_time.png](img/cpu_statistics_child_time.png)
