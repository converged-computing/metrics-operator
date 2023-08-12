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
metrics_yaml = os.path.join(tests, "app-amg", "metrics.yaml")

plt.style.use("bmh")


def get_parser():
    parser = argparse.ArgumentParser(
        description="Run AMG Metric and Get Output",
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
        print(f"Running AMG iteration {i}")
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


def split_seconds(t):
    """
    Split a time string into float seconds
    """
    assert "seconds" in t
    return float(t.split("seconds")[0])


def plot_results(output):
    """
    Plot results to a histogram and matrix heatmap
    """
    # Directory for plotting results
    img = os.path.join(here, "img")
    if not os.path.exists(img):
        os.makedirs(img)

    # Parse into data frame, stuff that might be important!
    df = pandas.DataFrame(
        columns=[
            "iter",
            "spatial_operator_seconds",
            "rhs_and_initial_guess_seconds",
            "problem_pcg_setup_seconds",
            "solve_time_seconds",
        ]
    )

    # I'm not actually sure what is important here
    idx = 0
    for iter, data in enumerate(output):
        res = data["data"][0]

        # Generate matrix time for spatial operator (seconds)
        spatial_op = split_seconds(
            res["generate_matrix"]["spatial_operator"]["wall_clock_time"]
        )

        # RHS and initial guess (seconds)
        guess = split_seconds(
            res["vector_setup"]["rhs_and_initial_guess"]["wall_clock_time"]
        )

        # problem setup - pcg setup (seconds)
        pcg_setup = split_seconds(res["problem_setup"]["pcg_setup"]["wall_clock_time"])

        # This is probably the one we care about? (seconds)
        solve_time = split_seconds(res["solve_time"]["pcg_solve"]["wall_clock_time"])

        df.loc[idx, :] = [iter, spatial_op, guess, pcg_setup, solve_time]
        idx += 1

    df.to_csv(os.path.join(img, "amg-times.csv"))

    # Make a plot for each time
    for ydim in df.columns:
        if ydim == "iter":
            continue

        # for sty in plt.style.available:
        ax = sns.boxplot(data=df, y=ydim)
        title = " ".join([x.capitalize() for x in ydim.split("_")])
        plt.title(title)
        ax.set_ylabel("time in seconds", fontsize=16)
        ax.set_xlabel(ydim.replace("_seconds", ""), fontsize=16)
        ax.set_xticklabels(ax.get_xmajorticklabels(), fontsize=14)
        ax.set_yticklabels(ax.get_yticks(), fontsize=14)

        # Make sure we round a bit so it's prettier
        ylabels = ["{:,.2f}".format(x) for x in ax.get_yticks()]
        ax.set_yticklabels(ylabels)
        plt.tight_layout()

        # Save to file
        plt.savefig(os.path.join(img, f"{ydim}.png"))
        plt.clf()
        plt.close()


if __name__ == "__main__":
    main()
