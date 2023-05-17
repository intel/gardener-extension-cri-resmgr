package lifecycle_test

import (
	"bytes"
	"context"

	"github.com/intel/gardener-extension-cri-resmgr/mocks"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/controller/lifecycle"

	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigMapToAllExtensionMapper tests", func() {
	var (
		ctx      context.Context
		log      logr.Logger
		reader   *mocks.MockReader // client.Reader
		requests []reconcile.Request
		mockCtrl *gomock.Controller
		buffer   *bytes.Buffer
	)

	BeforeEach(func() {
		ctx = context.TODO()

		buffer = bytes.NewBuffer(nil)
		klog.SetOutput(buffer)
		klog.LogToStderr(false)
		log = klogr.New()

		mockCtrl = gomock.NewController(GinkgoT())
		reader = mocks.NewMockReader(mockCtrl)
	})

	It("non Extension objects", func() {
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-configMap",
				Namespace: "test-namespace",
			},
		}
		reader.
			EXPECT().
			List(gomock.Any(), gomock.Any()).
			Return(nil)

		requests = lifecycle.ConfigMapToAllExtensionMapper(ctx, log, reader, configMap)
		Expect(requests).Should(BeEmpty())
	})

	It("some Extension objects", func() {
		configMap := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-configMap",
				Namespace: "test-namespace",
			},
		}
		extensionList := extensionsv1alpha1.ExtensionList{
			Items: []extensionsv1alpha1.Extension{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-extension-1",
						Namespace: "test-namespace",
					},
					Spec: extensionsv1alpha1.ExtensionSpec{
						DefaultSpec: extensionsv1alpha1.DefaultSpec{
							Type: consts.ExtensionType,
						},
					},
					Status: extensionsv1alpha1.ExtensionStatus{
						DefaultStatus: extensionsv1alpha1.DefaultStatus{
							LastOperation: &gardencorev1beta1.LastOperation{
								Description:    "Processing",
								LastUpdateTime: metav1.Time{},
								Progress:       0,
								State:          "Processing",
								Type:           "",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-extension-2",
						Namespace: "test-namespace",
					},
					Spec: extensionsv1alpha1.ExtensionSpec{
						DefaultSpec: extensionsv1alpha1.DefaultSpec{
							Type: consts.ExtensionType,
						},
					},
					Status: extensionsv1alpha1.ExtensionStatus{
						DefaultStatus: extensionsv1alpha1.DefaultStatus{
							LastOperation: &gardencorev1beta1.LastOperation{
								Description:    "Succeeded",
								LastUpdateTime: metav1.Time{},
								Progress:       0,
								State:          "Succeeded",
								Type:           "",
							},
						},
					},
				},
			},
		}
		reader.
			EXPECT().
			List(gomock.Any(), gomock.Any()).
			SetArg(1, extensionList).
			AnyTimes()

		requests = lifecycle.ConfigMapToAllExtensionMapper(ctx, log, reader, configMap)
		Expect(requests).Should(HaveLen(1))
		Expect(buffer.String()).To(ContainSubstring("ignore extension"))
		Expect(buffer.String()).To(ContainSubstring("configs: configMap changed so reconcile following extensions"))
	})
},
)
