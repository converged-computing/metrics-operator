# Development

## Writing Metric Containers

- They should contain wget if they need to download the wait script in the entrypoint

## Creation

```bash
mkdir metrics-operator
cd metrics-operator/
operator-sdk init --domain flux-framework.org --repo github.com/converged-computing/metrics-operator
operator-sdk create api --version v1alpha1 --kind MetricSet --resource --controller
```

