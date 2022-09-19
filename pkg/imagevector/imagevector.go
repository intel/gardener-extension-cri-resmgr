package imagevector

import (
	"strings"

	"github.com/gardener/gardener/pkg/utils/imagevector"
	"github.com/intel/gardener-extension-cri-resmgr/charts"
)

var imageVector imagevector.ImageVector

func init() {
	var err error
	imageVector, err = imagevector.Read(strings.NewReader(charts.ImagesYAML))
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
