package controllers

import (
	"encoding/json"
	"ethtool/beego4eth/module"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
)

/*
 */
type TFunsMap map[string]reflect.Value

/*
 */
type EthController struct {
	beego.Controller
}

/*
 */
type EthTx struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	Nonce    string `json:"nonce"`
}

/*
 */
type EthContract struct {
	From     string `json:"from"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	Nonce    string `json:"nonce"`
}

//func (e *EthController) Options() {
//	e.Data["json"] = map[string]interface{}{"status": 200, "message": "ok", "moreinfo": ""}
//	e.ServeJSON()
//}

/*
 */
func (e *EthController) Get() {
	logs.Debug("Get")
	e.Data["NodeAddr"] = beego.AppConfig.String("nodeaddr")
	e.TplName = "index.html"

	//common info
	//e.callFun("Commoninfo")

	e.callFun(e.GetString("action"))
	e.Data["Action"] = e.GetString("action")

	//e.Render()
}

func (e *EthController) callFun(action string) {
	logs.Debug("callfun", action)
	var method string
	fmap := make(TFunsMap, 0)
	vf := reflect.ValueOf(e)
	vft := vf.Type()
	mNum := vf.NumMethod()

	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		fmap[mName] = vf.Method(i)
	}

	//params := []reflect.Value{reflect.ValueOf(ec)}
	//fmap[action].Call(params)
	if action == "NewAccount" {
		method = "Personal" + action
	} else {
		method = "Eth" + action
	}

	fmap[method].Call(nil)
}

func (e *EthController) EthCommoninfo() {
	var commoninfo = []string{"blockNumber", "coinbase"}

	logs.Debug("get common info")
	for i := 0; i < len(commoninfo); i++ {
		if res, err := commRPCCall("eth_"+commoninfo[i], nil); err == nil {
			if i == 0 {
				snum := res.Result.(string)
				e.Data[commoninfo[i]] = module.Hex2dec(snum)
				beego.AppConfig.Set("blockNumber", string(module.Hex2dec(snum)))
			} else {
				e.Data[commoninfo[i]] = res.Result
				v,_:=res.Result.(string)
				beego.AppConfig.Set("coinbase",  v)
			}
		}
	}

	//write to config file


}

func (e *EthController) EthGetBlockByNumber() {
	var endn,tn int
	var content string
	//var params []interface{}
	var blockinfo = []string{"getBlockTransactionCountByNumber", "getBlockByNumber"}
	blocknumstr := e.GetString("blocknum")
	narray := strings.Split(blocknumstr, "-")
	startn,_ := strconv.Atoi(narray[0])
	if len(narray) == 1 {
		endn = startn
	} else {
		endn,_ = strconv.Atoi(narray[1])
	}
	//params = append(params, module.Dec2hex(e.GetString("blocknum")))
	for j:=startn;j<=endn;j++{
		for i := 0; i < len(blockinfo); i++ {
			var params []interface{}
			params = append(params, module.Dec2hex(j))
			if i == 1 {
				params = append(params, true)
			}

			if res, err := commRPCCall("eth_"+blockinfo[i], params); err == nil {
				if i == 1 {
					jsonbyte, _ := json.Marshal(res.Result)
					if content != "" {
						content += ","
					}
					content += string(jsonbyte)
				} else {
					tn += int(module.Hex2dec(res.Result.(string)))
				}
			}
		}

	}
	e.Data["Content"] = "["+content+"]"
	e.Data["getBlockTransactionCountByNumber"] = tn
}

func (e *EthController) EthGetTransactionByHash() {
	var params []interface{}
	var content string
	params = append(params, e.GetString("txhash"))

	if res, err := commRPCCall("eth_getTransactionByHash", params); err == nil {
		jsonbyte, _ := json.Marshal(res.Result)
		content = string(jsonbyte)
		e.Data["Content"] = string(jsonbyte)
	}
	if res, err := commRPCCall("eth_getTransactionReceipt", params); err == nil {
		jsonbyte, _ := json.Marshal(res.Result)
		content += string(jsonbyte)
		e.Data["Content1"] = string(jsonbyte)
	}
}

func (e *EthController) EthAccounts() {
	if res, err := commRPCCall("eth_accounts", nil); err == nil {
		e.Data["Content"] = res.Result
	}
}

func (e *EthController) EthAccountinfo() {
	var blockn string
	var accountinfo = []string{"getBalance", "getTransactionCount"}
	var params []interface{}
	account := e.GetString("account")
	if account == "" {
		account = e.Data["coinbase"].(string)
	}
	params = append(params, account)
	if e.GetString("blocknum") == "" {
		blockn = "latest"
	} else {
		blockn = module.Dec2hex(e.GetString("blocknum"))
	}
	params = append(params, blockn)

	logs.Debug("get account info for block: %s, account: %s", blockn, account)
	for i := 0; i < len(accountinfo); i++ {
		if res, err := commRPCCall("eth_"+accountinfo[i], params); err == nil {
			snum := res.Result.(string)
			e.Data[accountinfo[i]] = module.Hex2dec(snum)
		}
	}
}

func (e *EthController) EthSendTransaction() {
	var res module.JsonRes
	var params []interface{}
	var err error

	account := e.GetString("accountfrom")
	if account == "" {
		account = e.Data["Coinbase"].(string)
	}
	accountto := e.GetString("accountto")
	value := module.Dec2hex(e.GetString("value"))
	contract := e.GetString("contract")

	if value == "" && contract == "" {
		e.Data["Content"] = "缺少参数，转账金额或合约字节码"
		return
	}

	if accountto != "" && value == "" {
		e.Data["Content"] = "转账金额不能为空"
		return
	}

	// unlock the account
	acparam := []interface{}{account, beego.AppConfig.String("password"), 300}
	if res, err = commRPCCall("personal_unlockAccount", acparam); err != nil {
		e.Data["Content"] = "解锁帐户失败: " + err.Error()
		return
	} else if res.Error != nil {
		logs.Warn("unlock account error: ", res.Error, res.Result)
		e.Data["Content"] = res.Result
		return

	}

	// get nonce
	nonceparams := []interface{}{account, "latest"}
	if res, err = commRPCCall("eth_getTransactionCount", nonceparams); err != nil {
		logs.Warn("unlock account error: ", err.Error())
		e.Data["Content"] = "获取nonce失败: " + err.Error()
		return
	}
	intv := module.Hex2dec(res.Result.(string))
	nonce := module.Dec2hex(intv)
	logs.Debug("tx nonce: ", intv, nonce)

	// send tx
	var txparams EthTx
	if accountto == "" {
		txparams = EthTx{From: account, To: "0x0000000000000000000000000000000000000000", Nonce: nonce, Data: contract, Gas: beego.AppConfig.String("gasc")}

	} else {
		txparams = EthTx{From: account, To: accountto, Value: value, Nonce: nonce, Data: contract, Gas: beego.AppConfig.String("gas")}
	}
	params = append(params, txparams)
	if res, err = commRPCCall("eth_sendTransaction", params); err != nil {
		logs.Warn("unlock account error: ", err.Error())
		e.Data["Content"] = "解锁帐户失败: " + err.Error()
		return
	}
	if res.Result != nil {
		e.Data["Content"] = res.Result.(string)
	} else if res.Error != nil {
		e.Data["Content"] = res.Error
	}
}

func (e *EthController) PersonalNewAccount() {
	var params []interface{}
	params = append(params, beego.AppConfig.String("password"))

	if res, err := commRPCCall("personal_newAccount", params); err == nil {
		jsonbyte, _ := json.Marshal(res.Result)
		e.Data["Content"] = string(jsonbyte)
	}
}

func (e *EthController) EthCreateContract() {
	var res module.JsonRes
	var params []interface{}
	var err error

	//compiler contract
	cmd:=exec.Command("solc","--optimize" ,"--combined-json","abi,bin,interface" ,"sandbox.sol" )
	output, _ :=cmd.Output()
	stroutput := string(output)
	pos := strings.Index(stroutput,"bin")
	end := strings.Index(stroutput[pos+6:],"\"")
	binstr := "0x"+stroutput[pos+6:pos+6+end]

	// send tx
	// unlock the account
	account := beego.AppConfig.String("coinbase")
	acparam := []interface{}{account, beego.AppConfig.String("password"), 300}
	if res, err = commRPCCall("personal_unlockAccount", acparam); err != nil {
		e.Data["Content"] = "解锁帐户失败: " + err.Error()
		return
	} else if res.Error != nil {
		logs.Warn("unlock account error: ", res.Error, res.Result)
		e.Data["Content"] = res.Result
		return

	}

	// get nonce
	nonceparams := []interface{}{account, "latest"}
	if res, err = commRPCCall("eth_getTransactionCount", nonceparams); err != nil {
		logs.Warn("unlock account error: ", err.Error())
		e.Data["Content"] = "获取nonce失败: " + err.Error()
		return
	}
	intv := module.Hex2dec(res.Result.(string))
	nonce := module.Dec2hex(intv)
	logs.Debug("tx nonce: ", intv, nonce)
	var txparams EthContract
	txparams = EthContract{From: account, Nonce: nonce, Data: binstr, Gas: beego.AppConfig.String("gascc")}
	params = append(params, txparams)
	if res, err = commRPCCall("eth_sendTransaction", params); err != nil {
		logs.Warn("unlock account error: ", err.Error())
		e.Data["Content"] = "解锁帐户失败: " + err.Error()
		return
	}
	if res.Result != nil {
		e.Data["Content"] = res.Result.(string)
	} else if res.Error != nil {
		e.Data["Content"] = res.Error
	}

}


func commRPCCall(method string, params []interface{}) (module.JsonRes, error) {
	var res module.JsonRes
	var err error
	var str string

	reqstr := module.NewJsonReq(method, params)
	req := httplib.Post(beego.AppConfig.String("nodeaddr"))
	// 跨域
	//req.Header("Access-Control-Allow-Origin", "*")
	req.JSONBody(reqstr)
	if str, err = req.String(); err != nil {
		logs.Warn("req error: ", err)
	} else {
		json.Unmarshal([]byte(str), &res)
	}
	logs.Debug("req: ", method, params)
	logs.Debug("res: ", res)
	return res, err
}
