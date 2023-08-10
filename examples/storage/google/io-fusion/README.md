# Fusion with Google Storage Example

For these experiments, we will create a small cluster on Google Cloud, and attempt
to demonstrate measuring IO using the [FIO metric](https://converged-computing.github.io/metrics-operator/getting_started/metrics.html#fio) of the Metrics operator
with [Fusion](https://seqera.io/fusion/).

## Usage

### 1. Create the Cluster

For fusion I needed to make a new cluster and [follow the basic instructions here](https://flux-framework.org/flux-operator/deployment/google/fusion.html?h=fusion#create-cluster).

```bash
GOOGLE_PROJECT=myproject
```
```bash
gcloud container clusters create flux-cluster --project $GOOGLE_PROJECT \
    --zone us-central1-a --machine-type n1-standard-2 --cluster-version 1.25 \
    --num-nodes=2 --enable-network-policy --tags=flux-cluster --enable-intra-node-visibility \
    --ephemeral-storage-local-ssd count=1 --workload-pool=${GOOGLE_PROJECT}.svc.id.goog \
    --workload-metadata=GKE_METADATA
```

### 2. Permissions

Create a service account for the cluster:

```bash
kubectl create serviceaccount metrics-operator-sa
```

Note that instructions for this are mentioned [here](https://flux-framework.org/flux-operator/deployment/google/fusion.html?h=fusion#create-cluster).

```bash
# This is a Google service account with permission to stoage
# List: gcloud iam service-accounts list
GOOGLE_SERVICE_ACCOUNT=GSA_NAME@GSA_PROJECT.iam.gserviceaccount.com
NAMESPACE="default"
KSA_NAME="metrics-operator-sa"
gcloud iam service-accounts add-iam-policy-binding ${GOOGLE_SERVICE_ACCOUNT} \
    --role roles/iam.workloadIdentityUser \
    --member "serviceAccount:${GOOGLE_PROJECT}.svc.id.goog[${NAMESPACE}/${KSA_NAME}]"
```
The above is granting permission from the Google cloud service account to our Kubernetes service account.
We then "annotate" the Kubernetes service account with the email address of the IAM service account.

```bash
kubectl annotate serviceaccount ${KSA_NAME} \
    --namespace ${NAMESPACE} \
    iam.gke.io/gcp-service-account=${GOOGLE_SERVICE_ACCOUNT}
```

Now, permission wise, we are ready to go!

### 3. Install JobSet

We then need to install JobSet

```bash
VERSION=v0.2.0
kubectl apply --server-side -f https://github.com/kubernetes-sigs/jobset/releases/download/$VERSION/manifests.yaml
```

### 4. Install the Metrics Operator

Let's next install the operator. You can [choose one of the options here](https://converged-computing.github.io/metrics-operator/getting_started/user-guide.html). 
E.g., to deploy from the cloned repository:

```bash
$ kubectl apply -f https://raw.githubusercontent.com/converged-computing/metrics-operator/main/examples/dist/metrics-operator.yaml
```

Check logs to make sure they look OK:

```bash
kubectl logs -n metrics-system  metrics-controller-manager-59cb7fb47b-7rgzz -f
```

### 5. Install the Device Manager

Fusion is cool because it binds from the pod. We will first need to apply a Daemonset that enables a device manager so fuse to be used:

```bash
$ kubectl apply -f daemonset.yaml
```

A few things I learned for this daemonset! If you create it in a namespace other than kube-system (and have the priority class defined) you’ll get an error about quotas. The fix (for now) is to deploy to the kube-system namespace. I also adjusted the config.yaml at the bottom to equal the number of nodes in my cluster (2 for this demo) and I’m not sure how much that matters (the previous setting was at 20).
And then label the nodes so they can use the manager:

```bash
for n in $(kubectl get nodes | tail -n +2 | cut -d' ' -f1 ); do
    kubectl label node $n smarter-device-manager=enabled
done 
```

After this, you should have two daemonset pods running:

```bash
$ kubectl get -n kube-system pods | grep device
```
```console
smarter-device-manager-qv7xf                             1/1     Running   0          39m
smarter-device-manager-trtbl                             1/1     Running   0          39m
```

### 6. Run Storage Metric!

Let's do a test manually.

```bash
kubectl apply -f metrics.yaml
```

You should see the pod - wait until it's running, and then look at the log.

```bash
kubectl logs metricset-sample-m-0-qvsj6 -f
```

When that looks okay, let's run this in a small loop! You'll need the metricsoperator sdk:

```bash
pip install metricsoperator
```

Now let's run a series of experiments so we have more than one timepoint.

```bash
mkdir -p ./data
for i in 0 1 2 3 4 5; do
  echo "Running iteration $i"
  python run-metric.py --out ./data/metrics-fusion-$i.json
  kubectl delete -f metrics.yaml
done
```

When you are done:

```bash
$ gcloud container clusters delete flux-cluster
```

Note that we have the data results and plots [in this repository](https://github.com/converged-computing/metrics-operator-experiments/tree/main/google/storage#results) 
for comparison with a few other metrics.