# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

from .network import network_netmark, network_osu_benchmark
from .perf import perf_sysstat
from .storage import io_sysstat

metrics = {
    "io-sysstat": io_sysstat,
    "network-osu-benchmark": network_osu_benchmark,
    "network-netmark": network_netmark,
    "perf-sysstat": perf_sysstat,
}


def get_metric(name):
    """
    Get a named metric parser.
    """
    metric = metrics.get(name)
    if not metric:
        raise ValueError(f"Metric {name} does not have a known parser")
    return metric
