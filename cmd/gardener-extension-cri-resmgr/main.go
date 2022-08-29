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
	"context"
	"fmt"
	"os"
	"time"

	// Gardener
	extensionsconfig "github.com/gardener/gardener/extensions/pkg/apis/config"
	"github.com/gardener/gardener/extensions/pkg/controller"
	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	controllercmd "github.com/gardener/gardener/extensions/pkg/controller/cmd"
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	"github.com/gardener/gardener/extensions/pkg/controller/healthcheck"
	"github.com/gardener/gardener/extensions/pkg/controller/healthcheck/general"
	"github.com/gardener/gardener/extensions/pkg/util"
	gardenercorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	resourcemanagerv1alpha1 "github.com/gardener/gardener/pkg/apis/resources/v1alpha1"
	"github.com/gardener/gardener/pkg/logger"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	managedresources "github.com/gardener/gardener/pkg/utils/managedresources"

	// Other
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	runtimelog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

const (
	ExtensionName = "cri-resmgr"
	ExtensionType = "cri-resmgr-extension"

	ControllerName = "cri-resmgr-controller"
	ActuatorName   = "cri-resmgr-actuator"

	ManagedResourceName = "extension-runtime-cri-resmgr"
	ConfigKey           = "config.yaml"

	ChartPath               = "charts/cri-resmgr-installation/"
	ChartPathRemoval        = "charts/cri-resmgr-removal"
	InstallationImageName   = "installation_image_name" // TODO: to be replaced with proper "gardener-extension-cri-resmgr-" when ready
	InstallationReleaseName = "cri-resmgr-installation"
)

func RegisterHealthChecks(mgr manager.Manager) error {
	defaultSyncPeriod := time.Second * 30
	opts := healthcheck.DefaultAddArgs{
		HealthCheckConfig: extensionsconfig.HealthCheckConfig{SyncPeriod: metav1.Duration{Duration: defaultSyncPeriod}},
	}
	return healthcheck.DefaultRegistration(
		ExtensionType,
		extensionsv1alpha1.SchemeGroupVersion.WithKind(extensionsv1alpha1.ExtensionResource),
		func() client.ObjectList { return &extensionsv1alpha1.ExtensionList{} },
		func() extensionsv1alpha1.Object { return &extensionsv1alpha1.Extension{} },
		mgr,
		opts,
		nil,
		[]healthcheck.ConditionTypeToHealthCheck{
			{
				ConditionType: string(gardenercorev1beta1.ShootSystemComponentsHealthy),
				HealthCheck:   general.CheckManagedResource(ManagedResourceName),
			},
		},
	)
}

type Options struct {
	restOptions       *controllercmd.RESTOptions
	managerOptions    *controllercmd.ManagerOptions
	controllerOptions *controllercmd.ControllerOptions
	reconcileOptions  *controllercmd.ReconcilerOptions
}

// ---------------------------------------------------------------------------------------
// -                                        Main                                         -
// ---------------------------------------------------------------------------------------

func main() {
	runtimelog.SetLogger(logger.ZapLogger(true))

	ctx := signals.SetupSignalHandler()

	options := &Options{
		restOptions: &controllercmd.RESTOptions{},
		managerOptions: &controllercmd.ManagerOptions{
			LeaderElection:          false,
			LeaderElectionID:        controllercmd.LeaderElectionNameID(ExtensionName),
			LeaderElectionNamespace: os.Getenv("LEADER_ELECTION_NAMESPACE"),
		},
		controllerOptions: &controllercmd.ControllerOptions{
			MaxConcurrentReconciles: 1,
		},
		reconcileOptions: &controllercmd.ReconcilerOptions{},
	}

	optionAggregator := controllercmd.NewOptionAggregator(
		options.restOptions,
		options.managerOptions,
		options.controllerOptions,
		options.reconcileOptions,
	)

	cmd := &cobra.Command{
		Use:   "cri-resmgr-controller-manager",
		Short: "CRI Resource manager Controller manages components which install CRI-Resource-Manager as CRI proxy.",

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := optionAggregator.Complete(); err != nil {
				return fmt.Errorf("error completing options: %s", err)
			}

			mgrOpts := options.managerOptions.Completed().Options()
			mgrOpts.MetricsBindAddress = "0"

			// mgrOpts.ClientDisableCacheFor = []client.Object{
			// 	&corev1.Secret{},    // TODO: resolve race condition with small rsync time
			// }

			mgr, err := manager.New(options.restOptions.Completed().Config, mgrOpts)
			if err != nil {
				return fmt.Errorf("could not instantiate controller-manager: %s", err)
			}
			scheme := mgr.GetScheme()
			if err := extensionscontroller.AddToScheme(scheme); err != nil {
				return err
			}
			if err := resourcemanagerv1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
				return err
			}

			// enable healthcheck
			if err := RegisterHealthChecks(mgr); err != nil {
				return err
			}

			// For development purposes.
			ignoreOperationAnnotation := false

			if err := extension.Add(mgr, extension.AddArgs{
				Actuator:                  NewActuator(),
				ControllerOptions:         options.controllerOptions.Completed().Options(),
				Name:                      ControllerName,
				FinalizerSuffix:           ExtensionType,
				Resync:                    60 * time.Minute, // was 60 // FIXME: with 1 second resync we have race condition during deletion
				Predicates:                extension.DefaultPredicates(ignoreOperationAnnotation),
				Type:                      ExtensionType,
				IgnoreOperationAnnotation: ignoreOperationAnnotation,
			}); err != nil {
				return fmt.Errorf("error configuring actuator: %s", err)
			}

			if err := mgr.Start(ctx); err != nil {
				return fmt.Errorf("error running manager: %s", err)
			}

			return nil
		},
	}

	optionAggregator.AddFlags(cmd.Flags())

	if err := cmd.ExecuteContext(ctx); err != nil {
		runtimelog.Log.Error(err, "error executing the main controller command")
		os.Exit(1)
	}

}

