// Copyright 2022 Intel Corporation. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package imagevector

import (
	"fmt"
	"strings"

	"github.com/gardener/gardener/pkg/utils/imagevector"
	"github.com/intel/gardener-extension-cri-resmgr/charts"
)

var imageVector imagevector.ImageVector

func init() {
	var err error
	imageVector, err = imagevector.Read(strings.NewReader(charts.ImagesYAML))
	if err != nil {
		err_message := fmt.Errorf("cannot read images.yaml: %s", err)
		fmt.Println("An error message:", err_message)
	}

	imageVector, err = imagevector.WithEnvOverride(imageVector)
	if err != nil {
		err_message := fmt.Errorf("could not override or read environment variable: %s", err)
		fmt.Println("An error message:", err_message)
		return
	}
}

// ImageVector contains all images from charts/images.yaml
func ImageVector() imagevector.ImageVector {
	return imageVector
}
