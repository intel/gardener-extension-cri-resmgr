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

package main

import (
	"os"

	"github.com/intel/gardener-extension-cri-resmgr/cmd/gardener-extension-cri-resmgr/app"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"

	// Gardener
	"github.com/gardener/gardener/pkg/logger"

	// Other
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// ---------------------------------------------------------------------------------------
// -                                        Main                                         -
// ---------------------------------------------------------------------------------------

func main() {
	zapLogger, err := logger.NewZapLogger(logger.InfoLevel, logger.FormatText)
	if err != nil {
		log.Log.Error(err, "error creating NewZapLogger")
		os.Exit(1)
	}
	log.SetLogger(zapLogger)

	ctx := signals.SetupSignalHandler()
	log.Log.Info("Build", "Version", consts.Version, "Commit", consts.Commit)

	cmd := app.NewExtensionControllerCommand(ctx)
	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Log.Error(err, "error executing the main controller command")
		os.Exit(1)
	}
}
