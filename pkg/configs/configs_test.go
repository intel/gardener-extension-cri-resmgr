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
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"

	// Gardener
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/logger"

	// Other
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("cri-resource-manager extension configs reading", func() {

	var (
		extensions  []v1beta1.Extension
		baseConfigs map[string]string
		log         logr.Logger
	)
	BeforeEach(func() {
		log = logger.ZapLogger(true)
	})
	It("installation chart with zero configs provided", func() {
		extensions := []v1beta1.Extension{}
		configs, err := configs.MergeConfigs(log, map[string]string{}, extensions)
		Expect(err).NotTo(HaveOccurred())
		Expect(configs).Should(Equal(map[string]string{}))
	})

	Describe("with empty extensions and empty baseConfigs", func() {
		BeforeEach(func() {

			extensions = []v1beta1.Extension{
				{
					Type: consts.ExtensionType,
					ProviderConfig: &runtime.RawExtension{
						Raw: []byte("{}"),
					},
				},
			}

			baseConfigs = map[string]string{}
		})
		It("installation chart with just configs provided from shoot", func() {
			configs, err := configs.MergeConfigs(log, baseConfigs, extensions)
			Expect(err).NotTo(HaveOccurred())
			Expect(configs).Should(Equal(map[string]string{}))
		})
	})

	Describe("with not empty extensions", func() {
		BeforeEach(func() {
			extensions = []v1beta1.Extension{
				{
					Type: consts.ExtensionType,
					ProviderConfig: &runtime.RawExtension{
						Raw: []byte(`{"configs": {"foo":"bar"}}`),
					},
				},
			}
			baseConfigs = map[string]string{}
		})
		It("installation chart with just configs provided from shoot", func() {
			configs, err := configs.MergeConfigs(log, baseConfigs, extensions)
			Expect(err).NotTo(HaveOccurred())
			Expect(configs).Should(Equal(map[string]string{"foo": "bar"}))
		})
		Describe("with some baseConfigs, merge but do not override", func() {
			It("installation chart with just configs provided from shoot", func() {
				configs, err := configs.MergeConfigs(log, map[string]string{"bar": "baz\n"}, extensions)
				Expect(err).NotTo(HaveOccurred())
				Expect(configs).Should(Equal(map[string]string{"foo": "bar", "bar": "baz\n"}))
			})
		})
		Describe("with some baseConfigs, merge and override", func() {
			It("installation chart with just configs provided from shoot", func() {
				configs, err := configs.MergeConfigs(log, map[string]string{"foo": "old"}, extensions)
				Expect(err).NotTo(HaveOccurred())
				Expect(configs).Should(Equal(map[string]string{"foo": "bar"}))
			})
		})
	})
})
