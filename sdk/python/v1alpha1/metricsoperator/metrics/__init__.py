# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

import metricsoperator.metrics.app as apps
import metricsoperator.metrics.base as base
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


def get_metric(name=None):
    """
    Get a named metric parser.
    """
    metric = metrics.get(name)
    # If we don't have a matching metric, return base (for raw logs)
    if not metric:
        return base.MetricBase
    return metric
