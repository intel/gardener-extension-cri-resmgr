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

package lifecycle

import (
	"context"
	"fmt"
	"time"

	// Local
	"github.com/go-logr/logr"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/consts"
	"github.com/intel/gardener-extension-cri-resmgr/pkg/options"

	// Gardener
	"github.com/gardener/gardener/extensions/pkg/controller/extension"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/controllerutils/mapper"

	// Other
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// configMapToAllExtensionsMapper maps creates reconciliation requests for extensions based on dedicate configMap of cri-resmgr extension.
func configMapToAllExtensionMapper(ctx context.Context, log logr.Logger, reader client.Reader, obj client.Object) []reconcile.Request {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	configMap, ok := obj.(*corev1.ConfigMap)
	if !ok {
		log.Info("WARNING: expected to get ConfigMap but got something different in configMapMapper to Extension!", "module", "configs")
		return nil
	}

	extensionList := &extensionsv1alpha1.ExtensionList{}
	if err := reader.List(ctx, extensionList); err != nil {
		log.Info("WARNING: can not read list of Extension from Kubernetes", "error", err, "module", "configs")
		return nil
	}

	var requests []reconcile.Request
	extensionsFound := []string{}
	for _, extension := range extensionList.Items {
		if extension.Spec.Type == consts.ExtensionType {
			isOk := false
			// Assume, there is only one condition and it is is Ok "True", then add this extension to requests for reconciliation
			for _, condition := range extension.Status.Conditions {
				isOk = (condition.Status == "True")
				break
			}
			if isOk {
				requests = append(requests, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: extension.Namespace,
						Name:      extension.Name,
					},
				})
				extensionsFound = append(extensionsFound, fmt.Sprintf("%s/%s", extension.Namespace, extension.Name))
			}
		}
	}

	log.Info("found configMap so start reconciliation of all healthy extensions", "module", "configs", "configMap", configMap, "extensions", extensionsFound)
	return requests
}

// AddToManager creates controller that watches Extension object and deploys necessary objects to Shoot cluster.
func AddToManager(mgr manager.Manager, options *options.Options, ignoreOperationAnnotation bool) error {

	return extension.Add(mgr, extension.AddArgs{
		Actuator:                  NewActuator(consts.ActuatorName),
		ControllerOptions:         options.ControllerOptions.Completed().Options(),
		Name:                      consts.ControllerName,
		FinalizerSuffix:           consts.ExtensionType,
		Resync:                    60 * time.Minute,
		Type:                      consts.ExtensionType, // to be used for TypePredicate
		Predicates:                extension.DefaultPredicates(ignoreOperationAnnotation),
		IgnoreOperationAnnotation: ignoreOperationAnnotation,
	})
}

// AddConfigMapWatchingControllerToManager creates controller that watches cri-resmgr-extension ConfigMap object and reconciles everything on Shoot clusters.
func AddConfigMapWatchingControllerToManager(mgr manager.Manager, options *options.Options) error {

	// Create another instance of options - this time for "configMap2Extensions reconciler"
	controllerOptions := options.ControllerOptions.Completed().Options()
	configReconcilerArgs := extension.AddArgs{
		Actuator:        NewActuator(consts.ActuatorName + consts.ConfigsSuffix),
		Resync:          60 * time.Minute,
		FinalizerSuffix: consts.ExtensionType, // We're using the same finalizer as the original controller on purpose to "delete" only once without a need to wait for another "configs" controller
	}
	controllerOptions.Reconciler = extension.NewReconciler(configReconcilerArgs)
	controllerOptions.RecoverPanic = true

	controllerName := consts.ControllerName + consts.ConfigsSuffix
	ctrl, err := controller.New(controllerName, mgr, controllerOptions)
	if err != nil {
		return err
	}

	// only Watch for configMaps with properName and where its resourceVersionChanges
	matchingLabelSelectorPredicate, err := predicate.LabelSelectorPredicate(
		metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/name":              "gardener-extension-cri-resmgr",
				"resources.gardener.cloud/managed-by": "gardener",
			},
		},
	)
	if err != nil {
		return err
	}

	// Predicates to watch over my configMap
	predicates := []predicate.Predicate{matchingLabelSelectorPredicate, predicate.ResourceVersionChangedPredicate{}}

	return ctrl.Watch(
		&source.Kind{Type: &corev1.ConfigMap{}},
		mapper.EnqueueRequestsFrom(
			mapper.MapFunc(configMapToAllExtensionMapper),
			mapper.UpdateWithNew,
			mgr.GetLogger().WithName(controllerName),
		),
		predicates...,
	)
}
