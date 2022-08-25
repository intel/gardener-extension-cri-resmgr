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
			ExistingShootName:  "first", // "Name of an existing shoot to use instead of creating a new one."
			ProjectNamespace:   projectNamespace,
			GardenerKubeconfig: kubeconfigPath, // KUBECONFIG
			SkipAccessingShoot: false,          // if set to true then the test does not try to access the shoot via its kubeconfig
			CommonConfig:       &framework.CommonConfig{LogLevel: "debug"},
		},
	})
	f.Shoot = getShoot()

	ginkgo.It("Create Shoot, Enable cri-rm Extension, Disable cri-rm Extension, Delete Shoot", ginkgo.Label("good-case"), func() {
		ginkgo.By("Create Shoot")
		ctx, cancel := context.WithTimeout(backgroundCtx, fiveteenMinutes)
		defer cancel()
		gomega.Expect(f.CreateShootAndWaitForCreation(ctx, false)).To(gomega.Succeed())
		f.Verify()

		ginkgo.By("Enable cri-resmgr extension")
		ctx, cancel = context.WithTimeout(backgroundCtx, fiveteenMinutes)
		defer cancel()
		gomega.Expect(f.UpdateShoot(ctx, f.Shoot, enableCriResmgr)).To(gomega.Succeed())

		ginkgo.By("Disable cri-resmgr extension")
		ctx, cancel = context.WithTimeout(backgroundCtx, fiveteenMinutes)
		defer cancel()
		gomega.Expect(f.UpdateShoot(ctx, f.Shoot, disableCriResmgr)).To(gomega.Succeed())

		ginkgo.By("Delete Shoot")
		ctx, cancel = context.WithTimeout(backgroundCtx, fiveteenMinutes)
		defer cancel()
		gomega.Expect(f.DeleteShootAndWaitForDeletion(ctx, f.Shoot)).To(gomega.Succeed())
	})

})
