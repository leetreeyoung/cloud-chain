package controllers

import (
	"cloudchain-backend/models/dto"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *MainController) TestModel() {
	var testModel dto.TestModel
	err := json.Unmarshal(c.Ctx.Input.RequestBody,&testModel)
	if err != nil {
		logs.Error(err)
	}
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
