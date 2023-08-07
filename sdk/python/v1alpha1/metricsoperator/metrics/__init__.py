# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

from .network import network_osu_benchmark
from .storage import io_sysstat

metrics = {
    "io-sysstat": io_sysstat,
    "network-osu-benchmark": network_osu_benchmark,
}


def get_metric(name):
    """
    Get a named metric parser.
    """
    metric = metrics.get(name)
    if not metric:
        raise ValueError(f"Metric {name} does not have a known parser")
    return metric
