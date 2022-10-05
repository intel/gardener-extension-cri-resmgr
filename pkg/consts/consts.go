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
package consts

import _ "embed"

//go:embed VERSION
var Version string

//go:embed COMMIT
var Commit string

const (
	ExtensionName = "cri-resmgr"
	ExtensionType = "cri-resmgr-extension"

	ControllerName = "cri-resmgr-controller"
	ActuatorName   = "cri-resmgr-actuator"

	ManagedResourceName      = "extension-runtime-cri-resmgr"
	ConfigMapName            = "gardener-extension-cri-resmgr-configs"
	ConfigMapNamespaceEnvKey = "EXTENSION_CONFIGMAP_NAMESPACE"
	ConfigKey                = "config.yaml"

	ChartPath               = "charts/internal/cri-resmgr-installation/"
	InstallationImageName   = "gardener-extension-cri-resmgr-installation"
	AgentImageName          = "gardener-extension-cri-resmgr-agent"
	InstallationReleaseName = "cri-resmgr-installation"
	InstallationSecretKey   = "installation_chart"
)
