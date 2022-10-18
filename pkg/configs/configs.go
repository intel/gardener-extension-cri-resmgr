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

package configs

import (
	"context"
	"encoding/json"
	"os"

	// Local

	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"

	// Gardener
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"

	// Other
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CriResMgrConfig is a providerConfig specific type for CRI-res-mgr extension.
type CriResMgrConfig struct {
	// Configs is a map of name of config file for cri-resource-manager and its contents.
	Configs map[string]string `json:"configs,omitempty"`
}

// MergeConfigs merges base configs and values from Shoot.spec.extensions.providerConfig.
// Result is then used for helm installation charts to be rendered and additional configmaps for cri-resource-manager.
func MergeConfigs(logger logr.Logger, configs map[string]string, extensions []v1beta1.Extension) (map[string]string, error) {

	// Get and parse provideConfig data from Cluster.Extension (it is a copy from within Shoot.spec.extensions.providerConfig).
	var providerConfig *runtime.RawExtension
	var criResMgrConfig *CriResMgrConfig

	// If providerConfig were specified in Shoot spec.extensions then merge it with configs.
	for _, extension := range extensions {
		if extension.Type == consts.ExtensionType {
			providerConfig = extension.ProviderConfig
		}
	}
	if providerConfig != nil {
		if err := json.Unmarshal(providerConfig.Raw, &criResMgrConfig); err != nil {
			logger.Error(err, "configs: ERROR unmarshalling providerConfig", "providerConfig", string(providerConfig.Raw))
			return nil, err
		}
		configKeys := []string{}
		for configName, configContents := range criResMgrConfig.Configs {
			configs[configName] = configContents
			configKeys = append(configKeys, configName)
		}
		logger.Info("configs: from shoot.providerConfig configs", "types", configKeys)
	}
	return configs, nil
}

// GetBaseConfigsFromConfigMap reads extension ConfigMap and get its "configs" as baseConfigs
func GetBaseConfigsFromConfigMap(ctx context.Context, logger logr.Logger, k8sClient client.Client) (map[string]string, error) {

	baseConfigs := map[string]string{}
	extensionConfigMapNamespace := os.Getenv(consts.ConfigMapNamespaceEnvKey)
	if extensionConfigMapNamespace != "" {
		configMap := &corev1.ConfigMap{}
		err := k8sClient.Get(ctx, client.ObjectKey{Namespace: extensionConfigMapNamespace, Name: consts.ConfigMapName}, configMap)
		if err != nil {
			logger.Info("configs: cannot find configMap with base configs, return empty map")
			return baseConfigs, nil
		}
		baseConfigs = configMap.Data

		// Just for logging to not print all the contents
		baseConfigsKeys := []string{}
		for key := range configMap.Data {
			baseConfigsKeys = append(baseConfigsKeys, key)
		}
		logger.Info("configs: from configMap use as baseConfigs", "types", baseConfigsKeys)
	}
	return baseConfigs, nil
}

// PrepareConfigTypes merges baseConfigs and configs found in extensions.providerConfig and split that two "static" and "dynamic" types.
func PrepareConfigTypes(logger logr.Logger, baseConfigs map[string]string, extensions []v1beta1.Extension) (map[string]map[string]string, error) {

	// Get configs either from configMap (initial) and override with values from Shot.Spec.Extensions.providerConfig
	configs, err := MergeConfigs(logger, baseConfigs, extensions)
	if err != nil {
		return nil, err
	}

	configTypes := map[string]map[string]string{
		"static":  {},
		"dynamic": {},
	}
	configTypes["static"] = map[string]string{}
	for configName, configContent := range configs {
		if configName == "fallback" || configName == "force" || configName == "EXTRA_OPTIONS" {
			// static configs
			configTypes["static"][configName] = configContent
		} else {
			// dynamic configs
			configTypes["dynamic"][configName] = configContent
		}
	}
	return configTypes, nil
}
