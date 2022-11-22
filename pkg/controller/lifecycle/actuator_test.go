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

package lifecycle_test

import (
	"context"

	// Local
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"
	actuator "github.com/intel/gardener-extension-cri-resmgr/pkg/controller/lifecycle"

	// Gardener
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/logger"

	// Other
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cri-resource-manager extension actuator tests", func() {
	It("rendering charts installation chart without configs", func() {
		configs := map[string]map[string]string{
			// this should generate one ConfigMap with two keys
			"static": {
				"fallback": "FALLBACK_BODY:1",
				"force":    "FORCE_BODY:1",
			},
			"dynamic": {
				// this should generate
				// ConfigMap with name "cri-resmgr-config.default"
				"default": "CONFIG_BODY_OF_DEFAULT: 1",
				// ConfigMap with name "cri-resmgr-config.nodeFoo"
				"nodeFoo": "CONFIG_BODY_OF_NODEFOO: 1",
			},
		}
		// TODO: consider using mock instead of real rendered - not enough logic inside golang code yet!
		// unused but useful for future
		// "github.com/golang/mock/gomock"
		a := actuator.NewActuator("mock").(*actuator.Actuator)
		ctx := context.TODO()
		log := logger.ZapLogger(true)

		ex := &extensionsv1alpha1.Extension{}
		secret, err := a.GenerateSecretData(log, ctx, ex, consts.Charts, consts.ChartPath, "foo_namespace", "v1.0.0", configs)
		Expect(err).NotTo(HaveOccurred())

		Expect(secret).Should(HaveKey(consts.InstallationSecretKey))

		// check static
		Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring(`name: "cri-resmgr-static-configs"`))
		Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring("FALLBACK_BODY:1")) // notice no space between is passed as is
		Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring("FORCE_BODY:1"))

		// check dynamic (first level is unpacked) and rest becomes multi string
		Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring(`name: "cri-resmgr-config.default"`))
		Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring("CONFIG_BODY_OF_DEFAULT: |"))
		Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring(`name: "cri-resmgr-config.nodeFoo"`))
		Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring("CONFIG_BODY_OF_NODEFOO: |"))
	})
})
