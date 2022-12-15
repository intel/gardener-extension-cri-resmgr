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

package configs_test

import (

	// Local

	"github.com/go-logr/logr"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/configs"

	// Gardener

	"github.com/gardener/gardener/pkg/logger"

	// Other
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cri-resource-manager extension configs reading", func() {

	var (
		baseConfigs     map[string]string
		providerConfigs map[string]string
		log             logr.Logger
	)
	BeforeEach(func() {
		log = logger.ZapLogger(true)
	})
	It("non configs provided", func() {
		configs, err := configs.PrepareConfigTypes(log, map[string]string{}, map[string]string{})
		Expect(err).NotTo(HaveOccurred())
		Expect(configs).Should(Equal(map[string]map[string]string{"static": {}, "dynamic": {}}))
	})

	Describe("with not empty configs with all dynamics types", func() {
		It("but just configs provided from shoot", func() {
			baseConfigs = map[string]string{}
			providerConfigs = map[string]string{"foo": "bar"}
			configs, err := configs.PrepareConfigTypes(log, baseConfigs, providerConfigs)
			Expect(err).NotTo(HaveOccurred())
			Expect(configs).Should(Equal(map[string]map[string]string{"dynamic": {"foo": "bar"}, "static": {}}))
		})
		Describe("but just configs from configmap (baseConfigs)", func() {
			It("result in just dynamic from baseConfigs", func() {
				baseConfigs = map[string]string{"bar": "baz"}
				providerConfigs = map[string]string{}
				configs, err := configs.PrepareConfigTypes(log, baseConfigs, providerConfigs)
				Expect(err).NotTo(HaveOccurred())
				Expect(configs).Should(Equal(map[string]map[string]string{"dynamic": {"bar": "baz"}, "static": {}}))
			})
		})
		It("with some baseConfigs from configMap and some providerConfigs", func() {
			baseConfigs = map[string]string{"foo": "old", "bar": "baz"}
			providerConfigs = map[string]string{"foo": "new"}
			configs, err := configs.PrepareConfigTypes(log, baseConfigs, providerConfigs)
			Expect(err).NotTo(HaveOccurred())
			Expect(configs).Should(Equal(map[string]map[string]string{"dynamic": {"foo": "new", "bar": "baz"}, "static": {}}))
		})
		It("with some baseConfigs from configMap and some providerConfigs but types of static", func() {
			baseConfigs = map[string]string{"fallback": "old", "EXTRA_OPTIONS": "baz"}
			providerConfigs = map[string]string{"fallback": "new", "force": "force"}
			configs, err := configs.PrepareConfigTypes(log, baseConfigs, providerConfigs)
			Expect(err).NotTo(HaveOccurred())
			Expect(configs).Should(Equal(map[string]map[string]string{"dynamic": {}, "static": {"fallback": "new", "EXTRA_OPTIONS": "baz", "force": "force"}}))
		})
	})
})
