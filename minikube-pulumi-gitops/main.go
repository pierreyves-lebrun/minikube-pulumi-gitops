package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"os"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Use existing KUBECONFIG from environment
		kubeconfig := pulumi.String(os.Getenv("KUBECONFIG"))
		ctx.Export("kubeconfig", kubeconfig)

		// Create a Kubernetes Provider with the kubeconfig.
		k8sProvider, err := kubernetes.NewProvider(ctx, "k8sProvider", &kubernetes.ProviderArgs{
			Kubeconfig: kubeconfig,
		})
		if err != nil {
			return err
		}

		// Deploy ArgoCD and its repository secret
		argocd, secret, err := deployArgoCD(ctx, k8sProvider)
		if err != nil {
			return err
		}

		// Deploy the app-of-apps with dependencies on ArgoCD and the secret
		err = deployAppOfApps(ctx, k8sProvider, []pulumi.Resource{argocd, secret})
		if err != nil {
			return err
		}

		return nil
	})
}
