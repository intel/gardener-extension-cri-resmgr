package main_test

import (
	"context"

	main "github.com/intel/gardener-extension-cri-resmgr/cmd/gardener-extension-cri-resmgr"

	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

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
			a := main.NewActuator().(*main.Actuator)
			ctx := context.TODO()

			ex := &extensionsv1alpha1.Extension{}
			secret, err := a.GenerateSecretData(ctx, ex, main.ChartPath, "foo_namespace", "v1.0.0", configs)
			Expect(err).NotTo(HaveOccurred())

			Expect(secret).Should(HaveKey(main.InstallationSecretKey))

			Expect(string(secret[main.InstallationSecretKey])).Should(ContainSubstring(`name: "cri-resmgr-config.default"`))
			Expect(string(secret[main.InstallationSecretKey])).Should(ContainSubstring("policy:  THIS_WILL_CONFIG_BODY_OF_DEFAULT"))
			Expect(string(secret[main.InstallationSecretKey])).Should(ContainSubstring(`name: "cri-resmgr-config.nodeFoo"`))
			Expect(string(secret[main.InstallationSecretKey])).Should(ContainSubstring("policy:  THIS_WILL_CONFIG_BODY_OF_NODEFOO"))
		})
	})
})
