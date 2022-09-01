package imagevector

import (
	"strings"

	"github.com/gardener/gardener-extension-cri-resmgr/charts"
	"github.com/gardener/gardener/pkg/utils/imagevector"
	"k8s.io/apimachinery/pkg/util/runtime"
)

var imageVector imagevector.ImageVector

func init() {
	var imageVector, err = imagevector.Read(strings.NewReader(charts.ImagesYAML))
	if err != nil {
		return err
	}

	imageVector, err = imagevector.WithEnvOverride(imageVector)
	if err != nil {
		return err
	}
}

// ImageVector contains all images from charts/images.yaml
func ImageVector() imagevector.ImageVector {
	return imageVector
}
