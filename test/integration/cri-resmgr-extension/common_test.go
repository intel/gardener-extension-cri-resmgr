package cri_resmgr_extension

import (
	"context"
	"os"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var (
	backgroundCtx    = context.Background()
	ExtensionType    = "cri-resmgr-extension"
	projectNamespace = "garden"
	kubeconfigPath   = os.Getenv("KUBECONFIG")
	fiveteenMinutes  = 15 * time.Minute
)

func enableCriResmgr(shoot *gardencorev1beta1.Shoot) error {
	for i, extension := range shoot.Spec.Extensions {
		if extension.Type == ExtensionType {
			if extension.Disabled != nil {
				shoot.Spec.Extensions[i].Disabled = pointer.Bool(false)
			}
		}
	}
	shoot.Spec.Extensions = append(shoot.Spec.Extensions, gardencorev1beta1.Extension{
		Type:     ExtensionType,
		Disabled: pointer.Bool(false),
	})
	return nil
}

func disableCriResmgr(shoot *gardencorev1beta1.Shoot) error {
	for i, extension := range shoot.Spec.Extensions {
		if extension.Type == ExtensionType {
			if extension.Disabled != nil {
				shoot.Spec.Extensions[i].Disabled = pointer.Bool(true)
			}
		}
	}
	shoot.Spec.Extensions = append(shoot.Spec.Extensions, gardencorev1beta1.Extension{
		Type:     ExtensionType,
		Disabled: pointer.Bool(true),
	})
	return nil
}

func getShoot() *gardencorev1beta1.Shoot {
	return &gardencorev1beta1.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "local",
			Namespace:    projectNamespace,
			Annotations:  map[string]string{},
		},
		Spec: gardencorev1beta1.ShootSpec{
			Region:            "local",
			SecretBindingName: "local",
			CloudProfileName:  "local",
			SeedName:          pointer.String("local"),
			Kubernetes: gardencorev1beta1.Kubernetes{
				Version: "1.23.1",
			},
			Networking: gardencorev1beta1.Networking{
				Type: "local",
			},
			Provider: gardencorev1beta1.Provider{
				Type: "local",
				Workers: []gardencorev1beta1.Worker{{
					Name: "local",
					Machine: gardencorev1beta1.Machine{
						Type: "local",
					},
					CRI: &gardencorev1beta1.CRI{
						Name: gardencorev1beta1.CRINameContainerD,
					},
					Minimum: 1,
					Maximum: 1,
				}},
			},
		},
	}
}
