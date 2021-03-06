package perm

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/perm"
)

type PermRoleController struct {
	beego.Controller
}

// 权限管理入口页面
// @router /perm/role/list [get]
func (this *PermRoleController) PermRoleList() {
	this.TplName = "perm/role/list.html"
}

// 权限管理添加页面
// @router /perm/role/add [get]
func (this *PermRoleController) PermRoleAdd() {
	id := this.GetString("RoleId")
	update := perm.CloudPermRole{}
	// 更新操作
	if id != "" {
		searchMap := sql.GetSearchMap("RoleId", *this.Ctx)
		sql.Raw(sql.SearchSql(perm.CloudPermRole{}, perm.SelectCloudPermRole, searchMap)).QueryRow(&update)
	}
	this.Data["data"] = update
	this.TplName = "perm/role/add.html"
}

// 获取权限数据
// 2018-02-06 8:56
// router /api/perm/role [get]
func (this *PermRoleController) PermRoleData() {
	// 权限数据
	data := []perm.CloudPermRole{}
	q := sql.SearchSql(perm.CloudPermRole{}, perm.SelectCloudPermRole, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setPermRoleJson(this, data)
}

// string
// 权限保存
// @router /api/perm/role [post]
func (this *PermRoleController) PermRoleSave() {
	d := perm.CloudPermRole{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)
	
	q := sql.InsertSql(d, perm.InsertCloudPermRole)
	if d.RoleId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("RoleId", d.RoleId)
		q = sql.UpdateSql(d, perm.UpdateCloudPermRole, searchMap, "CreateTime,CreatePermRole")
	}
	_, err = sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存权限配置 "+msg, d.RoleName)
	setPermRoleJson(this, data)
}

// 获取权限数据
// 2018-02-06 08:36
// router /api/perm/role/name [get]
func (this *PermRoleController) PermRoleDataName() {
	// 权限数据
	data := []perm.CloudPermRole{}
	q := sql.SearchSql(perm.CloudPermRole{}, perm.SelectCloudPermRole, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setPermRoleJson(this, data)
}

// 权限数据
// @router /api/perm/role [get]
func (this *PermRoleController) PermRoleDatas() {
	data := []perm.CloudPermRole{}
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("RoleId", id)
	}
	searchSql := sql.SearchSql(perm.CloudPermRole{}, perm.SelectCloudPermRole, searchMap)
	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += " where 1=1 and (description like \"%" + key + "%\")"
	}
	num, _ := sql.OrderByPagingSql(searchSql, "role_id", *this.Ctx.Request, &data, perm.CloudPermRole{})
    r := util.ResponseMap(data, sql.Count("cloud_perm_role", int(num), key), this.GetString("draw"))
	setPermRoleJson(this, r)
}

// json
// 删除权限
// 2018-02-06 08:36
// @router /api/perm/role/:id:int [delete]
func (this *PermRoleController) PermRoleDelete() {
	searchMap := sql.GetSearchMap("RoleId", *this.Ctx)
	permData := perm.CloudPermRole{}
	sql.Raw(sql.SearchSql(permData, perm.SelectCloudPermRole, searchMap)).QueryRow(&permData)
	r, err := sql.Raw(sql.DeleteSql(perm.DeleteCloudPermRole, searchMap)).Exec()
	data := util.DeleteResponse(err, *this.Ctx, "删除权限"+permData.RoleName, this.GetSession("username"), permData.CreateUser, r)
	setPermRoleJson(this, data)
}

func setPermRoleJson(this *PermRoleController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}