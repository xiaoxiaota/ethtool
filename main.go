package main

import (
	_ "ethtool/routers"
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/context"
)

//func init() {
//	//跨域设置
//	var FilterGateWay = func(ctx *context.Context) {
//		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
//		//允许访问源
//		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
//		//允许post访问
//		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,ContentType,Authorization,accept,accept-encoding, authorization, content-type") //header的类型
//		ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
//		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
//	}
//	beego.InsertFilter("*", beego.BeforeRouter, FilterGateWay)
//}

func main() {
	beego.Run()
}
