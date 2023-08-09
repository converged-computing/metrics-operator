#!/usr/bin/env python3

import argparse
import os
import json
import time
from io import StringIO
from metricsoperator import MetricsOperator
import metricsoperator.utils as utils
import seaborn as sns
import matplotlib.pyplot as plt
import pandas

here = os.path.abspath(os.path.dirname(__file__))
metrics_yaml = os.path.join(here, "metrics.yaml")

plt.style.use("bmh")

def get_parser():
    parser = argparse.ArgumentParser(
        description="Run Netmark Metric and Get Output",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "--out",
        help="json file to save results",
        default=os.path.join(here, "metrics.json"),
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

    print(f'Sleeping {args.sleep} seconds so container can pull...')
    time.sleep(args.sleep)
    for output in m.watch():
        print(json.dumps(output, indent=4))
        utils.write_json(output, args.out)
        plot_results(output)

def plot_results(output):
    """
    Plot results to a histogram and matrix heatmap
    """
    # Directory for plotting results
    img = os.path.join(here, "img")
    if not os.path.exists(img):
        os.makedirs(img)
   
    for result in output['data']:
        data = StringIO(result['RTT.csv'])
        df = pandas.read_csv(data, sep=",", index_col=False, header=None)
        df.to_csv(os.path.join(img, "RTT.csv"))
        tasks = df.shape[0]

        # Save a heatmap of the particular run
        plt.figure(figsize=(36, 36))
        sns.heatmap(df, cmap="BrBG", annot=True)
        plt.title(f"MPI Connection Times (microseconds) for {tasks} Tasks")
        plt.savefig(
            os.path.join(img, f"netmark-{tasks}-tasks-heatmap.png"), dpi=300, bbox_inches="tight"
        )
        plt.clf()

        # And distribution of times
        plt.figure(figsize=(12, 12))
        flat = df.to_numpy().flatten()
        sns.histplot(pandas.DataFrame(flat, columns=['time']), x="time")
        plt.title(f"MPI Connection Times for {tasks} Tasks")
        plt.savefig(
            os.path.join(img, f"netmark-{tasks}-tasks-hist.png"), dpi=300, bbox_inches="tight"
        )
        plt.clf()
        plt.close()

if __name__ == "__main__":
    main()
