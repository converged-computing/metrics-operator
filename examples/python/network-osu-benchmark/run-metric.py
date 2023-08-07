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
        description="Run Storage Metric and Get Output",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument(
        "--out",
        help="json file to save results",
        default=os.path.join(here, "metrics.json"),
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

    # Directory for plotting results
    img = os.path.join(here, "img")
    if not os.path.exists(img):
        os.makedirs(img)

    print('Sleeping a minute so container can pull...')
    time.sleep(60)
    for output in m.watch():
        print(json.dumps(output, indent=4))
        utils.write_json(output, args.out)

    # Plot each metric to image. We only have one output above so
    # this outside loop thing is OK to do :)
    # This could be adjusted to handle more than one run of the metric
    # in different environments, but we just have one for now!
    for result in output['data']:
        df = pandas.DataFrame(columns=result['columns'])
        idx = 0
        for datum in result['matrix']:
            df.loc[idx, :] = datum
            idx +=1
        
        # Separate x and y - latency (y) is a function of size (x)
        x = result['columns'][0]
        y = result['columns'][1]

        # Save to data file
        title = result['header'][0].replace('#', '').strip()
        slug = title.replace(' ', '-')
        df.to_csv(os.path.join(img, f"{slug}.csv"))

        # for sty in plt.style.available:
        ax = sns.lineplot(data=df, x=x, y=y, markers=True, dashes=True)
        plt.title(title)
        ax.set_xlabel(x + " logscale", fontsize=16)
        ax.set_ylabel(y + " logscale", fontsize=16)
        ax.set_xticklabels(ax.get_xmajorticklabels(), fontsize=14)
        ax.set_yticklabels(ax.get_yticks(), fontsize=14)
        plt.xscale('log')
        plt.yscale('log')
        plt.tight_layout()
        plt.savefig(os.path.join(img, f"{slug}.png"))
        plt.clf()

if __name__ == "__main__":
    main()
