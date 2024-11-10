package model

import (
	"encoding/json"
	"regexp"
)

const (
	EvnBeta = "beta"
)

var (
	reVolume = regexp.MustCompile(`^[+]{0,1}(\d+)$`)
)

// Copy 从一个结构体复制到另一个结构体
func Copy(to, from interface{}) (err error) {
	b, err := json.Marshal(from)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, to)
	if err != nil {
		return
	}
	return nil
}

// IsETHAddress 是否为ETH地址
func IsETHAddress(address string) bool {
	// re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return reETHAddr.MatchString(address)
}

// IsVolumeLegal 判断volume（正整数）
func IsVolumeLegal(number string) bool {
	// re := regexp.MustCompile("^[+]{0,1}(\\d+)$")
	if number == "0" {
		return false
	}
	return reVolume.MatchString(number)
}

// IsNonNegativeInteger 判断volume（非负整数）
func IsNonNegativeInteger(number string) bool {
	// re := regexp.MustCompile("^[+]{0,1}(\\d+)$")
	return reVolume.MatchString(number)
}
