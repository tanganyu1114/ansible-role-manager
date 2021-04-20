package version

import (
	"github.com/tanganyu1114/ansible-role-manager/app/admin/models/system"

	//"github.com/tanganyu1114/ansible-role-manager/app/admin/models"
	"gorm.io/gorm"
	"runtime"

	"github.com/tanganyu1114/ansible-role-manager/cmd/migrate/migration"
	common "github.com/tanganyu1114/ansible-role-manager/common/models"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1613697961697Test)
}

func _1613697961697Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		//修改字段类型
		err := tx.Migrator().AlterColumn(&system.SysMenu{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&system.SysMenu{}, "update_by")
		if err != nil {
			return err
		}

		err = tx.Migrator().AlterColumn(&system.SysRole{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&system.SysRole{}, "update_by")
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
