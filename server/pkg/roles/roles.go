package roles

import (
	"github.com/tanganyu1114/ansible-role-manager/config"
	"os"
	"path/filepath"
	"sync"
)

var (
	singletonRolesIns        Roles
	onceForSingletonRolesIns = sync.Once{}
)

type Roles interface {
	ImportRoleData(roleName string, compressedData []byte) error
	ExportRoleData(roleName string) ([]byte, error)
	RemoveRole(roleName string) error
}

type roles struct {
	workspace string
	archiver  Archiver
}

func newRoles(workspace string, archiver Archiver) Roles {
	r := &roles{
		workspace: workspace,
		archiver:  archiver,
	}
	return Roles(r)
}

func GetSingletonRolesIns() Roles {
	onceForSingletonRolesIns.Do(func() {
		if singletonRolesIns == nil {
			wsDir := filepath.Join(config.ExtConfig.Ansible.BaseDir, config.ExtConfig.Ansible.RoleDir)
			arc := newArchiver()
			singletonRolesIns = newRoles(wsDir, arc)
		}
	})
	return singletonRolesIns
}

func (r roles) ImportRoleData(roleName string, compressedData []byte) error {
	exDir := filepath.Join(r.workspace, roleName)
	/*

		TODO: 判断该role是否已存在

		info, err := os.Stat(exDir)
		if err != nil {

		}*/

	// TODO: 递归最多5层，判断解压后的role文件是否在子级目录。将role文件从子级目录移动到role文件目录层，并删除子级目录
	return r.archiver.Decompress(exDir, compressedData)
}

func (r roles) ExportRoleData(roleName string) ([]byte, error) {
	exDir := filepath.Join(r.workspace, roleName)
	return r.archiver.Compress(exDir)
}

func (r roles) RemoveRole(roleName string) error {
	// TODO: 防止删除roles上级目录及其他同级role
	err := os.Chdir(r.workspace)
	if err != nil {
		return err
	}
	return os.RemoveAll(roleName)
}
