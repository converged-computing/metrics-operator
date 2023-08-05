# metrics-operator

![docs/images/metrics-operator-banner.png](docs/images/metrics-operator-banner.png)

Developing metrics and a catalog of applications to assess different kinds of Kubernetes performance.
We likely will choose different metrics that are important for HPC.
Note that I haven't started the operator yet because I'm [testing ideas for the design](hack/test).

View our ⭐️ [Documentation](https://converged-computing.github.io/metrics-operator/) ⭐️

## Dinosaur TODO

- Find better logging library for logging outside of controller
- Need a strategy for storing metrics output / logs
- For services we are measuring, we likely need to be able to kill after N seconds (to complete job) or to specify the success policy on the metrics containers instead of the application
- Start releasing earlier this time :)
- TBA

## License

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614
