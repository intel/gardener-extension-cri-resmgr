// Copyright 2022 Intel Corporation. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lifecycle

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	// Local
	"github.com/intel/gardener-extension-cri-resmgr/pkg/configs"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"

	// Gardener
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	"github.com/gardener/gardener/extensions/pkg/util"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	managedresources "github.com/gardener/gardener/pkg/utils/managedresources"

	// Other
	"github.com/go-logr/logr"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/imagevector"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CriResMgrConfig is a providerConfig specific type for CRI-res-mgr extension.
type CriResMgrConfig struct {
	// Configs is a map of name of config file for cri-resource-manager and its contents.
	Configs map[string]string `json:"configs,omitempty"`
	// nodeSelector
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// GetProviderConfig return CriResMgrConfig.
func GetProviderConfig(logger logr.Logger, extensions []v1beta1.Extension) (bool, CriResMgrConfig, error) {
	// Get and parse provideConfig data from Cluster.Extension (it is a copy from within Shoot.spec.extensions.providerConfig).
	var providerConfig *runtime.RawExtension
	var criResMgrConfig *CriResMgrConfig

	// Get providerConfig for our extension
	for _, extension := range extensions {
		if extension.Type == consts.ExtensionType {
			providerConfig = extension.ProviderConfig
		}
	}
	// If found, then parse and unmarshal
	if providerConfig != nil {
		if err := json.Unmarshal(providerConfig.Raw, &criResMgrConfig); err != nil {
			logger.Error(err, "ERROR unmarshalling providerConfig", "providerConfig", string(providerConfig.Raw))
			return false, CriResMgrConfig{}, err
		}
		return true, *criResMgrConfig, nil
	}
	return false, CriResMgrConfig{}, nil
}

// ---------------------------------------------------------------------------------------
// -                                        Actuator                                     -
// ---------------------------------------------------------------------------------------

// NewActuator return new Actuator.
func NewActuator(c client.Client, name string) extension.Actuator {
	return &Actuator{
		client:               c,
		ChartRendererFactory: extensionscontroller.ChartRendererFactoryFunc(util.NewChartRendererForShoot),
		logger:               log.Log.WithName(name),
	}
}

// NewActuatorWithSuffix return new Actuator with suffix.
func NewActuatorWithSuffix(c client.Client, nameSuffix string) extension.Actuator {
	return &Actuator{
		client:               c,
		ChartRendererFactory: extensionscontroller.ChartRendererFactoryFunc(util.NewChartRendererForShoot),
		logger:               log.Log.WithName(consts.ActuatorName + nameSuffix),
	}
}

// Actuator type.
type Actuator struct {
	client client.Client
	//config               *rest.Config
	ChartRendererFactory extensionscontroller.ChartRendererFactory
	//decoder              runtime.Decoder
	logger logr.Logger
}

// GenerateSecretData return byte map which is k8s secret with data.
func (a *Actuator) GenerateSecretData(logger logr.Logger, charts embed.FS, chartPath string,
	_ string, k8sVersion string, configs map[string]map[string]string, nodeSelector map[string]string) (map[string][]byte, error) {
	emptyMap := map[string][]byte{}
	// Depending on shoot, chartRenderer will have different capabilities based on K8s version.
	chartRenderer, err := a.ChartRendererFactory.NewChartRendererForShoot(k8sVersion)
	if err != nil {
		return emptyMap, err
	}

	// Check if config was not empty
	if nodeSelector == nil {
		nodeSelector = map[string]string{}
	}
	// Only run on containerd nodes
	nodeSelector[extensionsv1alpha1.CRINameWorkerLabel] = string(extensionsv1alpha1.CRINameContainerD)

	imageVector := imagevector.ImageVector()
	if len(imageVector) > 0 {
		for _, imageSource := range imageVector {
			logger.Info(fmt.Sprintf("images: found imageVector[name=%s]", imageSource.Name), "imageSource", (*imageSource.ToImage(&k8sVersion)).String())
		}
	}
	// TODO k8sVersion can be used to extend FindImage FindOptions(targetVersion)
	// to choose different version of image depending of target shoot Kubernetes. Not needed for now.
	balloonsImage, err := imageVector.FindImage(consts.BalloonsImageName)
	if err != nil {
		return emptyMap, err
	}
	chartValues := map[string]interface{}{
		"images": map[string]string{
			consts.BalloonsImageName: balloonsImage.String(),
		},
		"configs":      configs,
		"nodeSelector": nodeSelector,
	}

	release, err := chartRenderer.RenderEmbeddedFS(charts, chartPath, consts.InstallationReleaseName, metav1.NamespaceSystem, chartValues)

	if err != nil {
		return emptyMap, err
	}
	// Put chart into secret
	secretData := map[string][]byte{consts.InstallationSecretKey: release.Manifest()}
	return secretData, nil
}

// GenerateSecretDataToMonitoringManagedResource return byte map which is prepared config to monitoring.
func (a *Actuator) GenerateSecretDataToMonitoringManagedResource(namespace string) map[string][]byte {
	// Replace marker in namespace field to true namespace.
	yamlStringConfigNameWithNamespace := regexp.MustCompile(`{{ namespace }}`).ReplaceAllString(string(consts.MonitoringYaml), namespace)

	return map[string][]byte{"data": []byte(yamlStringConfigNameWithNamespace)}
}

// Reconcile the Extension resource.
func (a *Actuator) Reconcile(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()

	if a.client == nil {
		panic("a.client is nil!")
	}

	// Find what shoot cluster we dealing with.
	// to find k8s version for chart renderer
	// and get providerConfig for configurations for CRI-resource-manager configmaps
	// with imageVector support this would allow to choose different version of extensions depending on k8s version
	cluster, err := extensionscontroller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return err
	}

	// Get the baseConfigs from extension ConfigMap.
	baseConfigs, err := configs.GetBaseConfigsFromConfigMap(ctx, a.logger, a.client)
	if err != nil {
		return err
	}
	foundProviderConfig, criResMgrConfig, err := GetProviderConfig(logger, cluster.Shoot.Spec.Extensions)
	if err != nil {
		return err
	}
	var providerConfigs map[string]string
	nodeSelector := map[string]string{}
	if foundProviderConfig {
		providerConfigs = criResMgrConfig.Configs
		nodeSelector = criResMgrConfig.NodeSelector
	}

	// Merge baseConfigs and providerConfig.configs from Shoot.Spec and split into types "static","dynamic".
	configTypes, err := configs.PrepareConfigTypes(a.logger, baseConfigs, providerConfigs)
	if err != nil {
		return err
	}

	// Generate secret data that will be used by reference by ManagedResource to deploy.
	secretData, err := a.GenerateSecretData(a.logger, consts.Charts, consts.ChartPath, namespace, cluster.Shoot.Spec.Kubernetes.Version, configTypes, nodeSelector)
	if err != nil {
		return err
	}

	// Reconcile managedresource and secret for shoot.
	origin := "gardener"
	if err := managedresources.CreateForShoot(ctx, a.client, namespace, consts.ManagedResourceName, origin, false, secretData); err != nil {
		return err
	}

	//  Generate secret data that will be used by reference by ManagedResource to deploy.
	secretDataForMonitoring := a.GenerateSecretDataToMonitoringManagedResource(namespace)

	// Reconcile managedresource and secret for seed.
	err = managedresources.CreateForSeed(ctx, a.client, namespace, consts.MonitoringManagedResourceName, false, secretDataForMonitoring)
	return err
}

// ForceDelete forcefully deletes the Extension resource.
func (a *Actuator) ForceDelete(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Delete(ctx, logger, ex)
}

// Delete the Extension resource.
func (a *Actuator) Delete(ctx context.Context, _ logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()
	cluster, err := extensionscontroller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return err
	}
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

	if err := managedresources.DeleteForSeed(ctx, a.client, namespace, consts.MonitoringManagedResourceName); err != nil {
		return err
	}

	err = managedresources.WaitUntilDeleted(timeoutShootCtx, a.client, namespace, consts.MonitoringManagedResourceName)
	return err
}

// Restore the Extension resource.
func (a *Actuator) Restore(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Reconcile(ctx, logger, ex)
}

// Migrate the Extension resource.
func (a *Actuator) Migrate(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Delete(ctx, logger, ex)
}
