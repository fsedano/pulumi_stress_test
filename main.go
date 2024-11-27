package main

import (
	"fmt"
	"log"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(run)
}

func run(ctx *pulumi.Context) error {
	pod, err := createPod(ctx, 10)
	if err != nil {
		log.Printf("An error!!")
		ctx.Log.Error("an error happened", nil)
	}
	for i := 0; i < 50; i++ {

		createDeployment(ctx, i, pod)

	}

	return nil
}

func createDeployment(ctx *pulumi.Context, n int, pod *corev1.Pod) error {
	name := fmt.Sprintf("n1%d", n)
	labels := pulumi.StringMap{"app": pulumi.String(name)}
	_, err := appsv1.NewDeployment(ctx, name, &appsv1.DeploymentArgs{
		Metadata: &v1.ObjectMetaArgs{
			Labels:    labels,
			Name:      pulumi.String(name),
			Namespace: pulumi.String("pul"),
		},
		Spec: appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: labels,
			},
			Replicas: pulumi.Int(1),
			Template: corev1.PodTemplateSpecArgs{

				Metadata: &v1.ObjectMetaArgs{
					Namespace: pulumi.String("pul"),
					Labels:    labels,
				},
				Spec: corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:  pulumi.String(name),
							Image: pulumi.String("nginx"),
						},
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{pod}))
	return err
}

func createPod(ctx *pulumi.Context, i int) (*corev1.Pod, error) {
	pod, err := corev1.NewPod(ctx, fmt.Sprintf("pod%d", i), &corev1.PodArgs{
		Metadata: &v1.ObjectMetaArgs{
			Namespace: pulumi.String("pul"),
		},
		Spec: corev1.PodSpecArgs{
			Containers: corev1.ContainerArray{
				corev1.ContainerArgs{
					Name:  pulumi.String("nginx"),
					Image: pulumi.String("nginx2"),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return pod, err

}
