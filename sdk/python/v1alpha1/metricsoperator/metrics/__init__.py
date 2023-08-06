# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

from .storage import io_sysstat

metrics = {"io-sysstat": io_sysstat}


def get_metric(name):
    """
    Get a named metric parser.
    """
    metric = metrics.get(name)
    if not metric:
        raise ValueError("Metric %s does not have a known parser")
    return metric
