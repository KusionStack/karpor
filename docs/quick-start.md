# karpor Quick Start Guide

This guide will walk you through the installation process of karpor, a cloud-native multi-cluster search and insights software. The installation process consists of three steps: creating a cluster using kind, installing karpor's manifest. Finally, you can access karpor's dashboard.

## Prerequisites

* Ensure [kubectl](https://kubernetes.io/docs/tasks/tools/) is installed.
* Ensure [helm](https://helm.sh/docs/intro/install/) is installed.
* Ensure [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation/) is installed for local testing.

## Step 1: Create Cluster

First, you need to create a Kubernetes cluster in your local environment. We will use the `kind` tool to create the cluster. Follow these steps:

1. Create a cluster. You can create a cluster named karpor-cluster using the following command:
   ```shell
   kind create cluster --name karpor-cluster
   ```
   This will create a new Kubernetes cluster in your local Docker environment. Wait for a moment until the cluster creation is complete.
2. Verify that the cluster is running properly by executing the command:
   ```shell
   kubectl cluster-info
   ```
   If everything is set up correctly, you'll see information about your Kubernetes cluster.

## Step 2: Install karpor

After creating the cluster, proceed with the installation of karpor:
1. Clone the Karpor project using git and navigate to the charts/karpor directory:
   ```shell
   git clone https://github.com/KusionStack/karpor.git
   cd charts
   ```
2. Run the following command to install Karpor using Helm:
   ```shell
   helm install karpor ./karpor
   ```
3. Wait for the installation of karpor to complete. You can check the status of karpor installation by running the following command:
   ```shell
   kubectl get pods -n karpor
   ```
   When all the karpor pods are in the Running state, it means the installation is complete.

## Step 3: Access karpor Dashboard

1. Run the following command to forward the karpor server port:
   ```shell
   kubectl -n karpor port-forward service/karpor-server 7443:7443
   ```
   This will create a port forward from your local machine to the karpor server.
2. Open your browser and enter the following URL:
   ```shell
   https://127.0.0.1:7443
   ```
   This will take you to the karpor dashboard.

Congratulations! You have successfully installed karpor. Now you can start using karpor for multi-cluster search and insights.

Please note that this guide only provides a quick start for karpor, and you may need to refer to additional documentation and resources to configure and use other features of karpor.

If you have any questions or concerns, feel free to consult the official documentation of karpor or seek relevant support.