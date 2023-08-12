# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)


import re

import metricsoperator.utils as utils
from metricsoperator.metrics.base import MetricBase


def get_first_int(match):
    return int(match.group().strip().split(" ")[0])


class app_lammps(MetricBase):
    """
    Parse LAMMPS into nice JSON :)
    """

    container_name = "launcher"

    @property
    def pod_prefix(self):
        return f"{self.name}-l-0"

    def parse_row(self, row):
        """
        Given a row of two values with spaces, parse.
        """
        row = row.split(" ", 1)
        return [x.strip() for x in row]

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
            results.append(parse_lammps(section))
        return {"data": results, "metadata": metadata, "spec": self.spec}


def parse_lammps(lines):
    """
    Parse a raw lammps log and return data
    """
    if not isinstance(lines, list):
        lines = lines.split("\n")

    entry = {}
    while lines:
        line = lines.pop(0)

        # Figure out ranks and threads (in same line)
        match = re.search("[0-9]+ MPI tasks", line)
        if match:
            entry["ranks"] = get_first_int(match)
            continue

        match = re.search("[0-9]+ OpenMP threads", line)
        if match:
            entry["threads"] = get_first_int(match)

        match = re.search("(?P<percentage>[0-9]+[.][0-9]+)[%] CPU use with", line)
        if match:
            entry["percentage_cpu"] = float(match.groupdict()["percentage"])
            continue

        if re.search("^(NLocal|Nghost|Neighs)", line):
            category = line.split(":")[0].lower()
            line = line.split(":", 1)[-1].strip()
            ave, line = line.split("ave", 1)
            maxval, line = line.split("max", 1)
            minval, line = line.split("min", 1)
            entry.update(
                {
                    f"{category}_avg": float(ave),
                    f"{category}_min": float(minval),
                    f"{category}_max": float(maxval),
                }
            )
            hist = lines.pop(0).split(":")[-1].strip().split(" ")
            entry[f"{category}_hist"] = [int(x) for x in hist]
            continue

        if line.startswith("Total # of neighbors"):
            entry["neighbors"] = float(line.split("=")[-1].strip())
            continue

        if line.startswith("Ave neighs/atom"):
            entry["average_neighbors_per_atom"] = float(line.split("=")[-1].strip())
            continue

        if line.startswith("Neighbor list builds"):
            entry["neighbor_list_builds"] = int(line.split("=")[-1].strip())
            continue

        # This is total wall time as reported by lammps
        if line.startswith("Total wall time"):
            rawtime = line.split(":", 1)[-1].strip()
            entry["total_wall_time_raw"] = rawtime
            entry["total_wall_time_seconds"] = utils.timestr2seconds(rawtime)
            continue

        # Dimension of molecular matrix maybe?
        match = re.search(
            "(?P<x>[0-9]+) by (?P<y>[0-9]+) by (?P<z>[0-9]+) MPI processor grid",
            line,
        )

        if match:
            if "processor_grids" not in entry:
                entry["processor_grids"] = []
            new_grid = [{k: int(v)} for k, v in match.groupdict().items()]
            entry["processor_grids"].append(new_grid)
            continue

        # Number of atoms / velocities
        match = re.search("[0-9]+ atoms", line)
        if match and "atoms" not in entry:
            entry["atoms"] = get_first_int(match)
            continue
        match = re.search("[0-9]+ velocities", line)
        if match:
            entry["velocities"] = get_first_int(match)
            continue

        # reading data from CPU
        match = re.search("read_data CPU = (?P<cpu>[0-9]+[.][0-9]+) seconds", line)
        if match:
            entry["read_data_cpu_seconds"] = float(match.groupdict()["cpu"])
            continue

        match = re.search(
            "bounding box extra memory = (?P<mem>[0-9]+[.][0-9]+) MB", line
        )
        if match:
            entry["bounding_box_extra_memory_mb"] = float(match.groupdict()["mem"])
            continue

        match = re.search("replicate CPU = (?P<seconds>[0-9]+[.][0-9]+) seconds", line)
        if match:
            entry["replicate_cpu_seconds"] = float(match.groupdict()["seconds"])
            continue

        if line.startswith("Unit style"):
            entry["unit_style"] = line.split(":")[-1].strip()
            continue

        if line.startswith("Time step"):
            entry["time_step"] = float(line.split(":")[-1].strip())
            continue

        match = re.search(
            "Per MPI rank memory allocation (min/avg/max) = (?P<min_mpi_rank_memory_allocation_mb>[0-9]+[.][0-9]+) [|] (?P<avg_mpi_rank_memory_allocation_mb>[0-9]+[.][0-9]+) [|] (?P<max_mpi_rank_memory_allocation_mb>[0-9]+[.][0-9]+) Mbytes",
            line,
        )
        if match:
            [entry.update({k: int(v)}) for k, v in match.groupdict().items()]
            continue

        # If we find the embedded table with steps
        if line.startswith("Step"):
            header = [x.strip() for x in line.split(" ") if x.strip()]
            matrix = []
            line = lines.pop(0)
            while not line.startswith("Loop"):
                matrix.append([float(x.strip()) for x in line.split(" ") if x.strip()])
                line = lines.pop(0)
            entry["steps"] = {"matrix": matrix, "columns": header}

            # Here we have the loop line
            match = re.search(
                "Loop time of (?P<loop_time>[0-9]+[.][0-9]+) on (?P<loop_procs>[0-9]+) procs for (?P<loop_steps>[0-9]+) steps with (?P<loop_atoms>[0-9]+) atoms",
                line,
            )
            if match:
                values = match.groupdict()
                entry.update(
                    {
                        "loop_time": float(values["loop_time"]),
                        "loop_procs": int(values["loop_procs"]),
                        "loop_atoms": int(values["loop_atoms"]),
                    }
                )
                continue

        match = re.search(
            "Performance: (?P<performance_ns_per_day>[0-9]+[.][0-9]+) ns/day, (?P<performance_hours_per_ns>[0-9]+[.][0-9]+) hours/ns, (?P<timesteps_per_second>[0-9]+[.][0-9]+) timesteps/s",
            line,
        )
        if match:
            [entry.update({k: float(v)}) for k, v in match.groupdict().items()]
            continue

        # If we find the embedded table with times
        if line.startswith("Section"):
            header = [x.strip() for x in line.split("|") if x.strip()]

            # Line with full ----------------------
            lines.pop(0)
            matrix = []
            columns = header[1:]
            line = lines.pop(0)
            while line and line.strip():
                parts = [x.strip() for x in line.split("|")]
                _, rest = parts[0], parts[1:]
                if rest:
                    matrix.append(rest)
                if not lines:
                    break
                line = lines.pop(0)

            entry["times"] = {"matrix": matrix, "header": columns}

    return entry
