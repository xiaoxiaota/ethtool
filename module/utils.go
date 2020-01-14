package module

import (
	"fmt"
	"strconv"
)

func Hex2dec(hex string) int64 {
	intv, _ := strconv.ParseInt(hex[2:], 16, 32)
	return intv
}

func Dec2hex(n interface{}) (ret string) {
	if v, ok := n.(string); ok {
		intv, _ := strconv.Atoi(v)
		ret = fmt.Sprintf("0x%X", intv)
	} else if v, ok := n.(int); ok {
		ret = fmt.Sprintf("0x%X", v)
	} else if v, ok := n.(int64); ok {
		ret = fmt.Sprintf("0x%X", v)
	}

	return
}
