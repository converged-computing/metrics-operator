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
        description="Run PID Stat Metric and Get Output",
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

    # Ensure we cleanup!
    m.delete()


def plot_results(output):
    """
    Plot results to a histogram and matrix heatmap
    """
    # Directory for plotting results
    img = os.path.join(here, "img")
    if not os.path.exists(img):
        os.makedirs(img)
   
    # Assemble a data frame with all the result types
    columns = set()
    matrix = []
    for i, timepoint in enumerate(output['data']):
        
        # Assemble a row for the timepoint based on name
        row = {}
        for metric_type, listing in timepoint.items():

            # This is always a single dictionary
            for metric_name in listing[0]:
                metric_value = listing[0][metric_name]

                # The name of the column
                column = f"{metric_type}_{metric_name}"
                
                # Things we can plot
                if isinstance(metric_value, (int, float)):
                    columns.add(column)
                    row[column] = metric_value
            if row:
                row['timepoint'] = i
                matrix.append(row)

    df = pandas.DataFrame(columns=list(columns))
    idx = 0
    for row in matrix:
        df.loc[idx, list(row.keys())] = list(row.values())
        idx +=1

    # Save data frame to csv
    df.to_csv(os.path.join(img, 'metrics.csv'))

    # Only plot those that are defined
    for column in df.columns:
        # Don't plot the pid!
        if column.endswith('pid'):
            continue
        if column == "timepoint":
            continue
        if df[column].sum() == 0:
            print(f"Metric {column} is all zeros, skipping.")
            continue

        # Save to data file
        title = f"Result for {column}"
        slug = column.replace(' ', '-')
        
        # for sty in plt.style.available:
        ax = sns.lineplot(data=df, x="timepoint", y=column, markers=True, dashes=True)
        plt.title(title)
        ax.set_xlabel("timepoint", fontsize=16)
        ax.set_ylabel(column, fontsize=16)
        ax.set_xticklabels(ax.get_xmajorticklabels(), fontsize=14)
        ax.set_yticklabels(ax.get_yticks(), fontsize=14)
        plt.tight_layout()
        plt.savefig(os.path.join(img, f"{slug}.png"))
        plt.clf()

        # And distribution
        sns.histplot(df, x=column)
        plt.title(f"Histogram for result {column}")
        plt.savefig(
            os.path.join(img, f"{slug}-hist.png"), dpi=300, bbox_inches="tight"
        )
        plt.clf()


if __name__ == "__main__":
    main()
