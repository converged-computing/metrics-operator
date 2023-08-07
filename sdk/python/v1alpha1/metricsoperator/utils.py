# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

import json

import yaml


def read_yaml(filename):
    """
    Read yaml into dict
    """
    with open(filename, "r") as file:
        configuration = yaml.safe_load(file)
    return configuration


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
