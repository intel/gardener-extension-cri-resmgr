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
	"encoding/json"

	// Local
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"

	// Gardener
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"

	// Other
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
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

	for _, extension := range extensions {
		if extension.Type == consts.ExtensionType {
			providerConfig = extension.ProviderConfig
		}
	}
	// If providerConfig were specified in Shoot spec.extensions then merge it with configs.
	if providerConfig != nil {
		if err := json.Unmarshal(providerConfig.Raw, &criResMgrConfig); err != nil {
			logger.Error(err, "error unmarshalling providerConfig", "providerConfig", string(providerConfig.Raw))
			return nil, err
		}
		logger.Info("configs: from cluster.extensions.providerConfig", "criResMgrConfig", criResMgrConfig)
		for configName, configContents := range criResMgrConfig.Configs {
			configs[configName] = configContents
		}
	}

	return configs, nil
}
