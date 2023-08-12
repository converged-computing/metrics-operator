# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

import metricsoperator.metrics.app as apps
import metricsoperator.metrics.network as network
import metricsoperator.metrics.perf as perf
import metricsoperator.metrics.storage as storage

metrics = {
    "io-sysstat": storage.io_sysstat,
    "network-osu-benchmark": network.network_osu_benchmark,
    "network-netmark": network.network_netmark,
    "perf-sysstat": perf.perf_sysstat,
    "io-fio": storage.io_fio,
    "app-lammps": apps.app_lammps,
    "app-amg": apps.app_amg,
}


def get_metric(name):
    """
    Get a named metric parser.
    """
    metric = metrics.get(name)
    if not metric:
        raise ValueError(f"Metric {name} does not have a known parser")
    return metric
