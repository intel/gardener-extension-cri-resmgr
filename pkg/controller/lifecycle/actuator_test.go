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

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/logger"
	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"
	actuator "github.com/intel/gardener-extension-cri-resmgr/pkg/controller/lifecycle"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cri-resource-manager extension actuator tests", func() {
	var (
		log logr.Logger
	)
	BeforeEach(func() {
		var err error
		log, err = logger.NewZapLogger(logger.InfoLevel, logger.FormatText)
		if err != nil {
			log.Error(err, "error creating NewZapLogger")
		}
	})

	Describe("can extract data from providerConfig", func() {
		It("when there is not extensions, should not be found", func() {
			extensions := []v1beta1.Extension{}
			found, _, err := actuator.GetProviderConfig(log, extensions)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})
		It("when is not my type, should not be found", func() {
			extensions := []v1beta1.Extension{
				{
					Type: "notMyType",
					ProviderConfig: &runtime.RawExtension{
						Raw: []byte("{}"),
					},
				},
			}
			found, _, err := actuator.GetProviderConfig(log, extensions)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
		})

		It("when is empty, should be empty", func() {
			extensions := []v1beta1.Extension{
				{
					Type: consts.ExtensionType,
					ProviderConfig: &runtime.RawExtension{
						Raw: []byte("{}"),
					},
				},
			}
			found, criResMgrConfig, err := actuator.GetProviderConfig(log, extensions)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(criResMgrConfig).Should(Equal(actuator.CriResMgrConfig{}))
		})

		It("when non empty, then non empty", func() {
			extensions := []v1beta1.Extension{
				{
					Type: consts.ExtensionType,
					ProviderConfig: &runtime.RawExtension{
						Raw: []byte(`{"configs": {"foo":"bar"}}`),
					},
				},
			}
			found, criResMgrConfig, err := actuator.GetProviderConfig(log, extensions)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(criResMgrConfig).Should(
				Equal(
					actuator.CriResMgrConfig{Configs: map[string]string{"foo": "bar"}},
				),
			)
		})

	})
	Describe("rendering charts installation chart with configs", func() {
		nodeSelector := map[string]string{}
		configTypes := map[string]map[string]string{
			// this should generate one ConfigMap with "dynamic" key
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

		It("generate properly with expected bodies inside", func() {
			secret, err := a.GenerateSecretData(log, ctx, consts.Charts, consts.ChartPath, "foo_namespace", "v1.0.0", configTypes, nodeSelector)
			Expect(err).NotTo(HaveOccurred())

			Expect(secret).Should(HaveKey(consts.InstallationSecretKey))

			// check dynamic (first level is unpacked) and rest becomes multi string
			Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring(`name: "cri-resmgr-config.default"`))
			Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring("CONFIG_BODY_OF_DEFAULT: |"))
			Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring(`name: "cri-resmgr-config.nodeFoo"`))
			Expect(string(secret[consts.InstallationSecretKey])).Should(ContainSubstring("CONFIG_BODY_OF_NODEFOO: |"))
		})
	})

	Describe("rendering monitoring chart with GenerateSecretDataToMonitoringManagedResource", func() {
		It("generate correct config with replaced namespace", func() {
			a := actuator.NewActuator("mock").(*actuator.Actuator)

			output := a.GenerateSecretDataToMonitoringManagedResource("test")

			Expect(string(output["data"])).To(ContainSubstring("test"))
			Expect(string(output["data"])).NotTo(ContainSubstring("{{ namespace }}"))
		})
	})

})
