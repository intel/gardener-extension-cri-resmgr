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

package actuator_test

import (
	"context"

	// Local
	"github.com/intel/gardener-extension-cri-resmgr/pkg/actuator"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"

	// Gardener
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	// Other
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cri-resource-manager extension main tests", func() {
	Describe("rendering charts", func() {

		It("installation chart without configs", func() {
			configs := map[string]string{
				// this should generate
				// ConfigMap with name "cri-resmgr-config.default"
				// with data.policy "THIS_WILL_CONFIG_BODY_OF_DEFAULT"
				"default": "THIS_WILL_CONFIG_BODY_OF_DEFAULT",
				// nodeFoo -> "THIS_WILL_CONFIG_BODY_OF_NODEFOO"
				"nodeFoo": "THIS_WILL_CONFIG_BODY_OF_NODEFOO",
			}
			// TODO: consider using mock instead of real rendered - not enough logic inside golang code yet!
			// unused but usefull for future
			// "github.com/golang/mock/gomock"
			a := actuator.NewActuator().(*actuator.Actuator)
			ctx := context.TODO()

			ex := &extensionsv1alpha1.Extension{}
			secret, err := a.GenerateSecretData(ctx, ex, consts.ChartPath, "foo_namespace", "v1.0.0", configs)
			Expect(err).NotTo(HaveOccurred())

			Expect(secret).Should(HaveKey(consts.InstallationSecretKey))

			Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring(`name: "cri-resmgr-config.default"`))
			Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring("policy:  THIS_WILL_CONFIG_BODY_OF_DEFAULT"))
			Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring(`name: "cri-resmgr-config.nodeFoo"`))
			Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring("policy:  THIS_WILL_CONFIG_BODY_OF_NODEFOO"))
		})
	})
})
