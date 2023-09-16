# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

from metricsoperator.metrics.base import MetricBase


class network_netmark(MetricBase):
    """
    Parse the netmark terminal output.
    """

    container_name = "launcher"

    @property
    def pod_prefix(self):
        return f"{self.name}-n-0"

    def parse_row(self, row):
        """
        Given a row of netmark output, parse!
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

        # Netmark just has one long run (one section)
        for section in sections:
            section = section.split("\n")
            section = [x.strip() for x in section if x.strip()]

            # The first lines up to SETUP have information about ranks
            ranks = []
            while "SETUP" not in section[0]:
                ranks.append(section.pop(0))

            # Next is setup section (we have this in our metadata)
            section.pop(0)
            setup = []
            while "======" not in section[0]:
                setup.append(section.pop(0))

            # Pop end of setup line
            section.pop(0)

            # Find the RTT.csv between netmark.start and end lines
            section = "\n".join(section)
            netmark_data = section.split("NETMARK RTT.CSV START", 1)[-1]
            netmark_data = netmark_data.split("NETMARK RTT.CSV END", 1)[0]

            # Add rest of ranks from original data
            section = section.split("\n")
            while section:
                if section[0].startswith("size"):
                    ranks.append(section.pop(0))
                else:
                    section.pop(0)

            datum = {
                "RTT.csv": netmark_data.strip(),
                "ranks": ranks,
                "setup": setup,
            }
            results.append(datum)

        return {"data": results, "metadata": metadata, "spec": self.spec}
