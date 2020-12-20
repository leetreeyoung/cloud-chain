// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/sinksmell/bee-crontab/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/job",
			beego.NSRouter("/save", &controllers.JobController{}, "post:Save"),
			beego.NSRouter("/edit", &controllers.JobController{}, "post:Save"),
			beego.NSRouter("/list", &controllers.JobController{}, "get:List"),
			beego.NSRouter("/kill", &controllers.JobController{}, "post:Kill"),
			beego.NSRouter("/delete", &controllers.JobController{}, "post:Delete"),
			beego.NSRouter("/log/:id", &controllers.JobController{}, "get:Log"),
		),
		beego.NSNamespace("/worker",
			beego.NSRouter("/list", &controllers.WorkerController{}, "get:List"),
		),
	)
	beego.AddNamespace(ns)
}
