#!/usr/bin/env python3

import argparse
import os
import json
import time

from metricsoperator import MetricsOperator
import metricsoperator.utils as utils

here = os.path.abspath(os.path.dirname(__file__))
metrics_yaml = os.path.join(here, "metrics.yaml")

def get_parser():
    parser = argparse.ArgumentParser(
        description="Run Storage Metric and Get Output",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "--out",
        help="json file to save results",
        default=os.path.join(here, "metrics.json"),
    )
    parser.add_argument(
        "--sleep",
        help="seconds to sleep (for container to pull)",
        default=60,
        type=int,
    )
    return parser


def main():
    """
    Run a job.
    """
    parser = get_parser()
    args, _ = parser.parse_known_args()

    # Create a metrics operator with our metrics.yaml
    m = MetricsOperator(metrics_yaml)
    m.create()

    # Give pods time to create
    print("Sleeping to give containers time to pull...")
    time.sleep(args.sleep)
    for output in m.watch():
        print(json.dumps(output, indent=4))
        utils.write_json(output, args.out)


if __name__ == "__main__":
    main()
