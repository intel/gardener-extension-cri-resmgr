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

package app

import (
	controllercmd "github.com/gardener/gardener/extensions/pkg/controller/cmd"
)

type Options struct {
	restOptions       *controllercmd.RESTOptions       // kubeconfig / MasterURL
	controllerOptions *controllercmd.ControllerOptions // MaxConcurrentReconciles
	reconcileOptions  *controllercmd.ReconcilerOptions // IgnoreOperationAnnotation
	optionAggregator  controllercmd.OptionAggregator
}

func NewOptions() *Options {

	options := &Options{
		restOptions: &controllercmd.RESTOptions{},
		controllerOptions: &controllercmd.ControllerOptions{
			MaxConcurrentReconciles: 1,
		},
		reconcileOptions: &controllercmd.ReconcilerOptions{},
	}

	options.optionAggregator = controllercmd.NewOptionAggregator(
		options.restOptions,
		// options.managerOptions, // disabled until leader/webhooks or metrics/healthchecks are required to configure
		options.controllerOptions,
		options.reconcileOptions,
	)
	return options
}
