module github.com/intel/gardener-extension-cri-resmgr

go 1.19

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

require (
	github.com/gardener/gardener v1.56.0
	github.com/go-logr/logr v1.2.3
	github.com/onsi/ginkgo/v2 v2.3.0
	github.com/onsi/gomega v1.22.1
	github.com/spf13/cobra v1.4.0
	k8s.io/api v0.25.0
	k8s.io/apimachinery v0.25.0
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed
	sigs.k8s.io/controller-runtime v0.13.0
// github.com/intel/gardener-extension-cri-resmgr/pkg/imagevector v0.0.0-00010101000000-000000000000
// github.com/intel/gardener-extension-cri-resmgr/charts v0.0.0-00010101000000-000000000000 // indirect
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gardener/etcd-druid v0.12.3 // indirect
	github.com/gardener/hvpa-controller/api v0.5.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/zapr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kubernetes-csi/external-snapshotter/v2 v2.1.4 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mholt/archiver v3.1.1+incompatible // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nwaples/rardecode v1.1.2 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.12.2 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.0.0-20220516162934-403b01795ae8 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220502173005-c8bf987b8c21 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	istio.io/api v0.0.0-20220512181135-e8ec1e1d89de // indirect
	istio.io/client-go v1.14.0 // indirect
	k8s.io/apiextensions-apiserver v0.25.0 // indirect
	k8s.io/apiserver v0.25.0 // indirect
	k8s.io/autoscaler/vertical-pod-autoscaler v0.11.0 // indirect
	k8s.io/component-base v0.25.0 // indirect
	k8s.io/helm v2.16.1+incompatible // indirect
	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/kube-aggregator v0.25.0 // indirect
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1 // indirect
	k8s.io/metrics v0.25.0 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.25.0

replace (
	// Required for local development - remove before merging
	github.com/intel/gardener-extension-cri-resmgr/charts => ./charts
	github.com/intel/gardener-extension-cri-resmgr/pkg/imagevector => ./pkg/imagevector
)
