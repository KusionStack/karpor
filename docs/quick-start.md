# karbour Quick Start Guide

This guide will walk you through the installation process of karbour, a cloud-native multi-cluster search and insights software. The installation process consists of three steps: creating a cluster using kind, installing karbour's manifest. Finally, you can access karbour's dashboard.

## Prerequisites

* Ensure [kubectl](https://kubernetes.io/docs/tasks/tools/) is installed.
* Ensure [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation/) is installed.

## Step 1: Create Cluster

First, you need to create a Kubernetes cluster in your local environment. We will use the `kind` tool to create the cluster. Follow these steps:

1. Create a cluster. You can create a cluster named karbour-cluster using the following command:
   ```shell
   kind create cluster --name karbour-cluster
   ```
   This will create a new Kubernetes cluster in your local Docker environment. Wait for a moment until the cluster creation is complete.
2. Verify that the cluster is running properly by executing the command:
   ```shell
   kubectl cluster-info
   ```
   If everything is set up correctly, you'll see information about your Kubernetes cluster.

## Step 2: Install karbour

After creating the cluster, proceed with the installation of karbour:

1. Create a namespace for karbour by running the following command:
   ```shell
   kubectl create ns karbour
   ```
   This will create a namespace named karbour to install karbour resources.
2. Install karbour by applying the manifest files to the cluster:
   ```shell
   kubectl apply -f https://github.com/KusionStack/karbour/tree/main/manifests
   ```
   This command will apply the deployments and services defined in the karbour manifest files to the karbour namespace.
3. Wait for the installation of karbour to complete. You can check the status of karbour installation by running the following command:
   ```shell
   kubectl get pods -n karbour
   ```
   When all the karbour pods are in the Running state, it means the installation is complete.

## Step 3: Access karbour Dashboard

1. Run the following command to forward the karbour server port:
   ```shell
   kubectl -n karbour port-forward service/karbour-server 7443:7443
   ```
   This will create a port forward from your local machine to the karbour server.
2. Open your browser and enter the following URL:
   ```shell
   https://127.0.0.1:7443
   ```
   This will take you to the karbour dashboard.

Congratulations! You have successfully installed karbour. Now you can start using karbour for multi-cluster search and insights.

Please note that this guide only provides a quick start for karbour, and you may need to refer to additional documentation and resources to configure and use other features of karbour.

If you have any questions or concerns, feel free to consult the official documentation of karbour or seek relevant support.