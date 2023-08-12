#!/bin/bash

# Usage: /bin/bash script/test.sh $name 30
HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ROOT=$(dirname ${HERE})
cd ${ROOT}

set -eEu -o pipefail

name=${1}
jobtime=${2:-30}

echo "   Name: ${name}"
echo "Jobtime: ${jobtime}"

# Output and error files
out="${ROOT}/examples/tests/${name}/${name}-log.out"
err="${ROOT}/examples/tests/${name}/${name}-log.err"

# Quick helper script to run a test
make clean >> /dev/null 2>&1
make run > ${out} 2> ${err} &
pid=$!
echo "PID for running cluster is ${pid}"

/bin/bash examples/tests/${name}/pre-run.sh || true
kubectl apply -f examples/tests/${name}/metrics.yaml
echo "Sleeping for ${jobtime} seconds to allow job to complete 😴️."
sleep ${jobtime}

# The job should be completed
type=$(kubectl get jobset -o json | jq -r .items[0].status.conditions[0].type)
echo "JobSet status type is ${type}"
status=$(kubectl get jobset -o json | jq -r .items[0].status.conditions[0].status)

if [[ "${status}" != "True" ]] || [[ "${type}" != "Completed" ]]; then
    echo "Issue with running job ${name}"
    /bin/bash examples/tests/${name}/post-run.sh || true
    echo "LOGS for Metrics Operator pod"
    metrics_pod=$(kubectl get -n metrics-system pods -o json | jq -r .items[0].metadata.name)
    kubectl logs -n metrics-system ${metrics_pods}
    echo "LOGS for JobSet Operator pod"
    jobset_pod=$(kubectl get -n jobset-system pods -o json | jq -r .items[0].metadata.name)
    kubectl logs -n jobset-system ${jobset_pod}
    exit 1
fi

kill ${pid} || true
kill $(lsof -t -i:8080) || true
/bin/bash examples/tests/${name}/post-run.sh || true
