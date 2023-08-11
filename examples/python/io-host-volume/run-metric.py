#!/usr/bin/env python3

import argparse
import os
import json
import time
from metricsoperator import MetricsOperator
import metricsoperator.utils as utils

here = os.path.abspath(os.path.dirname(__file__))
examples = os.path.dirname(os.path.dirname(here))
tests = os.path.join(examples, "tests")
metrics_yaml = os.path.join(tests, "io-host-volume", "metrics.yaml")

def get_parser():
    parser = argparse.ArgumentParser(
        description="Run Storage Metric and Get Output",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "--out",
        help="json file to save results",
        default=os.path.join(here, "io-host-volume.json"),
    )
    parser.add_argument(
        "--sleep",
        help="seconds to sleep allowing for pull and run",
        type=int,
        default=60,
    )
    parser.add_argument(
        "--test",
        help="run metrics.json assertions to check run",
        action="store_true",
        default=False,
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

    print(f'Sleeping {args.sleep} seconds to allow container to pull...')
    time.sleep(args.sleep)
    for output in m.watch():
        print(json.dumps(output, indent=4))
        utils.write_json(output, args.out)

    # Ensure we cleanup!
    m.delete()

if __name__ == "__main__":
    main()
