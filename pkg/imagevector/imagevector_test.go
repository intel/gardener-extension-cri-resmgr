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

package imagevector_test

import (
	. "github.com/onsi/ginkgo/v2"
	
	. "github.com/onsi/gomega"

	"github.com/intel/gardener-extension-cri-resmgr/pkg/imagevector"
)

var _ = Describe("cri-resource-manager imagevector test", func() {

	It("imagevector has images", func() {

		imageVector := imagevector.ImageVector()

		Expect(imageVector).To(ContainElement(HaveField("Name", "gardener-extension-cri-resmgr-agent")))
		Expect(imageVector).To(ContainElement(HaveField("Name", "gardener-extension-cri-resmgr-installation")))

	})

})
