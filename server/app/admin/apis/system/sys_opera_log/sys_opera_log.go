package sys_opera_log

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"

	"github.com/tanganyu1114/ansible-role-manager/app/admin/models/system"
	"github.com/tanganyu1114/ansible-role-manager/app/admin/service"
	"github.com/tanganyu1114/ansible-role-manager/app/admin/service/dto"
	"github.com/tanganyu1114/ansible-role-manager/common/apis"
)

type SysOperaLog struct {
	apis.Api
}

func (e *SysOperaLog) GetSysOperaLogList(c *gin.Context) {
	log := e.GetLogger(c)
	d := new(dto.SysOperaLogSearch)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//查询列表
	err = d.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	list := make([]system.SysOperaLog, 0)
	var count int64
	serviceStudent := service.SysOperaLog{}
	serviceStudent.Log = log
	serviceStudent.Orm = db
	err = serviceStudent.GetSysOperaLogPage(d, &list, &count)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.PageOK(c, list, int(count), d.GetPageIndex(), d.GetPageSize(), "查询成功")
}

func (e *SysOperaLog) GetSysOperaLog(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.SysOperaLogById)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//查看详情
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	var object system.SysOperaLog

	serviceSysOperlog := service.SysOperaLog{}
	serviceSysOperlog.Log = log
	serviceSysOperlog.Orm = db
	err = serviceSysOperlog.GetSysOperaLog(control, &object)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "查询失败")
		return
	}

	e.OK(c, object, "查看成功")
}

func (e *SysOperaLog) InsertSysOperaLog(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.SysOperaLogControl)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//新增操作
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	// 设置创建人
	object.SetCreateBy(user.GetUserId(c))

	serviceSysOperaLog := service.SysOperaLog{}
	serviceSysOperaLog.Orm = db
	serviceSysOperaLog.Log = log
	err = serviceSysOperaLog.InsertSysOperaLog(object)
	if err != nil {
		log.Error(err)
		e.Error(c, http.StatusInternalServerError, err, "创建失败")
		return
	}

	e.OK(c, object.GetId(), "创建成功")
}

func (e *SysOperaLog) UpdateSysOperaLog(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.SysOperaLogControl)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//更新操作
	err = control.Bind(c)
	if err != nil {
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}
	object, err := control.Generate()
	if err != nil {
		e.Error(c, http.StatusInternalServerError, err, "模型生成失败")
		return
	}
	object.SetUpdateBy(user.GetUserId(c))

	serviceSysOperaLog := service.SysOperaLog{}
	serviceSysOperaLog.Orm = db
	serviceSysOperaLog.Log = log
	err = serviceSysOperaLog.UpdateSysOperaLog(object)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, object.GetId(), "更新成功")
}

func (e *SysOperaLog) DeleteSysOperaLog(c *gin.Context) {
	log := e.GetLogger(c)
	control := new(dto.SysOperaLogById)
	db, err := e.GetOrm(c)
	if err != nil {
		log.Error(err)
		return
	}

	//删除操作
	err = control.Bind(c)
	if err != nil {
		log.Errorf("Bind error: %s", err)
		e.Error(c, http.StatusUnprocessableEntity, err, "参数验证失败")
		return
	}

	serviceSysOperaLog := service.SysOperaLog{}
	serviceSysOperaLog.Orm = db
	serviceSysOperaLog.Log = log
	err = serviceSysOperaLog.RemoveSysOperaLog(control)
	if err != nil {
		log.Error(err)
		return
	}
	e.OK(c, control.GetId(), "删除成功")
}
