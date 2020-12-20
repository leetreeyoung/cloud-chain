package controllers

import (
	"context"
	"encoding/json"

	"github.com/sinksmell/bee-crontab/models/common"
	"github.com/sinksmell/bee-crontab/models/master"
	log "github.com/sirupsen/logrus"

	"github.com/astaxie/beego"
)

type JobController struct {
	beego.Controller
}

//func (c *JobController) URLMapping() {
//	c.Mapping("Save", c.Save)     // update or create job
//	c.Mapping("Delete", c.Delete) // delete job
//	c.Mapping("List", c.List)     // get all of the jobs
//	c.Mapping("Kill", c.Kill)     // kill job
//	c.Mapping("Log", c.Log)
//}

// @Title SaveJob
// @Description create jobs or update jobs
// @Param	body		body 	models.Job	true		"body for Job content"
// @Success 200 {int}
// @Failure 403 body is empty
// @router /save [post]
func (c *JobController) Save() {
	var (
		job  common.Job
		resp common.Response
		err  error
	)
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &job); err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
		c.Data["json"] = resp
		return
	}
	log.Infof("request is :%+v\n", job)
	if len(job.ID) == 0 {
		job.ID = master.NewID()
	}
	if oldJob, err := master.MJobManager.SaveJob(&job); err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
	} else {
		resp = common.NewResponse(0, "success", oldJob)
	}

	c.Data["json"] = &resp
	c.ServeJSON()
}

// @Title DeleteJob
// @Description delete job
// @Param	body		body 	models.Job	true		"body for Job content"
// @Success 200 {int}
// @Failure 403 body is empty
// @router /delete [post]
func (c *JobController) Delete() {
	var (
		job  common.Job
		resp common.Response
	)
	json.Unmarshal(c.Ctx.Input.RequestBody, &job)
	if oldJob, err := master.MJobManager.DeleteJob(&job); err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
	} else {
		resp = common.NewResponse(0, "success", oldJob)
	}
	c.Data["json"] = &resp
	c.ServeJSON()
}

// @Title ListJob
// @Description get all of the jobs
// @Success 200 {object} models.Job
// @router /list [get]
func (c *JobController) List() {
	var (
		resp common.Response
	)
	if jobs, err := master.MJobManager.ListJobs(); err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
	} else {
		resp = common.NewResponse(0, "success", jobs)
	}
	c.Data["json"] = &resp
	c.ServeJSON()
}

// @Title KillJob
// @Description Kill job by  name of job
// @Param	body		body 	models.Job	true		"body for Job content"
// @Success 200 {int}
// @Failure 403 body is empty
// @router /kill [post]
func (c *JobController) Kill() {
	var (
		resp common.Response
		job  common.Job
	)
	json.Unmarshal(c.Ctx.Input.RequestBody, &job)
	if err := master.MJobManager.KillJob(&job); err != nil {
		resp = common.NewResponse(-1, err.Error(), nil)
	} else {
		resp = common.NewResponse(0, "success", nil)
	}

	c.Data["json"] = &resp
	c.ServeJSON()
}

// @Title GetJobLog
// @Description get job execute log by job name
// @Param	name		path 	string	true		"The key for staticblock"
// @Success 200
// @router /log/:name [get]
func (c *JobController) Log() {
	var (
		logs []*common.HTTPJobLog
		resp common.Response
		err  error
	)

	jobID := c.GetString(":id")
	log.Info("job_id: ", jobID)
	if jobID != "" {
		if logs, err = master.ReadLog(context.TODO(), jobID); err != nil {
			log.Errorf("read log err job_id:%s err:%v", jobID, err)
			resp = common.NewResponse(-1, err.Error(), nil)
		} else {
			resp = common.NewResponse(0, "success", logs)
		}
	}
	c.Data["json"] = &resp
	c.ServeJSON()
}
