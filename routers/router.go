package routers

import (
	"ethtool/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//ns := beego.NewNamespace("/v1",
	//	//  用于跨域请求
	//	beego.NSRouter("*", &controllers.EthController{}, "OPTIONS:Options"))
	//beego.AddNamespace(ns)

	beego.Router("/", &controllers.MainController{})
	beego.Router("/get", &controllers.EthController{})

}
