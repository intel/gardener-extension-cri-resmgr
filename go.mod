module github.com/intel/gardener-extension-cri-resmgr

go 1.16

require (
        github.com/gardener/gardener v1.45.0
        github.com/go-logr/logr v1.2.0
        github.com/onsi/ginkgo/v2 v2.1.3
        github.com/onsi/gomega v1.18.0
        github.com/spf13/cobra v1.2.1
        k8s.io/apimachinery v0.23.3
        k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
        k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9
        sigs.k8s.io/controller-runtime v0.11.1
)

replace (
        github.com/gardener/gardener-resource-manager/api => github.com/gardener/gardener-resource-manager/api v0.25.0
        github.com/gardener/hvpa-controller => github.com/gardener/hvpa-controller v0.4.0
        github.com/gardener/hvpa-controller/api => github.com/gardener/hvpa-controller/api v0.4.0
        github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.5.5
        k8s.io/api => k8s.io/api v0.23.3
        k8s.io/apimachinery => k8s.io/apimachinery v0.23.3
        k8s.io/apiserver => k8s.io/apiserver v0.23.3
        k8s.io/autoscaler => k8s.io/autoscaler v0.0.0-20201008123815-1d78814026aa // translates to k8s.io/autoscaler/vertical-pod-autoscaler@v0.9.0
        k8s.io/autoscaler/vertical-pod-autoscaler => k8s.io/autoscaler/vertical-pod-autoscaler v0.9.0
        k8s.io/client-go => k8s.io/client-go v0.23.3
        k8s.io/code-generator => k8s.io/code-generator v0.23.3
        k8s.io/component-base => k8s.io/component-base v0.23.3
        k8s.io/helm => k8s.io/helm v2.13.1+incompatible
)
