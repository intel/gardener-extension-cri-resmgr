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
	"context"
	"os"

	mock_client "github.com/intel/gardener-extension-cri-resmgr/mocks"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/configs"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"

	// Gardener
	"github.com/gardener/gardener/pkg/logger"

	// Other
	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
)

var _ = Describe("cri-resource-manager extension configs reading", func() {

	var (
		baseConfigs     map[string]string
		providerConfigs map[string]string
		log             logr.Logger
	)
	BeforeEach(func() {
		var err error
		log, err = logger.NewZapLogger(logger.InfoLevel, logger.FormatText)
		if err != nil {
			log.Error(err, "error creating NewZapLogger")
		}
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

	Describe("GetBaseConfigsFromConfigMap reads extension ConfigMap and get its configs as baseConfigs", func() {
		var (
			ctx           context.Context
			log           logr.Logger
			mockCtrl      *gomock.Controller
			mockk8sClient *mock_client.MockClient
		)
		BeforeEach(func() {
			ctx = context.TODO()

			var err error
			log, err = logger.NewZapLogger(logger.InfoLevel, logger.FormatText)
			if err != nil {
				log.Error(err, "error creating NewZapLogger")
			}

			mockCtrl = gomock.NewController(GinkgoT())
			mockk8sClient = mock_client.NewMockClient(mockCtrl)
		})
		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("get empty config", func() {

			configs, err := configs.GetBaseConfigsFromConfigMap(ctx, log, mockk8sClient)

			Expect(err).NotTo(HaveOccurred())
			Expect(configs).Should(Equal(map[string]string{}))
		})

		It("get not empty config", func() {
			os.Setenv(consts.ConfigMapNamespaceEnvKey, "smoots")

			mockk8sClient.
				EXPECT().
				Get(gomock.Any(), gomock.Any(), gomock.Any()).
				SetArg(2, corev1.ConfigMap{Data: map[string]string{"test": "test"}}).
				AnyTimes()

			configs, err := configs.GetBaseConfigsFromConfigMap(ctx, log, mockk8sClient)

			Expect(err).NotTo(HaveOccurred())
			Expect(configs).Should(Equal(map[string]string{"test": "test"}))
		})

	})
})
