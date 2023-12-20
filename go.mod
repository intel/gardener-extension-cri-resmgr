module github.com/intel/gardener-extension-cri-resmgr

go 1.21

// TO BE REMOVED when concluded unnessesary ! :)
// github.com/gardener/gardener-resource-manager/api => github.com/gardener/gardener-resource-manager/api v0.25.0
// github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.12.1 // keep this value in sync with sigs.k8s.io/controller-runtime
// k8s.io/api => k8s.io/api v0.24.3
// k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.24.3
// k8s.io/apimachinery => k8s.io/apimachinery v0.24.3
// k8s.io/apiserver => k8s.io/apiserver v0.24.3
// k8s.io/autoscaler => k8s.io/autoscaler v0.0.0-20220531185024-cc90d57b7fe1 // translates to k8s.io/autoscaler/vertical-pod-autoscaler@v0.11.0
// k8s.io/autoscaler/vertical-pod-autoscaler => k8s.io/autoscaler/vertical-pod-autoscaler v0.11.0
// k8s.io/client-go => k8s.io/client-go v0.24.3
// k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.24.3
// k8s.io/code-generator => k8s.io/code-generator v0.24.3
// k8s.io/component-base => k8s.io/component-base v0.24.3
// k8s.io/helm => k8s.io/helm v2.16.1+incompatible
// k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.24.3
// sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.12.1
// github.com/go-logr/logr v1.2.4
// github.com/golang/mock v1.6.0
// github.com/onsi/ginkgo/v2 v2.13.0
// github.com/onsi/gomega v1.29.0
// github.com/spf13/cobra v1.7.0
// k8s.io/api v0.28.3
// k8s.io/apimachinery v0.28.3
// k8s.io/client-go v0.28.3
// k8s.io/klog/v2 v2.100.1
// k8s.io/utils v0.0.0-20230505201702-9f6742963106
// sigs.k8s.io/controller-runtime v0.16.3
// github.com/intel/gardener-extension-cri-resmgr/pkg/imagevector v0.0.0-00010101000000-000000000000
// github.com/intel/gardener-extension-cri-resmgr/charts v0.0.0-00010101000000-000000000000 // indirect

require (
	github.com/gardener/gardener v1.86.0
	github.com/go-logr/logr v1.2.4
	github.com/onsi/ginkgo/v2 v2.13.0
	github.com/onsi/gomega v1.29.0
	github.com/spf13/cobra v1.7.0
	go.uber.org/mock v0.3.0
	k8s.io/api v0.28.3
	k8s.io/apimachinery v0.28.3
	k8s.io/client-go v0.28.3
	k8s.io/klog/v2 v2.100.1
	k8s.io/utils v0.0.0-20230406110748-d93618cff8a2
	sigs.k8s.io/controller-runtime v0.16.3
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fluent/fluent-operator/v2 v2.2.0 // indirect
	github.com/frankban/quicktest v1.14.5 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gardener/etcd-druid v0.21.0 // indirect
	github.com/gardener/hvpa-controller/api v0.5.0 // indirect
	github.com/gardener/machine-controller-manager v0.50.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/zapr v1.2.4 // indirect
	github.com/go-openapi/errors v0.20.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.8.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/term v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230526161137-0005af68ea54 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230525234035-dd9d682886f9 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	istio.io/api v1.19.2-0.20231011000955-f3015ebb5bd4 // indirect
	istio.io/client-go v1.19.3 // indirect
	k8s.io/apiextensions-apiserver v0.28.3 // indirect
	k8s.io/autoscaler/vertical-pod-autoscaler v1.0.0 // indirect
	k8s.io/component-base v0.28.3 // indirect
	k8s.io/helm v2.17.0+incompatible // indirect
	k8s.io/kube-aggregator v0.28.3 // indirect
	k8s.io/kube-openapi v0.0.0-20230717233707-2695361300d9 // indirect
	k8s.io/kubelet v0.28.3 // indirect
	k8s.io/metrics v0.28.3 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.3.0 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	// Required for local development - remove before merging
	github.com/intel/gardener-extension-cri-resmgr/charts => ./charts
	github.com/intel/gardener-extension-cri-resmgr/pkg/imagevector => ./pkg/imagevector
)
