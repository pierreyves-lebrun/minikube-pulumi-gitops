package main

import (
	apiextensions "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
)

func deployAppOfApps(ctx *pulumi.Context, k8sProvider *kubernetes.Provider, dependencies []pulumi.Resource) error {
	_, err := apiextensions.NewCustomResource(ctx, "app-of-apps",
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
		pulumi.DependsOn(dependencies))

	return err
} 