# Development

## Creation

```bash
mkdir metrics-operator
cd metrics-operator/
operator-sdk init --domain flux-framework.org --repo github.com/converged-computing/metrics-operator
operator-sdk create api --version v1alpha1 --kind MetricSet --resource --controller
```