// ---------------------------------------------------------------------------------------
// -                                        Actuator                                     -
// ---------------------------------------------------------------------------------------

func NewActuator() extension.Actuator {
	return &actuator{
		chartRendererFactory: extensionscontroller.ChartRendererFactoryFunc(util.NewChartRendererForShoot),
		logger:               log.Log.WithName(ActuatorName),
	}
}

type actuator struct {
	client               client.Client
	config               *rest.Config
	chartRendererFactory extensionscontroller.ChartRendererFactory
	decoder              runtime.Decoder
	logger               logr.Logger
}

func (a *actuator) generateSecretData(ctx context.Context, ex *extensionsv1alpha1.Extension, chartPath string, namespace string) (map[string][]byte, error) {
	emptyMap := map[string][]byte{}
	// Find what shoot cluster we dealing with.
	cluster, err := controller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return emptyMap, err
	}
	// Depending on shoot, chartredner will have different capabilities based on K8s version.
	chartRenderer, err := a.chartRendererFactory.NewChartRendererForShoot(cluster.Shoot.Spec.Kubernetes.Version)
	if err != nil {
		return emptyMap, err
	}
	chartValues := map[string]interface{}{
		"images": map[string]string{
			InstallationImageName: "foo", // TODO: imagevector.FindImage(InstallationImageName),
		},
	}
	release, err := chartRenderer.Render(chartPath, InstallationReleaseName, metav1.NamespaceSystem, chartValues)
	if err != nil {
		return emptyMap, err
	}
	// Put chart into secret
	secretData := map[string][]byte{ConfigKey: release.Manifest()}
	return secretData, nil
}

func (a *actuator) deployDaemonsetToUninstallCriResMgr(ctx context.Context, ex *extensionsv1alpha1.Extension) error {
	a.logger.Info("Uninstalling CRI-Resource-Manager")
	namespace := ex.GetNamespace()
	secretData, err := a.generateSecretData(ctx, ex, ChartPathRemoval, namespace)
	if err != nil {
		return err
	}
	// Reconcile managedresource and secret for shoot.
	if err := managedresources.CreateForShoot(ctx, a.client, namespace, ManagedResourceName, false, secretData); err != nil {
		return err
	}
	// Sleep to give daemonset a time to remove cri-resmgr
	// TODO: detect if the script is finished
	a.logger.Info("Sleep for 120 seconds to make sure remove script is done.")
	time.Sleep(120 * time.Second)
	return nil
}

func (a *actuator) Reconcile(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()
	a.logger.Info("Reconcile: checking extension...") // , "shoot", cluster.Shoot.Name, "namespace", cluster.Shoot.Namespace)

	mr := &resourcemanagerv1alpha1.ManagedResource{}
	if err := a.client.Get(ctx, kutil.Key(namespace, ManagedResourceName), mr); err != nil {
		// Continue only if not found.
		if !apierrors.IsNotFound(err) {
			return err
		}
		secretData, err := a.generateSecretData(ctx, ex, ChartPath, namespace)
		if err != nil {
			return err
		}
		// Reconcile managedresource and secret for shoot.
		if err := managedresources.CreateForShoot(ctx, a.client, namespace, ManagedResourceName, false, secretData); err != nil {
			return err
		}
	} else {
		a.logger.Info("Reconcile: extension is already installed. Ignoring.") //, "shoot", cluster.Shoot.Name, "namespace", cluster.Shoot.Namespace)
	}

	return nil
}

func (a *actuator) Delete(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	namespace := ex.GetNamespace()
	cluster, err := controller.GetCluster(ctx, a.client, namespace)
	if err != nil {
		return err
	}

	if err := a.deployDaemonsetToUninstallCriResMgr(ctx, ex); err != nil {
		return err
	}

	a.logger.Info("Delete: deleting extension managedresources in shoot", "shoot", cluster.Shoot.Name, "namespace", cluster.Shoot.Namespace)

	twoMinutes := 2 * time.Minute

	timeoutShootCtx, cancelShootCtx := context.WithTimeout(ctx, twoMinutes)
	defer cancelShootCtx()

	if err := managedresources.DeleteForShoot(ctx, a.client, namespace, ManagedResourceName); err != nil {
		return err
	}

	if err := managedresources.WaitUntilDeleted(timeoutShootCtx, a.client, namespace, ManagedResourceName); err != nil {
		return err
	}

	return nil
}

func (a *actuator) Restore(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Reconcile(ctx, logger, ex)
}

func (a *actuator) Migrate(ctx context.Context, logger logr.Logger, ex *extensionsv1alpha1.Extension) error {
	return a.Delete(ctx, logger, ex)
}

func (a *actuator) InjectConfig(config *rest.Config) error {
	a.config = config
	return nil
}

func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

func (a *actuator) InjectScheme(scheme *runtime.Scheme) error {
	a.decoder = serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder()
	return nil
}
