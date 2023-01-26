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
	"context"
	"fmt"
	"os"

	// Local

	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/controller/healthcheck"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/controller/lifecycle"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/options"

	// Gardener
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	controllercmd "github.com/gardener/gardener/extensions/pkg/controller/cmd"
	"github.com/gardener/gardener/extensions/pkg/controller/heartbeat"
	heartbeatcmd "github.com/gardener/gardener/extensions/pkg/controller/heartbeat/cmd"
	resourcemanagerv1alpha1 "github.com/gardener/gardener/pkg/apis/resources/v1alpha1"
	"github.com/gardener/gardener/pkg/logger"

	// Other
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewExtensionControllerCommand(ctx context.Context) *cobra.Command {

	options := options.NewOptions()
	restOpts := &controllercmd.RESTOptions{}
	reconcileOpts := &controllercmd.ReconcilerOptions{
		IgnoreOperationAnnotation: true,
	}
	mgrOpts := &controllercmd.ManagerOptions{
		LeaderElection:     false,
		MetricsBindAddress: "0",
		LogLevel:           logger.InfoLevel,
		LogFormat:          logger.FormatText,
	}

	healthCheckOpts := &controllercmd.ControllerOptions{
		MaxConcurrentReconciles: 5,
	}

	heartbeatOpts := &heartbeatcmd.Options{
		ExtensionName:        consts.ExtensionName,
		RenewIntervalSeconds: 30,
		Namespace:            os.Getenv("EXTENSION_CONFIGMAP_NAMESPACE"),
	}

	allOpts := controllercmd.NewOptionAggregator(
		mgrOpts,
		healthCheckOpts,
		heartbeatOpts,
		restOpts,
		reconcileOpts,
	)

	cmd := &cobra.Command{
		Use:   "cri-resmgr-controller-manager",
		Short: "CRI Resource manager Controller manages components which install CRI-Resource-Manager as CRI proxy.",

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.OptionAggregator.Complete(); err != nil {
				return fmt.Errorf("error completing options: %s", err)
			}

			if err := allOpts.Complete(); err != nil {
				return err
			}

			if err := heartbeatOpts.Validate(); err != nil {
				return err
			}

			mgr, err := manager.New(restOpts.Completed().Config, mgrOpts.Completed().Options())
			if err != nil {
				return fmt.Errorf("could not instantiate controller-manager: %s", err)
			}
			scheme := mgr.GetScheme()
			if err := extensionscontroller.AddToScheme(scheme); err != nil {
				return err
			}
			if err := resourcemanagerv1alpha1.AddToScheme(scheme); err != nil {
				return err
			}

			// mgrOpts.ClientDisableCacheFor = []client.Object{
			// 	&corev1.ConfigMap{}, // applied for ManagedResources
			// }
			heartbeatOpts.Completed().Apply(&heartbeat.DefaultAddOptions)
			// Enable healthcheck.
			// "Registration" adds additional controller that watches over Extension/Cluster.
			// TODO: ENABLE before merging!!!
			if err := healthcheck.RegisterHealthChecks(mgr); err != nil {
				return err
			}

			if err := heartbeat.AddToManager(mgr); err != nil {
				return err
			}

			ignoreOperationAnnotation := reconcileOpts.Completed().IgnoreOperationAnnotation
			// if true:
			//		predicates: only observe "generation change" predicate (oldObject.generation != newObject.generation)
			// 		watches:  watch Cluster (additionally and map to extensions) and Extension
			//
			// if false (default):
			//      predicates: (defaultControllerPredicates) watches for "operation annotation" to be reconcile/migrate/restore
			//					or deletionTimestamp is set or lastOperation is not successful state (on Extension object)
			// 		watches: only Extension
			log.Log.Info("Reconciler options", "ignoreOperationAnnotation", ignoreOperationAnnotation)

			// I. This is the primary controller that watches over Extension (and possible Cluster based on ignoreOperationAnnotation)
			if err := lifecycle.AddToManager(mgr, options, ignoreOperationAnnotation); err != nil {
				return fmt.Errorf("error configuring controller with extensions actuator: %s", err)
			}
			// II. Create another controller for watching over specific configMap and
			// reconciling all Extensions that all only in Succeeded state to prevent race over Extension reconciliation
			if err := lifecycle.AddConfigMapWatchingControllerToManager(mgr, options); err != nil {
				return fmt.Errorf("error configuring configMap controller with extensions actuator: %s", err)
			}

			if err := mgr.Start(ctx); err != nil {
				return fmt.Errorf("error running manager: %s", err)
			}

			return nil
		},
	}

	allOpts.AddFlags(cmd.Flags())

	return cmd
}
