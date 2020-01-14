package module

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type JsonReq struct {
	Jsonrpc string        `json: jsonrpc`
	Method  string        `json: method`
	Params  []interface{} `json: params`
	Id      int           `json: id`
}

type JsonRes struct {
	Id      int         `json: id`
	Jsonrpc string      `json: jsonrpc`
	Result  interface{} `json: result`
	Error   interface{} `json:error`
}

type Contract struct {

}

type ABIModule struct {
	Contracts interface {} `json: contracts`
	Version string `json: string`
}

func NewJsonReq(method string, params []interface{}) *JsonReq {
	id, err := beego.AppConfig.Int(method)
	if err != nil {
		logs.Warn("method [%s] not exist.", method)
		return nil
	}

	return &JsonReq{"2.0", method, params, id}
}
