package roles

import svc "github.com/tanganyu1114/ansible-role-manager/pkg/roles"

type Roles interface {
	ImportRoleData(roleName string, compressedData []byte) error
	ExportRoleData(roleName string) ([]byte, error)
	RemoveRole(roleName string) error
}

type roles struct {
}

func newRoles() Roles {
	r := new(roles)
	return Roles(r)
}

func (r roles) ImportRoleData(roleName string, compressedData []byte) error {
	bo := svc.GetSingletonRolesIns()
	return bo.ImportRoleData(roleName, compressedData)
}

func (r roles) ExportRoleData(roleName string) ([]byte, error) {
	bo := svc.GetSingletonRolesIns()
	return bo.ExportRoleData(roleName)
}

func (r roles) RemoveRole(roleName string) error {
	bo := svc.GetSingletonRolesIns()
	return bo.RemoveRole(roleName)
}
