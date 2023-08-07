import json
import time

from kubernetes import client, config, watch
from kubernetes.client.api import core_v1_api
from kubernetes.client.exceptions import ApiException
from kubernetes.client.models.v1_pod_list import V1PodList


class MetricBase:
    separator = "METRICS OPERATOR TIMEPOINT"
    collection_start = "METRICS OPERATOR COLLECTION START"
    collection_end = "METRICS OPERATOR COLLECTION END"
    metadata_start = "METADATA START"
    metadata_end = "METADATA END"

    def __init__(self, spec=None, **kwargs):
        """
        Create a persistent client to interact with a MiniCluster

        This currently assumes the namespace exists.
        """
        self.spec = spec
        self._core_v1 = kwargs.get("core_v1_api")

        # Load kubeconfig on Metricbase init only
        if self.spec is not None:
            config.load_kube_config()

    @property
    def namespace(self):
        if not self.spec:
            return
        return self.spec["metadata"].get("namespace") or "default"

    @property
    def name(self):
        if not self.spec:
            return
        return self.spec["metadata"]["name"]

    @property
    def classname(self):
        return self.__class__.__name__

    def container(self):
        return self.classname.replace("-", "_")

    def parse(self, pod, container):
        """
        Retrieve logs output and call parsing function
        """
        lines = self.stream_output(
            name=pod.metadata.name, namespace=self.namespace, container=container.name
        )
        return self.parse_log(lines)

    @property
    def core_v1(self):
        """
        Instantiate a core_v1 api (if not done yet)
        """
        if self._core_v1 is not None:
            return self._core_v1

        self.c = client.Configuration.get_default_copy()
        self.c.assert_hostname = False
        client.Configuration.set_default(self.c)
        self._core_v1 = core_v1_api.CoreV1Api()
        return self._core_v1

    def logging_containers(self, namespace=None, states=None, retry_seconds=5):
        """
        Return list of containers intended to get logs from
        """
        containers = []
        pods = self.wait(namespace, states, retry_seconds)
        container_name = getattr(self, "container_name", self.container)
        print(f"Looking for container name {container_name}...")
        for pod in pods.items:
            for container in pod.spec.containers:
                print(f"Assessing {container.name}")
                if container.name == container_name:
                    print(f"Found logging container {container.name}")
                    containers.append(
                        (
                            pod,
                            container,
                        )
                    )
        return containers

    def wait(self, namespace=None, states=None, retry_seconds=5):
        """
        Wait for one or more pods of interest to be done
        """
        namespace = namespace or self.namespace
        print(f"Looking for prefix {self.pod_prefix} in namespace {namespace}")
        pod_list = self.get_pods(namespace, self.pod_prefix)
        size = len(pod_list.items)

        # We only want logs when they are completed
        states = states or ["Completed", "Succeeded"]
        if not isinstance(states, list):
            states = [states]

        ready = set()
        while len(ready) != size:
            print(f"{len(ready)} pods are ready, out of {size}")
            pod_list = self.get_pods(name=self.pod_prefix, namespace=namespace)

            for pod in pod_list.items:
                print(f"{pod.metadata.name} is in phase {pod.status.phase}")
                if pod.status.phase not in states:
                    time.sleep(retry_seconds)
                    continue

                if pod.status.phase not in ["Terminating"]:
                    ready.add(pod.metadata.name)

        states = '" or "'.join(states)
        print(f'All pods are in states "{states}"')
        return pod_list

    def _filter_pods(self, pods, name):
        """
        Filter a set of pods (associated with a job) to a name prefix.
        """
        filtered = []
        for pod in pods.items:
            if pod.metadata.name.startswith(name):
                filtered.append(pod)
        pods.items = filtered
        return pods

    def get_pods(self, namespace=None, name=None):
        """
        Get namespaced pods metadata, either scoped to a name or entire namespace.
        """
        namespace = namespace or self.namespace
        try:
            req = self.core_v1.list_namespaced_pod(namespace, async_req=True)
            pods = req.get()

            # If name is present, filter down to pods with that prefix
            if name is not None:
                pods = self._filter_pods(pods, name)
            return pods

        # Not found - it was deleted
        except ApiException:
            return V1PodList(items=[])
        except Exception:
            time.sleep(2)
            return self.get_pods(namespace, name)

    def get_log_metadata(self, lines):
        """
        Given a log dump, split based on the known separators.
        """
        if self.metadata_start not in lines or self.metadata_end not in lines:
            print("Cannot find expected collection start or end lines, cannot parse")
            return {}
        metadata = lines.split(self.metadata_start, 1)[-1]
        metadata = metadata.split(self.metadata_end, 1)[0]
        return json.loads(metadata)

    def get_log_sections(self, lines):
        """
        Given a log dump, split into data sections
        """
        if self.collection_start not in lines or self.collection_end not in lines:
            print("Cannot find expected metadata start or end lines, cannot parse")
            return {}
        data = lines.split(self.collection_start, 1)[1:]
        data = "\n".join(data).split(self.collection_end, 1)[0]
        return data.split(self.separator)

    def stream_output(
        self,
        name,
        namespace,
        stdout=True,
        filename=None,
        container=None,
        timestamps=False,
    ):
        """
        Stream output, optionally printing also to stdout.

        Also return the output to the user.
        """
        watcher = watch.Watch()

        out = None
        if filename:
            out = open(filename, "w")

        # Stream output to file and return it if desired!
        lines = []
        for line in watcher.stream(
            self.core_v1.read_namespaced_pod_log,
            name=name,
            namespace=namespace,
            timestamps=timestamps,
            container=container,
            follow=True,
        ):
            if out:
                # Lines end with /r and we need to strip and add a newline
                out.write(line.strip() + "\n")
            if stdout:
                print(line)
            lines.append(line)

        if out:
            out.close()

        # The parser needs to split, so we expect a single cohesive string
        return "\n".join(lines)
