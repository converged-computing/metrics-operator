#!/bin/bash

echo "Cleaning up /tmp/workflow in minikube"
minikube ssh -- sudo rm -rf /tmp/workflow