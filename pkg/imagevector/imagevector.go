package imagevector

import (
	"github.com/gardener/gardener/pkg/utils/imagevector"
	"github.com/intel/gardener-extension-cri-resmgr/charts"
	"strings"
)

var imageVector imagevector.ImageVector

func init() {
	var imageVector, err = imagevector.Read(strings.NewReader(charts.ImagesYAML))
	if err != nil {
		return
	}

	imageVector, err = imagevector.WithEnvOverride(imageVector)
	if err != nil {
		return
	}
}

// ImageVector contains all images from charts/images.yaml
func ImageVector() imagevector.ImageVector {
	return imageVector
}
