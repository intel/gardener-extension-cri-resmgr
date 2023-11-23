// Copyright 2022 Intel Corporation. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2etest

import (
	"context"
	"os"
	"time"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/test/framework"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
)

var (
	// commonLabel      = "cri-rm"
	backgroundCtx   = context.Background()
	fiveteenMinutes = 15 * time.Minute //nolint:all
	ExtensionType   = "cri-resmgr-extension"

	// Those options can be overridden with arguments like: -verbose, -disable-dump, -existing-shoot-name -kubecfg -project-namespace
	projectNamespace   = "garden-local"
	kubeconfigPath     = os.Getenv("KUBECONFIG")
	skipAccessingShoot = true // if set to true then the test does not try to access the shoot via its kubeconfig
	commonConfig       = &framework.CommonConfig{}

	kubernetesVersion = "1.24.8"
)

func enableOrDisableCriResmgr(shoot *gardencorev1beta1.Shoot, disabled bool) error {
	for i, extension := range shoot.Spec.Extensions {
		if extension.Type == ExtensionType {
			if extension.Disabled != nil {
				shoot.Spec.Extensions[i].Disabled = pointer.Bool(disabled)
			}
		}
	}
	if len(shoot.Spec.Extensions) == 0 {
		shoot.Spec.Extensions = append(shoot.Spec.Extensions, gardencorev1beta1.Extension{
			Type:     ExtensionType,
			Disabled: pointer.Bool(disabled),
		})
	}
	return nil
}

func enableCriResmgr(shoot *gardencorev1beta1.Shoot) error {
	return enableOrDisableCriResmgr(shoot, false)
}
func disableCriResmgr(shoot *gardencorev1beta1.Shoot) error {
	return enableOrDisableCriResmgr(shoot, true)
}

func getShoot() *gardencorev1beta1.Shoot {
	secretBindingName := "local"
	networkingType := "calico"
	networking := gardencorev1beta1.Networking{
		Type:           &networkingType,
		ProviderConfig: &runtime.RawExtension{Raw: []byte(`{"apiVersion":"calico.networking.extensions.gardener.cloud/v1alpha1","kind":"NetworkConfig","typha":{"enabled":false},"backend":"none"}`)},
	}
	return &gardencorev1beta1.Shoot{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "", // will be overridden anyway
			Namespace:    projectNamespace,
			Annotations: map[string]string{
				// Should speed up shoot cleanup: https://github.com/gardener/gardener/blob/5f62609530c035caa4115f80ae16d0e52f7b0d14/pkg/apis/core/v1beta1/constants/types_constants.go#L533
				v1beta1constants.AnnotationShootInfrastructureCleanupWaitPeriodSeconds: "0",
				v1beta1constants.AnnotationShootCloudConfigExecutionMaxDelaySeconds:    "0",
			},
		},
		Spec: gardencorev1beta1.ShootSpec{
			Region:            "local",
			SecretBindingName: &secretBindingName,
			CloudProfileName:  "local",
			SeedName:          pointer.String("local"),
			Kubernetes: gardencorev1beta1.Kubernetes{
				Version: kubernetesVersion,
			},
			Networking: &networking,
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
