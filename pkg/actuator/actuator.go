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

package actuator

import (
	"context"
	"fmt"
	"time"

	// Local
	"github.com/intel/gardener-extension-cri-resmgr/pkg/configs"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"

	// Gardener
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	"github.com/gardener/gardener/extensions/pkg/util"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	managedresources "github.com/gardener/gardener/pkg/utils/managedresources"

	// Other
	"github.com/go-logr/logr"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/imagevector"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ---------------------------------------------------------------------------------------
// -                                        Actuator                                     -
// ---------------------------------------------------------------------------------------

func NewActuator() extension.Actuator {
	return &Actuator{
		ChartRendererFactory: extensionscontroller.ChartRendererFactoryFunc(util.NewChartRendererForShoot),
		logger:               log.Log.WithName(consts.ActuatorName),
	}
}

type Actuator struct {
	client               client.Client
	config               *rest.Config
	ChartRendererFactory extensionscontroller.ChartRendererFactory
	decoder              runtime.Decoder
	logger               logr.Logger
}

func (a *Actuator) GenerateSecretData(logger logr.Logger, ctx context.Context, ex *extensionsv1alpha1.Extension,
	chartPath string, namespace string, k8sversion string, configs map[string]string) (map[string][]byte, error) {
	emptyMap := map[string][]byte{}
	// Depending on shoot, chartredner will have different capabilities based on K8s version.
	chartRenderer, err := a.ChartRendererFactory.NewChartRendererForShoot(k8sversion)
	if err != nil {
		return emptyMap, err
	}
	imageVector := imagevector.ImageVector()
	if len(imageVector) > 0 {
		for i, imageSource := range imageVector {

			logger.Info(fmt.Sprintf("imageVector[%d].imageSource", i), "imageSource", *imageSource)
		}
	}

	installationImage, err := imageVector.FindImage(consts.InstallationImageName)
	if err != nil {
		return emptyMap, err
	}
	agentImage, err := imageVector.FindImage(consts.AgentImageName)
	if err != nil {
		return emptyMap, err
	}
	chartValues := map[string]interface{}{
		"images": map[string]string{
			consts.InstallationImageName: installationImage.String(),
			consts.AgentImageName:        agentImage.String(),
		},
		"configs": configs,
	}
	release, err := chartRenderer.Render(chartPath, consts.InstallationReleaseName, metav1.NamespaceSystem, chartValues)
	//release, err := chartRenderer.RenderEmbeddedFS(chartPath, InstallationReleaseName, metav1.NamespaceSystem, chartValues)

	if err != nil {
		return emptyMap, err
	}
	// Put chart into secret
	secretData := map[string][]byte{consts.InstallationSecretKey: release.Manifest()}
	return secretData, nil
}

func (a *Actuator) Reconcile(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()
	a.logger.Info("Reconcile: checking extension...") // , "shoot", cluster.Shoot.Name, "namespace", cluster.Shoot.Namespace)

	// Find what shoot cluster we dealing with.
	// to find k8s version for chart renderer
	// and get providerConfig for configurations for CRI-resource-manager configmaps
	// with imageVector support this would allow to choose different version of extensions depending on k8s version
	cluster, err := extensionscontroller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return err
	}

	// Get configs either from providerConfig and merged it with provided files
	configs, err := configs.GetConfigs(a.logger, cluster.Shoot.Spec.Extensions)
	if err != nil {
		panic(err)
		// return err
	}

	secretData, err := a.GenerateSecretData(logger, ctx, ex, consts.ChartPath, namespace, cluster.Shoot.Spec.Kubernetes.Version, configs)
	if err != nil {
		panic(err)
		// return err
	}

	// Reconcile managedresource and secret for shoot.
	if err := managedresources.CreateForShoot(ctx, a.client, namespace, consts.ManagedResourceName, false, secretData); err != nil {
		return err
	}

	// mr := &resourcemanagerv1alpha1.ManagedResource{}
	// if err := a.client.Get(ctx, kutil.Key(namespace, ManagedResourceName), mr); err != nil {
	// 	// Continue only if not found.
	// 	if !apierrors.IsNotFound(err) {
	// 		return err
	// 	}
	// } else {
	// 	a.logger.Info("Reconcile: extension is already installed. Ignoring.") //, "shoot", cluster.Shoot.Name, "namespace", cluster.Shoot.Namespace)
	// }

	return nil
}

func (a *Actuator) Delete(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()
	cluster, err := extensionscontroller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return err
	}

	// if err := a.deployDaemonsetToUninstallCriResMgr(ctx, ex); err != nil {
	// 	return err
	// }

	a.logger.Info("Delete: deleting extension managedresources in shoot", "shoot", cluster.Shoot.Name, "namespace", cluster.Shoot.Namespace)

	twoMinutes := 1 * time.Minute

	timeoutShootCtx, cancelShootCtx := context.WithTimeout(ctx, twoMinutes)
	defer cancelShootCtx()

	if err := managedresources.DeleteForShoot(ctx, a.client, namespace, consts.ManagedResourceName); err != nil {
		return err
	}

	if err := managedresources.WaitUntilDeleted(timeoutShootCtx, a.client, namespace, consts.ManagedResourceName); err != nil {
		return err
	}

	return nil
}

func (a *Actuator) Restore(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Reconcile(ctx, logger, ex)
}

func (a *Actuator) Migrate(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Delete(ctx, logger, ex)
}

func (a *Actuator) InjectConfig(config *rest.Config) error {
	a.config = config
	return nil
}

func (a *Actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

func (a *Actuator) InjectScheme(scheme *runtime.Scheme) error {
	a.decoder = serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()
	return nil
}
