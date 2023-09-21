# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

import json
import re

from .base import MetricBase

# Known section headers
headers = [
    "CPU STATISTICS TASK",
    "CPU STATISTICS CHILD",
    "IO STATISTICS",
    "POLICY",
    "PAGEFAULTS TASK",
    "PAGEFAULTS CHILD",
    "STACK UTILIZATION",
    "THREADS TASK",
    "THREADS CHILD",
    "KERNEL TABLES",
    "TASK SWITCHING",
]


header_regex = "(%s)" % "|".join(headers)


class perf_sysstat(MetricBase):
    container_name = "perf-sysstat"

    @property
    def pod_prefix(self):
        return f"{self.name}-m-0"

    def parse_log(self, lines):
        """
        Given lines of output, parse and return json
        """
        # Get the log metadata
        metadata = self.get_log_metadata(lines)

        # Get command being monitored
        command = lines.split("PIDSTAT COMMAND START", 1)[-1]
        command = command.split("PIDSTAT COMMAND END", 1)[0].strip()

        # Split lines by section separator
        results = []
        sections = self.get_log_sections(lines)
        for section in sections:
            if not section.strip():
                continue

            # These will be parts for one timepoint
            parts = section.strip().split("\n")
            timepoint = {}
            while parts:
                part = parts.pop(0)
                title = part.replace(" ", "_").lower().strip()
                if re.search(header_regex, part) and parts:
                    jsondata = parts.pop(0)
                    try:
                        data = json.loads(jsondata)
                    except Exception:
                        print(f"Issue parsing {part}")
                        continue
                    if data:
                        timepoint[title] = data

            # Only add timepoint if we collected data
            if timepoint:
                results.append(timepoint)

        return {
            "data": results,
            "metadata": metadata,
            "command": command,
            "spec": self.spec,
        }
