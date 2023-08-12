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

here = os.path.abspath(os.path.dirname(__file__))
examples = os.path.dirname(os.path.dirname(here))
tests = os.path.join(examples, "tests")
metrics_yaml = os.path.join(tests, "network-osu-benchmark", "metrics.yaml")

plt.style.use("bmh")


def get_parser():
    parser = argparse.ArgumentParser(
        description="Run OSU Benchmarks Metric and Get Output",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "--out",
        help="json file to save results",
        default=os.path.join(here, "metrics-multiple.json"),
    )
    parser.add_argument(
        "--sleep",
        help="seconds to sleep allowing for pull and run",
        type=int,
        default=60,
    )
    parser.add_argument(
        "--iter",
        help="number of iterations to run",
        type=int,
        default=5,
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
    Run multiple OSU benchmark jobs on a Kind cluster.
    """
    parser = get_parser()
    args, _ = parser.parse_known_args()
    img = os.path.join(here, "img", "multiple")
    if not os.path.exists(img):
        os.makedirs(img)

    # Save listing of results
    results = []

    # Create a metrics operator with our metrics.yaml
    m = MetricsOperator(metrics_yaml)
    for i in range(args.iter):
        print(f"Running OSU Benchmarks iteration {i}")
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
    plot_results(results, img)


def plot_results(results, img):
    """
    Plot result images to file
    """
    import IPython

    IPython.embed()

    # Create a data frame for each result type
    dfs = {}
    idxs = {}
    columns = {}
    for entry in results:
        for result in entry["data"]:
            title = result["header"][0].replace("#", "").strip()
            slug = title.replace(" ", "-")
            if slug not in dfs:
                dfs[slug] = pandas.DataFrame(columns=result["columns"])
                idxs[slug] = 0
                columns[slug] = {"x": result["columns"][0], "y": result["columns"][1]}
            for datum in result["matrix"]:
                dfs[slug].loc[idxs[slug], :] = datum
                idxs[slug] += 1

    # Save each completed data frame to file and plot!
    for slug, df in dfs.items():
        print(f"Preparing plot for {slug}")

        # Save to data file
        df.to_csv(os.path.join(img, f"{slug}.csv"))

        # Separate x and y - latency (y) is a function of size (x)
        x = columns[slug]["x"]
        y = columns[slug]["y"]

        # Save to data file
        title = slug.replace("-", " ")

        # for sty in plt.style.available:
        ax = sns.lineplot(
            data=df, x=x, y=y, markers=True, dashes=True, errorbar=("ci", 95)
        )
        plt.title(title)
        ax.set_xlabel(x + " logscale", fontsize=16)
        ax.set_ylabel(y + " logscale", fontsize=16)
        ax.set_xticklabels(ax.get_xmajorticklabels(), fontsize=14)
        ax.set_yticklabels(ax.get_yticks(), fontsize=14)
        plt.xscale("log")
        plt.yscale("log")
        plt.tight_layout()
        plt.savefig(os.path.join(img, f"{slug}.png"))
        plt.clf()
        plt.close()


if __name__ == "__main__":
    main()
