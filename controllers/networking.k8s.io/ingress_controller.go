/*
Copyright 2021 xzzpig.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package networkingk8sio

import (
	"context"
	"fmt"

	v1 "github.com/xzzpig/alidns-ingress-updater/apis/network/v1"
	"github.com/xzzpig/alidns-ingress-updater/utils"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const finalizerName = "account.xzzpig.com/finalizer"

// IngressReconciler reconciles a Ingress object
type IngressReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Eventer record.EventRecorder
}

//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/finalizers,verbs=update
//+kubebuilder:rbac:groups=network.xzzpig.com,resources=alidnsaccounts,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Ingress object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var ingress netv1.Ingress
	if err := r.Get(ctx, req.NamespacedName, &ingress); err != nil {
		log.V(1).Info("unable to fetch Ingress")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var accounts v1.AliDnsAccountList
	if err := r.List(ctx, &accounts); err != nil {
		log.Error(err, "unable to list all AliDnsAccount")
		return ctrl.Result{}, err
	}

	for _, account := range accounts.Items {
		accountFullName := account.Namespace + "/" + account.Name
		publicIp := account.Status.LastIp
		if ingress.DeletionTimestamp.IsZero() {
			utils.HandleUpdate(&account, &ingress, publicIp, utils.Handler{
				InfoHandler: func(logMsg string, eventMsg string, keysAndValues ...interface{}) {
					if eventMsg == "" {
						eventMsg = logMsg
					}
					log.WithValues("account", accountFullName).Info(logMsg, keysAndValues...)
					r.Eventer.Event(&account, "Normal", "INFO", eventMsg)
				},
				ErrorHandler: func(err error, logMsg string, eventMsg string, keysAndValues ...interface{}) {
					if eventMsg == "" {
						eventMsg = logMsg
					}
					log.WithValues("account", accountFullName).Error(err, logMsg, keysAndValues...)
					r.Eventer.Event(&account, "Warning", "ERROR", eventMsg)
				},
				DebugHandler: func(logMsg string, keysAndValues ...interface{}) {
					log.WithValues("account", accountFullName).V(1).Info(logMsg, keysAndValues...)
				},
				SuccessHandler: func(host string) {
					if !utils.ContainsString(ingress.Finalizers, finalizerName) {
						ingress.Finalizers = append(ingress.Finalizers, finalizerName)
						if err := r.Update(ctx, &ingress); err != nil {
							log.Error(err, "unable to add Ingress Finalizers")
						} else {
							log.Info("add Finalizers to Ingress")
							r.Eventer.Event(&account, "Normal", "INFO", fmt.Sprintf("add Finalizers to Ingress %v", req.NamespacedName))
						}
					}
				},
			})
		} else {
			utils.HandleDelete(&account, &ingress, utils.Handler{
				InfoHandler: func(logMsg string, eventMsg string, keysAndValues ...interface{}) {
					if eventMsg == "" {
						eventMsg = logMsg
					}
					log.WithValues("account", accountFullName).Info(logMsg, keysAndValues...)
					r.Eventer.Event(&account, "Normal", "INFO", eventMsg)
				},
				ErrorHandler: func(err error, logMsg string, eventMsg string, keysAndValues ...interface{}) {
					if eventMsg == "" {
						eventMsg = logMsg
					}
					log.WithValues("account", accountFullName).Error(err, logMsg, keysAndValues...)
					r.Eventer.Event(&account, "Warning", "ERROR", eventMsg)
				},
				DebugHandler: func(logMsg string, keysAndValues ...interface{}) {
					log.WithValues("account", accountFullName).V(1).Info(logMsg, keysAndValues...)
				},
			})
			if utils.ContainsString(ingress.Finalizers, finalizerName) {
				ingress.Finalizers = utils.RemoveString(ingress.Finalizers, finalizerName)
				if err := r.Update(ctx, &ingress); err != nil {
					log.Error(err, "unable to add Ingress Finalizers")
					return ctrl.Result{}, err
				}
				log.Info("remove Finalizers from Ingress")
				r.Eventer.Event(&account, "Normal", "INFO", fmt.Sprintf("remove Finalizers from Ingress %v", req.NamespacedName))
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		For(&netv1.Ingress{}).
		Complete(r)
}
