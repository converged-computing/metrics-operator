#!/usr/bin/env python3

import argparse
import os
import sys
import metricsoperator.utils as utils
import metricsoperator.metrics as mutils

here = os.path.abspath(os.path.dirname(__file__))


def get_parser():
    parser = argparse.ArgumentParser(
        description="Parse MPI Trace Output",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "indir",
        help="input directory with files",
    )
    parser.add_argument(
        "--out",
        help="json file to save results",
        default=os.path.join(here, "metrics.json"),
    )
    parser.add_argument(
        "--prefix",
        help="prefix for mpi trace process files (defaults to mpi_profile)",
        default="mpi_profile",
    )
    return parser


def main():
    """
    Run a job.
    """
    parser = get_parser()
    args, _ = parser.parse_known_args()

    # Save listing of results
    results = []

    if not args.indir or not os.path.exists(args.indir):
        sys.exit("Input directory {args.indir} does not exist")

    # Create an instance of the metrics parser
    metric = mutils.get_metric("addon-mpitrace")()
    indir = os.path.abspath(args.indir)

    for filename in utils.recursive_find(indir, pattern=args.prefix):
        print(f"Parsing filename {filename}")
        log = utils.read_file(filename)
        results.append(metric.parse_log(log))
    utils.write_json(results, args.out)


if __name__ == "__main__":
    main()
