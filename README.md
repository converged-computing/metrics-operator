# metrics-operator

![docs/images/metrics-operator-banner.png](docs/images/metrics-operator-banner.png)

Developing metrics and a catalog of applications to assess different kinds of Kubernetes performance.
We likely will choose different metrics that are important for HPC.
Note that I haven't started the operator yet because I'm [testing ideas for the design](hack/test).
To learn more:

- ⭐️ [Documentation](https://converged-computing.github.io/metrics-operator/) ⭐️
- 🐯️ [Python Module](https://pypi.org/project/metricsoperator/) 🐯️

## Dinosaur TODO

- Find better logging library for logging outside of controller
- For larger metric collections, we should have a log streaming mode (and not wait for Completed/Successful)
- For services we are measuring, we likely need to be able to kill after N seconds (to complete job) or to specify the success policy on the metrics containers instead of the application
- Python function to save entire spec to yaml (for MetricSet and JobSet)?
- Netmark / OSU need resources set to ensure 1 pod/node
- Metrics parsers to do (need to add separators, formatting, python parser):
  - perf-sysstat
- Plotting examples needed for
  - perf-sysstat
  - io-sysstat

## License

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614
