package main

import (
	apiextensions "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
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

		// Install the ArgoCD chart.
		argocd, err := helmv3.NewRelease(ctx, "argocd", &helmv3.ReleaseArgs{
			Name:  pulumi.String("argocd"),
			Chart: pulumi.String("argo-cd"),
			RepositoryOpts: &helmv3.RepositoryOptsArgs{
				Repo: pulumi.String("https://argoproj.github.io/argo-helm"),
			},
			Version:         pulumi.String("8.0.16"),
			Namespace:       pulumi.String("argocd"),
			CreateNamespace: pulumi.Bool(true),
		}, pulumi.Provider(k8sProvider))

		if err != nil {
			return err
		}

		// Create the repository secret after ArgoCD is installed
		secret, err := corev1.NewSecret(ctx, "dummy-helm-repo-secret", &corev1.SecretArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("dummy-helm-repo"),
				Namespace: pulumi.String("argocd"),
				Labels: pulumi.StringMap{
					"argocd.argoproj.io/secret-type": pulumi.String("repository"),
				},
			},
			StringData: pulumi.StringMap{
				"enableOCI": pulumi.String("true"),
				"url":      pulumi.String("registry-1.docker.io/pierreyveslebrun"),
				"type":     pulumi.String("helm"),
			},
		}, pulumi.Provider(k8sProvider),
			// Make sure ArgoCD is installed before creating the secret
			pulumi.DependsOn([]pulumi.Resource{argocd}))

		if err != nil {
			return err
		}

		// Create the ArgoCD Application after both ArgoCD and the secret are ready
		_, err = apiextensions.NewCustomResource(ctx, "app-of-apps",
			&apiextensions.CustomResourceArgs{
				ApiVersion: pulumi.String("argoproj.io/v1alpha1"),
				Kind:       pulumi.String("Application"),
				Metadata: &metav1.ObjectMetaArgs{
					Name:      pulumi.String("apps"),
					Namespace: pulumi.String("argocd"),
				},
				OtherFields: map[string]interface{}{
					"spec": map[string]interface{}{
						"project": "default",
						"syncPolicy": map[string]interface{}{
							"automated": map[string]interface{}{
								"selfHeal": true,
								"prune":    false,
							},
							"syncOptions": []string{
								"CreateNamespace=true",
							},
						},
						"destination": map[string]interface{}{
							"name":      "in-cluster",
							"namespace": "argocd",
						},
						"source": map[string]interface{}{
							"repoURL":        "https://github.com/pierreyves-lebrun/minikube-pulumi-gitops",
							"targetRevision": "HEAD",
							"path":           "apps",
						},
					},
				},
			},
			pulumi.Provider(k8sProvider),
			// Make sure both ArgoCD and the secret are ready before creating the Application
			pulumi.DependsOn([]pulumi.Resource{argocd, secret}))

		if err != nil {
			return err
		}

		return nil
	})
}
