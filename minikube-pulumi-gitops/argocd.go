package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
)

func deployArgoCD(ctx *pulumi.Context, k8sProvider *kubernetes.Provider) (*helmv3.Release, *corev1.Secret, error) {
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
		return nil, nil, err
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
		return nil, nil, err
	}

	return argocd, secret, nil
} 