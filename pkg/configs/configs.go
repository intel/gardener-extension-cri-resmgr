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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// Local
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"

	// Gardener
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"

	// Other
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ConfigsOverrideEnv = "CONFIGS_OVERWRITE"
)

// CriResMgrConfig is a providerConfig specific type for CRI-res-mgr extension.
type CriResMgrConfig struct {
	// Configs is a map of name of config file for cri-resource-manager and its contents.
	Configs map[string]string `json:"configs,omitempty"`
}

// GetConfigs gets and merges configs values from Shoot.spec.extensions.providerConfig and
// from files found in directory defined by ConfigsOverrideEnv.
// Path defined by ConfigsOverrideEnv is first validated (is directory) then all files are read
// and passed to helm installation charts to be rendered and additional configmaps for cri-resource-manager.
func GetConfigs(logger logr.Logger, extensions []v1beta1.Extension) (map[string]string, error) {
	configs := map[string]string{}

	// I. Configs are read from directory provided by ConfigsOverrideEnv
	configsOverwritePath := os.Getenv(ConfigsOverrideEnv)
	if len(configsOverwritePath) > 0 {
		path, err := os.Open(configsOverwritePath)
		if err != nil {
			return nil, err
		}
		defer path.Close()
		fileStat, err := path.Stat()
		if err != nil {
			return nil, fmt.Errorf("cannot stat path provided by ConfigsOverrideEnv")
		}
		if !fileStat.IsDir() {
			return nil, fmt.Errorf("provided %s from is not a directory", ConfigsOverrideEnv)
		}
		dirInfo, err := path.ReadDir(-1)
		if err != nil {
			return nil, fmt.Errorf("cannot ReadDir %w", err)
		}
		for _, dirEntry := range dirInfo {
			configName := dirEntry.Name()
			fullPath := filepath.Join(path.Name(), dirEntry.Name())
			// ignore entries starting with dot (hidden or directories create by Kubernetes when mounting configMaps)
			if strings.HasPrefix(configName, ".") || strings.HasPrefix(configName, "..") {
				continue
			}
			configContents, err := os.ReadFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("cannot read file of config file: %w", err)
			}
			configs[configName] = string(configContents)
		}
		logger.Info("configs: from env provided directory", "configs", configs)
	}

	// II. Parse provideConfig data from Cluster.Extension (it is a copy from within Shoot.spec.extensions.providerConfig).
	var providerConfig *runtime.RawExtension
	var criResMgrConfig *CriResMgrConfig

	for _, extension := range extensions {
		if extension.Type == consts.ExtensionType {
			providerConfig = extension.ProviderConfig
		}
	}
	// If providerConfig were specified in Shoot spec.extensions then merge it values with those found in filesystem
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
