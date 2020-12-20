package controllers

import (
	"github.com/astaxie/beego"
	"github.com/sinksmell/bee-crontab/models/common"
	"github.com/sinksmell/bee-crontab/models/master"
)

type WorkerController struct {
	beego.Controller
}

//func (w *WorkerController) URLMapping() {
//	w.Mapping("List", w.List)
//}

// @Title List Workers Node
// @Description get all of the workers ip
// @Success 200
// @router /list [get]
func (w *WorkerController) List() {

	var (
		resp common.Response
	)
	if ips, err := master.WorkerManager.ListWorkers(); err == nil {
		resp = common.NewResponse(0, "success", ips)
	} else {
		resp = common.NewResponse(-1, err.Error(), nil)
	}
	w.Data["json"] = &resp
	w.ServeJSON()
}
