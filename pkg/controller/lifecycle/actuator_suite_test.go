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

package lifecycle_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

func TestCharts(t *testing.T) {
	// because we output very large charts
	format.MaxLength = 0
	RegisterFailHandler(Fail)
	// Go up three times pkg/controller/lifecycle to project root directory
	err := os.Chdir("../../..")
	if err != nil {
		fmt.Printf("ERROR: Cannot access requested directory: %s", err)
		os.Exit(1)
	}
	RunSpecs(t, "CRI-resource-manager extension test suite")
}
