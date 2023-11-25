# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

import re

from metricsoperator.metrics.base import MetricBase


class mpitrace(MetricBase):
    """
    Parse mpitrace (raw text) logs

    MPI trace doesn't stream to a log, so instead we take the raw lines
    for a single log. Each log should be one process.
    """

    def parse_log(self, lines):
        """
        Given lines of output, parse and return json
        """
        result = {}

        # change avg. bytes average_bytes, and others more readable
        lines = lines.replace("avg. bytes", "average_bytes")
        lines = lines.replace("#calls", "number_calls")
        lines = lines.replace("time(sec)", "time_seconds")
        lines = lines.strip().split("\n")
        while lines:
            line = lines.pop(0)

            # The MPI rank we are parsing
            if line.startswith("Data for MPI rank"):
                _, rank, _, total_ranks = line.strip(":").rsplit(" ", 3)
                result["rank"] = int(rank)
                result["total_ranks"] = int(total_ranks)
                continue

            if line.startswith("MPI Routine"):
                result["times_mpi_init_to_mpi_finalize"] = parse_init_finalize_table(
                    lines
                )
                # The next line is a description about communication time
                result["communication_time_description"] = lines.pop(0)

                # Next lines are key value pairs
                line = lines.pop(0)
                while "=" in line:
                    key, value = line.split("=", 1)
                    result[key.strip()] = value.strip()
                    line = lines.pop(0)
                    continue

            # Message size distribution table
            if line.startswith("Message size distributions:"):
                result["message_size_distributions"] = parse_message_sizes(lines)
                continue

            if line.startswith("Summary for all tasks"):
                result["summary_tasks"] = parse_summary(lines)
                continue

            if line.startswith("MPI timing summary"):
                result["mpi_timing_summary"] = parse_mpi_timing_summary(lines)

        return result


def parse_mpi_timing_summary(lines):
    line = lines.pop(0)
    columns = re.split("\\s+", line.strip())
    rows = [columns]
    while lines:
        line = lines.pop(0)
        if not line:
            continue
        row = re.split("\\s+", line.strip())
        row = [
            int(row[0]),
            row[1],
            int(row[2]),
            float(row[3]),
            float(row[4]),
            float(row[5]),
            float(row[6]),
            float(row[7]),
            int(row[8]),
        ]
        rows.append(row)
    return rows


def parse_summary(lines):
    summary = {}
    count = 0
    while True and count < 3:
        line = lines.pop(0)
        if not line:
            count += 1
            continue
        sep = ":" if ":" in line else "="
        key, value = line.split(sep)
        summary[key.strip()] = value.strip()

    return summary


def parse_message_sizes(lines):
    """
    Parse the times from the Message size Distributions tables
    """
    # pop top line that is empty
    lines.pop(0)

    tables = {}

    # columns should not change
    rows = []
    routine = None
    while True and lines:
        line = lines.pop(0)

        # Last line of multiple sections
        if "----" in line:
            break

        # this is a row of the same table
        if line.startswith(" "):
            row = re.split("\\s+", line.strip())
            rows.append([int(row[0]), float(row[1]), float(row[2])])

        # This is the row of a new table
        else:
            if routine:
                tables[routine] = rows
            rows = []
            # This is the header
            row = re.split("\\s+", line.strip())
            routine, row = row[0], row[1:]
            rows.append(row)

    return tables


def parse_init_finalize_table(lines):
    """
    Parse the Times from MPI_Init() to MPI_Finalize() table
    """
    # pop top line that is a divider
    lines.pop(0)

    # columns should not change
    columns = ["mpi_routine", "number_calls", "average_bytes", "time_seconds"]
    rows = [columns]

    while True:
        line = lines.pop(0)
        # Last line
        if "----" in line:
            break
        row = re.split("\\s+", line)
        row = [row[0], int(row[1]), float(row[2]), float(row[3])]
        rows.append(row)
    return rows
