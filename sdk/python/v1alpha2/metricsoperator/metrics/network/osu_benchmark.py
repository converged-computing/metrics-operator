# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

import os
import re

from metricsoperator.metrics.base import MetricBase

# prepare more consistent / formatted columns
latency_size_header = ["Size", "Latency(us)"]
average_latency_header = ["Size", "Avg Latency(us)"]
bandwidth_size_header = ["Size", "Bandwidth (MB/s)"]
column_lookup = {
    "osu_ibarrier": ["Overall(us)", "Compute(us)", "Pure Comm.(us)", "Overlap(%)"],
    "osu_barrier": ["Avg Latency(us)"],
    "osu_mbw_mr": ["Size", "MB/s", "Messages/s"],
    "osu_init": ["nprocs", "min", "max", "avg"],
    "osu_multi_lat": latency_size_header,
    "osu_get_acc_latency": latency_size_header,
    "osu_latency": latency_size_header,
    "osu_fop_latency": latency_size_header,
    "osu_put_latency": latency_size_header,
    "osu_get_latency": latency_size_header,
    "osu_acc_latency": latency_size_header,
    "osu_mt_latency": latency_size_header,
    "osu_get_latency": latency_size_header,
    "osu_cas_latency": latency_size_header,
    "osu_latency_mp": latency_size_header,
    "osu_latency_mt": latency_size_header,
    "osu_bibw": bandwidth_size_header,
    "osu_get_bw": bandwidth_size_header,
    "osu_put_bw": bandwidth_size_header,
    "osu_put_bibw": bandwidth_size_header,
    "osu_bw": bandwidth_size_header,
    "osu_allgather": average_latency_header,
    "osu_allreduce": average_latency_header,
}


def parse_commented_header(section):
    """
    Parse a commented (#) header
    """
    # Each section has some number of header lines (with #)
    header = []
    while section and section[0].startswith("#"):
        header.append(section.pop(0))
    return header


def parse_row(row):
    """
    Given a row of two values with spaces, parse.
    """
    row = row.split(" ", 2)
    return [x.strip() for x in row if x]


def parse_multi_value_row(row):
    """
    Parse a row with multiple values.
    """
    row = row.split(" ")
    return [x.strip() for x in row if x]


def parse_hello_section(section):
    """
    Parse the osu hello section

    This only has a print line for output
    """
    # Command is the first entry
    command = section.pop(0)
    header = parse_commented_header(section)

    # The next row is just a print line
    message = section.pop(0)
    timed = parse_timed_section(section)
    result = {
        "matrix": [[message]],
        "columns": ["message"],
        "header": header,
        "command": command,
    }
    if timed:
        result["timed"] = timed
    return result


def parse_init_section(section):
    """
    Parse the osu init section

    This section has one column, and all the values there!
    """
    # Command is the first entry
    command = section.pop(0)
    header = parse_commented_header(section)

    # The next row has all the data!
    row = section.pop(0)
    values = [x.strip() for x in row.split(",")]
    data = {}
    for entry in values:
        field, value = [x.strip() for x in entry.split(":")]
        # Do we have a unit?
        unit = ""
        if " " in value:
            value, unit = [x.strip() for x in value.split(" ")]
        if unit:
            field = f"{field}-{unit}"
        data[field] = float(value)

    # If we have additional sections (possibly with times)
    timed = parse_timed_section(section)
    result = {
        "matrix": [list(data.values())],
        "columns": list(data.keys()),
        "header": header,
        "command": command,
    }
    if timed:
        result["timed"] = timed
    return result


def parse_timed_section(section):
    """
    If the remainder is wrapped in time, parse it.
    """
    timed = {}
    for line in section:
        if line and re.search("^(real|user|sys\t)", line):
            time_type, ts = line.strip().split("\t")
            timed[time_type] = ts
    return timed


def parse_barrier_section(section):
    """
    Parse a barrier section.

    This section is unique in that it has two columns
    but the header is not preceded with a #
    """
    # Command is the first entry
    command = section.pop(0)
    header = parse_commented_header(section)

    # The columns are the last row of the header
    section.pop(0)
    result = parse_value_matrix(section)
    result.update({"header": header, "command": command})
    return result


def parse_multi_section(section):
    """
    Parse a multi-value section.

    This section has standard format, but >2 values in the matrix
    """
    # Command is the first entry
    command = section.pop(0)
    header = parse_commented_header(section)

    # The columns are the last row of the header
    header.pop()
    result = parse_value_matrix(section)
    result.update({"header": header, "command": command})
    return result


def parse_value_matrix(section):
    """
    Parse a matrix of times
    """
    # The remainder is data, again always two points
    # If there are real / user / sys at the end, we ran with timed:true
    data = []
    timed = {}
    for line in section:
        if not line:
            continue

        # We found a time
        if re.search("^(real|user|sys\t)", line):
            time_type, ts = line.strip().split("\t")
            timed[time_type] = ts
            continue

        row = parse_multi_value_row(line)
        row = [float(x) for x in row]
        data.append(row)

    result = {"matrix": data}
    if timed:
        result["timed"] = timed
    return result


def run_parsing_function(section):
    """
    Parsing functions for different sections
    """
    # The command is the first line
    command = os.path.basename(section[0])
    result = None
    if command in ["osu_ibarrier"]:
        result = parse_barrier_section(section)
    elif command in [
        "osu_bw",
        "osu_bibw",
        "osu_barrier",
        "osu_get_bw",
        "osu_put_bw",
        "osu_put_bibw",
        "osu_mbw_mr",
        "osu_multi_lat",
        "osu_allgather",
        "osu_latency",
        "osu_cas_latency",
        "osu_put_latency",
        "osu_get_latency",
        "osu_latency_mp",
        "osu_latency_mt",
        "osu_fop_latency",
        "osu_acc_latency",
        "osu_get_acc_latency",
        "osu_allreduce",
    ]:
        result = parse_multi_section(section)

    # Special snowflakes
    elif command == "osu_init":
        result = parse_init_section(section)

    # This only potentially has a time :)
    elif command == "osu_hello":
        result = parse_hello_section(section)

    # Columns aren't predictible, so we ensure they are more consistent this way
    # Some parsers do their own columns
    if "columns" not in result:
        result["columns"] = column_lookup[command]

    return result


class network_osu_benchmark(MetricBase):
    """
    Parse the OSU benchmarks output into data!

    For pair to pair we had a common format (two values with size and another field)
    but adding on the multi- benchmarks, we have a new challenge that there is slight
    variance in format, so I needed to extend the class to be more specific to parsing.
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
            section = [x.strip() for x in section if x.strip()]

            # Parse the section. If this fails, we want to know
            datum = run_parsing_function(section)
            results.append(datum)

        return {"data": results, "metadata": metadata, "spec": self.spec}
