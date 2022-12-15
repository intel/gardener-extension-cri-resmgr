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
	"os"

	// Local

	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"

	// Gardener

	// Other
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

// PrepareConfigTypes merges baseConfigs and providerConfigs (found in extensions.providerConfig) and split that two "static" and "dynamic" types.
func PrepareConfigTypes(logger logr.Logger, baseConfigs map[string]string, providerConfigs map[string]string) (map[string]map[string]string, error) {

	// Merge result is used for helm installation charts to be rendered and additional configmaps for cri-resource-manager.
	// overwrite baseConfigs with those read from providerConfig.Configs
	configKeys := []string{}
	for configName, configContents := range providerConfigs {
		baseConfigs[configName] = configContents
		configKeys = append(configKeys, configName)
	}
	logger.Info("configs: from shoot.providerConfig configs", "types", configKeys)

	// Split into types
	configTypes := map[string]map[string]string{
		"static":  {},
		"dynamic": {},
	}
	configTypes["static"] = map[string]string{}
	for configName, configContent := range baseConfigs {
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
