/**/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	alidnsaccountv1 "github.com/xzzpig/alidns-ingress-updater/pkg/apis/alidnsaccount/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAliDnsAccounts implements AliDnsAccountInterface
type FakeAliDnsAccounts struct {
	Fake *FakeXzzpigV1
}

var alidnsaccountsResource = schema.GroupVersionResource{Group: "xzzpig.com", Version: "v1", Resource: "alidnsaccounts"}

var alidnsaccountsKind = schema.GroupVersionKind{Group: "xzzpig.com", Version: "v1", Kind: "AliDnsAccount"}

// Get takes name of the aliDnsAccount, and returns the corresponding aliDnsAccount object, and an error if there is any.
func (c *FakeAliDnsAccounts) Get(ctx context.Context, name string, options v1.GetOptions) (result *alidnsaccountv1.AliDnsAccount, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(alidnsaccountsResource, name), &alidnsaccountv1.AliDnsAccount{})
	if obj == nil {
		return nil, err
	}
	return obj.(*alidnsaccountv1.AliDnsAccount), err
}

// List takes label and field selectors, and returns the list of AliDnsAccounts that match those selectors.
func (c *FakeAliDnsAccounts) List(ctx context.Context, opts v1.ListOptions) (result *alidnsaccountv1.AliDnsAccountList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(alidnsaccountsResource, alidnsaccountsKind, opts), &alidnsaccountv1.AliDnsAccountList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &alidnsaccountv1.AliDnsAccountList{ListMeta: obj.(*alidnsaccountv1.AliDnsAccountList).ListMeta}
	for _, item := range obj.(*alidnsaccountv1.AliDnsAccountList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested aliDnsAccounts.
func (c *FakeAliDnsAccounts) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(alidnsaccountsResource, opts))
}

// Create takes the representation of a aliDnsAccount and creates it.  Returns the server's representation of the aliDnsAccount, and an error, if there is any.
func (c *FakeAliDnsAccounts) Create(ctx context.Context, aliDnsAccount *alidnsaccountv1.AliDnsAccount, opts v1.CreateOptions) (result *alidnsaccountv1.AliDnsAccount, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(alidnsaccountsResource, aliDnsAccount), &alidnsaccountv1.AliDnsAccount{})
	if obj == nil {
		return nil, err
	}
	return obj.(*alidnsaccountv1.AliDnsAccount), err
}

// Update takes the representation of a aliDnsAccount and updates it. Returns the server's representation of the aliDnsAccount, and an error, if there is any.
func (c *FakeAliDnsAccounts) Update(ctx context.Context, aliDnsAccount *alidnsaccountv1.AliDnsAccount, opts v1.UpdateOptions) (result *alidnsaccountv1.AliDnsAccount, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(alidnsaccountsResource, aliDnsAccount), &alidnsaccountv1.AliDnsAccount{})
	if obj == nil {
		return nil, err
	}
	return obj.(*alidnsaccountv1.AliDnsAccount), err
}

// Delete takes name of the aliDnsAccount and deletes it. Returns an error if one occurs.
func (c *FakeAliDnsAccounts) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(alidnsaccountsResource, name), &alidnsaccountv1.AliDnsAccount{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAliDnsAccounts) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(alidnsaccountsResource, listOpts)

	_, err := c.Fake.Invokes(action, &alidnsaccountv1.AliDnsAccountList{})
	return err
}

// Patch applies the patch and returns the patched aliDnsAccount.
func (c *FakeAliDnsAccounts) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *alidnsaccountv1.AliDnsAccount, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(alidnsaccountsResource, name, pt, data, subresources...), &alidnsaccountv1.AliDnsAccount{})
	if obj == nil {
		return nil, err
	}
	return obj.(*alidnsaccountv1.AliDnsAccount), err
}
