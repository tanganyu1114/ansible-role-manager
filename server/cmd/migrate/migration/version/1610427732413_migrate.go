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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _1610427732413Test)
}

func _1610427732413Test(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {

		// TODO: 这里开始写入要变更的内容
		user := system.SysUser{}
		err := tx.Model(&user).Where("status = ?", 0).Update("status", 2).Error
		if err != nil {
			return err
		}

		return tx.Create(&common.Migration{
			Version: version,
		}).Error
	})
}
