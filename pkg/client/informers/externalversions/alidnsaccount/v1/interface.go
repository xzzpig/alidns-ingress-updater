/**/
// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	internalinterfaces "github.com/xzzpig/alidns-ingress-updater/pkg/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// AliDnsAccounts returns a AliDnsAccountInformer.
	AliDnsAccounts() AliDnsAccountInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// AliDnsAccounts returns a AliDnsAccountInformer.
func (v *version) AliDnsAccounts() AliDnsAccountInformer {
	return &aliDnsAccountInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
