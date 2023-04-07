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

package cri_resmgr_extension

import (
	"context"

	"github.com/gardener/gardener/test/framework"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("cri-resmgr enable tests", ginkgo.Label("enable"), func() {
	f := framework.NewShootCreationFramework(&framework.ShootCreationConfig{
		GardenerConfig: &framework.GardenerConfig{
			ProjectNamespace:   projectNamespace,
			GardenerKubeconfig: kubeconfigPath,
			SkipAccessingShoot: skipAccessingShoot,
			CommonConfig:       commonConfig,
		},
	})
	f.Shoot = getShoot()
	f.Shoot.Name = "e2e-default"

	var a = "get controllerregistrations.core.gardener.cloud cri-resmgr-extension"
	var b = "get controllerdeployments.core.gardener.cloud cri-resmgr-extension"
	var c = "get controllerinstallation.core.gardener.cloud"
	var d = "describe pod -l app.kubernetes.io/name=gardener-extension-cri-resmgr --all-namespaces"
	var e = "describe pod --all-namespaces"

	ginkgo.It("Create Shoot, Enable cri-rm Extension, Delete Shoot", func() {
		kubectl(a)
		kubectl(b)
		kubectl(c)
		kubectl(d)
		kubectl(e)

		ginkgo.By("Create Shoot")
		ctx, cancel := context.WithTimeout(backgroundCtx, fiveteenMinutes)
		defer cancel()
		test := f.CreateShootAndWaitForCreation(ctx, false)
		kubectl(d)
		kubectl(e)
		gomega.Expect(test).To(gomega.Succeed())
		f.Verify()

		ginkgo.By("Enable cri-resmgr extension")
		ctx, cancel = context.WithTimeout(backgroundCtx, fiveteenMinutes)
		defer cancel()
		gomega.Expect(f.UpdateShoot(ctx, f.Shoot, enableCriResmgr)).To(gomega.Succeed())

		kubectl(d)

		ginkgo.By("Delete Shoot")
		ctx, cancel = context.WithTimeout(backgroundCtx, fiveteenMinutes)
		defer cancel()
		gomega.Expect(f.DeleteShootAndWaitForDeletion(ctx, f.Shoot)).To(gomega.Succeed())
	})
})
