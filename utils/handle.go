package utils

import (
	"fmt"
	"strings"

	networkv1 "github.com/xzzpig/alidns-ingress-updater/apis/network/v1"
	netv1 "k8s.io/api/networking/v1"
)

type InfoHandler = func(string, string, ...interface{})
type DebugHandler = func(string, ...interface{})
type ErrorHandler = func(error, string, string, ...interface{})
type SuccessHandler = func(host string)

type Handler struct {
	InfoHandler
	DebugHandler
	ErrorHandler
	SuccessHandler
}

func HandleUpdate(account *networkv1.AliDnsAccount, ingress *netv1.Ingress, ip string, handler Handler) {
	handler.DebugHandler("start handle ingress")
	if ingress == nil {
		handler.DebugHandler("ingress is nil")
		return
	}
	if ingress.ObjectMeta.Annotations["xzzpig.com/alidns-ignore"] == "true" {
		handler.DebugHandler("ingress is ignored")
		return
	}

	for _, rule := range ingress.Spec.Rules {
		host := rule.Host
		if host == "" {
			continue
		}

		domainName := "." + account.Spec.DomainName
		if !strings.HasSuffix(host, domainName) {
			continue
		}
		dnsUtils, err := NewAliDnsUtils(account.Spec)
		if err != nil {
			handler.ErrorHandler(err, "unable to create alidns utils", "")
			continue
		}
		rr := strings.ReplaceAll(host, domainName, "")
		record, err := dnsUtils.FindRecordByRR(rr)
		if err != nil {
			handler.ErrorHandler(err, "unable to find record by rr", "", "rr", rr)
			continue
		}
		if record == nil {
			_, err = dnsUtils.CreateRecord(rr, ip, "A")
			if err != nil {
				handler.ErrorHandler(err, "dns record create failed", fmt.Sprintf("dns record %v create as ip %v failed", host, ip), "host", host, "ip", ip)
				continue
			}
			handler.InfoHandler("dns record created", fmt.Sprintf("dns record %v created as ip %v", host, ip), "host", host, "ip", ip)
			handler.SuccessHandler(host)
			continue
		}
		if *record.Type == "A" && *record.Value == ip {
			handler.InfoHandler("dns record not updated because ip not changed", fmt.Sprintf("dns record %v not updated because ip %v not changed", host, ip), "host", host, "ip", ip, "oldIP", *record.Value, "ip", ip)
			continue
		}
		err = dnsUtils.UpdateRecord(*record.RecordId, rr, ip, "A")
		if err != nil {
			handler.ErrorHandler(err, "dns record update failed", fmt.Sprintf("dns record %v update from %v to %v failed", host, *record.Value, ip), "host", host, "ip", ip, "oldIP", *record.Value)
			continue
		}
		handler.InfoHandler("dns record updated", fmt.Sprintf("dns record %v updated from %v to %v", host, *record.Value, ip), "host", host, "ip", ip, "oldIP", *record.Value)
		handler.SuccessHandler(host)
	}
}

func HandleDelete(account *networkv1.AliDnsAccount, ingress *netv1.Ingress, handler Handler) {
	handler.DebugHandler("start handle ingress")
	if ingress == nil {
		handler.DebugHandler("ingress is nil")
		return
	}
	if ingress.ObjectMeta.Annotations["xzzpig.com/alidns-ignore"] == "true" {
		handler.DebugHandler("ingress is ignored")
		return
	}

	for _, rule := range ingress.Spec.Rules {
		host := rule.Host
		if host == "" {
			continue
		}

		domainName := "." + account.Spec.DomainName
		if !strings.HasSuffix(host, domainName) {
			continue
		}
		dnsUtils, err := NewAliDnsUtils(account.Spec)
		if err != nil {
			handler.ErrorHandler(err, "unable to create alidns utils", "")
			continue
		}
		rr := strings.ReplaceAll(host, domainName, "")
		record, err := dnsUtils.FindRecordByRR(rr)
		if err != nil {
			handler.ErrorHandler(err, "unable to find record by rr", "", "rr", rr)
			continue
		}
		if record == nil || record.RecordId == nil {
			handler.DebugHandler("dns record not need to delete (not exists)", fmt.Sprintf("dns record %v not need to delete (not exists)", host), "host", host)

		}
		err = dnsUtils.DeleteRecord(*record.RecordId)
		if err != nil {
			handler.ErrorHandler(err, "dns record delete failed", fmt.Sprintf("dns record %v delete failed", host), "host", host)
			continue
		}
		handler.InfoHandler("dns record deleted", fmt.Sprintf("dns record %v deleted", host), "host", host)
	}
}
