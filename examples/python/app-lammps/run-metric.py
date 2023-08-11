#!/usr/bin/env python3

import argparse
import os
import json
import time
from metricsoperator import MetricsOperator
import metricsoperator.utils as utils
import seaborn as sns
import matplotlib.pyplot as plt
import pandas
import collections

here = os.path.abspath(os.path.dirname(__file__))
examples = os.path.dirname(os.path.dirname(here))
tests = os.path.join(examples, "tests")
metrics_yaml = os.path.join(tests, "app-lammps", "metrics.yaml")

plt.style.use("bmh")


def get_parser():
    parser = argparse.ArgumentParser(
        description="Run LAMMPS Metric and Get Output",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "--out",
        help="json file to save results",
        default=os.path.join(here, "metrics.json"),
    )
    parser.add_argument(
        "--iter",
        help="number of iterations to run",
        type=int,
        default=5,
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

    # Save listing of results
    results = []

    # Create a metrics operator with our metrics.yaml
    m = MetricsOperator(metrics_yaml)
    for i in range(args.iter):
        print(f"Running LAMMPS iteration {i}")
        m.create()
        if i == 0:
            print(f"Sleeping {args.sleep} seconds so container can pull...")
            time.sleep(args.sleep)
        else:
            time.sleep(5)
        for output in m.watch():
            print(json.dumps(output, indent=4))
            results.append(output)
        m.delete()
    utils.write_json(results, args.out)
    plot_results(results)


def plot_results(output):
    """
    Plot results to a histogram and matrix heatmap
    """
    # Directory for plotting results
    img = os.path.join(here, "img")
    if not os.path.exists(img):
        os.makedirs(img)

    # Parse into data frame
    df = pandas.DataFrame(columns=["ranks", "iter", "time"])

    idx = 0
    for iter, data in enumerate(output):
        loop_time = data["data"][0]["loop_time"]
        ranks = data["data"][0]["ranks"]
        df.loc[idx, :] = [ranks, iter, loop_time]
        idx += 1

    df.to_csv(os.path.join(img, "lammps-times.csv"))

    # Plot each!
    colors = sns.color_palette("hls", 8)
    hexcolors = colors.as_hex()
    types = list(df.ranks.unique())

    # ALWAYS double check this ordering, this
    # is almost always wrong and the colors are messed up
    palette = collections.OrderedDict()
    for t in types:
        palette[t] = hexcolors.pop(0)

    make_plot(
        df,
        title="LAMMPS Times",
        tag="lammps",
        ydimension="time",
        xdimension="ranks",
        palette=palette,
        outdir=img,
        ext="png",
        plotname="lammps",
        hue="ranks",
        plot_type="bar",
        xlabel="MPI Ranks",
        ylabel="Time (seconds)",
    )


def make_plot(
    df,
    title,
    tag,
    ydimension,
    xdimension,
    palette,
    xlabel,
    ylabel,
    ext="pdf",
    plotname="lammps",
    plot_type="violin",
    hue="ranks",
    outdir="img",
):
    """
    Helper function to make common plots.
    """
    plotfunc = sns.boxplot
    if plot_type == "violin":
        plotfunc = sns.violinplot

    ext = ext.strip(".")
    plt.figure(figsize=(12, 12))
    sns.set_style("dark")
    ax = plotfunc(
        x=xdimension, y=ydimension, hue=hue, data=df, whis=[5, 95], palette=palette
    )
    plt.title(title)
    ax.set_xlabel(xlabel, fontsize=16)
    ax.set_ylabel(ylabel, fontsize=16)
    ax.set_xticklabels(ax.get_xmajorticklabels(), fontsize=14)
    ax.set_yticklabels(ax.get_yticks(), fontsize=14)
    plt.savefig(os.path.join(outdir, f"{tag}_{plotname}.{ext}"))
    plt.clf()


if __name__ == "__main__":
    main()
