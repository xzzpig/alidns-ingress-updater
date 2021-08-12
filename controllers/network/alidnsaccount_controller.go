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

package network

import (
	"context"
	"fmt"
	"time"

	networkv1 "github.com/xzzpig/alidns-ingress-updater/apis/network/v1"
	"github.com/xzzpig/alidns-ingress-updater/utils"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const finalizerName = "account.xzzpig.com/finalizer"

// AliDnsAccountReconciler reconciles a AliDnsAccount object
type AliDnsAccountReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Eventer record.EventRecorder
}

//+kubebuilder:rbac:groups=network.xzzpig.com,resources=alidnsaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=network.xzzpig.com,resources=alidnsaccounts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=network.xzzpig.com,resources=alidnsaccounts/finalizers,verbs=update
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingress,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AliDnsAccount object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *AliDnsAccountReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var account networkv1.AliDnsAccount
	if err := r.Get(ctx, req.NamespacedName, &account); err != nil {
		log.V(1).Info("unable to fetch AliDnsAccount")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// TODO
	// account.Status.LastIp = ""

	var publicIp string
	if account.DeletionTimestamp.IsZero() {
		if ip, err := utils.GetPublicIP(); err != nil {
			log.Error(err, "unable to get Public IP")
			r.Eventer.Event(&account, "Warning", "ERROR", "unable to get Public IP")
			return ctrl.Result{}, err
		} else {
			if account.Status.LastIp == ip {
				log.WithValues("ip", ip).V(1).Info("ip not changed")
				return ctrl.Result{RequeueAfter: time.Minute}, nil
			}
			publicIp = ip
		}
	}

	var ingresses netv1.IngressList
	if err := r.List(ctx, &ingresses); err != nil {
		log.Error(err, "unable to list all Ingress")
		return ctrl.Result{}, err
	}

	for _, ingress := range ingresses.Items {
		ingressFullName := ingress.Namespace + "/" + ingress.Name
		if account.DeletionTimestamp.IsZero() {
			utils.HandleUpdate(&account, &ingress, publicIp, utils.Handler{
				InfoHandler: func(logMsg string, eventMsg string, keysAndValues ...interface{}) {
					if eventMsg == "" {
						eventMsg = logMsg
					}
					log.WithValues("ingress", ingressFullName).Info(logMsg, keysAndValues...)
					r.Eventer.Event(&account, "Normal", "INFO", eventMsg)
				},
				ErrorHandler: func(err error, logMsg string, eventMsg string, keysAndValues ...interface{}) {
					if eventMsg == "" {
						eventMsg = logMsg
					}
					log.WithValues("ingress", ingressFullName).Error(err, logMsg, keysAndValues...)
					r.Eventer.Event(&account, "Warning", "ERROR", eventMsg)
				},
				DebugHandler: func(logMsg string, keysAndValues ...interface{}) {
					log.WithValues("ingress", ingressFullName).V(1).Info(logMsg, keysAndValues...)
				},
				SuccessHandler: func(host string) {
					if !utils.ContainsString(ingress.Finalizers, finalizerName) {
						ingress.Finalizers = append(ingress.Finalizers, finalizerName)
						if err := r.Update(ctx, &ingress); err != nil {
							log.Error(err, "unable to add Ingress Finalizers")
						} else {
							log.Info("add Finalizers to Ingress")
							r.Eventer.Event(&account, "Normal", "INFO", fmt.Sprintf("add Finalizers to Ingress %v", ingressFullName))
						}
					}
				},
			})

			if !utils.ContainsString(account.Finalizers, finalizerName) {
				account.Finalizers = append(account.Finalizers, finalizerName)
				if err := r.Update(ctx, &account); err != nil {
					log.Error(err, "unable to add AliDnsAccount Finalizers")
					return ctrl.Result{}, err
				}
				log.Info("add Finalizers to AliDnsAccount")
				r.Eventer.Event(&account, "Normal", "INFO", fmt.Sprintf("add Finalizers to AliDnsAccount %v", req.NamespacedName))
			}

		} else {
			utils.HandleDelete(&account, &ingress, utils.Handler{
				InfoHandler: func(logMsg string, eventMsg string, keysAndValues ...interface{}) {
					if eventMsg == "" {
						eventMsg = logMsg
					}
					log.WithValues("ingress", ingressFullName).Info(logMsg, keysAndValues...)
					r.Eventer.Event(&account, "Normal", "INFO", eventMsg)
				},
				ErrorHandler: func(err error, logMsg string, eventMsg string, keysAndValues ...interface{}) {
					if eventMsg == "" {
						eventMsg = logMsg
					}
					log.WithValues("ingress", ingressFullName).Error(err, logMsg, keysAndValues...)
					r.Eventer.Event(&account, "Warning", "ERROR", eventMsg)
				},
				DebugHandler: func(logMsg string, keysAndValues ...interface{}) {
					log.WithValues("ingress", ingressFullName).V(1).Info(logMsg, keysAndValues...)
				},
			})
			if utils.ContainsString(account.Finalizers, finalizerName) {
				account.Finalizers = utils.RemoveString(account.Finalizers, finalizerName)
				if err := r.Update(ctx, &account); err != nil {
					log.Error(err, "unable to add AliDnsAccount Finalizers")
					return ctrl.Result{}, err
				}
				log.Info("remove Finalizers from AliDnsAccount")
				r.Eventer.Event(&account, "Normal", "INFO", fmt.Sprintf("remove Finalizers from AliDnsAccount %v", req.NamespacedName))
			}
		}
	}

	account.Status.LastIp = publicIp
	r.Status().Update(ctx, &account)

	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AliDnsAccountReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkv1.AliDnsAccount{}).
		Complete(r)
}
