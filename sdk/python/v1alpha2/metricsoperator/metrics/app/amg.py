# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)


import re

import metricsoperator.utils as utils
from metricsoperator.metrics.base import MetricBase


class app_amg(MetricBase):
    """
    Parse AMG into nice JSON :)
    """

    container_name = "launcher"

    @property
    def pod_prefix(self):
        return f"{self.name}-l-0"

    def parse_log(self, lines):
        """
        Given lines of output, parse and return json
        """
        # Get the log metadata
        metadata = self.get_log_metadata(lines)

        # Split lines by separator
        results = []
        sections = self.get_log_sections(lines)
        for section in sections:
            if not section.strip():
                continue
            section = section.split("\n")
            # Don't lose indentation, meaningful!
            section = [x for x in section if x.strip()]
            results.append(parse_amg(section))
        return {"data": results, "metadata": metadata, "spec": self.spec}


def parse_key_value_pair(line, sep="="):
    """
    Parse a key value pair, ensuring it is slugified too.
    """
    key, value = line.split(sep, 1)
    key = utils.slugify(key.strip())
    if not key:
        return
    return key, value.strip()


def parse_params(lines):
    """
    Parse a section.

    This is a generic parsing strategy that works with the test case.
    If we need to create problem-specific parsers we can do that.
    """
    # The title here is the first section header
    params = {}
    section = None
    while lines:
        line = lines.pop(0)

        # if we find an = with a space prefix, add directly as a data attribute
        if "=" in line and line.startswith(" "):
            pair = parse_key_value_pair(line)
            if not pair:
                continue
            if section:
                params[section][pair[0]] = pair[1]
            else:
                params[pair[0]] = pair[1]

        # If we find an = and the line is on equal level, it's a standalone attribute
        # OR an indented parameter of the same section
        elif ("=" in line and not line.startswith(" ")) or (
            ":" in line and line.startswith("   ")
        ):
            pair = parse_key_value_pair(line)
            if not pair:
                continue
            print(f"Adding {pair} to section {section}")
            params[pair[0]] = pair[1]

        # If we find a colon, it's a section header
        elif ":" in line:
            parts = line.strip().split(":", 1)
            parts = [x.strip() for x in parts if x.strip()]

            # If we have one piece, it's a section name
            if len(parts) == 1:
                section = parts[0]

                # Slugify!
                section = utils.slugify(section)
                params[section] = {}
                print(f"Adding new section {section}")

            # If we have two pieces, it's a separate key/value pair
            elif len(parts) > 1:
                pair = parse_key_value_pair(line, sep=":")
                if not pair:
                    continue
                params[pair[0]] = pair[1]
                print(f"Adding {pair} to section {section}")

    return params


separator = "============================================="


def parse_amg(lines):
    """
    Parse the AMG section. This is notably different, but original credit for parsing this goes to:
    https://github.com/flux-framework/flux-k8s/blob/canopie22-artifacts/canopie22-artifacts/amg/process_amg.py
    with credit to Dan Milroy.
    """
    entry = {}

    # Prepare to split into subset
    def get_subset(lines):
        subset = []
        while lines and "======" not in lines[0]:
            subset.append(lines.pop(0))
        return subset

    while lines:
        line = lines.pop(0)

        # This is the first section, driver parameters
        if re.search("Running with these driver", line, re.IGNORECASE):
            # Get the subset up to the next section
            section = get_subset(lines)

            # Each subsection header has TITLE:
            entry["driver_params"] = parse_params(section)
            continue

        if re.search("Generate Matrix", line, re.IGNORECASE):
            lines.pop(0)  # get rid of extra header
            section = get_subset(lines)
            entry["generate_matrix"] = parse_params(section)
            continue

        if re.search("Vector Setup", line, re.IGNORECASE):
            lines.pop(0)  # get rid of extra header
            section = get_subset(lines)
            entry["vector_setup"] = parse_params(section)
            continue

        if re.search("Setup Time", line, re.IGNORECASE):
            lines.pop(0)  # get rid of extra header
            section = get_subset(lines)
            entry["problem_setup"] = parse_params(section)
            continue

        if re.search("Solve Time", line, re.IGNORECASE):
            lines.pop(0)  # get rid of extra header
            section = get_subset(lines)
            entry["solve_time"] = parse_params(section)
            continue

    return entry
