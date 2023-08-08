#!/usr/bin/env python3

import argparse
import os
import time
import re
import json
from metricsoperator import MetricsOperator
import metricsoperator.utils as utils
import seaborn as sns
import matplotlib.pyplot as plt
import pandas

here = os.path.abspath(os.path.dirname(__file__))
examples = os.path.dirname(os.path.dirname(here))
experiments = os.path.join(examples, "experiments")
metrics_yaml = os.path.join(experiments, "flux-operator", "metrics.yaml")

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

    pods = {}

    # We actually have multiple containers here!
    for i, output in enumerate(m.watch()):
        print(json.dumps(output, indent=4))
        pods[i] = output

    # This output structure is a little different
    # it is combined single metrics.json into one
    utils.write_json(pods, args.out)
    plot_results(pods)

def plot_results(ranks):
    """
    Plot results to a histogram and matrix heatmap
    """
    # Directory for plotting results
    img = os.path.join(here, "img")
    if not os.path.exists(img):
        os.makedirs(img)

    # Assemble a data frame with all the result types
    # AND in addition we want to group based on the pod (rank)
    matrix = []
    columns = set()
    for rank, output in ranks.items():
        for i, timepoint in enumerate(output['data']):        

            for metric_type, listing in timepoint.items():

                # Assemble a row for the timepoint based on name
                row = {}

                # This is a list with different child names
                for item in listing:
                    cmd = item.get('command') or "global"
                    cmd = cmd.replace('|__', '')
                    for metric_name in item:
                        metric_value = item[metric_name]

                        # The name of the column (specific to command)
                        column = f"{metric_type}_{metric_name}"
                    
                        # Damn, some of the parsed values come back as strings...
                        try:
                            metric_value = float(metric_value)
                        except:
                            continue

                        # Things we can plot
                        if isinstance(metric_value, (int, float)):
                            columns.add(column)
                            row[column] = metric_value

                if row:
                    row['timepoint'] = i
                    row['rank'] = rank
                    row['process'] = cmd
                    matrix.append(row)

    df = pandas.DataFrame(columns=list(columns))
    idx = 0
    for row in matrix:
        df.loc[idx, list(row.keys())] = list(row.values())
        idx +=1

    # Save data frame to csv
    df.to_csv(os.path.join(img, 'metrics.csv'))

    # Break into dataframe with just process names and global
    global_df = df[df.process=="global"]
    global_df = global_df.dropna(axis=1, how='all')
    process_df = df[df.process!="global"]

    plot_df(global_df, "global", img)
    plot_df(process_df, "process", img, "process")

def plot_df(df, identifier, img, hue="rank"):
    """
    Make plots for a specific dataframe
    """
    # Only plot those that are defined
    for column in df.columns:
        # Don't plot the pid / uid / gid
        if re.search('(_pid|_uid)', column):
            continue
        if column in ["timepoint", "rank"]:
            continue
        if df[column].sum() == 0:
            print(f"Metric {column} is all zeros, skipping.")
            continue

        # Save to data file
        title = f"Result for {column} in {identifier} space"
        slug = column.replace(' ', '-').replace('/', '-')
        
        # for sty in plt.style.available:
        plt.figure(figsize=(12, 12))
        ax = sns.lineplot(data=df, x="timepoint", hue=hue, y=column, markers=True, dashes=True)
        plt.title(title)
        ax.set_xlabel("timepoint", fontsize=16)
        ax.set_ylabel(column, fontsize=16)
        ax.set_xticklabels(ax.get_xmajorticklabels(), fontsize=14)
        ax.set_yticklabels(ax.get_yticks(), fontsize=14)
        plt.tight_layout()
        plt.savefig(os.path.join(img, f"{identifier}-{slug}.png"))
        plt.clf()

        # And distribution
        plt.figure(figsize=(12, 12))
        sns.histplot(df, x=column, hue=hue, bins=20)
        plt.title(f"Histogram for result {column} across ranks")
        plt.savefig(
            os.path.join(img, f"{identifier}-{slug}-hist.png"), dpi=300, bbox_inches="tight"
        )
        plt.clf()
        plt.close()


if __name__ == "__main__":
    main()
