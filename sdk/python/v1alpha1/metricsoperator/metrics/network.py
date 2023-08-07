# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

from .base import MetricBase


class network_osu_benchmark(MetricBase):
    """
    Parse the OSU benchmarks output into data!

    We will eventually want to plot these via means from:
    https://github.com/ULHPC/tutorials/tree/devel/parallel/mpi/OSU_MicroBenchmarks/plots
    But likely I'll provide Python equivalents to do this.
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

    def parse(self, pod, container):
        """
        Parse osu benchmark output log and return json
        """
        lines = self.stream_output(
            name=pod.metadata.name, namespace=self.namespace, container=container.name
        )

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

            # Command is the first entry
            command = section.pop(0)

            # Each section has some number of header lines (with #)
            header = []
            while section[0].startswith("#"):
                header.append(section.pop(0))

            # Last row of the header are the column names
            columns = header.pop()
            columns = columns.replace("# ", "").strip()
            columns = self.parse_row(columns)

            # The remainder is data, again always two points
            data = []
            for line in section:
                if not line:
                    continue
                row = self.parse_row(line)
                row = [float(x) for x in row]
                data.append(row)

            datum = {
                "matrix": data,
                "columns": columns,
                "header": header,
                "command": command,
            }
            results.append(datum)

        return {"data": results, "metadata": metadata, "spec": self.spec}
