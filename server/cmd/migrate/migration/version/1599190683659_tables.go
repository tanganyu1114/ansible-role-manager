package version

import (
	"runtime"

	"github.com/tanganyu1114/ansible-role-manager/app/admin/models/system"
	"gorm.io/gorm"

	"github.com/tanganyu1114/ansible-role-manager/app/admin/models"
	"github.com/tanganyu1114/ansible-role-manager/app/admin/models/tools"
	"github.com/tanganyu1114/ansible-role-manager/cmd/migrate/migration"
	common "github.com/tanganyu1114/ansible-role-manager/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1599190683659Tables)
}

func _1599190683659Tables(db *gorm.DB, version string) error {
	err := db.Debug().Migrator().AutoMigrate(
		new(system.CasbinRule),
		new(system.SysDept),
		new(system.SysConfig),
		new(tools.SysTables),
		new(tools.SysColumns),
		new(system.SysMenu),
		new(system.SysLoginLog),
		new(system.SysOperaLog),
		new(models.RoleMenu),
		new(system.SysRoleDept),
		new(system.SysUser),
		new(system.SysRole),
		new(Post),
		new(DictData),
		new(DictType),
		new(models.SysJob),
		new(system.SysConfig),
		new(models.SysSetting),
		new(models.SysFileDir),
		new(models.SysFileInfo),
		new(models.SysCategory),
		new(models.SysContent),
	)
	if err != nil {
		return err
	}
	return db.Create(&common.Migration{
		Version: version,
	}).Error
}
