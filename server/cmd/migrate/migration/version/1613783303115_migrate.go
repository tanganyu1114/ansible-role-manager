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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1613783303115Test)
}

func _1613783303115Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		//修改字段类型
		err := tx.Model(&system.SysDictType{}).Select("create_by").Not("create_by  > 0").Update("create_by", "0").Error
		if err != nil {
			return err
		}
		err = tx.Model(&system.SysDictType{}).Select("update_by").Not("update_by > 0").Update("update_by", "0").Error
		if err != nil {
			return err
		}
		err = tx.Model(&system.SysDictData{}).Select("create_by").Not("create_by > 0").Update("create_by", "0").Error
		if err != nil {
			return err
		}
		err = tx.Model(&system.SysDictData{}).Select("update_by").Not("update_by > 0").Update("update_by", "0").Error
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&system.SysDictType{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&system.SysDictType{}, "update_by")
		if err != nil {
			return err
		}

		err = tx.Migrator().AlterColumn(&system.SysDictData{}, "create_by")
		if err != nil {
			return err
		}
		err = tx.Migrator().AlterColumn(&system.SysDictData{}, "update_by")
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
