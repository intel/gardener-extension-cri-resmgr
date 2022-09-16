package main_test

import (
	"context"

	main "github.com/intel/gardener-extension-cri-resmgr/cmd/gardener-extension-cri-resmgr"

	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	// unused but usefull for future
	// "github.com/gardener/gardener/pkg/chartrenderer"
	// "github.com/gardener/gardener/pkg/chartrenderer"
	// "github.com/golang/mock/gomock"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/helm/pkg/manifest"
)

var _ = Describe("cri-resource-manager extension main tests", func() {
	Describe("rendering charts", func() {

		It("installation chart without configs", func() {
			configs := map[string]string{}
			a := main.NewActuator().(*main.Actuator)
			ctx := context.TODO()

			ex := &extensionsv1alpha1.Extension{}
			secret, err := a.GenerateSecretData(ctx, ex, main.ChartPath, "foo_namespace", "v1.0.0", configs)
			Expect(err).NotTo(HaveOccurred())

			Expect(secret).Should(HaveKey(main.InstallationSecretKey))
			Expect(string(secret[main.InstallationSecretKey])).Should(ContainElement("foo_namespace"))
		})
	})
})
