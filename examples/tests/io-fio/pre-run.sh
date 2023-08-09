#!/bin/bash

HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TESTS=$(dirname ${HERE})

echo "Creating local volume in minikube"

# The "data" volume will be mounted at /mnt/data
minikube ssh -- mkdir -p /tmp/workflow