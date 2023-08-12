# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

import json
import os
import re

import yaml


def read_file(filename):
    """
    Read file into text blob
    """
    with open(filename, "r") as fd:
        text = fd.read()
    return text


def read_yaml(filename):
    """
    Read yaml into dict
    """
    with open(filename, "r") as file:
        configuration = yaml.safe_load(file)
    return configuration


def slugify(section):
    """
    Slugify into lowercase and underscores
    """
    return section.lower().replace(" ", "_").replace("-", "_")


def write_file(content, filename):
    """
    Write dict to json
    """
    with open(filename, "w") as fd:
        fd.write(content)


def write_json(obj, filename):
    """
    Write dict to json
    """
    with open(filename, "w") as fd:
        fd.write(json.dumps(obj, indent=4))


def read_json(filename):
    """
    Read json into dict
    """
    with open(filename, "r") as fd:
        content = json.loads(fd.read())
    return content


def recursive_find(base, pattern="^(laamps[.]out|lammps.*[.]out|log[.]out)$"):
    """
    Recursively find lammps output files.
    """
    for root, _, filenames in os.walk(base):
        for filename in filenames:
            if re.search(pattern, filename):
                yield os.path.join(root, filename)


def read_lines(filename):
    """
    Read lines of a file into a list.
    """
    with open(filename, "r") as fd:
        lines = fd.readlines()
    return lines


def timestr2seconds(timestr):
    """
    Given a timestring in two formats, return seconds (float).
    """
    # Minutes and seconds, MM:SS.mm
    if timestr.count(":") == 1:
        minutes, seconds = timestr.split(":")
        return (int(minutes) * 60) + float(seconds)

    # hours, minutes, seconds HH:MM:SS.mm
    elif timestr.count(":") == 2:
        hours, minutes, seconds = timestr.split(":")
        return (int(hours) * 360) + (int(minutes) * 60) + float(seconds)
    raise ValueError(f"Unrecognized time format {timestr}")
