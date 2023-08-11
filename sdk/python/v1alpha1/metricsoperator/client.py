# Copyright 2023 Lawrence Livermore National Security, LLC
# (c.f. AUTHORS, NOTICE.LLNS, COPYING)

# SPDX-License-Identifier: MIT

import os

from kubernetes import client, config

import metricsoperator.metrics as mutils
import metricsoperator.utils as utils


class MetricsOperator:
    def __init__(self, yaml_file):
        """
        Given a YAML file with one or more metrics, apply
        to create it and stream logs for each metric of interest.
        """
        self._core_v1 = None
        self.yaml_file = os.path.abspath(yaml_file)
        self.spec = utils.read_yaml(self.yaml_file)
        config.load_kube_config()

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
        return api.create_namespaced_custom_object(
            group=self.group,
            version=self.version,
            namespace=self.namespace,
            plural=self.plural,
            body=self.spec,
        )

    @property
    def group(self):
        return self.spec["apiVersion"].split("/", 2)[0]

    @property
    def version(self):
        return self.spec["apiVersion"].split("/", 2)[1]

    @property
    def plural(self):
        return self.spec["kind"].lower() + "s"

    @property
    def namespace(self):
        return self.spec["metadata"].get("namespace") or "default"

    @property
    def name(self):
        return self.spec["metadata"]["name"]

    def delete(self):
        """
        Delete the associated YAML file.
        """
        api = client.CustomObjectsApi()
        return api.delete_namespaced_custom_object(
            group=self.group,
            version=self.version,
            namespace=self.namespace,
            plural=self.plural,
            name=self.name,
        )
