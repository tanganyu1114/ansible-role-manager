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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1606579402398Test)
}

func _1606579402398Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		dicData := DictData{}
		err := tx.Model(&dicData).Where("dict_code = ?", 1).Update("dict_value", 2).Error
		if err != nil {
			return err
		}

		dicType := DictType{}
		err = tx.Model(&dicType).Where("status = ?", 0).Update("status", 2).Error
		if err != nil {
			return err
		}

		user := system.SysUser{}
		err = tx.Model(&user).Where("user_id = ?", 1).Update("status", 2).Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
