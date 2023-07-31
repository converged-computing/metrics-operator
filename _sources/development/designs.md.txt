# Design Thinking

## Database for Metric Storage

I want to try creating a consistent database that can be used to store metrics across runs. In the space of an operator,
this means we can't clean it up when the specific metric is deleted, but rather it should be owned by the namespace.
I'm not sure how to do that but will think about ideas. Worst case, we have the user deploy the database in the same namespace
separately. Best case, we can manage it for them, or (better) not require it at all.
I don't want anything complicated (I don't want to re-create prometheus or a monitoring service!)

## Kubernetes Objects

JobSet gives us a lot of flexibility to deploy different kinds of applications or services alongside one another. I'm wondering if we can
have some design where a replicated job corresponds to one metric, and then one run can include one or more metrics.
I also like the design of the JobSet, so I'm going to design something similar, a `MetricSet` that might allow
for several metrics to be run across an application or context of choice (e.g., storage).

### Metrics

The following metrics are interesting, and here is how we might measure them with an operator:

### Performance

More generally, I wonder if we can add the SYS_PTRACE capability to containers in the same pod and then be able to monitor processes from one container into another? If we are able to know one or more processes of interest, and find tools that can give meaningful 
metrics from the processes, that could be a cool setup. I [tested the shared process namespace](https://vsoch.github.io/2023/shared-process-namespace/) and think this is a good idea, at least to start. For the current implementation, we allow
a "perf" flavored metric to be deployed alongside an application of interest, and then get access to the PID of the running process. We are currently running the monitoring tool under the PID is no longer found, and then finishing.
Metric output is in the pod logs, and hopefully we can improve upon this. In addition to performance, it would be nice to have a simple means to measure the timing of the application.

### Storage

Setting up storage, typically by way of a persistent volume claim that turns into a persistent volume, is complex. To start, I'm going to require that the user create the PVC on their own,
and then provide information about it to the operator. The operator will then request a volume, measure something on it for some rate and length of time, and then clean up.

### Others

There are likely others (and I need to think about it)