package routers

import (
	"cloudchain-backend/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/test-model", &controllers.MainController{},"post:TestModel")
}
