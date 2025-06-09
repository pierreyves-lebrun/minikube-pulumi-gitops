# minikube-pulumi-gitops

## Pre-requisties

- Install [Golang](https://go.dev/doc/install)
- Install [Pulumi](https://www.pulumi.com/docs/install/)
  - `mkdir -p ~/.pulumi-local` - Create a local folder to retain pulumi state
  - `pulumi login file://~/.pulumi-local` - Use local filsystem as pulumi backend
- Install [Minikube](https://minikube.sigs.k8s.io/docs/start/)
  - `minikube start` - Initialize a minikube cluster 

## Usage

- `cd minikube-pulumi-gitops` - Navigate to pulumi code directory
- `pulumi up` - Create the resources in the stack
- `export PULUMI_CONFIG_PASSPHRASE=""` - Use empty passphrase for secrets
- `pulumi stack init dev --non-interactive` - Create a new stack named 'dev'
- `pulumi up --yes --non-interactive --skip-preview` - Deploy the stack automatically: skip confirmation prompts and preview

## Design

This repository demonstrates a GitOps approach to Kubernetes cluster management using [Pulumi](https://www.pulumi.com/docs/install/) for cluster boostrapping and [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) for continuous deployment.

The setup consists of two main parts:

1. **Bootstrap Phase (Pulumi)**:
   - Boostraps a Minikube cluster from `KUBECONFIG`
   - Installs and configures ArgoCD
   - Creates an initial ArgoCD "app of apps" pattern for managing subsequent deployments

2. **Continuous Deployment (ArgoCD)**:
   - Deploys the External Secrets Operator
   - Sets up a [fake ClusterSecretStore provider](./apps/external-secrets/clustersecretstore.yaml) for demonstration purposes
   - Deploys a sample [dummy-helm](https://github.com/pierreyves-lebrun/dummy-helm) chart that showcases secret management through External Secrets

Note: In a production environment, you would replace the fake provider with a secure secrets backend such as [AWS Secrets Manager](https://external-secrets.io/latest/provider/aws-secrets-manager/) or [Google Cloud Secret Manager](https://external-secrets.io/latest/provider/google-secrets-manager/), etc.

## Validation

To access the dummy-helm service locally:
```bash
kubectl -n dummy-helm port-forward svc/dummy-helm 5678:5678
```

You can then access the service at `http://localhost:5678`
