# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

# SPDX-License-Identifier: MIT

import os

from kubernetes import client, config

import metricsoperator.metrics as mutils
import metricsoperator.utils as utils

config.load_kube_config()


class MetricsOperator:
    def __init__(self, yaml_file):
        """
        Given a YAML file with one or more metrics, apply
        to create it and stream logs for each metric of interest.
        """
        self._core_v1 = None
        self.yaml_file = os.path.abspath(yaml_file)
        self.spec = utils.read_yaml(self.yaml_file)

    def watch(self):
        """
        Wait for (and yield parsed) metric logs.
        """
        for metric in self.spec["spec"]["metrics"]:
            parser = mutils.get_metric(metric["name"])(self.spec)
            print("Watching %s" % metric["name"])
            for pod, container in parser.logging_containers():
                yield parser.parse(pod=pod, container=container)

    def create(self):
        """
        Create the associated YAML file.
        """
        api = client.CustomObjectsApi()
        group, version = self.spec["apiVersion"].split("/", 2)
        plural = self.spec["kind"].lower() + "s"
        namespace = self.spec["metadata"].get("namespace") or "default"
        return api.create_namespaced_custom_object(
            group=group,
            version=version,
            namespace=namespace,
            plural=plural,
            body=self.spec,
        )
