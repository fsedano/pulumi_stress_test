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

type WaveDef map[int][]pulumi.Resource

func run(ctx *pulumi.Context) error {
	// pod, err := createPod(ctx, 10)
	// if err != nil {
	// 	log.Printf("An error!!")
	// 	ctx.Log.Error("an error happened", nil)
	// }
	waveList := make(WaveDef)
	for wave := 0; wave < 7; wave++ {
		waveList[wave] = []pulumi.Resource{}

		for i := 0; i < 10; i++ {
			depends := []pulumi.Resource{}
			if wave > 0 {
				depends = waveList[wave-1]
			}
			//log.Printf("Depends is %#v for wave=%d", depends, wave)
			deplName := fmt.Sprintf("wave%di%d", wave, i)
			res, err := createDeployment(ctx, deplName, depends)
			if err == nil {
				waveList[wave] = append(waveList[wave], res)
			}
		}

	}

	return nil
}

func createDeployment(ctx *pulumi.Context, name string, depends []pulumi.Resource) (*appsv1.Deployment, error) {
	log.Printf("Creating %s that depends on %#v", name, depends)
	labels := pulumi.StringMap{"app": pulumi.String(name)}
	deployment, err := appsv1.NewDeployment(ctx, name, &appsv1.DeploymentArgs{
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
	}, pulumi.DependsOn(depends))
	return deployment, err
}

// func createPod(ctx *pulumi.Context, i int) (*corev1.Pod, error) {
// 	pod, err := corev1.NewPod(ctx, fmt.Sprintf("pod%d", i), &corev1.PodArgs{
// 		Metadata: &v1.ObjectMetaArgs{
// 			Namespace: pulumi.String("pul"),
// 		},
// 		Spec: corev1.PodSpecArgs{
// 			Containers: corev1.ContainerArray{
// 				corev1.ContainerArgs{
// 					Name:  pulumi.String("nginx"),
// 					Image: pulumi.String("nginx2"),
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return pod, err

// }
