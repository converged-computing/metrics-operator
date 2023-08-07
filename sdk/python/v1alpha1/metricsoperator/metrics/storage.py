# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

import json

from .base import MetricBase


class io_sysstat(MetricBase):
    container_name = "io-sysstat"

    @property
    def pod_prefix(self):
        return f"{self.name}-m-0"

    def parse_log(self, lines):
        """
        Given lines of output, parse and return json
        """
        # Get the log metadata
        metadata = self.get_log_metadata(lines)

        # Split lines by IOSTAT TIMEPOINT
        results = []
        sections = self.get_log_sections(lines)
        for section in sections:
            if not section.strip():
                continue
            results.append(json.loads(section))
        return {"data": results, "metadata": metadata, "spec": self.spec}
