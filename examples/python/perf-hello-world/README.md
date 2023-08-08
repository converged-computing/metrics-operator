# Performance Sysstat (pidstat) Hello World in Python

I'm trying to understnd what I'm looking at for output of sysstat (pidstat) metrics.
This assumes you have a running cluster with JobSet and the Metrics Operator installed,

```bash
$ python run-metric.py
```

The script will use the metricsoperator library to submit the [metrics.yaml](metrics.yaml)
and then wait for the pod to complete and parse the output in the log.

