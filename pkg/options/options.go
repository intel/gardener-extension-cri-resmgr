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

package options

import (
	controllercmd "github.com/gardener/gardener/extensions/pkg/controller/cmd"
)

type Options struct {
	RestOptions       *controllercmd.RESTOptions       // kubeconfig / MasterURL
	ControllerOptions *controllercmd.ControllerOptions // MaxConcurrentReconciles
	ReconcileOptions  *controllercmd.ReconcilerOptions // IgnoreOperationAnnotation
	OptionAggregator  controllercmd.OptionAggregator
}

func NewOptions() *Options {

	options := &Options{
		RestOptions: &controllercmd.RESTOptions{},
		ControllerOptions: &controllercmd.ControllerOptions{
			MaxConcurrentReconciles: 1,
		},
		ReconcileOptions: &controllercmd.ReconcilerOptions{},
	}

	options.OptionAggregator = controllercmd.NewOptionAggregator(
		options.RestOptions,
		options.ControllerOptions,
		options.ReconcileOptions,
	)
	return options
}
