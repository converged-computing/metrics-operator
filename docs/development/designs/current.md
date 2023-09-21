# Current Design

For this second design,m we can more easily say:

> A Metric Set is a collection of metrics to measure IO, performance, or networking that can be customized with addons.

The original design was a good first shot, but was flawed in several ways:

1. I could not combined metrics into one. E.g., if I wanted to use a launcher jobset design combined with HPCToolkit, another metric, I could not.
2. The top level set types - standalone, application, and storage, didn't have much meaning.
3. The use of Storage, Application, and Volume was messy at best (external entities to add to a metric set)

For this second design, the "MetricSet" is still mirroring the design of a JobSet, but it is more generic, and of one type. There are no longer different
flavors of metric sets. Rather, we allow metrics to generate replicated jobs. For the "extras" that we need to integrate to supplement those jobs - e.g., applications, volumes/storage, or
even extra containers that add logic, these are now called metric addons.  More specifically, an addon can:

 - Add extra containers (and config maps for their entrypoints)
 - Add custom logic to entrypoints for specific jobs and/or containers
 - Add additional volumes that range the gamut from empty to persistent disk.

The current design allows only one JobSet per metrics.yaml, and this was an explicit choice after realizing that it's unlikely to want more than one.

## Kubernetes Abstractions

We use a JobSet on the top level with Replica set to 1, and within that set, each metric is allowed to create one or more ReplcatedJob. We can easily customize the style of the replicated job based
on interfacs. E.g.,:

- The `LauncherWorker` is a typical design that might have a launcher and MPI hostlist written, and a main command run there to then interact with the workers.
- The `SingleApplication` is a basic design that expects one or more pods in an indexed job, and also shares the process namespace.
- The `StorageGeneric` is almost the same, but doesn't share a process namespace.

I haven't found a need for another kind of design yet (most are the launcher worker type) but can easily add them if needed.
There is no longer any distinction between MetricSet types, as there is only one MetricSet that serves as a shell from the metric.

## Output Options

### Logging Parser

For the simplest start, I've decided to allow for metrics to have their own custom output (indeed it would be hard to standardize this between so many different tools) but have the operator
provide structure to that, meaning separators to distinguish sections, and a consistent way to output metadata. As an example, here is what the top level metadata and sections (with some custom output data between)
would look like:

```console
METADATA START {"pods":1,"completions":1,"storageVolumePath":"/workflow","storageVolumeHostPath":"/tmp/workflow","metricName":"io-sysstat","metricDescription":"statistics for Linux tasks (processes) : I/O, CPU, memory, etc.","metricType":"storage","metricOptions":{"completions":2,"human":"false","rate":10}}
METADATA END
METRICS OPERATOR COLLECTION START
METRICS OPERATOR TIMEPOINT
...custom data output here for timepoint 1...
METRICS OPERATOR TIMEPOINT
...custom data output here for timepoint 2...
METRICS OPERATOR TIMEPOINT
...custom data output here for timepoint N...
METRICS OPERATOR COLLECTION END
```

In the above, we can parse the metadata for the run from the first line (a subset of flattened, important features dumped in json) and then clearly mark the start and end of collection,
along with separation between timepoints. This is the most structure we can provide, as each metric output looks different. It's up to the Python module parser from the "metricsoperator"
module to know how to parse (and possibly plot) any specific output type.