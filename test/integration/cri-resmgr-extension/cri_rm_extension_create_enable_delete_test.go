package cri_resmgr_extension

import (
	"context"

	"github.com/gardener/gardener/test/framework"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)


var _ = ginkgo.Describe("cri-rm Extension Tests", ginkgo.Label("CRI-RM"), func() {
	f := framework.NewShootCreationFramework(&framework.ShootCreationConfig{
		GardenerConfig: &framework.GardenerConfig{
			ProjectNamespace:   projectNamespace,
			GardenerKubeconfig: kubeconfigPath,
			SkipAccessingShoot: true,
			CommonConfig:       &framework.CommonConfig{},
		},
	})
	f.Shoot = getShoot()
	ginkgo.It("Create Shoot, Enable cri-rm Extension, Delete Shoot", ginkgo.Label("good-case"), func() {
		ginkgo.By("Create Shoot")
		ctx, cancel := context.WithTimeout(backgroundCtx, fiveteenMinutes)
		defer cancel()
		gomega.Expect(f.CreateShootAndWaitForCreation(ctx, false)).To(gomega.Succeed())
		f.Verify()

		ginkgo.By("Enable cri-resmgr extension")
		ctx, cancel = context.WithTimeout(backgroundCtx, fiveteenMinutes)
		defer cancel()
		gomega.Expect(f.UpdateShoot(ctx, f.Shoot, enableCriResmgr)).To(gomega.Succeed())

		ginkgo.By("Delete Shoot")
		ctx, cancel = context.WithTimeout(backgroundCtx, fiveteenMinutes)
		defer cancel()
		gomega.Expect(f.DeleteShootAndWaitForDeletion(ctx, f.Shoot)).To(gomega.Succeed())
	})
})
