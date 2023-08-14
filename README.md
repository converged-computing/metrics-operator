# metrics-operator

![docs/images/metrics-operator-banner.png](docs/images/metrics-operator-banner.png)

Developing metrics and a catalog of applications to assess different kinds of Kubernetes performance.
We likely will choose different metrics that are important for HPC.
Note that I haven't started the operator yet because I'm [testing ideas for the design](hack/test).
To learn more:

- ⭐️ [Documentation](https://converged-computing.github.io/metrics-operator/) ⭐️
- 🐯️ [Python Module](https://pypi.org/project/metricsoperator/) 🐯️

## Dinosaur TODO

- We need a way for the entrypoint command to monitor (based on the container) to differ (potentially)
- make flux operator command generator
- For larger metric collections, we should have a log streaming mode (and not wait for Completed/Successful)
- For services we are measuring, we likely need to be able to kill after N seconds (to complete job) or to specify the success policy on the metrics containers instead of the application
- Look into pod affinity/anti-affintiy vs. topology constraint (which do we want)?
- Add assertions checking for python tests
- Plotting examples (python parsers) needed for
  - io-sysstat
  - app-kripke
  - app-quicksilver
  - app-pennant

## License

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614
