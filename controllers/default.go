package controllers

import (
	//"encoding/json"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["NodeAddr"] = beego.AppConfig.String("nodeaddr")
	c.TplName = "index.html"
	c.Render()
}
