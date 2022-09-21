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
	// Just for test
	Foo bool `json:"foo,omitempty"`
	// Configs is a map of name of config file for cri-resource-manager and its contents.
	Configs map[string]string `json:"configs,omitempty"`
}

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
			// ignore entries starting with dot (hidden or directories create by kuberntes when mounting configMaps)
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

	// II. Parse provideConfig from Cluster.Extension within Shoot defintion
	var providerConfig *runtime.RawExtension
	var criResMgrConfig *CriResMgrConfig

	for _, extension := range extensions {
		if extension.Type == consts.ExtensionType {
			providerConfig = extension.ProviderConfig
		}
	}

	// Has to be empty to allow helm values to merge
	if providerConfig != nil {
		if err := json.Unmarshal(providerConfig.Raw, &criResMgrConfig); err != nil {
			// gardencorev1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
			// conditionValid = gardencorev1beta1helper.UpdatedCondition(conditionValid, gardencorev1beta1.ConditionFalse, "ChartInformationInvalid", fmt.Sprintf("CRI-ResMgr Extension (providerConfig) connfig cannot be unmarshalled: %+v", err))
			panic(err)
			// logger.Error(err, "error unmarhasling providerConfig", "providerConfig", string(providerConfig.Raw))
			// return err
		}
		logger.Info("configs: from cluster.extensions.providerConfig", "criResMgrConfig", criResMgrConfig)
		for configName, configContents := range criResMgrConfig.Configs {
			configs[configName] = configContents
		}
	}

	return configs, nil
}
