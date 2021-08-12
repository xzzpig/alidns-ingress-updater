package utils

import (
	"sort"

	netv1 "k8s.io/api/networking/v1"
)

func RulesToStrings(rules []netv1.IngressRule) []string {
	if rules == nil {
		return nil
	}
	hosts := make([]string, len(rules))
	for _, rule := range rules {
		hosts = append(hosts, rule.Host)
	}
	sort.Strings(hosts)
	return hosts
}

func RulesEqual(lRules []netv1.IngressRule, rRules []netv1.IngressRule) bool {
	lHosts := RulesToStrings(lRules)
	rHosts := RulesToStrings(rRules)
	if lHosts == nil && rHosts == nil {
		return true
	}
	if len(lHosts) != len(rHosts) {
		return false
	}
	for i, lHost := range lHosts {
		rHost := rHosts[i]
		if lHost != rHost {
			return false
		}
	}
	return true
}
