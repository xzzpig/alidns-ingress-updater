package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type PublicIP struct {
	Address  string `json:"address"`
	Code     int64  `json:"code"`
	IP       string `json:"ip"`
	IsDomain int64  `json:"isDomain"`
	Rs       int64  `json:"rs"`
}

func GetPublicIP() (ip string, err error) {
	ip, err = GetPublicIP1()
	if err == nil {
		return
	}
	ip, err = GetPublicIP2()
	if err == nil {
		return
	}

	return
}

func GetPublicIP1() (string, error) {
	resp, err := http.Get("https://v6r.ipip.net/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("Error Status" + resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func GetPublicIP2() (string, error) {
	resp, err := http.Get("https://ip.cn/api/index?ip=&type=0")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("Error Status" + resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	// return string(body), nil
	var publicIp PublicIP
	json.Unmarshal(body, &publicIp)
	return publicIp.IP, nil
}
