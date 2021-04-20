package version

import (
	"gorm.io/gorm"

	"github.com/tanganyu1114/ansible-role-manager/app/admin/models/system"
	"github.com/tanganyu1114/ansible-role-manager/cmd/migrate/migration"
	common "github.com/tanganyu1114/ansible-role-manager/common/models"

	"runtime"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1603516925109Test)
}

func _1603516925109Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		_ = tx.Migrator().RenameTable("sys_operlog", "sys_opera_log")
		_ = tx.Migrator().RenameTable("sys_loginlog", "sys_login_log")

		if tx.Migrator().HasColumn(&system.SysLoginLog{}, "info_id") {
			err := tx.Migrator().RenameColumn(&system.SysLoginLog{}, "info_id", "id")
			if err != nil {
				return err
			}
		}

		if tx.Migrator().HasColumn(&system.SysOperaLog{}, "oper_id") {
			err := tx.Migrator().RenameColumn(&system.SysOperaLog{}, "oper_id", "id")
			if err != nil {
				return err
			}
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
