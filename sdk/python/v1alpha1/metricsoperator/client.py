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

    def watch(self, raw_logs=False, pod_prefix=None, container_name=None):
        """
        Wait for (and yield parsed) metric logs.
        """
        if raw_logs and not pod_prefix:
            raise ValueError("You must provide a pod_prefix to ask for raw logs.")

        for metric in self.spec["spec"]["metrics"]:
            if raw_logs:
                parser = mutils.get_metric()(self.spec, container_name=container_name)
            else:
                parser = mutils.get_metric(metric["name"])(
                    self.spec, container_name=container_name
                )
            print("Watching %s" % metric["name"])
            for pod, container in parser.logging_containers(pod_prefix=pod_prefix):
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

    def delete(self, pod_prefix=None):
        """
        Delete the associated YAML file.
        """
        api = client.CustomObjectsApi()
        result = api.delete_namespaced_custom_object(
            group=self.group,
            version=self.version,
            namespace=self.namespace,
            plural=self.plural,
            name=self.name,
        )
        self.wait_for_delete(pod_prefix)
        return result

    def wait_for_delete(self, pod_prefix=None):
        """
        Wait for pods to be gone (deleted)
        """
        for metric in self.spec["spec"]["metrics"]:
            parser = mutils.get_metric(metric["name"])(self.spec)
            print("Watching %s for deletion" % metric["name"])
            parser.wait_for_delete(pod_prefix=pod_prefix)
