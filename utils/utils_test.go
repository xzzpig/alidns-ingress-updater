package utils_test

import (
	"fmt"
	"testing"

	"github.com/xzzpig/alidns-ingress-updater/utils"
)

func TestGetIp(t *testing.T) {
	ip, err := utils.GetPublicIP()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ip)
}

func TestContainsString(t *testing.T) {
	strs := make([]string, 0)
	strs = append(strs, "aaa", "bbb", "ccc")
	fmt.Println(strs)
	fmt.Println(utils.ContainsString(strs, "aaa"))
	fmt.Println(utils.ContainsString(strs, "bbb"))
	fmt.Println(utils.ContainsString(strs, "bbc"))
}

func TestRemoveString(t *testing.T) {
	strs := make([]string, 0)
	strs = append(strs, "aaa", "bbb", "ccc", "bbb")
	fmt.Println(strs)
	strs = utils.RemoveString(strs, "bbb")
	fmt.Println(strs)
	fmt.Println(len(strs))
}
